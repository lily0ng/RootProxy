package tui

import (
	"context"
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/lily0ng/RootProxy/internal/proxy"
	"github.com/lily0ng/RootProxy/internal/rootproxy"
)

type screen int

const (
	screenDashboard screen = iota
	screenCertificates
	screenProfiles
	screenRouting
	screenChains
	screenMonitoring
	screenSecurity
	screenIntegrations
	screenAdvanced
	screenSettings
	screenProxyDashboard
)

type proxyTestMsg proxy.TestResult

type Model struct {
	app   *rootproxy.App
	theme Theme

	width  int
	height int

	screen screen

	activeProxyName string
	statusText      string
	latencyText     string

	helpVisible bool
}

func NewModel(app *rootproxy.App) Model {
	m := Model{
		app:             app,
		theme:           HTBDark(),
		screen:          screenDashboard,
		activeProxyName: "",
		statusText:      "DISCONNECTED",
		latencyText:     "-",
	}
	proxies := app.Proxies.List()
	if len(proxies) > 0 {
		m.activeProxyName = proxies[0].Name
		m.statusText = "CONNECTED"
	}
	return m
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil
	case tea.KeyMsg:
		return m.handleKey(msg)
	case proxyTestMsg:
		tr := proxy.TestResult(msg)
		if tr.OK {
			m.statusText = "CONNECTED"
			m.latencyText = tr.Latency.Round(time.Millisecond).String()
		} else {
			m.statusText = "DISCONNECTED"
			m.latencyText = "-"
		}
		return m, nil
	}
	return m, nil
}

func (m Model) View() string {
	base := lipgloss.NewStyle().Background(m.theme.Bg).Foreground(m.theme.Fg)

	head := m.renderHeader()
	body := m.renderBody()
	foot := m.renderFooter()

	out := lipgloss.JoinVertical(lipgloss.Left, head, body, foot)
	return base.Render(out)
}

func (m Model) handleKey(k tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch k.String() {
	case "q", "esc", "f10":
		return m, tea.Quit
	case "f1":
		m.helpVisible = !m.helpVisible
		return m, nil
	case "f4":
		return m, m.testActiveProxyCmd()
	case "ctrl+p":
		m.screen = screenProxyDashboard
		return m, nil
	case "ctrl+c":
		m.screen = screenCertificates
		return m, nil
	case "ctrl+r":
		m.screen = screenRouting
		return m, nil
	case "ctrl+m":
		m.screen = screenMonitoring
		return m, nil
	case "ctrl+s":
		m.screen = screenSettings
		return m, nil
	case "1":
		m.screen = screenProxyDashboard
		return m, nil
	case "2":
		m.screen = screenCertificates
		return m, nil
	case "3":
		m.screen = screenProfiles
		return m, nil
	case "4":
		m.screen = screenRouting
		return m, nil
	case "5":
		m.screen = screenChains
		return m, nil
	case "6":
		m.screen = screenMonitoring
		return m, nil
	case "7":
		m.screen = screenSecurity
		return m, nil
	case "8":
		m.screen = screenIntegrations
		return m, nil
	case "9":
		m.screen = screenAdvanced
		return m, nil
	case "0":
		m.screen = screenSettings
		return m, nil
	}
	return m, nil
}

func (m Model) testActiveProxyCmd() tea.Cmd {
	p, ok := m.app.Proxies.GetByName(m.activeProxyName)
	if !ok {
		return nil
	}
	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		return proxyTestMsg(proxy.TestConnectivity(ctx, p))
	}
}

func (m Model) renderHeader() string {
	title := lipgloss.NewStyle().Foreground(m.theme.Accent).Bold(true).Render("ROOTPROXY v1.0.0")
	themeBar := lipgloss.NewStyle().Foreground(m.theme.Success).Render("[HTB Theme: ■■■■□□]")

	left := title
	right := themeBar

	w := max(0, m.width-2)
	line := lipgloss.PlaceHorizontal(w, lipgloss.Left, left)
	line = overlayRight(line, right, w)

	border := lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, true, false).BorderForeground(m.theme.Border)
	return border.Render(" " + line + " ")
}

func (m Model) renderBody() string {
	if m.helpVisible {
		return m.renderHelp()
	}
	switch m.screen {
	case screenProxyDashboard:
		return renderProxyDashboard(m)
	case screenCertificates:
		return renderCertManager(m)
	case screenProfiles:
		return renderProfiles(m)
	case screenRouting:
		return renderRouting(m)
	case screenChains:
		return renderChains(m)
	case screenMonitoring:
		return renderMonitoring(m)
	case screenSecurity:
		return renderSecurity(m)
	case screenIntegrations:
		return renderIntegrations(m)
	case screenAdvanced:
		return renderAdvanced(m)
	case screenSettings:
		return renderSettings(m)
	default:
		return renderDashboard(m)
	}
}

func (m Model) renderFooter() string {
	foot := lipgloss.NewStyle().Foreground(m.theme.Muted)
	return foot.Render(" F1:Help  F2:QuickSwitch  F3:Logs  F4:Test  F10:Exit ")
}

func (m Model) renderHelp() string {
	panel := panelStyle(m.theme)
	text := "Quick Commands\n\n" +
		"F1  - Show help\n" +
		"F4  - Test current proxy\n" +
		"F10 - Exit\n\n" +
		"Hotkeys\n\n" +
		"Ctrl+P - Proxy manager\n" +
		"Ctrl+C - Cert manager\n" +
		"Ctrl+R - Routing rules\n" +
		"Ctrl+M - Monitoring\n" +
		"Ctrl+S - Settings\n"
	return panel.Render(text)
}

func renderDashboard(m Model) string {
	panel := panelStyle(m.theme)

	menu := "[1] Proxy Dashboard        [2] Certificate Manager\n" +
		"[3] Profile System         [4] Routing Rules\n" +
		"[5] Proxy Chains           [6] Monitoring\n" +
		"[7] Security Settings      [8] Integrations\n" +
		"[9] Advanced Tools         [0] Settings\n"

	activeProfile := fmt.Sprintf("Current Profile: %s", m.app.Profiles.Active())
	activeProxy := "Active Proxy: -"
	if m.activeProxyName != "" {
		if p, ok := m.app.Proxies.GetByName(m.activeProxyName); ok {
			activeProxy = fmt.Sprintf("Active Proxy: %s://%s", p.Type, p.Address())
		}
	}
	status := fmt.Sprintf("Status: %s (%s)", statusDot(m), m.latencyText)

	content := lipgloss.JoinVertical(lipgloss.Left,
		menu,
		"",
		activeProfile,
		activeProxy,
		status,
	)
	return panel.Render(content)
}

func statusDot(m Model) string {
	if m.statusText == "CONNECTED" {
		return lipgloss.NewStyle().Foreground(m.theme.Success).Render("● CONNECTED")
	}
	return lipgloss.NewStyle().Foreground(m.theme.Danger).Render("● DISCONNECTED")
}

func panelStyle(t Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Background(t.PanelBg)
}

func overlayRight(line, right string, width int) string {
	if width <= 0 {
		return line
	}
	rw := lipgloss.Width(right)
	pos := width - rw
	if pos < 0 {
		pos = 0
	}
	return lipgloss.PlaceHorizontal(width, lipgloss.Left, line[:min(len(line), pos)]) + right
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
