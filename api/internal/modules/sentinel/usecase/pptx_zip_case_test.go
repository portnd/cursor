package usecase

import (
	"archive/zip"
	"bytes"
	"strings"
	"testing"
)

func TestReadZipEntry_CaseInsensitivePath(t *testing.T) {
	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	w, err := zw.Create("PPT/Media/Photo.PNG")
	if err != nil {
		t.Fatal(err)
	}
	payload := bytes.Repeat([]byte{0x89, 0x50, 0x4e, 0x47}, 90) // > pptxMinImageBytes, fake png-ish
	if _, err := w.Write(payload); err != nil {
		t.Fatal(err)
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}

	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}
	got := readZipEntry(zr, "ppt/media/photo.png")
	if len(got) != len(payload) {
		t.Fatalf("readZipEntry: got len %d want %d", len(got), len(payload))
	}
}

func TestExtractImagesFromRels_TargetBeforeType(t *testing.T) {
	rels := []byte(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Target="../media/i.png" Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/image"/>
</Relationships>`)

	buf := new(bytes.Buffer)
	zw := zip.NewWriter(buf)
	add := func(name string, body []byte) {
		t.Helper()
		w, err := zw.Create(name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write(body); err != nil {
			t.Fatal(err)
		}
	}
	add("ppt/media/i.png", bytes.Repeat([]byte{0x89}, 300))
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatal(err)
	}

	entries := extractImagesFromRels(zr, "ppt/slides", rels)
	if len(entries) != 1 {
		t.Fatalf("entries: got %d want 1", len(entries))
	}
	if !strings.HasPrefix(entries[0].DataURL, "data:image/png;base64,") {
		t.Fatalf("unexpected data url prefix: %s", entries[0].DataURL[:40])
	}
}
