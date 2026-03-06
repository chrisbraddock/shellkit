package ui

import (
	"fmt"
	"image/color"
	"math"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/chrisbraddock/shellkit/internal/data"
)

// DashboardTab shows startup metrics, graphs, and system info.
type DashboardTab struct {
	viewport viewport.Model
	summary  data.MetricsSummary
	sysInfo  data.SystemInfo
	version  string
	width    int
	height   int
	styles   *Styles
}

// NewDashboardTab creates the dashboard tab.
func NewDashboardTab(entries []data.MetricEntry, info data.SystemInfo, version string, styles *Styles) DashboardTab {
	vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))
	t := DashboardTab{
		viewport: vp,
		summary:  data.ComputeSummary(entries),
		sysInfo:  info,
		version:  version,
		styles:   styles,
	}
	t.renderContent()
	return t
}

func (t *DashboardTab) SetStyles(s *Styles) {
	t.styles = s
	t.renderContent()
}

func (t *DashboardTab) SetSize(w, h int) {
	t.width = w
	t.height = h
	t.viewport.SetWidth(w)
	t.viewport.SetHeight(h)
	t.renderContent()
}

func (t *DashboardTab) SetMetrics(entries []data.MetricEntry) {
	t.summary = data.ComputeSummary(entries)
	t.renderContent()
}

func (t *DashboardTab) SetSysInfo(info data.SystemInfo) {
	t.sysInfo = info
	t.renderContent()
}

func (t *DashboardTab) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.viewport, cmd = t.viewport.Update(msg)
	return cmd
}

func (t *DashboardTab) View() string {
	return t.viewport.View()
}

// Summary returns a short string for the status bar.
func (t *DashboardTab) Summary() string {
	if t.summary.Count == 0 {
		return ""
	}
	return fmt.Sprintf("%.0fms avg · %d sessions", t.summary.Average, t.summary.Count)
}

func (t *DashboardTab) renderContent() {
	var b strings.Builder

	if t.summary.Count == 0 {
		b.WriteString(t.renderEmptyState())
	} else {
		b.WriteString(t.renderStatsBox())
		b.WriteString("\n")
		b.WriteString(t.renderSparkline())
		b.WriteString("\n")
		b.WriteString(t.renderDistribution())
	}

	b.WriteString("\n")
	b.WriteString(t.renderSystemInfo())

	t.viewport.SetContent(b.String())
}

func (t *DashboardTab) renderEmptyState() string {
	var b strings.Builder
	b.WriteString(t.styles.Title.Render("  Shell Startup"))
	b.WriteString("\n\n")
	b.WriteString(t.styles.Subtle.Render("  No metrics recorded yet."))
	b.WriteString("\n")
	b.WriteString(t.styles.Subtle.Render("  Startup timing will appear here after opening a few new shell sessions."))
	b.WriteString("\n\n")
	b.WriteString(t.styles.Subtle.Render("  Metrics are saved to ~/.local/share/shellkit/metrics.jsonl"))
	b.WriteString("\n")
	return b.String()
}

func (t *DashboardTab) renderStatsBox() string {
	var b strings.Builder

	b.WriteString(t.styles.Title.Render("  Shell Startup"))
	b.WriteString("\n\n")

	// Color-code the current value
	currentStyle := t.styles.StatusOK
	if t.summary.Current > 300 {
		currentStyle = t.styles.StatusFail
	} else if t.summary.Current > 150 {
		currentStyle = t.styles.Warning
	}

	label := t.styles.Subtle
	// Row 1
	b.WriteString(fmt.Sprintf("  %s  %s     %s  %s     %s  %s\n",
		label.Render("Current"),
		currentStyle.Render(fmt.Sprintf("%6.1fms", t.summary.Current)),
		label.Render("Average"),
		t.styles.Highlight.Render(fmt.Sprintf("%6.1fms", t.summary.Average)),
		label.Render("Median"),
		t.styles.Highlight.Render(fmt.Sprintf("%6.1fms", t.summary.Median)),
	))
	// Row 2
	b.WriteString(fmt.Sprintf("  %s  %s     %s  %s     %s  %s\n",
		label.Render("Min    "),
		t.styles.Highlight.Render(fmt.Sprintf("%6.1fms", t.summary.Min)),
		label.Render("Max    "),
		t.styles.Highlight.Render(fmt.Sprintf("%6.1fms", t.summary.Max)),
		label.Render("P95   "),
		t.styles.Highlight.Render(fmt.Sprintf("%6.1fms", t.summary.P95)),
	))

	return b.String()
}

func (t *DashboardTab) renderSparkline() string {
	var b strings.Builder
	blocks := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}

	entries := t.summary.Entries
	// Take last N entries that fit width
	maxEntries := 50
	if t.width > 4 {
		maxEntries = t.width - 6
	}
	if maxEntries > len(entries) {
		maxEntries = len(entries)
	}
	if maxEntries < 1 {
		return ""
	}

	recent := entries[len(entries)-maxEntries:]

	// Find min/max for normalization
	minVal, maxVal := recent[0].DurationMs, recent[0].DurationMs
	for _, e := range recent {
		if e.DurationMs < minVal {
			minVal = e.DurationMs
		}
		if e.DurationMs > maxVal {
			maxVal = e.DurationMs
		}
	}

	b.WriteString(t.styles.Subtle.Render(fmt.Sprintf("  Last %d startups", len(recent))))
	b.WriteString("\n  ")

	// Generate gradient colors from green to red
	colors := lipgloss.Blend1D(len(blocks), t.styles.MetricFast, t.styles.MetricSlow)

	rangeVal := maxVal - minVal
	if rangeVal < 1 {
		rangeVal = 1
	}

	for _, e := range recent {
		// Normalize to 0-7
		idx := int(math.Round((e.DurationMs - minVal) / rangeVal * float64(len(blocks)-1)))
		if idx < 0 {
			idx = 0
		}
		if idx >= len(blocks) {
			idx = len(blocks) - 1
		}

		b.WriteString(
			lipgloss.NewStyle().
				Foreground(colors[idx]).
				Render(string(blocks[idx])),
		)
	}
	b.WriteString("\n")

	return b.String()
}

func (t *DashboardTab) renderDistribution() string {
	var b strings.Builder

	type bucket struct {
		label string
		min   float64
		max   float64
	}
	buckets := []bucket{
		{"< 100ms  ", 0, 100},
		{"100-150ms", 100, 150},
		{"150-200ms", 150, 200},
		{"200-300ms", 200, 300},
		{"> 300ms  ", 300, math.MaxFloat64},
	}

	// Count entries per bucket
	counts := make([]int, len(buckets))
	total := len(t.summary.Entries)
	for _, e := range t.summary.Entries {
		for i, bkt := range buckets {
			if e.DurationMs >= bkt.min && e.DurationMs < bkt.max {
				counts[i]++
				break
			}
		}
	}

	b.WriteString(t.styles.Subtle.Render("  Distribution"))
	b.WriteString("\n")

	barWidth := t.width - 26
	if barWidth < 20 {
		barWidth = 20
	}

	bucketColors := []color.Color{
		lipgloss.Color("#73F59F"),
		lipgloss.Color("#A8E86C"),
		lipgloss.Color("#FFD93D"),
		lipgloss.Color("#FF9A3D"),
		lipgloss.Color("#FF6B6B"),
	}

	for i, bkt := range buckets {
		pct := 0.0
		if total > 0 {
			pct = float64(counts[i]) / float64(total)
		}

		bar := progress.New(
			progress.WithColors(bucketColors[i]),
			progress.WithWidth(barWidth),
			progress.WithoutPercentage(),
		)

		pctStr := fmt.Sprintf("%3.0f%%", pct*100)
		b.WriteString(fmt.Sprintf("  %s %s %s\n",
			t.styles.Subtle.Render(bkt.label),
			bar.ViewAs(pct),
			t.styles.Subtle.Render(pctStr),
		))
	}

	return b.String()
}

func (t *DashboardTab) renderSystemInfo() string {
	var b strings.Builder

	b.WriteString(t.styles.Title.Render("  System"))
	b.WriteString("\n\n")

	label := t.styles.Subtle

	// Count installed tools
	installed := 0
	for _, tool := range t.sysInfo.Tools {
		if tool.Installed {
			installed++
		}
	}
	totalTools := len(t.sysInfo.Tools)

	uptime := getUptime()

	// Row 1
	b.WriteString(fmt.Sprintf("  %s  %-20s  %s  %s\n",
		label.Render("OS        "),
		fmt.Sprintf("%s/%s", t.sysInfo.OS, t.sysInfo.Arch),
		label.Render("Uptime   "),
		uptime,
	))
	// Row 2
	b.WriteString(fmt.Sprintf("  %s  %-20s  %s  %s\n",
		label.Render("Shell     "),
		t.sysInfo.Shell,
		label.Render("Terminal "),
		t.sysInfo.Terminal,
	))
	// Row 3
	versionStr := t.version
	if versionStr == "" {
		versionStr = "dev"
	}
	toolsStr := fmt.Sprintf("%d/%d installed", installed, totalTools)
	if totalTools == 0 {
		toolsStr = "loading..."
	}
	b.WriteString(fmt.Sprintf("  %s  %-20s  %s  %s\n",
		label.Render("Shellkit  "),
		versionStr,
		label.Render("Tools    "),
		toolsStr,
	))

	return b.String()
}

// getUptime returns a human-readable uptime string.
func getUptime() string {
	switch runtime.GOOS {
	case "darwin":
		return getDarwinUptime()
	case "linux":
		return getLinuxUptime()
	default:
		return "unknown"
	}
}

func getDarwinUptime() string {
	out, err := exec.Command("sysctl", "-n", "kern.boottime").Output()
	if err != nil {
		return "unknown"
	}
	// Output format: { sec = 1709654400, usec = 0 } ...
	s := string(out)
	idx := strings.Index(s, "sec = ")
	if idx < 0 {
		return "unknown"
	}
	s = s[idx+6:]
	end := strings.Index(s, ",")
	if end < 0 {
		return "unknown"
	}
	bootSec, err := strconv.ParseInt(strings.TrimSpace(s[:end]), 10, 64)
	if err != nil {
		return "unknown"
	}
	uptime := time.Since(time.Unix(bootSec, 0))
	return formatDuration(uptime)
}

func getLinuxUptime() string {
	out, err := exec.Command("cat", "/proc/uptime").Output()
	if err != nil {
		return "unknown"
	}
	parts := strings.Fields(string(out))
	if len(parts) == 0 {
		return "unknown"
	}
	sec, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return "unknown"
	}
	return formatDuration(time.Duration(sec * float64(time.Second)))
}

func formatDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	mins := int(d.Minutes()) % 60

	if days > 0 {
		return fmt.Sprintf("%dd %dh %dm", days, hours, mins)
	}
	if hours > 0 {
		return fmt.Sprintf("%dh %dm", hours, mins)
	}
	return fmt.Sprintf("%dm", mins)
}
