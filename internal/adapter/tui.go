package adapter

import (
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	m "github.com/mouse-blink/gooze/internal/model"
	"golang.org/x/term"
)

const (
	// ANSI color codes for zero values (dark gray, faint).
	grayColor  = "\033[2;90m" // Faint + dark gray
	resetColor = "\033[0m"    // Reset
)

// TUI implements UI using Bubble Tea for interactive display.
type TUI struct {
	output io.Writer
}

// NewTUI creates a new TUI.
func NewTUI(output io.Writer) *TUI {
	return &TUI{output: output}
}

// Display shows source files using Bubble Tea TUI.
func (p *TUI) Display(sources []m.Source) error {
	model := newTUIModel(sources, false)

	// Get initial terminal size
	if f, ok := p.output.(*os.File); ok {
		width, height, err := term.GetSize(int(f.Fd()))
		if err == nil {
			model.height = height
			model.width = width
		}
	}

	// If list is small, just print and exit
	if !model.needsPagination() {
		_, err := fmt.Fprint(p.output, model.View())
		return err
	}

	program := tea.NewProgram(model, tea.WithOutput(p.output), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}

// ShowNotImplemented displays a "not implemented" message using TUI.
func (p *TUI) ShowNotImplemented(count int) error {
	model := newNotImplementedModel(count)

	// Not implemented message is always short, just print and exit
	_, err := fmt.Fprint(p.output, model.View())

	return err
}

// DisplayMutationEstimations shows pre-calculated mutation estimations using TUI.
func (p *TUI) DisplayMutationEstimations(estimations map[m.Path]int) error {
	// Convert map to sorted slice for display
	counts := make([]mutationCount, 0, len(estimations))
	total := 0

	// Sort paths for consistent output
	paths := make([]m.Path, 0, len(estimations))
	for path := range estimations {
		paths = append(paths, path)
	}

	// Simple sort by string comparison
	for i := 0; i < len(paths); i++ {
		for j := i + 1; j < len(paths); j++ {
			if string(paths[i]) > string(paths[j]) {
				paths[i], paths[j] = paths[j], paths[i]
			}
		}
	}

	for _, path := range paths {
		count := estimations[path]
		counts = append(counts, mutationCount{
			file:  string(path),
			count: count,
		})
		total += count
	}

	model := newMutationCountModel(counts, total, m.MutationArithmetic)

	// Get initial terminal size
	if f, ok := p.output.(*os.File); ok {
		width, height, err := term.GetSize(int(f.Fd()))
		if err == nil {
			model.height = height
			model.width = width
		}
	}

	// If list is small, just print and exit
	if !model.needsPagination() {
		_, err := fmt.Fprint(p.output, model.View())
		return err
	}

	program := tea.NewProgram(model, tea.WithOutput(p.output), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		return err
	}

	return nil
}

// mutationCount holds the mutation count for a single file.
type mutationCount struct {
	file  string
	count int
} // tuiModel represents the Bubble Tea model for displaying source files.
type tuiModel struct {
	sources        []m.Source
	notImplemented bool
	count          int
	height         int
	width          int
	offset         int // Current scroll offset
	quitting       bool
}

func newTUIModel(sources []m.Source, notImplemented bool) tuiModel {
	return tuiModel{
		sources:        sources,
		notImplemented: notImplemented,
		height:         0, // Will be set on first WindowSizeMsg
		width:          0,
		offset:         0,
		quitting:       false,
	}
}

func newNotImplementedModel(count int) tuiModel {
	return tuiModel{
		notImplemented: true,
		count:          count,
		height:         0,
		width:          0,
		offset:         0,
		quitting:       false,
	}
}

func (tm tuiModel) Init() tea.Cmd {
	return nil
}

func (tm tuiModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		tm.height = msg.Height
		tm.width = msg.Width

		return tm, nil

	case tea.KeyMsg:
		return tm.handleKeyPress(msg)
	}

	return tm, nil
}

//nolint:cyclop,exhaustive,dupl // Key handling requires multiple cases for UI navigation, duplication intentional for different models
func (tm tuiModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		tm.quitting = true
		return tm, tea.Quit
	default:
		// Handle other key types in the string switch below
	}

	switch msg.String() {
	case "q":
		tm.quitting = true
		return tm, tea.Quit

	case "down", "j":
		tm.offset++

		maxOffset := tm.maxOffset()
		if tm.offset > maxOffset {
			tm.offset = maxOffset
		}

		return tm, nil

	case "up", "k":
		tm.offset--
		if tm.offset < 0 {
			tm.offset = 0
		}

		return tm, nil

	case "g", "home":
		tm.offset = 0

		return tm, nil

	case "G", "end":
		tm.offset = tm.maxOffset()

		return tm, nil

	case "d", "pgdown":
		tm.offset += tm.itemsPerPage()

		maxOffset := tm.maxOffset()
		if tm.offset > maxOffset {
			tm.offset = maxOffset
		}

		return tm, nil

	case "u", "pgup":
		tm.offset -= tm.itemsPerPage()
		if tm.offset < 0 {
			tm.offset = 0
		}

		return tm, nil
	}

	return tm, nil
}

// itemsPerPage calculates how many items can fit on screen.
func (tm tuiModel) itemsPerPage() int {
	if tm.height == 0 {
		return 10 // Default
	}
	// Reserve space for:
	// - Header: 4 lines (box + empty)
	// - Title: 2 lines (count + empty)
	// - Footer: 3 lines (empty + page + help)
	// - Top margin: 1 line
	// Total: 10 lines
	reserved := 10

	available := tm.height - reserved
	if available < 1 {
		return 1
	}

	return available
}

// maxOffset returns the maximum scroll offset.
func (tm tuiModel) maxOffset() int {
	itemCount := len(tm.sources)

	perPage := tm.itemsPerPage()
	if perPage <= 0 {
		return 0
	}

	maxOff := itemCount - perPage
	if maxOff < 0 {
		return 0
	}

	return maxOff
}

// needsPagination returns true if the list is too large to fit on screen.
func (tm tuiModel) needsPagination() bool {
	if tm.notImplemented {
		return false
	}

	totalFiles := len(tm.sources)
	if totalFiles == 0 {
		return false
	}

	itemsPerPage := tm.itemsPerPage()

	return totalFiles > itemsPerPage && tm.height > 0
}

func (tm tuiModel) View() string {
	var b strings.Builder

	tm.renderHeader(&b)

	if tm.notImplemented {
		tm.renderNotImplemented(&b)
		return b.String()
	}

	if len(tm.sources) == 0 {
		b.WriteString("  📭 No source files found\n")
		return b.String()
	}

	tm.renderSourceList(&b)

	return b.String()
}

func (tm tuiModel) renderHeader(b *strings.Builder) {
	b.WriteString("╔════════════════════════════════════════════════════════════════╗\n")
	b.WriteString("║                    Gooze - Mutation Testing                    ║\n")
	b.WriteString("╚════════════════════════════════════════════════════════════════╝\n\n")
}

func (tm tuiModel) renderNotImplemented(b *strings.Builder) {
	fmt.Fprintf(b, "  📊 Found %d source file(s)\n\n", tm.count)
	b.WriteString("  ⚠️  Mutation testing not yet implemented.\n")
	b.WriteString("  💡 Use --list flag to see source files.\n\n")
}

func (tm tuiModel) renderSourceList(b *strings.Builder) {
	totalFiles := len(tm.sources)
	fmt.Fprintf(b, "  📁 Found %d source file(s):\n\n", totalFiles)

	// Calculate pagination
	itemsPerPage := tm.itemsPerPage()
	needsPagination := totalFiles > itemsPerPage && tm.height > 0

	start := tm.offset

	end := start + itemsPerPage
	if end > totalFiles {
		end = totalFiles
	}

	if start >= totalFiles {
		start = totalFiles - 1
		if start < 0 {
			start = 0
		}
	}

	// Show items for current page
	displaySources := tm.sources

	if needsPagination {
		displaySources = tm.sources[start:end]
	}

	for i, source := range displaySources {
		actualIndex := start + i + 1
		fmt.Fprintf(b, "  %2d. %s\n", actualIndex, source.Origin)
	}

	// Footer with navigation help
	b.WriteString("\n")

	if needsPagination {
		currentPage := (tm.offset / itemsPerPage) + 1
		totalPages := (totalFiles + itemsPerPage - 1) / itemsPerPage
		fmt.Fprintf(b, "  Page %d/%d | Showing %d-%d of %d\n",
			currentPage, totalPages, start+1, end, totalFiles)
		b.WriteString("  ↑/k: up | ↓/j: down | g: top | G: bottom | q: quit\n")
	}
}

// mutationCountModel represents the Bubble Tea model for displaying mutation counts.
type mutationCountModel struct {
	counts       []mutationCount
	total        int
	mutationType m.MutationType
	height       int
	width        int
	offset       int // Current scroll offset
	quitting     bool
}

func newMutationCountModel(counts []mutationCount, total int, mutationType m.MutationType) mutationCountModel {
	return mutationCountModel{
		counts:       counts,
		total:        total,
		mutationType: mutationType,
		height:       0, // Will be set on first WindowSizeMsg
		width:        0,
		offset:       0,
		quitting:     false,
	}
}

func (mcm mutationCountModel) Init() tea.Cmd {
	return nil
}

func (mcm mutationCountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		mcm.height = msg.Height
		mcm.width = msg.Width

		return mcm, nil

	case tea.KeyMsg:
		return mcm.handleKeyPress(msg)
	}

	return mcm, nil
}

//nolint:cyclop,exhaustive,dupl // Key handling requires multiple cases for UI navigation, duplication intentional for different models
func (mcm mutationCountModel) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		mcm.quitting = true
		return mcm, tea.Quit
	default:
		// Handle other key types in the string switch below
	}

	switch msg.String() {
	case "q":
		mcm.quitting = true
		return mcm, tea.Quit

	case "down", "j":
		mcm.offset++

		maxOffset := mcm.maxOffset()
		if mcm.offset > maxOffset {
			mcm.offset = maxOffset
		}

		return mcm, nil

	case "up", "k":
		mcm.offset--
		if mcm.offset < 0 {
			mcm.offset = 0
		}

		return mcm, nil

	case "g", "home":
		mcm.offset = 0

		return mcm, nil

	case "G", "end":
		mcm.offset = mcm.maxOffset()

		return mcm, nil

	case "d", "pgdown":
		mcm.offset += mcm.itemsPerPage()

		maxOffset := mcm.maxOffset()
		if mcm.offset > maxOffset {
			mcm.offset = maxOffset
		}

		return mcm, nil

	case "u", "pgup":
		mcm.offset -= mcm.itemsPerPage()
		if mcm.offset < 0 {
			mcm.offset = 0
		}

		return mcm, nil
	}

	return mcm, nil
}

// itemsPerPage calculates how many items can fit on screen.
func (mcm mutationCountModel) itemsPerPage() int {
	if mcm.height == 0 {
		return 10 // Default
	}
	// Reserve space for:
	// - Header: 4 lines (box + empty)
	// - Title: 2 lines (summary + empty)
	// - Total: 2 lines (empty + total)
	// - Footer: 3 lines (empty + page + help)
	// - Top margin: 1 line
	// Total: 12 lines
	reserved := 12

	available := mcm.height - reserved
	if available < 1 {
		return 1
	}

	return available
}

// maxOffset returns the maximum scroll offset.
func (mcm mutationCountModel) maxOffset() int {
	itemCount := len(mcm.counts)

	perPage := mcm.itemsPerPage()
	if perPage <= 0 {
		return 0
	}

	maxOff := itemCount - perPage
	if maxOff < 0 {
		return 0
	}

	return maxOff
}

// needsPagination returns true if the list is too large to fit on screen.
func (mcm mutationCountModel) needsPagination() bool {
	totalFiles := len(mcm.counts)
	if totalFiles == 0 {
		return false
	}

	itemsPerPage := mcm.itemsPerPage()

	return totalFiles > itemsPerPage && mcm.height > 0
}

func (mcm mutationCountModel) View() string {
	var b strings.Builder

	mcm.renderHeader(&b)

	if len(mcm.counts) == 0 {
		b.WriteString("  📭 No source files found\n")
		return b.String()
	}

	mcm.renderMutationCountList(&b)

	return b.String()
}

func (mcm mutationCountModel) renderHeader(b *strings.Builder) {
	b.WriteString("╔════════════════════════════════════════════════════════════════╗\n")
	b.WriteString("║                    Gooze - Mutation Testing                    ║\n")
	b.WriteString("╚════════════════════════════════════════════════════════════════╝\n\n")
}

func (mcm mutationCountModel) renderMutationCountList(b *strings.Builder) {
	totalFiles := len(mcm.counts)

	mutationTypeName := "arithmetic"
	if mcm.mutationType != m.MutationArithmetic {
		mutationTypeName = "unknown"
	}

	fmt.Fprintf(b, "  🔢 %s mutations summary:\n\n", mutationTypeName)

	// Calculate pagination
	itemsPerPage := mcm.itemsPerPage()
	needsPagination := totalFiles > itemsPerPage && mcm.height > 0

	start := mcm.offset

	end := start + itemsPerPage
	if end > totalFiles {
		end = totalFiles
	}

	if start >= totalFiles {
		start = totalFiles - 1
		if start < 0 {
			start = 0
		}
	}

	// Show items for current page
	displayCounts := mcm.counts

	if needsPagination {
		displayCounts = mcm.counts[start:end]
	}

	for _, mc := range displayCounts {
		if mc.count == 0 {
			// Gray out only the zero count value
			fmt.Fprintf(b, "  %s: %s%d mutations%s\n", mc.file, grayColor, mc.count, resetColor)
		} else {
			fmt.Fprintf(b, "  %s: %d mutations\n", mc.file, mc.count)
		}
	}

	// Total count
	b.WriteString("\n")
	fmt.Fprintf(b, "  📊 Total: %d %s mutations across %d file(s)\n", mcm.total, mutationTypeName, totalFiles)

	// Footer with navigation help
	if needsPagination {
		b.WriteString("\n")

		currentPage := (mcm.offset / itemsPerPage) + 1
		totalPages := (totalFiles + itemsPerPage - 1) / itemsPerPage
		fmt.Fprintf(b, "  Page %d/%d | Showing %d-%d of %d\n",
			currentPage, totalPages, start+1, end, totalFiles)
		b.WriteString("  ↑/k: up | ↓/j: down | g: top | G: bottom | q: quit\n")
	}
}
