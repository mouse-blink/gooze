package controller

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tickMsg time.Time

// Simple delegate for estimate list items.
type estimateDelegate struct {
	offset int
}

func (d estimateDelegate) Height() int  { return 1 }
func (d estimateDelegate) Spacing() int { return 0 }
func (d estimateDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd {
	return nil
}

func (d estimateDelegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	file, ok := item.(fileItem)
	if !ok {
		return
	}

	isSelected := index == m.Index()

	var pathStyle, countStyle lipgloss.Style

	var displayPath string

	width := m.Width() - 8 // Subtract count width (6) + spacing (2)

	if isSelected {
		pathStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("6")).
			Bold(true)
		countStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("0")).
			Background(lipgloss.Color("6")).
			Bold(true).
			Width(6).
			Align(lipgloss.Right)

		displayPath = animateScroll(file.path, width, d.offset)
	} else {
		pathStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
		countStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11")).
			Bold(true).
			Width(6).
			Align(lipgloss.Right)

		displayPath = truncateToWidth(file.path, width)
	}

	line := fmt.Sprintf("%s  %s",
		countStyle.Render(fmt.Sprintf("%d", file.count)),
		pathStyle.Render(displayPath),
	)
	_, _ = fmt.Fprint(w, line)
}

func animateScroll(text string, width int, offset int) string {
	if width <= 0 {
		return ""
	}

	textWidth := lipgloss.Width(text)
	if textWidth <= width {
		return text
	}

	// Gap between repeats
	gap := "   "

	// Initial pause before scrolling starts (in ticks)
	pause := 5

	if offset < pause {
		return truncateToWidth(text, width)
	}

	effectiveStep := offset - pause

	// Create the repeating pattern: text + gap
	// We work with runes to handle multi-byte characters correctly
	runes := []rune(text + gap)
	n := len(runes)

	if n == 0 {
		return ""
	}

	start := effectiveStep % n

	// Construct the window
	res := make([]rune, 0, width)
	for i := range width {
		idx := (start + i) % n
		res = append(res, runes[idx])
	}

	return string(res)
}

func truncateToWidth(text string, width int) string {
	if width <= 0 {
		return ""
	}

	if lipgloss.Width(text) <= width {
		return text
	}

	const ellipsis = "…"

	if width <= 1 {
		return ellipsis
	}

	maxWidth := width - lipgloss.Width(ellipsis)
	if maxWidth <= 0 {
		return ellipsis
	}

	currentWidth := 0

	result := make([]rune, 0, len(text))
	for _, r := range text {
		rWidth := lipgloss.Width(string(r))
		if currentWidth+rWidth > maxWidth {
			break
		}

		result = append(result, r)
		currentWidth += rWidth
	}

	return string(result) + ellipsis
}

// estimateModel is used for listing mutations without interactive testing.
type estimateModel struct {
	width        int
	height       int
	fileList     list.Model
	delegate     estimateDelegate
	total        int
	totalFiles   int
	rendered     bool
	animOffset   int
	lastSelected int
}

func newEstimateModel() estimateModel {
	delegate := estimateDelegate{}
	fileList := list.New([]list.Item{}, delegate, 80, 20)
	fileList.SetShowPagination(false)
	fileList.SetShowFilter(true)
	fileList.SetShowHelp(false)
	fileList.SetShowTitle(false)
	fileList.SetShowStatusBar(false)
	fileList.FilterInput.Placeholder = "Filter by path…"

	return estimateModel{
		fileList:     fileList,
		delegate:     delegate,
		lastSelected: -1,
	}
}

func (m estimateModel) Init() tea.Cmd {
	return tea.Tick(time.Second/2, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m estimateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.fileList.SetWidth(m.width)
		// ... adjustment logic handles in View ...

	case tickMsg:
		if m.fileList.FilterState() != list.Filtering && m.rendered {
			m.animOffset++
			m.delegate.offset = m.animOffset
			m.fileList.SetDelegate(m.delegate)

			return m, tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}

		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			// Pass all key events to the list
			var newList list.Model

			newList, cmd = m.fileList.Update(msg)
			m.fileList = newList

			// Detect selection change to reset animation
			if m.fileList.Index() != m.lastSelected {
				m.lastSelected = m.fileList.Index()
				m.animOffset = 0
				m.delegate.offset = 0
				m.fileList.SetDelegate(m.delegate)
			}

			return m, cmd
		}

	case estimationMsg:
		m = m.handleEstimationMsg(msg)
	}

	return m, cmd
}

func (m estimateModel) handleEstimationMsg(msg estimationMsg) estimateModel {
	m.total = msg.total
	m.totalFiles = msg.paths

	// Create sorted file items
	pathsList := make([]string, 0, len(msg.fileStats))
	for path := range msg.fileStats {
		pathsList = append(pathsList, path)
	}

	sort.Strings(pathsList)

	items := make([]list.Item, 0, len(pathsList))
	for _, path := range pathsList {
		items = append(items, fileItem{path: path, count: msg.fileStats[path]})
	}

	m.fileList.SetItems(items)
	m.rendered = true

	// Start animation loop if not started (Init calls it, but just in case)
	// Or ensure selection is tracked
	if len(items) > 0 && m.lastSelected == -1 {
		m.lastSelected = 0
	}

	return m
}

func (m estimateModel) View() string {
	if !m.rendered {
		return "Loading mutation list…\n"
	}

	// Styles
	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Padding(1, 0, 0, 2)

	summaryStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252")).
		Padding(0, 0, 1, 2)

	accentStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("6")) // Cyan

	// 1. Title
	title := titleStyle.Render("🧬 Gooze Mutation Estimate")

	// 2. Summary
	summary := summaryStyle.Render(fmt.Sprintf(
		"Total Mutations: %s   Files: %s",
		accentStyle.Render(fmt.Sprintf("%d", m.total)),
		accentStyle.Render(fmt.Sprintf("%d", m.totalFiles)),
	))

	// 3. Table with border
	table := m.renderTable()

	// 4. Footer
	footerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Align(lipgloss.Center).
		Width(m.width)

	footer := footerStyle.Render("↑/k up • ↓/j down • g/G top/bottom • / filter • q quit")

	return lipgloss.JoinVertical(lipgloss.Left,
		title,
		summary,
		table,
		footer,
	)
}

func (m estimateModel) renderTable() string {
	// List sizing
	// Display calculations:
	// Screen Height
	// - Title (2)
	// - Summary (2)
	// - Footer (1)
	// - Border (2)
	// - Padding/Headers (2)
	// = Left for list
	listHeight := m.height - 9
	if listHeight < 5 {
		listHeight = 5
	}

	// Widths:
	// Window Width
	// - Margin (2)
	// - Border (2)
	// - Padding (2)
	// = List Width
	listWidth := m.width - 6

	m.fileList.SetHeight(listHeight)
	m.fileList.SetWidth(listWidth)

	// Column Headers inside the table area
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("8")).
		Bold(true).
		Border(lipgloss.NormalBorder(), false, false, true, false).
		BorderForeground(lipgloss.Color("8")).
		Width(listWidth)

	headers := headerStyle.Render(fmt.Sprintf("%6s  %s", "Count", "File Path"))

	tableContainer := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("6")).
		Margin(0, 1). // Outer margin
		Padding(0, 1) // Inner padding

	return tableContainer.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			headers,
			m.fileList.View(),
		),
	)
}
