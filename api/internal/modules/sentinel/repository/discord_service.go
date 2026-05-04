package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/portnd/the-sentinel-core/internal/modules/sentinel/domain"
)

// DiscordService handles sending notifications to Discord via webhook
type DiscordService struct {
	webhookURL string
}

// NewDiscordService creates a new Discord notification service
func NewDiscordService(webhookURL string) *DiscordService {
	return &DiscordService{webhookURL: webhookURL}
}

// DiscordEmbed represents a Discord embed object
type DiscordEmbed struct {
	Title       string               `json:"title,omitempty"`
	Description string               `json:"description,omitempty"`
	Color       int                  `json:"color,omitempty"`
	Fields      []DiscordEmbedField  `json:"fields,omitempty"`
	Footer      *DiscordEmbedFooter  `json:"footer,omitempty"`
	Timestamp   string               `json:"timestamp,omitempty"`
}

// DiscordEmbedField represents a field in a Discord embed
type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// DiscordEmbedFooter represents a footer in a Discord embed
type DiscordEmbedFooter struct {
	Text string `json:"text"`
}

// DiscordWebhookPayload represents the payload sent to Discord webhook
type DiscordWebhookPayload struct {
	Content string         `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

// Ensure DiscordService implements the domain.DiscordNotifier interface
var _ domain.DiscordNotifier = (*DiscordService)(nil)

// SendTimeLogNotification sends a notification when a user logs time
func (d *DiscordService) SendTimeLogNotification(entries []domain.TimeLogEntryForDiscord, date string) error {
	if d.webhookURL == "" {
		return nil // Skip if webhook not configured
	}

	// Group entries by user
	userEntries := make(map[string][]domain.TimeLogEntryForDiscord)
	for _, entry := range entries {
		userEntries[entry.DisplayName] = append(userEntries[entry.DisplayName], entry)
	}

	// Create embeds for each user
	var embeds []DiscordEmbed
	
	for displayName, userLogs := range userEntries {
		var fields []DiscordEmbedField
		
		for _, log := range userLogs {
			// Format task info
			taskInfo := ""
			if log.TaskCode != "" {
				taskInfo = log.TaskCode
				if log.TaskTitle != "" {
					taskInfo += " - " + log.TaskTitle
				}
			} else if log.TaskTitle != "" {
				taskInfo = log.TaskTitle
			}
			
			// Format hours
			hours := float64(log.Minutes) / 60.0
			hoursStr := fmt.Sprintf("%.1f ชั่วโมง", hours)
			if log.Progress > 0 {
				hoursStr += fmt.Sprintf(" (%d%%)", log.Progress)
			}
			
			// Format description
			desc := log.Description
			if desc == "" {
				desc = "-"
			}
			
			// Truncate long values for Discord limits
			if len(taskInfo) > 100 {
				taskInfo = taskInfo[:97] + "..."
			}
			if len(desc) > 100 {
				desc = desc[:97] + "..."
			}
			
			fields = append(fields, DiscordEmbedField{
				Name:  fmt.Sprintf("[%s] %s", log.WorkType, taskInfo),
				Value: fmt.Sprintf("%s\n%s", desc, hoursStr),
			})
		}
		
		// Calculate total hours for this user
		totalMinutes := 0
		for _, log := range userLogs {
			totalMinutes += log.Minutes
		}
		totalHours := float64(totalMinutes) / 60.0
		
		embeds = append(embeds, DiscordEmbed{
			Title:       fmt.Sprintf("📝 %s", displayName),
			Description: fmt.Sprintf("**วันที่:** %s\n**รวม:** %.1f ชั่วโมง", date, totalHours),
			Color:       0x5865F2, // Discord blurple
			Fields:      fields,
			Footer: &DiscordEmbedFooter{
				Text: "KOMGRIP Work Log",
			},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		})
	}

	payload := DiscordWebhookPayload{
		Content: "📢 **แจ้งเตือน - มีการลง log การทำงาน**",
		Embeds:  embeds,
	}

	return d.sendWebhook(payload)
}

// SendMissingLogNotification sends a notification listing users who haven't logged time
func (d *DiscordService) SendMissingLogNotification(users []domain.UserWithoutLogForDiscord, date string) error {
	if d.webhookURL == "" {
		return nil // Skip if webhook not configured
	}

	if len(users) == 0 {
		return nil // No users to report
	}

	var userLines []string
	for _, u := range users {
		name := u.DisplayName
		if name == "" {
			name = "ไม่ทราบชื่อ"
		}
		userLines = append(userLines, fmt.Sprintf("• %s (%.0f ชั่วโมง)", name, u.TotalHours))
	}

	// Split into chunks if too long (Discord has 4096 char limit for description)
	description := strings.Join(userLines, "\n")
	if len(description) > 4000 {
		description = description[:3997] + "..."
	}

	payload := DiscordWebhookPayload{
		Content: "⚠️ **ประกาศรายชื่อที่บันทึกข้อมูลประจำวันไม่ครบถ้วน**",
		Embeds: []DiscordEmbed{
			{
				Title:       "📋 รายชื่อผู้ที่ยังไม่ได้ลง log",
				Description: description,
				Color:       0xFF6B6B, // Red-ish
				Footer: &DiscordEmbedFooter{
					Text: fmt.Sprintf("วันที่: %s", date),
				},
				Timestamp: time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	return d.sendWebhook(payload)
}

// SendLeaveNotification sends a notification listing employees on leave today
func (d *DiscordService) SendLeaveNotification(leaves []domain.LeaveEntryForDiscord, date string) error {
	if d.webhookURL == "" {
		return nil // Skip if webhook not configured
	}

	if len(leaves) == 0 {
		return nil // No leaves to report
	}

	// Format leave type in Thai
	formatLeaveType := func(leaveType string) string {
		switch leaveType {
		case "ANNUAL":
			return "ลาพักร้อน"
		case "SICK":
			return "ลาป่วย"
		case "PERSONAL":
			return "ลากิจส่วนตัว"
		case "UNPAID":
			return "ลาไม่รับเงินเดือน"
		default:
			return leaveType
		}
	}

	var lines []string
	for _, l := range leaves {
		leaveTypeStr := formatLeaveType(l.LeaveType)
		if l.IsHalfDay {
			leaveTypeStr += " (ครึ่งวัน)"
		}
		lines = append(lines, fmt.Sprintf("• %s | %s", l.DisplayName, leaveTypeStr))
	}

	// Split into chunks if too long (Discord has 4096 char limit for description)
	description := strings.Join(lines, "\n")
	if len(description) > 4000 {
		description = description[:3997] + "..."
	}

	payload := DiscordWebhookPayload{
		Content: fmt.Sprintf("📋 **ประกาศรายชื่อแจ้งลาประจำวันที่ %s มีจำนวน %d ท่าน**", date, len(leaves)),
		Embeds: []DiscordEmbed{
			{
				Description: description,
				Color:       0xF39C12, // Orange
				Timestamp:   time.Now().UTC().Format(time.RFC3339),
			},
		},
	}

	return d.sendWebhook(payload)
}

// sendWebhook sends the payload to Discord
func (d *DiscordService) sendWebhook(payload DiscordWebhookPayload) error {
	if d.webhookURL == "" {
		return nil
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal Discord payload: %w", err)
	}

	req, err := http.NewRequest("POST", d.webhookURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return fmt.Errorf("failed to create Discord request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send Discord webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Discord webhook returned status %d", resp.StatusCode)
	}

	return nil
}

// IsEnabled returns true if the webhook is configured
func (d *DiscordService) IsEnabled() bool {
	return d.webhookURL != ""
}
