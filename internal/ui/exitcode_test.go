package ui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func ptrInt(v int) *int { return &v }

func TestExitCodeStatus(t *testing.T) {
	cases := []struct {
		code *int
		want string
	}{
		{nil, "info"},
		{ptrInt(0), "success"},
		{ptrInt(-1), "secondary"},
		{ptrInt(75), "warning"},
		{ptrInt(127), "warning"},
		{ptrInt(1), "danger"},
		{ptrInt(255), "danger"},
	}
	for _, tc := range cases {
		if got := ExitCodeStatus(tc.code); got != tc.want {
			t.Errorf("ExitCodeStatus(%v) = %q, want %q", tc.code, got, tc.want)
		}
	}
}

func TestFormatExitCode(t *testing.T) {
	pending := FormatExitCode(nil)
	if !strings.Contains(pending, "⏳") {
		t.Errorf("pending should show hourglass, got %q", pending)
	}

	ok := FormatExitCode(ptrInt(0))
	if !strings.Contains(ok, "✅") || !strings.Contains(ok, "0") {
		t.Errorf("success format = %q", ok)
	}

	fail := FormatExitCode(ptrInt(1))
	if !strings.Contains(fail, "❌") || !strings.Contains(fail, "1") {
		t.Errorf("danger format = %q", fail)
	}

	warn := FormatExitCode(ptrInt(75))
	if !strings.Contains(warn, "⚠️") || !strings.Contains(warn, "75") {
		t.Errorf("warning format = %q", warn)
	}

	plain := FormatExitCodePlain(nil)
	if plain != "⏳ (pending)" {
		t.Errorf("plain pending = %q", plain)
	}
	plain0 := FormatExitCodePlain(ptrInt(0))
	if plain0 != "✅ 0" {
		t.Errorf("plain success = %q", plain0)
	}
}

func TestTruncateCellANSIAware(t *testing.T) {
	styled := FormatExitCode(ptrInt(0))
	padded := truncateCell(styled, 14)
	if lipgloss.Width(padded) != 14 {
		t.Errorf("padded visible width = %d, want 14 (raw %q)", lipgloss.Width(padded), padded)
	}
	if !strings.Contains(padded, "✅") {
		t.Errorf("padding dropped emoji: %q", padded)
	}
}
