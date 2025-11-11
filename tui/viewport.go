package tui

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	// Fixed height for UI elements
	headerHeight        = 3  // 空行 + 标题 + 空行
	statusBarHeight     = 2  // 空行 + 状态栏
	backToMenuHeight    = 1  // "Back to menu" 选项
	warningHeight       = 3  // 警告框（包括边框）
	scrollIndicatorHeight = 1 // 滚动指示器
)

// Scroll indicator style
var scrollIndicatorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("240")).
	Italic(true).
	Padding(0, 1)

// calculateListViewportHeight 计算列表视口的可用高度
// 注意：我们不预留滚动指示器的空间，因为滚动指示器和配置项共享显示区域
// 当需要显示滚动指示器时，会减少显示的配置项数量
func calculateListViewportHeight(windowHeight int, hasWarning bool, compact bool) int {
	if windowHeight <= 0 {
		// 默认高度10，但需要通过实际渲染验证是否适合当前窗口
		return 10
	}

	// 计算固定元素占用的总行数
	// header: 3行（空行 + 标题 + 空行）
	// statusBar: 2行（空行 + 状态栏）
	// backToMenu: 1行
	fixedHeight := headerHeight + statusBarHeight + backToMenuHeight

	// 如果有警告信息，加上警告框的高度
	if hasWarning {
		fixedHeight += warningHeight
	}

	// 计算每个配置项的高度
	itemHeight := 3 // compact 模式：1行内容 + 2行边框
	if !compact {
		// 非compact模式下，每个配置项显示多行信息 + 边框：
		// - 上边框：1行
		// - 内容：3行（Name、Provider、BaseURL）
		// - 下边框：1行
		// - MarginBottom：0行（在 itemBoxStyle 中设置）
		// 总计：5行
		itemHeight = 5
	}

	// 计算可用高度
	availableHeight := windowHeight - fixedHeight

	// 计算可以显示的配置项数量
	// 注意：当显示滚动指示器时，它们会占用配置项的空间
	// 例如：如果这里计算显示10个配置项，当需要显示底部滚动指示器时，
	// 实际可能只显示9个配置项 + 1行滚动指示器
	visibleItems := availableHeight / itemHeight

	if visibleItems < 3 {
		return 3 // 最少显示3个项目
	}

	if visibleItems > 20 {
		return 20 // 最多显示20个项目
	}

	return visibleItems
}

// updateCursorViewport updates the viewport to ensure the cursor is visible
func updateCursorViewport(cursor, totalItems, viewportSize int) (int, int) {
	if totalItems <= viewportSize {
		return 0, totalItems // All items fit, no scrolling needed
	}

	// If cursor is within the first viewport, show the first viewport
	if cursor < viewportSize {
		return 0, viewportSize
	}

	// If cursor is near the end, show the last viewport
	if cursor >= totalItems - viewportSize {
		start := totalItems - viewportSize
		return start, totalItems
	}

	// Otherwise, show cursor at the bottom of the viewport
	offset := cursor - viewportSize + 1
	return offset, offset + viewportSize
}