// External importers: Google Slides, Google Sheets, Canva, PPTX upload.
package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// --- Google Slides Import ---

func (h *SentinelHandler) PreviewGoogleSlides(c *gin.Context) {
	var req domain.PreviewGoogleSlidesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	result, err := h.usecase.PreviewGoogleSlides(&req, h.googleAPIKey)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") {
			statusCode = http.StatusBadRequest
		}
		if contains(err.Error(), "API error 403") || contains(err.Error(), "API error 401") {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": "Preview failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *SentinelHandler) ImportGoogleSlides(c *gin.Context) {
	creatorID := getUserIDFromContext(c)
	if creatorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.ImportGoogleSlidesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	result, err := h.usecase.ImportFromGoogleSlides(&req, h.googleAPIKey, creatorID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") {
			statusCode = http.StatusBadRequest
		}
		if contains(err.Error(), "API error 403") || contains(err.Error(), "API error 401") {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": "Import failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Imported %d tasks from \"%s\"", result.CreatedCount, result.PresentationTitle),
		"data":    result,
	})
}

// --- Google Sheets Import ---

func (h *SentinelHandler) PreviewGoogleSheets(c *gin.Context) {
	var req domain.PreviewGoogleSheetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	result, err := h.usecase.PreviewGoogleSheets(&req)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") || contains(err.Error(), "no importable") || contains(err.Error(), "no data rows") {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": "Preview failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *SentinelHandler) ImportGoogleSheets(c *gin.Context) {
	creatorID := getUserIDFromContext(c)
	if creatorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.ImportGoogleSheetsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	result, err := h.usecase.ImportFromGoogleSheets(&req, creatorID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") || contains(err.Error(), "at least one row") || contains(err.Error(), "parent task") || contains(err.Error(), "nested sub-task") {
			statusCode = http.StatusBadRequest
		}
		if domain.IsBadRequest(err) {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": "Import failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Imported %d tasks from \"%s\"", result.CreatedCount, result.SheetTitle),
		"data":    result,
	})
}

// --- Canva Import ---

func (h *SentinelHandler) PreviewCanva(c *gin.Context) {
	var req domain.PreviewCanvaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}
	result, err := h.usecase.PreviewCanva(&req, h.canvaAccessToken)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") || contains(err.Error(), "empty") || contains(err.Error(), "could not find") {
			statusCode = http.StatusBadRequest
		}
		if contains(err.Error(), "HTTP 401") || contains(err.Error(), "HTTP 403") {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": "Preview failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *SentinelHandler) ImportCanva(c *gin.Context) {
	creatorID := getUserIDFromContext(c)
	if creatorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req domain.ImportCanvaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request", "message": err.Error()})
		return
	}

	result, err := h.usecase.ImportFromCanva(&req, h.canvaAccessToken, creatorID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if contains(err.Error(), "invalid") || contains(err.Error(), "required") || contains(err.Error(), "at least one") || contains(err.Error(), "parent task") || contains(err.Error(), "nested sub-task") || contains(err.Error(), "no matching") {
			statusCode = http.StatusBadRequest
		}
		if domain.IsBadRequest(err) {
			statusCode = http.StatusBadRequest
		}
		if contains(err.Error(), "HTTP 401") || contains(err.Error(), "HTTP 403") {
			statusCode = http.StatusForbidden
		}
		c.JSON(statusCode, gin.H{"error": "Import failed", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Imported %d tasks from \"%s\"", result.CreatedCount, result.DesignTitle),
		"data":    result,
	})
}

// --- PPTX File Upload Import ---

// maxPPTXUploadBytes: Canva / rich decks with images often exceed 50 MB.
const maxPPTXUploadBytes = 250 * 1024 * 1024 // 250 MiB

func (h *SentinelHandler) PreviewPPTX(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file", "message": "multipart field 'file' is required"})
		return
	}
	if fileHeader.Size > maxPPTXUploadBytes {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File too large",
			"message": fmt.Sprintf("PPTX file must be under %d MB (got %.1f MB)", maxPPTXUploadBytes/(1024*1024), float64(fileHeader.Size)/(1024*1024)),
		})
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file", "message": err.Error()})
		return
	}
	defer f.Close()
	data, err := readAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file", "message": err.Error()})
		return
	}
	result, err := h.usecase.PreviewPPTX(data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid PPTX", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": result})
}

func (h *SentinelHandler) ImportPPTX(c *gin.Context) {
	creatorID := getUserIDFromContext(c)
	if creatorID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file", "message": "multipart field 'file' is required"})
		return
	}
	if fileHeader.Size > maxPPTXUploadBytes {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "File too large",
			"message": fmt.Sprintf("PPTX file must be under %d MB (got %.1f MB)", maxPPTXUploadBytes/(1024*1024), float64(fileHeader.Size)/(1024*1024)),
		})
		return
	}
	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file", "message": err.Error()})
		return
	}
	defer f.Close()
	data, err := readAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file", "message": err.Error()})
		return
	}
	payloadStr := c.PostForm("payload")
	if payloadStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing payload", "message": "form field 'payload' (JSON) is required"})
		return
	}
	var req domain.ImportPPTXRequest
	if err := json.Unmarshal([]byte(payloadStr), &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload", "message": err.Error()})
		return
	}
	result, err := h.usecase.ImportFromPPTX(data, &req, creatorID)
	if err != nil {
		statusCode := http.StatusInternalServerError
		if domain.IsBadRequest(err) || contains(err.Error(), "invalid") || contains(err.Error(), "required") || contains(err.Error(), "at least one") || contains(err.Error(), "parent task") || contains(err.Error(), "nested sub-task") || contains(err.Error(), "no matching") {
			statusCode = http.StatusBadRequest
		}
		c.JSON(statusCode, gin.H{"error": "Import failed", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": fmt.Sprintf("Imported %d tasks from \"%s\"", result.CreatedCount, result.Title),
		"data":    result,
	})
}
