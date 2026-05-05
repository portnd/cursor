package repository

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// --- DiscordService Tests ---

func TestDiscordService_IsEnabled_WithWebhook(t *testing.T) {
	svc := NewDiscordService("https://discord.com/api/webhooks/test")
	if !svc.IsEnabled() {
		t.Error("Should be enabled when webhook URL is set")
	}
}

func TestDiscordService_IsEnabled_WithoutWebhook(t *testing.T) {
	svc := NewDiscordService("")
	if svc.IsEnabled() {
		t.Error("Should be disabled when webhook URL is empty")
	}
}

func TestDiscordService_SendMissingLogNotification_NoWebhook(t *testing.T) {
	svc := NewDiscordService("")
	err := svc.SendMissingLogNotification([]domain.UserWithoutLogForDiscord{
		{DisplayName: "Alice", TotalHours: 0},
	}, "2025-05-04")
	if err != nil {
		t.Errorf("Should return nil when webhook not configured, got: %v", err)
	}
}

func TestDiscordService_SendMissingLogNotification_EmptyUsers(t *testing.T) {
	svc := NewDiscordService("https://discord.com/api/webhooks/test")
	err := svc.SendMissingLogNotification([]domain.UserWithoutLogForDiscord{}, "2025-05-04")
	if err != nil {
		t.Errorf("Should return nil for empty users list, got: %v", err)
	}
}

func TestDiscordService_SendMissingLogNotification_Success(t *testing.T) {
	var receivedPayload DiscordWebhookPayload
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		if err := json.NewDecoder(r.Body).Decode(&receivedPayload); err != nil {
			t.Errorf("Failed to decode payload: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewDiscordService(server.URL)
	err := svc.SendMissingLogNotification([]domain.UserWithoutLogForDiscord{
		{DisplayName: "Alice", TotalHours: 0},
		{DisplayName: "Bob", TotalHours: 0.5},
	}, "2025-05-04")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if receivedPayload.Content == "" {
		t.Error("Expected non-empty Content")
	}
	if len(receivedPayload.Embeds) != 1 {
		t.Fatalf("Expected 1 embed, got %d", len(receivedPayload.Embeds))
	}
	if receivedPayload.Embeds[0].Title == "" {
		t.Error("Expected non-empty embed Title")
	}
	if receivedPayload.Embeds[0].Footer == nil || receivedPayload.Embeds[0].Footer.Text == "" {
		t.Error("Expected footer with date")
	}
}

func TestDiscordService_SendLeaveNotification_NoWebhook(t *testing.T) {
	svc := NewDiscordService("")
	err := svc.SendLeaveNotification([]domain.LeaveEntryForDiscord{
		{DisplayName: "Alice", LeaveType: "ANNUAL"},
	}, "2025-05-05")
	if err != nil {
		t.Errorf("Should return nil when webhook not configured, got: %v", err)
	}
}

func TestDiscordService_SendLeaveNotification_EmptyLeaves(t *testing.T) {
	svc := NewDiscordService("https://discord.com/api/webhooks/test")
	err := svc.SendLeaveNotification([]domain.LeaveEntryForDiscord{}, "2025-05-05")
	if err != nil {
		t.Errorf("Should return nil for empty leaves list, got: %v", err)
	}
}

func TestDiscordService_SendLeaveNotification_Success(t *testing.T) {
	var receivedPayload DiscordWebhookPayload
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if err := json.NewDecoder(r.Body).Decode(&receivedPayload); err != nil {
			t.Errorf("Failed to decode payload: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewDiscordService(server.URL)
	err := svc.SendLeaveNotification([]domain.LeaveEntryForDiscord{
		{DisplayName: "Alice", LeaveType: "ANNUAL", StartDate: "2025-05-05", EndDate: "2025-05-05", IsHalfDay: false},
		{DisplayName: "Bob", LeaveType: "SICK", StartDate: "2025-05-05", EndDate: "2025-05-06", IsHalfDay: true},
	}, "2025-05-05")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if receivedPayload.Content == "" {
		t.Error("Expected non-empty Content")
	}
	if len(receivedPayload.Embeds) != 1 {
		t.Fatalf("Expected 1 embed, got %d", len(receivedPayload.Embeds))
	}
}

func TestDiscordService_SendMissingStandupNotification_NoWebhook(t *testing.T) {
	svc := NewDiscordService("")
	err := svc.SendMissingStandupNotification([]domain.UserWithoutStandupForDiscord{
		{DisplayName: "Alice"},
	}, "2025-05-05")
	if err != nil {
		t.Errorf("Should return nil when webhook not configured, got: %v", err)
	}
}

func TestDiscordService_SendMissingStandupNotification_EmptyUsers(t *testing.T) {
	svc := NewDiscordService("https://discord.com/api/webhooks/test")
	err := svc.SendMissingStandupNotification([]domain.UserWithoutStandupForDiscord{}, "2025-05-05")
	if err != nil {
		t.Errorf("Should return nil for empty users list, got: %v", err)
	}
}

func TestDiscordService_SendMissingStandupNotification_Success(t *testing.T) {
	var receivedPayload DiscordWebhookPayload
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if err := json.NewDecoder(r.Body).Decode(&receivedPayload); err != nil {
			t.Errorf("Failed to decode payload: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewDiscordService(server.URL)
	err := svc.SendMissingStandupNotification([]domain.UserWithoutStandupForDiscord{
		{DisplayName: "Alice"},
		{DisplayName: "Bob"},
	}, "2025-05-05")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if receivedPayload.Content == "" {
		t.Error("Expected non-empty Content")
	}
	if len(receivedPayload.Embeds) != 1 {
		t.Fatalf("Expected 1 embed, got %d", len(receivedPayload.Embeds))
	}
}

func TestDiscordService_SendTimeLogNotification_NoWebhook(t *testing.T) {
	svc := NewDiscordService("")
	err := svc.SendTimeLogNotification([]domain.TimeLogEntryForDiscord{
		{DisplayName: "Alice", Minutes: 120, WorkType: "DEV", TaskCode: "T-001", TaskTitle: "Test task"},
	}, "2025-05-05")
	if err != nil {
		t.Errorf("Should return nil when webhook not configured, got: %v", err)
	}
}

func TestDiscordService_SendTimeLogNotification_Success(t *testing.T) {
	var receivedPayload DiscordWebhookPayload
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if err := json.NewDecoder(r.Body).Decode(&receivedPayload); err != nil {
			t.Errorf("Failed to decode payload: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewDiscordService(server.URL)
	err := svc.SendTimeLogNotification([]domain.TimeLogEntryForDiscord{
		{DisplayName: "Alice", Minutes: 120, WorkType: "DEV", TaskCode: "T-001", TaskTitle: "Test task", Description: "Did stuff", Progress: 50},
	}, "2025-05-05")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if receivedPayload.Content == "" {
		t.Error("Expected non-empty Content")
	}
	if len(receivedPayload.Embeds) == 0 {
		t.Fatal("Expected at least 1 embed")
	}
	if len(receivedPayload.Embeds[0].Fields) == 0 {
		t.Fatal("Expected at least 1 field in embed")
	}
}

func TestDiscordService_WebhookError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	svc := NewDiscordService(server.URL)
	err := svc.SendMissingLogNotification([]domain.UserWithoutLogForDiscord{
		{DisplayName: "Alice", TotalHours: 0},
	}, "2025-05-04")
	if err == nil {
		t.Error("Expected error when webhook returns 500")
	}
}

func TestDiscordService_SendMissingLogNotification_LongDescription(t *testing.T) {
	var receivedPayload DiscordWebhookPayload
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		_ = json.NewDecoder(r.Body).Decode(&receivedPayload)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewDiscordService(server.URL)

	// Create many users to exceed 4000 char limit
	var users []domain.UserWithoutLogForDiscord
	for i := 0; i < 200; i++ {
		users = append(users, domain.UserWithoutLogForDiscord{
			DisplayName: "VeryLongNameThatRepeatedlyExceedsLimits",
			TotalHours:  0,
		})
	}
	err := svc.SendMissingLogNotification(users, "2025-05-04")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if len(receivedPayload.Embeds[0].Description) > 4000 {
		t.Errorf("Description should be truncated to <= 4000 chars, got %d", len(receivedPayload.Embeds[0].Description))
	}
}

func TestDiscordService_SendMissingLogNotification_EmptyDisplayName(t *testing.T) {
	var receivedPayload DiscordWebhookPayload
	var mu sync.Mutex

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		_ = json.NewDecoder(r.Body).Decode(&receivedPayload)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	svc := NewDiscordService(server.URL)
	err := svc.SendMissingLogNotification([]domain.UserWithoutLogForDiscord{
		{DisplayName: "", TotalHours: 0},
	}, "2025-05-04")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	// Should use "ไม่ทราบชื่อ" as fallback
	if receivedPayload.Embeds[0].Description == "" {
		t.Error("Expected non-empty description even with empty DisplayName")
	}
}
