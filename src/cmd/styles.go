package cmd

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/raci-sh/raci-CLI/src/internal/ui"
)

var (
	brandStyle = lipgloss.NewStyle().Bold(true).Foreground(ui.Primary)
	labelStyle = lipgloss.NewStyle().Foreground(ui.Secondary).Width(10)
	valueStyle = lipgloss.NewStyle().Foreground(ui.Accent)
	cmdStyle   = lipgloss.NewStyle().Bold(true).Foreground(ui.Primary)
	descStyle  = lipgloss.NewStyle().Foreground(ui.Secondary)
)
