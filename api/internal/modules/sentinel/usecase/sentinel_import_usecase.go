package usecase

// Google Slides and Google Sheets import (PPTX pipeline, optional Slides/Drive APIs, sheet CSV).
// Canva and file-upload PPTX remain in canva_import.go and pptx_upload_import.go (same usecase package).

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
	"gorm.io/datatypes"
)
// --- Google Slides Import ---

// slideInfo is the unified internal representation of a parsed slide
type slideInfo struct {
	Index        int
	Title        string
	Body         string
	Notes        string
	ThumbnailURL string   // URL (from Slides API) or empty
	Images       []string // base64 data URLs (from PPTX media, no API key needed)
	SlideObjID   string   // populated when using Slides API
}

// ---- PPTX parsing (no API key required for public presentations) ----

type pptxPresentation struct {
	Title    string
	Slides   []pptxSlideData
	SlideIDs []string // object IDs from presentation.xml, same order as Slides (for export?format=jpeg&slide=id.XXX)
}

type pptxSlideData struct {
	Index  int
	Title  string
	Body   string
	Notes  string
	Images []string // base64 data URLs, largest-first
	Hidden bool     // true when slide has show="0" in presentation.xml (hidden/skipped in show)
}

type pptxImageEntry struct {
	DataURL string
	Size    int
}

var (
	spreadsheetIDRegex  = regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	spreadsheetGIDRegex = regexp.MustCompile(`[#?&]gid=(\d+)`)
	presentationIDRegex = regexp.MustCompile(`/presentation/d/([a-zA-Z0-9_-]+)`)
	pptxSlideNumRegex   = regexp.MustCompile(`slide(\d+)\.xml$`)
	pptxTitlePhRe       = regexp.MustCompile(`<p:ph[^>]*type="(?:title|ctrTitle)"`)
	pptxSystemPhRe      = regexp.MustCompile(`<p:ph[^>]*type="(?:dt|ftr|sldNum|hdr)"`)
	pptxNotesBodyPhRe   = regexp.MustCompile(`<p:ph(?:[^>]*idx="1"|[^>]*type="body")`)
	pptxATextRe         = regexp.MustCompile(`<a:t(?:\s[^>]*)?>([^<]*)</a:t>`)
	pptxRelationshipTagRe = regexp.MustCompile(`(?i)<Relationship\s+([^>]+)/\s*>`)
	pptxRelTypeDq         = regexp.MustCompile(`(?i)\bType\s*=\s*"([^"]*)"`)
	pptxRelTypeSq         = regexp.MustCompile(`(?i)\bType\s*=\s*'([^']*)'`)
	pptxRelTargetDq       = regexp.MustCompile(`(?i)\bTarget\s*=\s*"([^"]*)"`)
	pptxRelTargetSq       = regexp.MustCompile(`(?i)\bTarget\s*=\s*'([^']*)'`)
	pptxDocTitleRe    = regexp.MustCompile(`<dc:title>([^<]*)</dc:title>`)
	pptxSldIdRe       = regexp.MustCompile(`<p:sldId[^>]*id="(\d+)"`)
	pptxSldShowAttrRe = regexp.MustCompile(`<p:sld\s*([^>]+)>`) // opening tag of slide: check for show="0" or show="false" (hidden)
)

const (
	pptxMaxImageBytes     = 12 * 1024 * 1024 // cap per embedded file (design exports can be large)
	pptxMinImageBytes     = 256               // skip tiny icons/bullets
	pptxMaxImagesPerSlide = 30                // เก็บได้หลายรูปต่อหน้า (รูปฝัง + รูป export)
)

func pptxMimeType(ext string) string {
	switch strings.ToLower(ext) {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	case ".svg":
		return "image/svg+xml"
	default:
		return "" // EMF, WMF, etc. — skip (no lightweight decoder)
	}
}

func normalizeZipEntryName(name string) string {
	return path.Clean(strings.ReplaceAll(strings.TrimSpace(name), "\\", "/"))
}

// pptxRelAttr reads Type or Target from a Relationship attribute fragment; order-independent.
func pptxRelAttr(attrs, key string) string {
	switch strings.ToLower(strings.TrimSpace(key)) {
	case "type":
		if m := pptxRelTypeDq.FindStringSubmatch(attrs); len(m) > 1 {
			return m[1]
		}
		if m := pptxRelTypeSq.FindStringSubmatch(attrs); len(m) > 1 {
			return m[1]
		}
	case "target":
		if m := pptxRelTargetDq.FindStringSubmatch(attrs); len(m) > 1 {
			return m[1]
		}
		if m := pptxRelTargetSq.FindStringSubmatch(attrs); len(m) > 1 {
			return m[1]
		}
	}
	return ""
}

func isPPTXImageRelationshipType(typ string) bool {
	t := strings.ToLower(strings.TrimSpace(typ))
	return strings.Contains(t, "relationships/image") || strings.HasSuffix(t, "/image")
}

// pptxFirstRelationshipTargetByTypeContains returns Target for the first Relationship whose Type contains hint (e.g. notesSlide, slideLayout).
func pptxFirstRelationshipTargetByTypeContains(relsData []byte, typeHint string) string {
	h := strings.ToLower(typeHint)
	for _, m := range pptxRelationshipTagRe.FindAllStringSubmatch(string(relsData), -1) {
		if len(m) < 2 {
			continue
		}
		typ := pptxRelAttr(m[1], "Type")
		if typ == "" || !strings.Contains(strings.ToLower(typ), h) {
			continue
		}
		if tgt := pptxRelAttr(m[1], "Target"); tgt != "" {
			return tgt
		}
	}
	return ""
}

// extractImagesFromRels reads a _rels file and extracts all image targets from relDir.
// relDir is the directory containing the .rels file (e.g. ppt/slides for slide1.xml.rels).
func extractImagesFromRels(zr *zip.Reader, relDir string, relsData []byte) []pptxImageEntry {
	seen := make(map[string]bool)
	var entries []pptxImageEntry

	for _, m := range pptxRelationshipTagRe.FindAllStringSubmatch(string(relsData), -1) {
		if len(m) < 2 {
			continue
		}
		typ := pptxRelAttr(m[1], "Type")
		target := pptxRelAttr(m[1], "Target")
		if target == "" || !isPPTXImageRelationshipType(typ) {
			continue
		}
		mediaPath := path.Clean(relDir + "/" + target)
		if seen[mediaPath] {
			continue
		}
		seen[mediaPath] = true

		ext := path.Ext(mediaPath)
		mime := pptxMimeType(ext)
		if mime == "" {
			continue
		}

		imgData := readZipEntry(zr, mediaPath)
		if len(imgData) < pptxMinImageBytes || len(imgData) > pptxMaxImageBytes {
			continue
		}

		b64 := base64.StdEncoding.EncodeToString(imgData)
		dataURL := fmt.Sprintf("data:%s;base64,%s", mime, b64)
		entries = append(entries, pptxImageEntry{DataURL: dataURL, Size: len(imgData)})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Size > entries[j].Size
	})
	return entries
}

// extractSlideImages extracts embedded images from a PPTX slide via its _rels file.
// If the slide has no images, tries the slide layout's _rels (background/placeholders).
// Returns base64 data URLs sorted by image size (largest first).
func extractSlideImages(zr *zip.Reader, slideFile string) []string {
	dir := path.Dir(slideFile)
	base := path.Base(slideFile)
	relsPath := dir + "/_rels/" + base + ".rels"

	relsData := readZipEntry(zr, relsPath)
	if relsData == nil {
		return nil
	}

	entries := extractImagesFromRels(zr, dir, relsData)

	// Fallback: if slide has no images, try the slide layout (title slides often have only layout art)
	if len(entries) == 0 {
		if layoutTarget := pptxFirstRelationshipTargetByTypeContains(relsData, "slideLayout"); layoutTarget != "" {
			layoutPath := path.Clean(dir + "/" + layoutTarget)
			layoutDir := path.Dir(layoutPath)
			layoutBase := path.Base(layoutPath)
			layoutRelsPath := layoutDir + "/_rels/" + layoutBase + ".rels"
			if layoutRelsData := readZipEntry(zr, layoutRelsPath); layoutRelsData != nil {
				entries = extractImagesFromRels(zr, layoutDir, layoutRelsData)
			}
		}
	}

	result := make([]string, 0, len(entries))
	for i, e := range entries {
		if i >= pptxMaxImagesPerSlide {
			break
		}
		result = append(result, e.DataURL)
	}
	return result
}

// fetchSlideExportImage gets the full-slide image from Google's export (no API key).
// slideID is the object ID from presentation.xml (e.g. "256"). Returns base64 data URL or empty.
func fetchSlideExportImage(presentationID, slideID string) string {
	if presentationID == "" || slideID == "" {
		return ""
	}
	// Google: /export?format=jpeg&slide=id.[id] — try both "id.256" and "256"
	for _, slideParam := range []string{"id." + slideID, slideID} {
		exportURL := fmt.Sprintf(
			"https://docs.google.com/presentation/d/%s/export?format=jpeg&slide=%s",
			presentationID, slideParam,
		)
		client := &http.Client{
			Timeout: 8 * time.Second,
			CheckRedirect: func(req *http.Request, _ []*http.Request) error {
				if strings.Contains(req.URL.Host, "accounts.google.com") {
					return fmt.Errorf("auth required")
				}
				return nil
			},
		}
		resp, err := client.Get(exportURL) //nolint:noctx
		if err != nil {
			continue
		}
		data, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil || resp.StatusCode != http.StatusOK {
			continue
		}
		// Must look like image (JPEG magic bytes or reasonable size)
		if len(data) < 500 {
			continue
		}
		b64 := base64.StdEncoding.EncodeToString(data)
		return fmt.Sprintf("data:image/jpeg;base64,%s", b64)
	}
	return ""
}

func extractPresentationID(rawURL string) (string, error) {
	m := presentationIDRegex.FindStringSubmatch(rawURL)
	if len(m) < 2 {
		return "", errors.New("invalid Google Slides URL: could not extract presentation ID")
	}
	return m[1], nil
}

var xmlEntities = strings.NewReplacer(
	"&amp;", "&", "&lt;", "<", "&gt;", ">", "&apos;", "'", "&quot;", `"`,
)

func pptxUnescapeXML(s string) string { return xmlEntities.Replace(s) }

// readZipEntry reads a zip member by logical path. Matching is case-insensitive so PPTX from
// macOS/Canva (mixed-case paths like PPT/Media/...) still resolves on Linux/Docker.
func readZipEntry(r *zip.Reader, name string) []byte {
	want := normalizeZipEntryName(name)
	var match *zip.File
	for _, f := range r.File {
		if strings.EqualFold(normalizeZipEntryName(f.Name), want) {
			match = f
			break
		}
	}
	if match == nil {
		return nil
	}
	rc, err := match.Open()
	if err != nil {
		return nil
	}
	defer rc.Close()
	data, err := io.ReadAll(rc)
	if err != nil {
		return nil
	}
	return data
}

func pptxSlideFiles(r *zip.Reader) []string {
	var names []string
	for _, f := range r.File {
		if strings.HasPrefix(f.Name, "ppt/slides/slide") && strings.HasSuffix(f.Name, ".xml") &&
			!strings.Contains(f.Name[len("ppt/slides/"):], "/") {
			names = append(names, f.Name)
		}
	}
	sort.Slice(names, func(i, j int) bool {
		ni, nj := 0, 0
		if m := pptxSlideNumRegex.FindStringSubmatch(names[i]); len(m) > 1 {
			ni, _ = strconv.Atoi(m[1])
		}
		if m := pptxSlideNumRegex.FindStringSubmatch(names[j]); len(m) > 1 {
			nj, _ = strconv.Atoi(m[1])
		}
		return ni < nj
	})
	return names
}

// pptxSlideIsHidden returns true if the slide XML has show="0" or show="false" on the root <p:sld> (OOXML: hidden/skip in show).
func pptxSlideIsHidden(slideXML []byte) bool {
	m := pptxSldShowAttrRe.FindSubmatch(slideXML)
	if len(m) < 2 {
		return false
	}
	attr := string(m[1])
	return strings.Contains(attr, `show="0"`) || strings.Contains(attr, `show='0'`) ||
		strings.Contains(attr, `show="false"`) || strings.Contains(attr, `show='false'`)
}

// splitPPTXShapes extracts individual <p:sp>...</p:sp> blocks from slide XML.
func splitPPTXShapes(content string) []string {
	var shapes []string
	rest := content
	for {
		start := strings.Index(rest, "<p:sp")
		if start == -1 {
			break
		}
		// Verify it's really a <p:sp> tag (not <p:spTree> etc.)
		if len(rest) > start+5 {
			c := rest[start+5]
			if c != ' ' && c != '>' && c != '\t' && c != '\n' && c != '\r' {
				rest = rest[start+5:]
				continue
			}
		}
		end := strings.Index(rest[start:], "</p:sp>")
		if end == -1 {
			break
		}
		shapes = append(shapes, rest[start:start+end+7])
		rest = rest[start+end+7:]
	}
	return shapes
}

// parsePPTXSlideText extracts title and body text from a PPTX slide XML.
func parsePPTXSlideText(data []byte) (title, body string) {
	content := string(data)
	shapes := splitPPTXShapes(content)

	var bodyParts []string
	for _, shape := range shapes {
		// Skip system placeholders (date, footer, slide number)
		if pptxSystemPhRe.MatchString(shape) {
			continue
		}
		matches := pptxATextRe.FindAllStringSubmatch(shape, -1)
		var texts []string
		for _, m := range matches {
			t := strings.TrimSpace(pptxUnescapeXML(m[1]))
			if t != "" {
				texts = append(texts, t)
			}
		}
		if len(texts) == 0 {
			continue
		}
		shapeText := strings.Join(texts, "")
		if pptxTitlePhRe.MatchString(shape) && title == "" {
			title = strings.TrimSpace(shapeText)
		} else {
			bodyParts = append(bodyParts, strings.TrimSpace(shapeText))
		}
	}
	body = strings.Join(bodyParts, "\n")
	return
}

// parsePPTXNotesText extracts speaker notes from a notesSlide XML.
func parsePPTXNotesText(data []byte) string {
	content := string(data)
	shapes := splitPPTXShapes(content)
	for _, shape := range shapes {
		if pptxSystemPhRe.MatchString(shape) {
			continue
		}
		// Notes body is typically idx="1" or type="body"
		if !pptxNotesBodyPhRe.MatchString(shape) {
			continue
		}
		matches := pptxATextRe.FindAllStringSubmatch(shape, -1)
		var texts []string
		for _, m := range matches {
			t := strings.TrimSpace(pptxUnescapeXML(m[1]))
			if t != "" {
				texts = append(texts, t)
			}
		}
		if len(texts) > 0 {
			return strings.TrimSpace(strings.Join(texts, "\n"))
		}
	}
	return ""
}

// findNotesFileInZip resolves the notes slide file for a given slide via its _rels file.
func findNotesFileInZip(r *zip.Reader, slideFile string) string {
	// e.g. slideFile = "ppt/slides/slide1.xml"
	// rels  = "ppt/slides/_rels/slide1.xml.rels"
	dir := path.Dir(slideFile)
	base := path.Base(slideFile)
	relsPath := dir + "/_rels/" + base + ".rels"

	relsData := readZipEntry(r, relsPath)
	if relsData == nil {
		return ""
	}
	notesTarget := pptxFirstRelationshipTargetByTypeContains(relsData, "notesSlide")
	if notesTarget == "" {
		return ""
	}
	// Target is relative to the slide directory, e.g. "../notesSlides/notesSlide1.xml"
	return path.Clean(dir + "/" + notesTarget)
}

func downloadAndParsePPTX(presentationID string) (*pptxPresentation, error) {
	exportURL := fmt.Sprintf(
		"https://docs.google.com/presentation/d/%s/export/pptx",
		presentationID,
	)
	client := &http.Client{
		CheckRedirect: func(req *http.Request, _ []*http.Request) error {
			if strings.Contains(req.URL.Host, "accounts.google.com") {
				return fmt.Errorf("presentation requires authentication — make sure it is shared as 'Anyone with the link can view'")
			}
			return nil
		},
	}
	resp, err := client.Get(exportURL) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("failed to download presentation: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("access denied (HTTP %d) — share the presentation as 'Anyone with the link can view'", resp.StatusCode)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download failed (HTTP %d)", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read PPTX data: %w", err)
	}

	return parsePPTX(data)
}

func parsePPTX(data []byte) (*pptxPresentation, error) {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, fmt.Errorf("invalid PPTX file: %w", err)
	}

	// Get presentation title from docProps/core.xml
	title := "Imported Presentation"
	if coreXML := readZipEntry(zr, "docProps/core.xml"); coreXML != nil {
		if m := pptxDocTitleRe.FindSubmatch(coreXML); len(m) > 1 {
			if t := strings.TrimSpace(pptxUnescapeXML(string(m[1]))); t != "" {
				title = t
			}
		}
	}

	// Slide object IDs from ppt/presentation.xml (same order as slide files)
	var slideIDs []string
	if presXML := readZipEntry(zr, "ppt/presentation.xml"); presXML != nil {
		for _, m := range pptxSldIdRe.FindAllStringSubmatch(string(presXML), -1) {
			if len(m) > 1 {
				slideIDs = append(slideIDs, m[1])
			}
		}
	}

	slideFiles := pptxSlideFiles(zr)
	var slides []pptxSlideData
	for i, sf := range slideFiles {
		slideXML := readZipEntry(zr, sf)
		if slideXML == nil {
			continue
		}
		// Hidden: OOXML stores show="0" or show="false" on <p:sld> in each slide XML (not in presentation.xml)
		hidden := pptxSlideIsHidden(slideXML)

		slideTitle, slideBody := parsePPTXSlideText(slideXML)

		var notes string
		notesFile := findNotesFileInZip(zr, sf)
		if notesFile != "" {
			if notesXML := readZipEntry(zr, notesFile); notesXML != nil {
				notes = parsePPTXNotesText(notesXML)
			}
		}

		images := extractSlideImages(zr, sf)

		slides = append(slides, pptxSlideData{
			Index:  i + 1,
			Title:  slideTitle,
			Body:   slideBody,
			Notes:  notes,
			Images: images,
			Hidden: hidden,
		})
	}

	return &pptxPresentation{Title: title, Slides: slides, SlideIDs: slideIDs}, nil
}

// ---- Google Slides REST API (optional — used when API key provided for thumbnails & comments) ----

type gSlidePresentation struct {
	PresentationID string   `json:"presentationId"`
	Title          string   `json:"title"`
	Slides         []gSlide `json:"slides"`
}

type gSlide struct {
	ObjectID        string           `json:"objectId"`
	PageElements    []gPageElement   `json:"pageElements"`
	SlideProperties gSlideProperties `json:"slideProperties"`
}

type gSlideProperties struct {
	NotesPage gNotesPage `json:"notesPage"`
}

type gNotesPage struct {
	PageElements []gPageElement `json:"pageElements"`
}

type gPageElement struct {
	ObjectID string  `json:"objectId"`
	Shape    *gShape `json:"shape"`
}

type gShape struct {
	ShapeType   string        `json:"shapeType"`
	Placeholder *gPlaceholder `json:"placeholder"`
	Text        *gTextContent `json:"text"`
}

type gPlaceholder struct {
	Type  string `json:"type"`
	Index int    `json:"index"`
}

type gTextContent struct {
	TextElements []gTextElement `json:"textElements"`
}

type gTextElement struct {
	TextRun *gTextRun `json:"textRun"`
}

type gTextRun struct {
	Content string `json:"content"`
}

type gThumbnail struct {
	ContentURL string `json:"contentUrl"`
}

type gDriveCommentsResponse struct {
	Comments []gDriveComment `json:"comments"`
}

type gDriveComment struct {
	ID       string        `json:"id"`
	Content  string        `json:"content"`
	Author   gDriveAuthor  `json:"author"`
	Resolved bool          `json:"resolved"`
	Replies  []gDriveReply `json:"replies"`
}

type gDriveAuthor struct {
	DisplayName string `json:"displayName"`
}

type gDriveReply struct {
	Content string       `json:"content"`
	Author  gDriveAuthor `json:"author"`
}

// googleAPIError is the JSON shape of Google API error responses.
type googleAPIError struct {
	Error struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Status  string `json:"status"`
	} `json:"error"`
}

func fetchGoogleSlidesAPI(presentationID, apiKey string) (*gSlidePresentation, error) {
	apiURL := fmt.Sprintf("https://slides.googleapis.com/v1/presentations/%s?key=%s", presentationID, apiKey)
	resp, err := http.Get(apiURL) //nolint:noctx
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		var apiErr googleAPIError
		if jsonErr := json.Unmarshal(body, &apiErr); jsonErr == nil && apiErr.Error.Message != "" {
			msg := apiErr.Error.Message
			// Google Slides API does not accept API keys; only OAuth2. Return a clear message.
			if strings.Contains(msg, "API keys are not supported") || strings.Contains(msg, "Expected OAuth2") {
				msg = "Google Slides API ไม่รองรับ API Key — ต้องใช้ OAuth2 (ล็อกอินด้วย Google). ตอนนี้ใช้โหมด PPTX ได้โดยไม่ต้องใส่ Key (ได้ text, notes, รูปฝังใน slide)."
			}
			return nil, fmt.Errorf("%s", msg)
		}
		return nil, fmt.Errorf("Google Slides API error %d: %s", resp.StatusCode, string(body))
	}
	var p gSlidePresentation
	if err := json.Unmarshal(body, &p); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &p, nil
}

func fetchSlideThumbnail(presentationID, pageObjectID, apiKey string) (string, error) {
	apiURL := fmt.Sprintf(
		"https://slides.googleapis.com/v1/presentations/%s/pages/%s/thumbnail?key=%s&thumbnailProperties.mimeType=PNG&thumbnailProperties.thumbnailSize=LARGE",
		presentationID, pageObjectID, apiKey,
	)
	resp, err := http.Get(apiURL) //nolint:noctx
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("thumbnail API error %d", resp.StatusCode)
	}
	var t gThumbnail
	if err := json.Unmarshal(body, &t); err != nil {
		return "", err
	}
	return t.ContentURL, nil
}

// downloadThumbnailAsDataURL fetches the thumbnail image from Google's contentUrl and returns a data URL
// so it can be stored in the task and displayed without expiry (contentUrl lifetime is ~30 min).
// The thumbnail includes the full rendered slide (shapes, lines, drawings).
// apiKey is optional; when set, it is appended to the URL so Google may allow the server-side download.
func downloadThumbnailAsDataURL(contentURL, apiKey string) (string, error) {
	if contentURL == "" {
		return "", fmt.Errorf("empty content URL")
	}
	downloadURL := contentURL
	if apiKey != "" {
		if strings.Contains(contentURL, "?") {
			downloadURL = contentURL + "&key=" + apiKey
		} else {
			downloadURL = contentURL + "?key=" + apiKey
		}
	}
	req, err := http.NewRequest(http.MethodGet, downloadURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "Sentinel-Slides-Import/1.0 (https://github.com/portnd/the-sentinel-core)")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[Slides Import] thumbnail download request failed: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		snippet := string(body)
		if len(snippet) > 300 {
			snippet = snippet[:300] + "..."
		}
		log.Printf("[Slides Import] thumbnail download failed: status=%d body=%q", resp.StatusCode, snippet)
		return "", fmt.Errorf("thumbnail download returned %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[Slides Import] thumbnail response read failed: %v", err)
		return "", err
	}
	b64 := base64.StdEncoding.EncodeToString(body)
	return "data:image/png;base64," + b64, nil
}

func fetchDriveComments(fileID, apiKey string) ([]gDriveComment, error) {
	apiURL := fmt.Sprintf(
		"https://www.googleapis.com/drive/v3/files/%s/comments?key=%s&fields=comments(id,content,anchor,author,resolved,replies)&pageSize=100",
		fileID, apiKey,
	)
	resp, err := http.Get(apiURL) //nolint:noctx
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Drive API error %d: %s", resp.StatusCode, string(body))
	}
	var result gDriveCommentsResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Comments, nil
}

func extractAPISlideTitle(slide gSlide) string {
	for _, el := range slide.PageElements {
		if el.Shape == nil || el.Shape.Text == nil {
			continue
		}
		if el.Shape.Placeholder != nil {
			pt := el.Shape.Placeholder.Type
			if pt == "CENTERED_TITLE" || pt == "TITLE" {
				var sb strings.Builder
				for _, te := range el.Shape.Text.TextElements {
					if te.TextRun != nil {
						sb.WriteString(te.TextRun.Content)
					}
				}
				if t := strings.TrimSpace(strings.ReplaceAll(sb.String(), "\n", " ")); t != "" {
					return t
				}
			}
		}
	}
	return ""
}

func extractAPISlideBody(slide gSlide) string {
	var parts []string
	for _, el := range slide.PageElements {
		if el.Shape == nil || el.Shape.Text == nil {
			continue
		}
		var sb strings.Builder
		for _, te := range el.Shape.Text.TextElements {
			if te.TextRun != nil {
				sb.WriteString(te.TextRun.Content)
			}
		}
		if t := strings.TrimSpace(sb.String()); t != "" {
			parts = append(parts, t)
		}
	}
	return strings.Join(parts, "\n")
}

func extractAPISpeakerNotes(slide gSlide) string {
	for _, el := range slide.SlideProperties.NotesPage.PageElements {
		if el.Shape == nil || el.Shape.Text == nil {
			continue
		}
		if el.Shape.Placeholder != nil && el.Shape.Placeholder.Type == "BODY" {
			var sb strings.Builder
			for _, te := range el.Shape.Text.TextElements {
				if te.TextRun != nil {
					sb.WriteString(te.TextRun.Content)
				}
			}
			if t := strings.TrimSpace(sb.String()); t != "" {
				return t
			}
		}
	}
	return ""
}

// ---- Preview: get slide list only (no thumbnails, no task creation) ----

func (u *sentinelUsecase) PreviewGoogleSlides(req *domain.PreviewGoogleSlidesRequest, serverAPIKey string) (*domain.PreviewGoogleSlidesResult, error) {
	presentationID, err := extractPresentationID(req.PresentationURL)
	if err != nil {
		return nil, err
	}
	apiKey := strings.TrimSpace(req.APIKey)
	if apiKey == "" {
		apiKey = serverAPIKey
	}

	title, items, importMode, apiKeyStatus, apiKeyErrMsg, err := getSlidesListOnly(presentationID, apiKey)
	if err != nil {
		return nil, err
	}
	alreadyImported, _ := u.repo.GetImportedSlideIndicesByPresentationID(presentationID)
	return &domain.PreviewGoogleSlidesResult{
		PresentationTitle:           title,
		PresentationID:               presentationID,
		Slides:                       items,
		AlreadyImportedSlideIndices:  alreadyImported,
		ImportMode:                   importMode,
		APIKeyStatus:                 apiKeyStatus,
		APIKeyError:                  apiKeyErrMsg,
	}, nil
}

// getSlidesListOnly returns presentation title, slide list, import mode, API key status, and optional error message. No thumbnails.
func getSlidesListOnly(presentationID, apiKey string) (title string, slides []domain.PreviewSlideItem, importMode, apiKeyStatus, apiKeyError string, err error) {
	var presentationTitle string
	apiKeyProvided := apiKey != ""

	pptxData, pptxErr := downloadAndParsePPTX(presentationID)
	if pptxErr != nil && !apiKeyProvided {
		return "", nil, "", "", "", fmt.Errorf("failed to download presentation: %w\nTip: ensure the presentation is shared as 'Anyone with the link can view'", pptxErr)
	}
	if pptxErr == nil {
		presentationTitle = pptxData.Title
		for _, s := range pptxData.Slides {
			t := s.Title
			if t == "" {
				t = fmt.Sprintf("Slide %d", s.Index)
			}
			slides = append(slides, domain.PreviewSlideItem{
				Index:              s.Index,
				Title:              t,
				SuggestedTaskTitle: suggestedTaskTitleFromSlideText(s.Body, s.Index),
				Hidden:             s.Hidden,
			})
		}
	}
	if apiKeyProvided {
		apiPresentation, apiErr := fetchGoogleSlidesAPI(presentationID, apiKey)
		if apiErr == nil {
			apiKeyStatus = "valid"
			if pptxErr != nil {
				importMode = "api_only"
				presentationTitle = apiPresentation.Title
				slides = nil
				for i, slide := range apiPresentation.Slides {
					t := extractAPISlideTitle(slide)
					if t == "" {
						t = fmt.Sprintf("Slide %d", i+1)
					}
					body := extractAPISlideBody(slide)
					slides = append(slides, domain.PreviewSlideItem{
						Index:              i + 1,
						Title:              t,
						SuggestedTaskTitle: suggestedTaskTitleFromSlideText(body, i+1),
						Hidden:             false,
					})
				}
			} else {
				importMode = "pptx_with_api"
			}
		} else {
			apiKeyStatus = "invalid"
			apiKeyError = apiErr.Error()
			if pptxErr != nil {
				return "", nil, "", "", "", fmt.Errorf("failed to download PPTX (%v) and Slides API failed: %w", pptxErr, apiErr)
			}
			importMode = "pptx_only"
		}
	} else {
		apiKeyStatus = "not_provided"
		importMode = "pptx_only"
	}
	if len(slides) == 0 {
		return "", nil, "", "", "", errors.New("no slides found in the presentation")
	}
	return presentationTitle, slides, importMode, apiKeyStatus, apiKeyError, nil
}

// ---- Main ImportFromGoogleSlides usecase ----

func (u *sentinelUsecase) ImportFromGoogleSlides(req *domain.ImportGoogleSlidesRequest, serverAPIKey string, creatorID uint) (*domain.ImportGoogleSlidesResult, error) {
	presentationID, err := extractPresentationID(req.PresentationURL)
	if err != nil {
		return nil, err
	}

	var sprintUUID *uuid.UUID
	if req.SprintID != "" {
		parsed, err := uuid.Parse(req.SprintID)
		if err != nil {
			return nil, fmt.Errorf("invalid sprint_id: %w", err)
		}
		sprintUUID = &parsed
	}
	var epicUUID *uuid.UUID
	if req.EpicID != "" {
		parsed, err := uuid.Parse(req.EpicID)
		if err != nil {
			return nil, fmt.Errorf("invalid epic_id: %w", err)
		}
		epicUUID = &parsed
	}
	projectUUID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: %w", err)
	}
	var parentUUID *uuid.UUID
	if req.ParentID != "" {
		parsed, err := uuid.Parse(req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		parentUUID = &parsed
	}

	// Build triage map: slide_index -> TriagedSlide for per-slide overrides.
	// Also derive SlideIndices from triage data when provided.
	triagedMap := make(map[int]domain.TriagedSlide)
	if len(req.Slides) > 0 {
		for _, ts := range req.Slides {
			triagedMap[ts.SlideIndex] = ts
		}
		// Derive SlideIndices from triage data so filtering works normally downstream.
		if len(req.SlideIndices) == 0 {
			for idx := range triagedMap {
				req.SlideIndices = append(req.SlideIndices, idx)
			}
		}
	}

	apiKey := strings.TrimSpace(req.APIKey)
	if apiKey == "" {
		apiKey = serverAPIKey
	}

	// Step 1: Get slide content.
	// Primary: PPTX export — works for any public "anyone with link" presentation, no API key needed.
	// If API key provided: also fetch thumbnails and Drive comments.
	var slides []slideInfo
	var presentationTitle string

	pptxData, pptxErr := downloadAndParsePPTX(presentationID)
	if pptxErr != nil && apiKey == "" {
		return nil, fmt.Errorf("failed to download presentation: %w\nTip: ensure the presentation is shared as 'Anyone with the link can view'", pptxErr)
	}

	if pptxErr == nil {
		// Successfully parsed PPTX
		presentationTitle = pptxData.Title
		for _, s := range pptxData.Slides {
			title := s.Title
			if title == "" {
				title = fmt.Sprintf("Slide %d", s.Index)
			}
			// ใช้เฉพาะรูปที่ฝังใน slide นั้น (จาก PPTX media) — ไม่เรียก export ต่อ slide เพราะมักได้รูป slide แรกซ้ำทุกหน้า
			slides = append(slides, slideInfo{
				Index:  s.Index,
				Title:  title,
				Body:   s.Body,
				Notes:  s.Notes,
				Images: s.Images,
			})
		}
	}

		// Step 2: If API key available, enhance with thumbnails (via Slides API) and comments (via Drive API).
	var allComments []domain.SlideComment
	if apiKey != "" {
		log.Printf("[Slides Import] API key present, calling Slides API for thumbnails/comments")
		apiPresentation, apiErr := fetchGoogleSlidesAPI(presentationID, apiKey)
		if apiErr == nil {
			log.Printf("[Slides Import] Slides API OK, slides=%d", len(apiPresentation.Slides))
			// If PPTX also failed, use the API data as the content source
			if pptxErr != nil {
				presentationTitle = apiPresentation.Title
				slides = nil
				for i, slide := range apiPresentation.Slides {
					title := extractAPISlideTitle(slide)
					if title == "" {
						title = fmt.Sprintf("Slide %d", i+1)
					}
					slides = append(slides, slideInfo{
						Index:      i + 1,
						Title:      title,
						Body:       extractAPISlideBody(slide),
						Notes:      extractAPISpeakerNotes(slide),
						SlideObjID: slide.ObjectID,
					})
				}
			} else {
				// Merge: fill in slideObjID from API data (same order)
				for i := range slides {
					if i < len(apiPresentation.Slides) {
						slides[i].SlideObjID = apiPresentation.Slides[i].ObjectID
					}
				}
			}

			// Fetch per-slide thumbnails and download as base64 so drawings/lines are persisted (contentUrl expires in ~30 min)
			withObjID := 0
			for i := range slides {
				if slides[i].SlideObjID != "" {
					withObjID++
				}
			}
			log.Printf("[Slides Import] fetching thumbnails: %d slides with SlideObjID", withObjID)
			for i := range slides {
				if slides[i].SlideObjID == "" {
					log.Printf("[Slides Import] slide %d: skip (no SlideObjID)", i+1)
					continue
				}
				url, err := fetchSlideThumbnail(presentationID, slides[i].SlideObjID, apiKey)
				if err != nil {
					log.Printf("[Slides Import] slide %d: thumbnail URL failed: %v", i+1, err)
					continue
				}
				slides[i].ThumbnailURL = url
				dataURL, dlErr := downloadThumbnailAsDataURL(url, apiKey)
				if dlErr != nil && apiKey != "" {
					dataURL, dlErr = downloadThumbnailAsDataURL(url, "") // fallback: CDN may not accept key param
				}
				if dlErr == nil {
					// Prepend so the first image shown is the full slide with shapes/lines (กรอบแดง, เส้นวาด)
					slides[i].Images = append([]string{dataURL}, slides[i].Images...)
					log.Printf("[Slides Import] slide %d: thumbnail OK (base64 prepended)", i+1)
				} else {
					log.Printf("[Slides Import] slide %d: thumbnail download failed (กรอบ/เส้น will be missing): %v", i+1, dlErr)
				}
			}

			// Fetch Drive comments (non-fatal)
			driveComments, _ := fetchDriveComments(presentationID, apiKey)
			for _, c := range driveComments {
				comment := domain.SlideComment{
					Content:  c.Content,
					Author:   c.Author.DisplayName,
					Resolved: c.Resolved,
				}
				for _, r := range c.Replies {
					comment.Content += fmt.Sprintf("\n  ↳ [%s]: %s", r.Author.DisplayName, r.Content)
				}
				allComments = append(allComments, comment)
			}
		} else if pptxErr != nil {
			// Both PPTX and API failed
			return nil, fmt.Errorf("failed to download PPTX (%v) and Slides API also failed: %w", pptxErr, apiErr)
		} else {
			log.Printf("[Slides Import] Slides API failed (thumbnails/comments skipped): %v", apiErr)
		}
		// If API key provided but API call failed, we still continue with PPTX data (no thumbnails/comments)
	} else {
		log.Printf("[Slides Import] no API key: thumbnails (กรอบ/เส้น) and Drive comments will not be fetched")
	}

	if len(slides) == 0 {
		return nil, errors.New("no slides found in the presentation")
	}

	// Step 3: Validate priority and story points
	priority := strings.ToUpper(strings.TrimSpace(req.Priority))
	if !map[string]bool{"CRITICAL": true, "HIGH": true, "MEDIUM": true, "LOW": true}[priority] {
		priority = "MEDIUM"
	}
	storyPoints := req.StoryPoints
	if storyPoints < 0 {
		storyPoints = 0
	}

	slug := "task"
	proj, err := u.repo.GetProjectByID(projectUUID, domain.CallerContext{Role: domain.RoleCEO})
	if err == nil && proj != nil {
		slug = slugify(proj.Name)
	}

	// Filter by selected slide indices if provided (1-based)
	if len(req.SlideIndices) > 0 {
		allowed := make(map[int]bool)
		for _, idx := range req.SlideIndices {
			allowed[idx] = true
		}
		filtered := slides[:0]
		for _, s := range slides {
			if allowed[s.Index] {
				filtered = append(filtered, s)
			}
		}
		slides = filtered
	}
	if len(slides) == 0 {
		return nil, errors.New("no slides selected to import")
	}

	// Next code numbers: use global max suffix so codes are unique across all projects (idx_tasks_code is global).
	maxSuffix, _ := u.repo.GetMaxTaskCodeSuffix(slug)

	// Step 4: Create one task per (filtered) slide
	var createdTasks []*domain.Task
	for i, slide := range slides {
		// Build description as HTML so images appear in Description (single place); no separate Slide Images section.
		var htmlParts []string
		if slide.Body != "" {
			htmlParts = append(htmlParts, "<p>"+html.EscapeString(slide.Body)+"</p>")
		}
		for _, imgSrc := range slide.Images {
			// Block-level <img> — not inside <p> (TipTap/ProseMirror block image cannot nest in paragraph).
			htmlParts = append(htmlParts, "<img src=\""+html.EscapeString(imgSrc)+"\" class=\"editor-image\" alt=\"Slide\" />")
		}
		if slide.Notes != "" {
			htmlParts = append(htmlParts, "<p><em>Speaker Notes:</em> "+html.EscapeString(slide.Notes)+"</p>")
		}
		description := strings.Join(htmlParts, "\n")
		if description == "" {
			description = "<p></p>"
		}

		var slideURL string
		if slide.SlideObjID != "" {
			// Prefer object ID so link stays correct if slides are reordered
			slideURL = fmt.Sprintf("https://docs.google.com/presentation/d/%s/edit#slide=id.%s", presentationID, slide.SlideObjID)
		} else {
			// Fallback: use 1-based slide index; Google rewrites to internal ID and navigates to that position
			slideURL = fmt.Sprintf("https://docs.google.com/presentation/d/%s/edit#slide=%d", presentationID, slide.Index)
		}

		// resource_urls: keep only metadata for "Open in Slides"; images are now in description
		resourceURLs := domain.SlideResourceURLs{
			ThumbnailURL:   "",   // no longer used; images in description
			Images:         nil,  // no duplicate; images in description
			SlideURL:       slideURL,
			Source:         "google_slides",
			SlideIndex:     slide.Index,
			PresentationID: presentationID,
			Comments:       allComments,
		}
		resourceURLsJSON, err := json.Marshal(resourceURLs)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal resource URLs for slide %d: %w", slide.Index, err)
		}

		code := fmt.Sprintf("%s-%03d", slug, maxSuffix+1+i)

		projectIDCopy := projectUUID

		// Apply per-slide triage overrides when provided; default task title from slide body text, not placeholder title.
		taskTitle := suggestedTaskTitleFromSlideText(slide.Body, slide.Index)
		taskPriority := priority
		taskEstimatedMinutes := 0
		var taskAssigneeID *uint
		if triage, ok := triagedMap[slide.Index]; ok {
			if strings.TrimSpace(triage.Title) != "" {
				taskTitle = strings.TrimSpace(triage.Title)
			}
			p := strings.ToUpper(strings.TrimSpace(triage.Priority))
			if map[string]bool{"CRITICAL": true, "HIGH": true, "MEDIUM": true, "LOW": true}[p] {
				taskPriority = p
			}
			if triage.EstimatedMinutes > 0 {
				taskEstimatedMinutes = triage.EstimatedMinutes
			}
			taskAssigneeID = triage.AssigneeID
		}

		task := &domain.Task{
			ID:               uuid.New(),
			Code:             code,
			Title:            taskTitle,
			Description:      description,
			TaskType:         string(domain.TaskTypeTask),
			CreatedBy:        &creatorID,
			Status:           "PENDING",
			Priority:         taskPriority,
			StoryPoints:      storyPoints,
			EstimatedMinutes: taskEstimatedMinutes,
			SprintID:         sprintUUID,
			EpicID:           epicUUID,
			ProjectID:        &projectIDCopy,
			ParentID:         parentUUID,
			ResourceURLs:     datatypes.JSON(resourceURLsJSON),
		}
		if taskAssigneeID != nil {
			task.AssignedTo = taskAssigneeID
		}

		if err := u.repo.CreateTask(task); err != nil {
			return nil, fmt.Errorf("failed to create task for slide %d: %w", slide.Index, err)
		}
		createdTasks = append(createdTasks, task)
	}

	return &domain.ImportGoogleSlidesResult{
		CreatedCount:      len(createdTasks),
		SlideCount:        len(slides),
		PresentationTitle: presentationTitle,
		Tasks:             createdTasks,
	}, nil
}

var thaiMonthAbbrevToMonth = map[string]time.Month{
	"ม.ค.": time.January, "ก.พ.": time.February, "มี.ค.": time.March, "เม.ย.": time.April,
	"พ.ค.": time.May, "มิ.ย.": time.June, "ก.ค.": time.July, "ส.ค.": time.August,
	"ก.ย.": time.September, "ต.ค.": time.October, "พ.ย.": time.November, "ธ.ค.": time.December,
}

var validSheetImportStatuses = map[string]bool{
	"PENDING": true, "IN_PROGRESS": true, "READY_FOR_TEST": true,
	"WAIT_FOR_DEPLOY": true, "READY_FOR_UAT": true, "COMPLETED": true, "CANCELLED": true, "BLOCKED": true,
}

func sheetCSVCell(row []string, col int) string {
	if col < 0 || col >= len(row) {
		return ""
	}
	return strings.TrimSpace(strings.ReplaceAll(row[col], "\u200b", ""))
}

func sheetTitleFromContentDisposition(cd string) string {
	cd = strings.TrimSpace(cd)
	if cd == "" {
		return ""
	}
	if m := regexp.MustCompile(`filename\*=UTF-8''([^;\s]+)`).FindStringSubmatch(cd); len(m) > 1 {
		if dec, err := url.PathUnescape(strings.Trim(m[1], `"`)); err == nil {
			return strings.TrimSuffix(strings.TrimSpace(dec), ".csv")
		}
		return strings.TrimSuffix(strings.TrimSpace(m[1]), ".csv")
	}
	if m := regexp.MustCompile(`filename="([^"]+)"`).FindStringSubmatch(cd); len(m) > 1 {
		return strings.TrimSuffix(m[1], ".csv")
	}
	return ""
}

func parseGoogleSheetURL(raw string) (sheetID, gid string, err error) {
	u := strings.TrimSpace(raw)
	if u == "" {
		return "", "", errors.New("empty sheet URL")
	}
	m := spreadsheetIDRegex.FindStringSubmatch(u)
	if len(m) < 2 {
		return "", "", errors.New("invalid Google Sheets URL: missing spreadsheet id")
	}
	sheetID = m[1]
	gid = "0"
	if gm := spreadsheetGIDRegex.FindStringSubmatch(u); len(gm) > 1 {
		gid = gm[1]
	}
	return sheetID, gid, nil
}

func fetchGoogleSheetCSVRecords(sheetID, gid string) (records [][]string, sheetTitle string, err error) {
	exportURL := fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s/export?format=csv&gid=%s", sheetID, gid)
	client := &http.Client{Timeout: 60 * time.Second}
	req, err := http.NewRequest(http.MethodGet, exportURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Sentinel/1.0)")
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download sheet CSV: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(io.LimitReader(resp.Body, 20<<20))
	if err != nil {
		return nil, "", err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("Google Sheets returned HTTP %d: ensure the spreadsheet is shared as \"Anyone with the link can view\"", resp.StatusCode)
	}
	sheetTitle = sheetTitleFromContentDisposition(resp.Header.Get("Content-Disposition"))
	r := csv.NewReader(bytes.NewReader(body))
	r.LazyQuotes = true
	r.FieldsPerRecord = -1
	records, err = r.ReadAll()
	if err != nil {
		return nil, sheetTitle, fmt.Errorf("invalid CSV: %w", err)
	}
	if len(records) > 0 && len(records[0]) > 0 {
		first := strings.TrimSpace(records[0][0])
		if strings.HasPrefix(first, "<!DOCTYPE") || strings.HasPrefix(first, "<html") {
			return nil, "", errors.New("sheet export returned HTML (wrong gid or spreadsheet not shared as \"Anyone with the link can view\")")
		}
	}
	return records, sheetTitle, nil
}

func parseThaiBuddhistShortDate(s string) (time.Time, bool) {
	parts := strings.Fields(strings.TrimSpace(s))
	if len(parts) < 3 {
		return time.Time{}, false
	}
	day, err1 := strconv.Atoi(parts[0])
	month, ok := thaiMonthAbbrevToMonth[parts[1]]
	yearBE, err2 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || !ok {
		return time.Time{}, false
	}
	if day < 1 || day > 31 {
		return time.Time{}, false
	}
	yearCE := yearBE - 543
	if yearCE < 1900 || yearCE > 2100 {
		return time.Time{}, false
	}
	return time.Date(yearCE, month, day, 0, 0, 0, 0, time.UTC), true
}

func parseSlashSheetDate(s string) (time.Time, bool) {
	parts := strings.Split(strings.TrimSpace(s), "/")
	if len(parts) != 3 {
		return time.Time{}, false
	}
	d, err1 := strconv.Atoi(parts[0])
	m, err2 := strconv.Atoi(parts[1])
	y, err3 := strconv.Atoi(parts[2])
	if err1 != nil || err2 != nil || err3 != nil || d < 1 || m < 1 || m > 12 {
		return time.Time{}, false
	}
	if y < 100 {
		y += 2000
	}
	if y < 1900 || y > 2100 {
		return time.Time{}, false
	}
	return time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.UTC), true
}

func parseSheetDueRaw(s string) string {
	s = strings.TrimSpace(strings.ReplaceAll(s, "\u200b", ""))
	if s == "" {
		return ""
	}
	if t, ok := parseThaiBuddhistShortDate(s); ok {
		return t.Format("2006-01-02")
	}
	if t, ok := parseSlashSheetDate(s); ok {
		return t.Format("2006-01-02")
	}
	return ""
}

func mapKGSheetStatus(raw string) string {
	raw = strings.TrimSpace(strings.ReplaceAll(raw, "\u200b", ""))
	if raw == "" {
		return "PENDING"
	}
	if idx := strings.IndexAny(raw, "\n\r"); idx >= 0 {
		raw = strings.TrimSpace(raw[:idx])
	}
	switch raw {
	case "แก้ไขแล้ว", "นำขึ้น Prod แล้ว":
		return "COMPLETED"
	case "ยังไม่แก้":
		return "PENDING"
	case "กำลังแก้ไข", "แก้แล้วแต่ไม่ถูกต้อง", "แก้ไขอีกครั้ง":
		return "IN_PROGRESS"
	case "ทดสอบอีกครั้ง":
		return "READY_FOR_TEST"
	case "รอนำขึ้น Prod":
		return "READY_FOR_UAT"
	default:
		return "PENDING"
	}
}

// sheetColumnMap holds the result of dynamic header detection.
type sheetColumnMap struct {
	layout           string // "iod_dynamic", "kg_ifp"
	dataStart        int    // first data row index
	colTitle         int    // column index for task title / ปัญหา / Detail (-1 = not found)
	colStatus        int    // column index for working Status (-1 = not found)
	colDue           int    // column index for due date (-1 = not found)
	colPriority      int    // column index for Priority (-1 = not found)
	colSection       int    // column index for Header / section (-1 = not found)
	colHeaderLink    int    // column index for Header Link / URL (-1 = not found)
	colRequestMethod int    // column index for Request Method (-1 = not found)
	colPayload       int    // column index for Payload (-1 = not found)
	colImage         int    // column index for Image (-1 = not found)
	notesCols        []int  // leftover columns to include in notes
}

// detectSheetLayout scans the first 3 rows looking for an IOD-style header row.
// It dynamically maps column positions so it works for any IOD sheet variant.
func detectSheetLayout(records [][]string) sheetColumnMap {
	limit := 3
	if len(records) < limit {
		limit = len(records)
	}
	for ri := 0; ri < limit; ri++ {
		row := records[ri]
		colTitle := -1
		colStatus := -1
		colDue := -1
		colPriority := -1
		colSection := -1
		colHeaderLink := -1
		colRequestMethod := -1
		colPayload := -1
		colImage := -1

		for ci, rawCell := range row {
			cell := strings.TrimSpace(strings.ReplaceAll(rawCell, "\u200b", ""))
			low := strings.ToLower(cell)
			switch {
			case strings.Contains(cell, "ปัญหา") || low == "detail":
				colTitle = ci
			case strings.Contains(low, "working") || low == "status":
				colStatus = ci
			case strings.Contains(cell, "วันที่กำหนด") || low == "due date" || low == "due_date":
				colDue = ci
			case strings.Contains(low, "priority") || strings.Contains(cell, "ความสำคัญ"):
				colPriority = ci
			case low == "header" || strings.Contains(cell, "สเตป") || strings.Contains(cell, "มอบหมาย"):
				colSection = ci
			case low == "header link" || low == "header_link" || low == "headerlink":
				colHeaderLink = ci
			case strings.Contains(low, "request method") || low == "method" || low == "request_method":
				colRequestMethod = ci
			case low == "payload":
				colPayload = ci
			case low == "image" || strings.Contains(low, "รูป") || strings.Contains(low, "ภาพ"):
				colImage = ci
			}
		}

		if colTitle >= 0 && colStatus >= 0 {
			knownCols := map[int]bool{
				colTitle: true, colStatus: true, colDue: true, colPriority: true,
				colSection: true, colHeaderLink: true, colRequestMethod: true,
				colPayload: true, colImage: true,
			}
			var notesCols []int
			for ci, rawCell := range row {
				if knownCols[ci] {
					continue
				}
				cell := strings.TrimSpace(strings.ReplaceAll(rawCell, "\u200b", ""))
				low := strings.ToLower(cell)
				if low == "no." || low == "" {
					continue
				}
				notesCols = append(notesCols, ci)
			}
			return sheetColumnMap{
				layout:           "iod_dynamic",
				dataStart:        ri + 1,
				colTitle:         colTitle,
				colStatus:        colStatus,
				colDue:           colDue,
				colPriority:      colPriority,
				colSection:       colSection,
				colHeaderLink:    colHeaderLink,
				colRequestMethod: colRequestMethod,
				colPayload:       colPayload,
				colImage:         colImage,
				notesCols:        notesCols,
			}
		}
	}
	return sheetColumnMap{layout: "kg_ifp"}
}

// extractBugTitle cleans up a Detail cell: strips leading URLs, joins remaining lines.
// If a section name is provided its first line is prepended as context.
func extractBugTitle(detail, section string) string {
	// Take only the first line of the section (it may contain login creds on following lines)
	sectionFirst := strings.TrimSpace(section)
	if idx := strings.IndexAny(sectionFirst, "\n\r"); idx >= 0 {
		sectionFirst = strings.TrimSpace(sectionFirst[:idx])
	}

	lines := strings.Split(detail, "\n")
	var textLines []string
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" || strings.HasPrefix(l, "http://") || strings.HasPrefix(l, "https://") {
			continue
		}
		textLines = append(textLines, l)
	}
	var title string
	if len(textLines) > 0 {
		title = strings.Join(textLines, " ")
	} else {
		// Only a URL was in the detail — use section or the raw first URL line
		if sectionFirst != "" {
			return sectionFirst
		}
		title = strings.TrimSpace(detail)
		if idx := strings.Index(title, "\n"); idx > 0 {
			title = title[:idx]
		}
	}
	if sectionFirst != "" && !strings.Contains(title, sectionFirst) {
		title = sectionFirst + ": " + title
	}
	if len([]rune(title)) > 220 {
		runes := []rune(title)
		title = string(runes[:220]) + "..."
	}
	return title
}

// mapIODWorkingStatus maps English / mixed working-status labels from IOD bug sheets.
func mapIODWorkingStatus(raw string) string {
	raw = strings.TrimSpace(strings.ReplaceAll(raw, "\u200b", ""))
	if raw == "" {
		return "PENDING"
	}
	if idx := strings.IndexAny(raw, "\n\r"); idx >= 0 {
		raw = strings.TrimSpace(raw[:idx])
	}
	low := strings.ToLower(raw)
	switch {
	case strings.Contains(low, "done") || strings.Contains(low, "completed"):
		return "COMPLETED"
	case strings.Contains(low, "in process") || strings.Contains(low, "in progress") || low == "processing":
		return "IN_PROGRESS"
	case strings.Contains(low, "pause"):
		return "BLOCKED"
	case low == "bug":
		return "READY_FOR_TEST"
	case strings.Contains(low, "uat") || strings.Contains(low, "prod"):
		return "READY_FOR_UAT"
	default:
		return "PENDING"
	}
}

func (u *sentinelUsecase) PreviewGoogleSheets(req *domain.PreviewGoogleSheetsRequest) (*domain.PreviewGoogleSheetsResult, error) {
	if req == nil || strings.TrimSpace(req.SheetURL) == "" {
		return nil, errors.New("sheet_url is required")
	}
	sheetID, gid, err := parseGoogleSheetURL(req.SheetURL)
	if err != nil {
		return nil, err
	}
	records, sheetTitle, err := fetchGoogleSheetCSVRecords(sheetID, gid)
	if err != nil {
		return nil, err
	}
	if len(records) < 2 {
		return nil, errors.New("sheet has no data rows (only header or empty)")
	}
	if sheetTitle == "" {
		sheetTitle = "Google Sheet"
	}
	colMap := detectSheetLayout(records)
	var rows []domain.SheetRowPreviewItem
	if colMap.layout == "iod_dynamic" {
		for i := colMap.dataStart; i < len(records); i++ {
			row := records[i]
			rawDetail := sheetCSVCell(row, colMap.colTitle)
			section := ""
			if colMap.colSection >= 0 {
				section = sheetCSVCell(row, colMap.colSection)
			}
			title := extractBugTitle(rawDetail, section)
			if title == "" {
				continue
			}

			statusRaw := sheetCSVCell(row, colMap.colStatus)

			dueStr := ""
			if colMap.colDue >= 0 {
				dueStr = parseSheetDueRaw(sheetCSVCell(row, colMap.colDue))
			}

			priority := "MEDIUM"
			if colMap.colPriority >= 0 {
				if p := strings.ToUpper(strings.TrimSpace(sheetCSVCell(row, colMap.colPriority))); validPriorities[p] {
					priority = p
				}
			}

			var noteParts []string
			for _, ci := range colMap.notesCols {
				if s := sheetCSVCell(row, ci); s != "" {
					noteParts = append(noteParts, s)
				}
			}
			notes := strings.Join(noteParts, " | ")

			// Collect IOD-specific extra fields
			headerLink := ""
			if colMap.colHeaderLink >= 0 {
				headerLink = sheetCSVCell(row, colMap.colHeaderLink)
			}
			// Extract ALL URLs embedded in the Detail cell.
			// Many IOD bug rows contain one or more page URLs mixed into the description text.
			var detailURLs []string
			for _, detailLine := range strings.Split(rawDetail, "\n") {
				detailLine = strings.TrimSpace(detailLine)
				if strings.HasPrefix(detailLine, "http://") || strings.HasPrefix(detailLine, "https://") {
					detailURLs = append(detailURLs, detailLine)
				}
			}
			// Use the first URL as headerLink when no explicit header link column value exists.
			if headerLink == "" && len(detailURLs) > 0 {
				headerLink = detailURLs[0]
				detailURLs = detailURLs[1:] // remaining URLs become DetailLinks
			}
			requestMethod := ""
			if colMap.colRequestMethod >= 0 {
				requestMethod = sheetCSVCell(row, colMap.colRequestMethod)
			}
			payload := ""
			if colMap.colPayload >= 0 {
				payload = sheetCSVCell(row, colMap.colPayload)
			}
			imageRef := ""
			if colMap.colImage >= 0 {
				imageRef = sheetCSVCell(row, colMap.colImage)
			}

			rawFirst := statusRaw
			if idx := strings.IndexAny(rawFirst, "\n\r"); idx >= 0 {
				rawFirst = strings.TrimSpace(rawFirst[:idx])
			}
			log.Printf("[IOD preview] row %d headerLink=%q detailURLs=%v", i+1, headerLink, detailURLs)
			rows = append(rows, domain.SheetRowPreviewItem{
				RowIndex:      i + 1,
				Title:         title,
				DueDate:       dueStr,
				Status:        mapIODWorkingStatus(statusRaw),
				RawStatus:     rawFirst,
				Notes:         notes,
				Header:        section,
				HeaderLink:    headerLink,
				RequestMethod: requestMethod,
				Payload:       payload,
				ImageRef:      imageRef,
				DetailLinks:   detailURLs,
			})
			_ = priority // exposed via triage in the frontend
		}
		if len(rows) == 0 {
			return nil, errors.New("no importable rows: title column is empty for all data rows")
		}
	} else {
		// kg_ifp legacy layout: col A=date, col B=title, col F=status, col K=notes
		for i := 1; i < len(records); i++ {
			row := records[i]
			title := sheetCSVCell(row, 1)
			if title == "" {
				continue
			}
			dueRaw := sheetCSVCell(row, 0)
			statusRaw := sheetCSVCell(row, 5)
			notes := sheetCSVCell(row, 10)
			dueStr := parseSheetDueRaw(dueRaw)
			rawFirst := statusRaw
			if idx := strings.IndexAny(rawFirst, "\n\r"); idx >= 0 {
				rawFirst = strings.TrimSpace(rawFirst[:idx])
			}
			rows = append(rows, domain.SheetRowPreviewItem{
				RowIndex:  i + 1,
				Title:     title,
				DueDate:   dueStr,
				Status:    mapKGSheetStatus(statusRaw),
				RawStatus: rawFirst,
				Notes:     notes,
			})
		}
		if len(rows) == 0 {
			return nil, errors.New("no importable rows: column B (รายละเอียด) is empty for all data rows")
		}
	}
	return &domain.PreviewGoogleSheetsResult{
		SheetTitle: sheetTitle,
		SheetID:    sheetID,
		Rows:       rows,
	}, nil
}

// buildIODTaskDescription creates a rich HTML description from a TriagedSheetRow.
// For IOD-style rows, it formats all available fields (Header, Detail, Header Link,
// Request Method, Payload, Image, Notes) as a structured bug report.
// For plain rows (only Notes), it falls back to a simple paragraph.
func buildIODTaskDescription(tr domain.TriagedSheetRow) string {
	hasIOD := tr.Header != "" || tr.HeaderLink != "" || tr.RequestMethod != "" || tr.Payload != "" || tr.ImageRef != "" || len(tr.DetailLinks) > 0
	if !hasIOD {
		n := strings.TrimSpace(tr.Notes)
		if n == "" {
			return ""
		}
		return "<p>" + html.EscapeString(n) + "</p>"
	}

	var b strings.Builder

	if h := strings.TrimSpace(tr.Header); h != "" {
		b.WriteString(`<p><strong>📌 Section / หน้า:</strong> `)
		b.WriteString(html.EscapeString(h))
		b.WriteString(`</p>`)
	}

	detail := strings.TrimSpace(tr.Title)
	if detail != "" {
		b.WriteString(`<p><strong>🐛 Bug Detail:</strong></p>`)
		b.WriteString(`<p>`)
		b.WriteString(strings.ReplaceAll(html.EscapeString(detail), "\n", "<br>"))
		b.WriteString(`</p>`)
	}

	if link := strings.TrimSpace(tr.HeaderLink); link != "" {
		b.WriteString(`<p><strong>🔗 URL: </strong>`)
		if strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://") {
			b.WriteString(`<a href="`)
			b.WriteString(html.EscapeString(link))
			b.WriteString(`" target="_blank" rel="noopener">`)
			b.WriteString(html.EscapeString(link))
			b.WriteString(`</a>`)
		} else {
			b.WriteString(html.EscapeString(link))
		}
		b.WriteString(`</p>`)
	}

	if len(tr.DetailLinks) > 0 {
		b.WriteString(`<p><strong>🔗 Related Links:</strong></p><ul>`)
		for _, dl := range tr.DetailLinks {
			dl = strings.TrimSpace(dl)
			if dl == "" {
				continue
			}
			b.WriteString(`<li><a href="`)
			b.WriteString(html.EscapeString(dl))
			b.WriteString(`" target="_blank" rel="noopener">`)
			b.WriteString(html.EscapeString(dl))
			b.WriteString(`</a></li>`)
		}
		b.WriteString(`</ul>`)
	}

	if m := strings.TrimSpace(tr.RequestMethod); m != "" {
		b.WriteString(`<p><strong>⚙️ Request Method:</strong> <code>`)
		b.WriteString(html.EscapeString(m))
		b.WriteString(`</code></p>`)
	}

	if p := strings.TrimSpace(tr.Payload); p != "" {
		b.WriteString(`<p><strong>📦 Payload / Response:</strong></p>`)
		b.WriteString(`<pre><code>`)
		b.WriteString(html.EscapeString(p))
		b.WriteString(`</code></pre>`)
	}

	if img := strings.TrimSpace(tr.ImageRef); img != "" {
		if strings.HasPrefix(img, "http://") || strings.HasPrefix(img, "https://") {
			b.WriteString(`<p><strong>🖼️ Screenshot:</strong></p>`)
			b.WriteString(`<img src="`)
			b.WriteString(html.EscapeString(img))
			b.WriteString(`" alt="screenshot">`)
		} else {
			b.WriteString(`<p><strong>🖼️ Image:</strong> `)
			b.WriteString(html.EscapeString(img))
			b.WriteString(`</p>`)
		}
	}

	if n := strings.TrimSpace(tr.Notes); n != "" {
		b.WriteString(`<p><strong>📝 Notes:</strong> `)
		b.WriteString(html.EscapeString(n))
		b.WriteString(`</p>`)
	}

	return b.String()
}

func (u *sentinelUsecase) ImportFromGoogleSheets(req *domain.ImportGoogleSheetsRequest, creatorID uint) (*domain.ImportGoogleSheetsResult, error) {
	if req == nil {
		return nil, errors.New("request is required")
	}
	if len(req.Rows) == 0 {
		return nil, errors.New("at least one row is required to import")
	}
	_, _, err := parseGoogleSheetURL(req.SheetURL)
	if err != nil {
		return nil, err
	}

	var sprintUUID *uuid.UUID
	if req.SprintID != "" {
		parsed, err := uuid.Parse(req.SprintID)
		if err != nil {
			return nil, fmt.Errorf("invalid sprint_id: %w", err)
		}
		sprintUUID = &parsed
	}
	var epicUUID *uuid.UUID
	if req.EpicID != "" {
		parsed, err := uuid.Parse(req.EpicID)
		if err != nil {
			return nil, fmt.Errorf("invalid epic_id: %w", err)
		}
		epicUUID = &parsed
	}
	projectUUID, err := uuid.Parse(req.ProjectID)
	if err != nil {
		return nil, fmt.Errorf("invalid project_id: %w", err)
	}
	var parentUUID *uuid.UUID
	if req.ParentID != "" {
		parsed, err := uuid.Parse(req.ParentID)
		if err != nil {
			return nil, fmt.Errorf("invalid parent_id: %w", err)
		}
		parentUUID = &parsed
	}

	if parentUUID != nil {
		parent, err := u.repo.GetTaskByID(*parentUUID)
		if err != nil || parent == nil {
			return nil, errors.New("parent task not found")
		}
		if parent.ParentID != nil {
			return nil, &domain.ErrBadRequest{Msg: "cannot attach sheet import under a nested sub-task"}
		}
		if parent.ProjectID == nil || *parent.ProjectID != projectUUID {
			return nil, &domain.ErrBadRequest{Msg: "parent task must belong to the same project"}
		}
	}

	slug := "task"
	if proj, err := u.repo.GetProjectByID(projectUUID, domain.CallerContext{Role: domain.RoleCEO}); err == nil && proj != nil {
		slug = slugify(proj.Name)
	}
	maxSuffix, err := u.repo.GetMaxTaskCodeSuffix(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get next task code: %w", err)
	}

	var created []*domain.Task
	for i, tr := range req.Rows {
		title := strings.TrimSpace(tr.Title)
		if title == "" {
			return nil, fmt.Errorf("row %d: title is required", tr.RowIndex)
		}
		estMins := tr.EstimatedMinutes
		if estMins < 0 {
			return nil, fmt.Errorf("row %d: estimated_minutes cannot be negative", tr.RowIndex)
		}
		priority := strings.ToUpper(strings.TrimSpace(tr.Priority))
		if priority == "" {
			priority = "MEDIUM"
		}
		if !validPriorities[priority] {
			return nil, fmt.Errorf("row %d: invalid priority %q", tr.RowIndex, tr.Priority)
		}
		st := strings.ToUpper(strings.TrimSpace(tr.Status))
		if st == "" {
			st = "PENDING"
		}
		if !validSheetImportStatuses[st] {
			return nil, fmt.Errorf("row %d: invalid status %q", tr.RowIndex, tr.Status)
		}

		log.Printf("[IOD import] row %d detail_links=%v", tr.RowIndex, tr.DetailLinks)
		desc := buildIODTaskDescription(tr)

		var duePtr *time.Time
		if ds := strings.TrimSpace(tr.DueDate); ds != "" {
			t, err := time.Parse("2006-01-02", ds)
			if err != nil {
				return nil, fmt.Errorf("row %d: invalid due_date %q (use YYYY-MM-DD)", tr.RowIndex, tr.DueDate)
			}
			duePtr = &t
		}

		projectIDCopy := projectUUID
		task := &domain.Task{
			ID:               uuid.New(),
			Code:             fmt.Sprintf("%s-%03d", slug, maxSuffix+1+i),
			Title:            title,
			Description:      desc,
			TaskType:         string(domain.TaskTypeBug),
			CreatedBy:        &creatorID,
			Status:           st,
			Priority:         priority,
			StoryPoints:      0,
			EstimatedMinutes: estMins,
			DueAt:            duePtr,
			SprintID:         sprintUUID,
			EpicID:           epicUUID,
			ProjectID:        &projectIDCopy,
			ParentID:         parentUUID,
		}

		if err := u.repo.CreateTask(task); err != nil {
			return nil, fmt.Errorf("failed to create task for sheet row %d: %w", tr.RowIndex, err)
		}
		created = append(created, task)
	}

	titleOut := strings.TrimSpace(req.SheetTitle)
	if titleOut == "" {
		titleOut = "Google Sheet"
	}
	return &domain.ImportGoogleSheetsResult{
		CreatedCount: len(created),
		SheetTitle:   titleOut,
		Tasks:        created,
	}, nil
}
