package ui

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"

	"github.com/chrisbraddock/shellkit/internal/config"
)

// DoctorTab shows health check results.
type DoctorTab struct {
	viewport viewport.Model
	cfg      config.Config
	width    int
	height   int
	styles   *Styles
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
}

func (t *DoctorTab) Update(msg tea.Msg) tea.Cmd {
	var cmd tea.Cmd
	t.viewport, cmd = t.viewport.Update(msg)
	return cmd
}

func (t *DoctorTab) View() string {
	return t.viewport.View()
}

type check struct {
	name   string
	passed bool
	detail string
}

func (t *DoctorTab) runChecks() {
	var checks []check

	// Core tools
	for _, tool := range []struct{ name, cmd string }{
		{"chezmoi", "chezmoi"},
		{"zsh", "zsh"},
		{"git", "git"},
		{"fzf", "fzf"},
		{"tmux", "tmux"},
		{"neovim", "nvim"},
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
			return "run: Ctrl-b I in tmux"
		}(),
	})

	// Render
	var b strings.Builder
	b.WriteString(t.styles.Title.Render("  Health Check"))
	b.WriteString("\n\n")

	for _, c := range checks {
		status := t.styles.StatusOK.Render("  ✓")
		if !c.passed {
			status = t.styles.StatusFail.Render("  ✗")
		}
		b.WriteString(fmt.Sprintf("%s  %-20s %s\n",
			status,
			c.name,
			t.styles.Subtle.Render(c.detail),
		))
	}

	t.viewport.SetContent(b.String())
}
