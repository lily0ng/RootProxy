package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderProxyDashboard(m Model) string {
	panel := panelStyle(m.theme)
	items := m.app.Proxies.List()
	var b strings.Builder
	b.WriteString("Proxy Dashboard\n\n")
	if len(items) == 0 {
		b.WriteString("No proxies configured.\n")
		return panel.Render(b.String())
	}
	for _, p := range items {
		marker := " "
		if p.Name == m.activeProxyName {
			marker = lipgloss.NewStyle().Foreground(m.theme.Success).Render("▶")
		}
		b.WriteString(fmt.Sprintf("%s %s  %s://%s\n", marker, p.Name, p.Type, p.Address()))
	}
	b.WriteString("\nTip: Press F4 to test the active proxy.\n")
	return panel.Render(b.String())
}

func renderCertManager(m Model) string {
	panel := panelStyle(m.theme)
	certs := m.app.Certs.List()
	var b strings.Builder
	b.WriteString("Certificate Manager\n\n")
	if len(certs) == 0 {
		b.WriteString("No certificates imported.\n")
		return panel.Render(b.String())
	}
	for _, c := range certs {
		b.WriteString("- " + c.Name + "\n")
	}
	return panel.Render(b.String())
}

func renderProfiles(m Model) string {
	panel := panelStyle(m.theme)
	profiles := m.app.Profiles.List()
	active := m.app.Profiles.Active()
	var b strings.Builder
	b.WriteString("Profile System\n\n")
	for _, p := range profiles {
		marker := " "
		if p.Name == active {
			marker = lipgloss.NewStyle().Foreground(m.theme.Success).Render("●")
		}
		b.WriteString(fmt.Sprintf("%s %s  chain=%v\n", marker, p.Name, p.Chain))
	}
	return panel.Render(b.String())
}

func renderRouting(m Model) string {
	panel := panelStyle(m.theme)
	return panel.Render("Routing Rules\n\nDomain-based / geo-based / app-based routing is scaffolded here.")
}

func renderChains(m Model) string {
	panel := panelStyle(m.theme)
	return panel.Render("Proxy Chains\n\nMulti-hop chain builder (up to 5 hops) is scaffolded here.")
}

func renderMonitoring(m Model) string {
	panel := panelStyle(m.theme)
	return panel.Render("Monitoring & Analytics\n\nTraffic graphs, logs, and latency charts are scaffolded here.")
}

func renderSecurity(m Model) string {
	panel := panelStyle(m.theme)
	return panel.Render("Security Settings\n\nDoH/DoT, leak protection, kill switch are scaffolded here.")
}

func renderIntegrations(m Model) string {
	panel := panelStyle(m.theme)
	return panel.Render("Integrations\n\nBurp, Nmap, Metasploit routing integrations are scaffolded here.")
}

func renderAdvanced(m Model) string {
	panel := panelStyle(m.theme)
	return panel.Render("Advanced Tools\n\nScraper / validator / bruteforcer stubs are scaffolded here.")
}

func renderSettings(m Model) string {
	panel := panelStyle(m.theme)
	return panel.Render("Settings\n\nTheme: " + m.app.Settings.Theme + "\n")
}

// munal code with go ( for rendercontriolshelp )

func renderControlsHelp(m Model) string {
	helpStyle := lipgloss.NewStyle().Foreground(m.theme.Faint)
	helpText := "Controls: ↑/↓ Navigate | Enter Select | F1 Dashboard | F2 Proxies | F3 Certs | F4 Test Proxy | F5 Profiles | F6 Routing | F7 Chains | F8 Monitoring | F9 Security | F10 Integrations | F11 Advanced | F12 Settings | q Quit"
	return helpStyle.Render(helpText)
}

func PanelStyleThemes(theme Theme) lipgloss.Style {
	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 2).
		Margin(1, 2)
}

// End Of Render Controls Help


