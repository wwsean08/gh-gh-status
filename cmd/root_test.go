package cmd

import (
	"strings"
	"testing"
	"time"

	"github.com/wwsean08/gh-gh-status/status"
)

func TestStripAnsiCodes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text without ANSI codes",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "text with simple ANSI color code",
			input:    "\x1b[32mGreen Text\x1b[0m",
			expected: "Green Text",
		},
		{
			name:     "text with multiple ANSI codes",
			input:    "\x1b[31mRed\x1b[0m and \x1b[32mGreen\x1b[0m",
			expected: "Red and Green",
		},
		{
			name:     "text with complex ANSI codes",
			input:    "\x1b[1;31;40mBold Red on Black\x1b[0m",
			expected: "Bold Red on Black",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only ANSI codes",
			input:    "\x1b[32m\x1b[0m",
			expected: "",
		},
		{
			name:     "ANSI code at start",
			input:    "\x1b[32mText",
			expected: "Text",
		},
		{
			name:     "ANSI code at end",
			input:    "Text\x1b[0m",
			expected: "Text",
		},
		{
			name:     "ANSI code in middle",
			input:    "Before\x1b[32mAfter",
			expected: "BeforeAfter",
		},
		{
			name:     "multiple consecutive ANSI codes",
			input:    "\x1b[1m\x1b[32m\x1b[40mText\x1b[0m",
			expected: "Text",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripAnsiCodes(tt.input)
			if result != tt.expected {
				t.Errorf("stripAnsiCodes(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPadLineToWidth(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		width    int
		expected string
	}{
		{
			name:     "plain text shorter than width",
			line:     "Hello",
			width:    10,
			expected: "Hello     ",
		},
		{
			name:     "plain text equal to width",
			line:     "HelloWorld",
			width:    10,
			expected: "HelloWorld",
		},
		{
			name:     "plain text longer than width",
			line:     "HelloWorldExtra",
			width:    10,
			expected: "HelloWorldExtra",
		},
		{
			name:     "text with ANSI codes shorter than width",
			line:     "\x1b[32mGreen\x1b[0m",
			width:    10,
			expected: "\x1b[32mGreen\x1b[0m     ",
		},
		{
			name:     "text with ANSI codes equal to width",
			line:     "\x1b[32mGreenText!\x1b[0m",
			width:    10,
			expected: "\x1b[32mGreenText!\x1b[0m",
		},
		{
			name:     "text with ANSI codes longer than width",
			line:     "\x1b[32mGreenTextExtra\x1b[0m",
			width:    10,
			expected: "\x1b[32mGreenTextExtra\x1b[0m",
		},
		{
			name:     "empty string with width",
			line:     "",
			width:    5,
			expected: "     ",
		},
		{
			name:     "only ANSI codes with width",
			line:     "\x1b[32m\x1b[0m",
			width:    5,
			expected: "\x1b[32m\x1b[0m     ",
		},
		{
			name:     "width of zero",
			line:     "Text",
			width:    0,
			expected: "Text",
		},
		{
			name:     "multiple ANSI codes in text",
			line:     "\x1b[31mRed\x1b[0m\x1b[32mGreen\x1b[0m",
			width:    10,
			expected: "\x1b[31mRed\x1b[0m\x1b[32mGreen\x1b[0m  ",
		},
		{
			name:     "ANSI codes at different positions",
			line:     "Start\x1b[32mMiddle\x1b[0mEnd",
			width:    20,
			expected: "Start\x1b[32mMiddle\x1b[0mEnd      ",
		},
		{
			name:     "exact fit with ANSI codes",
			line:     "\x1b[1m\x1b[32mHello\x1b[0m",
			width:    5,
			expected: "\x1b[1m\x1b[32mHello\x1b[0m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padLineToWidth(tt.line, tt.width)
			if result != tt.expected {
				t.Errorf("padLineToWidth(%q, %d) = %q, expected %q", tt.line, tt.width, result, tt.expected)
			}
			// Verify the visual length (without ANSI codes) matches or exceeds width
			strippedResult := stripAnsiCodes(result)
			if len(strippedResult) < tt.width && len(stripAnsiCodes(tt.line)) < tt.width {
				t.Errorf("padLineToWidth(%q, %d) resulted in visual length %d, expected at least %d",
					tt.line, tt.width, len(strippedResult), tt.width)
			}
		})
	}
}

func TestPadLineToWidthConsistency(t *testing.T) {
	// Test that padding is consistent - calling twice should give same result
	tests := []struct {
		line  string
		width int
	}{
		{"Hello", 10},
		{"\x1b[32mGreen\x1b[0m", 15},
		{"", 5},
	}

	for _, tt := range tests {
		result1 := padLineToWidth(tt.line, tt.width)
		result2 := padLineToWidth(tt.line, tt.width)
		if result1 != result2 {
			t.Errorf("padLineToWidth is not consistent: first call = %q, second call = %q", result1, result2)
		}
	}
}

func TestStripAnsiCodesIdempotent(t *testing.T) {
	// Test that stripping ANSI codes twice gives the same result as once
	tests := []string{
		"Plain text",
		"\x1b[32mGreen\x1b[0m",
		"\x1b[1m\x1b[31mBold Red\x1b[0m",
		"",
	}

	for _, input := range tests {
		result1 := stripAnsiCodes(input)
		result2 := stripAnsiCodes(result1)
		if result1 != result2 {
			t.Errorf("stripAnsiCodes is not idempotent for input %q: first = %q, second = %q", input, result1, result2)
		}
	}
}

func TestRenderUI_NilSummary(t *testing.T) {
	// Test rendering with nil summary
	lastUpdate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	result := renderUI(nil, false, "", lastUpdate, false)

	// Should contain last updated time
	if !strings.Contains(result, "Last Updated") {
		t.Error("Expected output to contain 'Last Updated'")
	}

	// Should not be empty
	if result == "" {
		t.Error("Expected non-empty output for nil summary")
	}
}

func TestRenderUI_WithError(t *testing.T) {
	// Test rendering with error state
	lastUpdate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	errMsg := "Failed to fetch status"
	result := renderUI(nil, true, errMsg, lastUpdate, false)

	// Should contain error message
	if !strings.Contains(result, errMsg) {
		t.Errorf("Expected output to contain error message %q", errMsg)
	}

	// Should contain last updated time
	if !strings.Contains(result, "Last Updated") {
		t.Error("Expected output to contain 'Last Updated'")
	}
}

func TestRenderUI_WithComponents(t *testing.T) {
	// Create test data with components
	lastUpdate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	summary := &status.SystemStatus{
		Components: []status.Components{
			{
				ID:        "comp1",
				Component: "Git Operations",
				Status:    status.COMPONENT_OPERATIONAL,
			},
			{
				ID:        "comp2",
				Component: "API Requests",
				Status:    status.COMPONENT_DEGREDADED_PERFORMANCE,
			},
		},
		Incidents: []status.Incidents{},
	}

	result := renderUI(summary, false, "", lastUpdate, false)

	// Should contain component names
	if !strings.Contains(result, "Git Operations") {
		t.Error("Expected output to contain 'Git Operations'")
	}
	if !strings.Contains(result, "API Requests") {
		t.Error("Expected output to contain 'API Requests'")
	}

	// Should contain status information
	if !strings.Contains(result, "Operational") {
		t.Error("Expected output to contain 'Operational'")
	}
	if !strings.Contains(result, "Degraded Performance") {
		t.Error("Expected output to contain 'Degraded Performance'")
	}

	// Should contain System Status box title
	if !strings.Contains(result, "System Status") {
		t.Error("Expected output to contain 'System Status' box title")
	}
}

func TestRenderUI_FilterIgnoredComponent(t *testing.T) {
	// Test that ignored component is filtered out
	lastUpdate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	summary := &status.SystemStatus{
		Components: []status.Components{
			{
				ID:        IGNORE_GHSTATUS_COMPONENTID,
				Component: "Visit www.githubstatus.com",
				Status:    status.COMPONENT_OPERATIONAL,
			},
			{
				ID:        "comp1",
				Component: "Git Operations",
				Status:    status.COMPONENT_OPERATIONAL,
			},
		},
		Incidents: []status.Incidents{},
	}

	result := renderUI(summary, false, "", lastUpdate, false)

	// Should NOT contain ignored component
	if strings.Contains(result, "Visit www.githubstatus.com") {
		t.Error("Expected output to NOT contain ignored component")
	}

	// Should contain the other component
	if !strings.Contains(result, "Git Operations") {
		t.Error("Expected output to contain 'Git Operations'")
	}
}

func TestRenderUI_AllComponentStatuses(t *testing.T) {
	// Test all different component status types
	lastUpdate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	summary := &status.SystemStatus{
		Components: []status.Components{
			{
				ID:        "comp1",
				Component: "Service A",
				Status:    status.COMPONENT_OPERATIONAL,
			},
			{
				ID:        "comp2",
				Component: "Service B",
				Status:    status.COMPONENT_DEGREDADED_PERFORMANCE,
			},
			{
				ID:        "comp3",
				Component: "Service C",
				Status:    status.COMPONENT_PARTIAL_OUTAGE,
			},
			{
				ID:        "comp4",
				Component: "Service D",
				Status:    status.COMPONENT_MAJOR_OUTAGE,
			},
		},
		Incidents: []status.Incidents{},
	}

	result := renderUI(summary, false, "", lastUpdate, false)

	// Check all status types are rendered
	statusChecks := map[string]string{
		"Service A": "Operational",
		"Service B": "Degraded Performance",
		"Service C": "Partial Outage",
		"Service D": "Major Outage",
	}

	for service, statusText := range statusChecks {
		if !strings.Contains(result, service) {
			t.Errorf("Expected output to contain %q", service)
		}
		if !strings.Contains(result, statusText) {
			t.Errorf("Expected output to contain status %q", statusText)
		}
	}
}

func TestRenderUI_WithIncidents(t *testing.T) {
	// Test rendering with incidents
	lastUpdate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	incidentTime := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

	summary := &status.SystemStatus{
		Components: []status.Components{
			{
				ID:        "comp1",
				Component: "Git Operations",
				Status:    status.COMPONENT_OPERATIONAL,
			},
		},
		Incidents: []status.Incidents{
			{
				ID:     "incident123",
				Status: "investigating",
				IncidentUpdates: []status.IncidentUpdate{
					{
						Status: "investigating",
						Update: "We are investigating degraded performance",
						Timestamp: &status.Time{
							Time: &incidentTime,
						},
					},
				},
			},
		},
	}

	result := renderUI(summary, false, "", lastUpdate, false)

	// Should contain incident URL
	if !strings.Contains(result, "githubstatus.com/incidents/incident123") {
		t.Error("Expected output to contain incident URL")
	}

	// Should contain incident update text
	if !strings.Contains(result, "investigating degraded performance") {
		t.Error("Expected output to contain incident update text")
	}

	// Should contain Incident Updates box title
	if !strings.Contains(result, "Incident Updates") {
		t.Error("Expected output to contain 'Incident Updates' box title")
	}

	// Should contain both System Status and Incident Updates boxes
	if !strings.Contains(result, "System Status") {
		t.Error("Expected output to contain 'System Status' box title")
	}
}

func TestRenderUI_EmptyComponents(t *testing.T) {
	// Test with empty components list
	lastUpdate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	summary := &status.SystemStatus{
		Components: []status.Components{},
		Incidents:  []status.Incidents{},
	}

	result := renderUI(summary, false, "", lastUpdate, false)

	// Should not crash and should contain basic elements
	if !strings.Contains(result, "Last Updated") {
		t.Error("Expected output to contain 'Last Updated'")
	}

	// Should contain System Status box even if empty
	if !strings.Contains(result, "System Status") {
		t.Error("Expected output to contain 'System Status' box title")
	}
}

func TestRenderUI_ContainsNewlines(t *testing.T) {
	// Test that output contains newlines for padding
	lastUpdate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	summary := &status.SystemStatus{
		Components: []status.Components{
			{
				ID:        "comp1",
				Component: "Test Component",
				Status:    status.COMPONENT_OPERATIONAL,
			},
		},
		Incidents: []status.Incidents{},
	}

	result := renderUI(summary, false, "", lastUpdate, false)

	// Should contain multiple newlines for terminal height padding
	newlineCount := strings.Count(result, "\n")
	if newlineCount < 5 {
		t.Errorf("Expected at least 5 newlines for padding, got %d", newlineCount)
	}
}

func TestRenderUI_ConsistentOutput(t *testing.T) {
	// Test that calling renderUI twice with same inputs gives consistent results
	lastUpdate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	summary := &status.SystemStatus{
		Components: []status.Components{
			{
				ID:        "comp1",
				Component: "Test",
				Status:    status.COMPONENT_OPERATIONAL,
			},
		},
		Incidents: []status.Incidents{},
	}

	result1 := renderUI(summary, false, "", lastUpdate, false)
	result2 := renderUI(summary, false, "", lastUpdate, false)

	if result1 != result2 {
		t.Error("renderUI should produce consistent output for same inputs")
	}
}

func TestRenderUI_WatchModeHelpText(t *testing.T) {
	// Test that watch mode includes help text at bottom
	lastUpdate := time.Date(2024, 1, 15, 14, 30, 0, 0, time.UTC)
	summary := &status.SystemStatus{
		Components: []status.Components{
			{
				ID:        "comp1",
				Component: "Test Component",
				Status:    status.COMPONENT_OPERATIONAL,
			},
		},
		Incidents: []status.Incidents{},
	}

	// Test without watch mode
	resultNoWatch := renderUI(summary, false, "", lastUpdate, false)
	if strings.Contains(resultNoWatch, "Press 'r' to refresh") {
		t.Error("Expected no help text when watch mode is disabled")
	}

	// Test with watch mode
	resultWatch := renderUI(summary, false, "", lastUpdate, true)
	if !strings.Contains(resultWatch, "Press 'r' to refresh") {
		t.Error("Expected help text to contain 'Press 'r' to refresh' in watch mode")
	}
	if !strings.Contains(resultWatch, "quit") {
		t.Error("Expected help text to mention quit option in watch mode")
	}
}
