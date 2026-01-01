package settings

import lg "github.com/charmbracelet/lipgloss"

func joinWithSpacer(width int, left string, right string) string {
	style := lg.NewStyle().Width(width - lg.Width(left) - lg.Width(right) - 2)
	return lg.JoinHorizontal(lg.Center, left, style.Render(" "), right)
}
