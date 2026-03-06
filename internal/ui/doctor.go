package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/progress"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"

	"github.com/chrisbraddock/shellkit/internal/config"
)

// DoctorTab shows health check results.
type DoctorTab struct {
	viewport viewport.Model
	cfg      config.Config
	width    int
	height   int
	styles   *Styles
	passed   int
	total    int
}

// NewDoctorTab creates the doctor tab.
func NewDoctorTab(cfg config.Config, styles *Styles) DoctorTab {
	vp := viewport.New(viewport.WithWidth(80), viewport.WithHeight(20))

	t := DoctorTab{
		viewport: vp,
		cfg:      cfg,
		styles:   styles,
	}
	t.runChecks()
	return t
}

func (t *DoctorTab) SetStyles(s *Styles) {
	t.styles = s
	t.runChecks()
}

func (t *DoctorTab) SetSize(w, h int) {
	t.width = w
	t.height = h
	t.viewport.SetWidth(w)
	t.viewport.SetHeight(h)
	t.runChecks()
}

func (t *DoctorTab) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.viewport, cmd = t.viewport.Update(msg)
	return cmd
}

func (t *DoctorTab) View() string {
	return t.viewport.View()
}

// Summary returns a short summary string for the status bar.
func (t *DoctorTab) Summary() string {
	if t.total == 0 {
		return ""
	}
	return fmt.Sprintf("%d/%d passing", t.passed, t.total)
}

type check struct {
	name   string
	passed bool
	detail string
	hint   string
}

func (t *DoctorTab) runChecks() {
	var checks []check

	// Core tools
	for _, tool := range []struct{ name, cmd, hint string }{
		{"chezmoi", "chezmoi", "brew install chezmoi"},
		{"zsh", "zsh", ""},
		{"git", "git", "brew install git"},
		{"fzf", "fzf", "brew install fzf"},
		{"tmux", "tmux", "brew install tmux"},
		{"neovim", "nvim", "brew install neovim"},
	} {
		_, err := exec.LookPath(tool.cmd)
		checks = append(checks, check{
			name:   tool.name,
			passed: err == nil,
			detail: func() string {
				if err == nil {
					return "installed"
				}
				return "not found"
			}(),
			hint: tool.hint,
		})
	}

	// Source directory
	if t.cfg.ChezmoiSrc != "" {
		_, err := os.Stat(t.cfg.ChezmoiSrc)
		checks = append(checks, check{
			name:   "chezmoi source",
			passed: err == nil,
			detail: t.cfg.ChezmoiSrc,
		})
	}

	// Zsh directory
	_, err := os.Stat(t.cfg.ZshDir)
	checks = append(checks, check{
		name:   "~/.zsh directory",
		passed: err == nil,
		detail: t.cfg.ZshDir,
	})

	// Alias files
	if entries, err := os.ReadDir(t.cfg.AliasDir); err == nil {
		count := 0
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".zsh") {
				count++
			}
		}
		checks = append(checks, check{
			name:   "alias files",
			passed: count > 0,
			detail: fmt.Sprintf("%d files", count),
		})
	}

	// Functions
	if entries, err := os.ReadDir(t.cfg.FunctionDir); err == nil {
		count := 0
		for _, e := range entries {
			if !strings.HasPrefix(e.Name(), ".") && !strings.HasPrefix(e.Name(), "_") {
				count++
			}
		}
		checks = append(checks, check{
			name:   "shell functions",
			passed: count > 0,
			detail: fmt.Sprintf("%d functions", count),
		})
	}

	// TPM
	tpmDir := filepath.Join(t.cfg.HomeDir, ".tmux", "plugins", "tpm")
	_, tpmErr := os.Stat(tpmDir)
	checks = append(checks, check{
		name:   "tmux TPM",
		passed: tpmErr == nil,
		detail: func() string {
			if tpmErr == nil {
				return "installed"
			}
			return "not found"
		}(),
		hint: "run: Ctrl-b I in tmux",
	})

	// Count passed/total
	t.passed = 0
	t.total = len(checks)
	for _, c := range checks {
		if c.passed {
			t.passed++
		}
	}

	// Render
	var b strings.Builder

	// Progress bar header
	pct := 0.0
	if t.total > 0 {
		pct = float64(t.passed) / float64(t.total)
	}

	headerLine := t.styles.Title.Render("  Health Check")
	statsText := t.styles.StatsBadge.Render(fmt.Sprintf("%d/%d passing", t.passed, t.total))

	barWidth := t.width - 4
	if barWidth < 20 {
		barWidth = 40
	}

	// Title + stats on same line
	titleWidth := lipgloss.Width(headerLine)
	statsWidth := lipgloss.Width(statsText)
	gap := barWidth - titleWidth - statsWidth
	if gap < 1 {
		gap = 1
	}
	b.WriteString(headerLine)
	b.WriteString(strings.Repeat(" ", gap))
	b.WriteString(statsText)
	b.WriteString("\n")

	// Progress bar using bubbles progress
	progressBar := progress.New(
		progress.WithColors(
			lipgloss.Color("#FF6B6B"), // red
			lipgloss.Color("#FFD93D"), // yellow
			lipgloss.Color("#73F59F"), // green
		),
		progress.WithWidth(barWidth),
		progress.WithoutPercentage(),
	)
	b.WriteString("  ")
	b.WriteString(progressBar.ViewAs(pct))
	b.WriteString("\n\n")

	// Individual checks
	for _, c := range checks {
		status := t.styles.StatusOK.Render("  ✓")
		if !c.passed {
			status = t.styles.StatusFail.Render("  ✗")
		}

		detail := t.styles.Subtle.Render(c.detail)
		line := fmt.Sprintf("%s  %-20s %s", status, c.name, detail)

		// Show remediation hint for failed checks
		if !c.passed && c.hint != "" {
			line += t.styles.Warning.Render(" — " + c.hint)
		}

		b.WriteString(line)
		b.WriteString("\n")
	}

	t.viewport.SetContent(b.String())
}
