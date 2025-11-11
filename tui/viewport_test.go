package tui

import (
	"fmt"
	"testing"
)

func TestUpdateCursorViewport(t *testing.T) {
	testCases := []struct {
		cursor       int
		totalItems   int
		viewportSize int
		wantStart    int
		wantEnd      int
		description  string
	}{
		{0, 5, 10, 0, 5, "Small list (5 items, cursor at 0, viewport 10)"},
		{2, 5, 10, 0, 5, "Small list (5 items, cursor at 2, viewport 10)"},
		{0, 15, 10, 0, 10, "Medium list (15 items, cursor at 0, viewport 10)"},
		{5, 15, 10, 0, 10, "Medium list (15 items, cursor at 5, viewport 10)"},
		{10, 15, 10, 5, 15, "Medium list (15 items, cursor at 10, viewport 10)"},
		{0, 25, 10, 0, 10, "Large list (25 items, cursor at 0, viewport 10)"},
		{10, 25, 10, 1, 11, "Large list (25 items, cursor at 10, viewport 10)"},
		{20, 25, 10, 15, 25, "Large list (25 items, cursor at 20, viewport 10)"},
		{24, 25, 10, 15, 25, "Large list (25 items, cursor at last item, viewport 10)"},
		{5, 20, 5, 1, 6, "Custom viewport (20 items, viewport 5, cursor at 5)"},
		{10, 20, 5, 6, 11, "Custom viewport (20 items, viewport 5, cursor at 10)"},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			start, end := updateCursorViewport(tc.cursor, tc.totalItems, tc.viewportSize)

			if start != tc.wantStart {
				t.Errorf("updateCursorViewport(%d, %d, %d) start = %d, want %d",
					tc.cursor, tc.totalItems, tc.viewportSize, start, tc.wantStart)
			}

			if end != tc.wantEnd {
				t.Errorf("updateCursorViewport(%d, %d, %d) end = %d, want %d",
					tc.cursor, tc.totalItems, tc.viewportSize, end, tc.wantEnd)
			}

			// Verify that the cursor is within the visible range
			if tc.cursor < start || tc.cursor >= end {
				t.Errorf("updateCursorViewport(%d, %d, %d) cursor %d not in visible range [%d, %d)",
					tc.cursor, tc.totalItems, tc.viewportSize, tc.cursor, start, end)
			}
		})
	}
}

func ExampleupdateCursorViewport() {
	// Test with a large list
	start, end := updateCursorViewport(15, 30, 10)
	fmt.Printf("With 30 items, cursor at position 15, viewport size 10:\n")
	fmt.Printf("Visible range: items %d to %d (showing %d items)\n", start, end-1, end-start)

	if start > 0 {
		fmt.Printf("↑ Scroll up indicator would be shown\n")
	}
	if end < 30 {
		fmt.Printf("↓ Scroll down indicator would be shown\n")
	}

	// Output:
	// With 30 items, cursor at position 15, viewport size 10:
	// Visible range: items 6 to 15 (showing 10 items)
	// ↑ Scroll up indicator would be shown
	// ↓ Scroll down indicator would be shown
}