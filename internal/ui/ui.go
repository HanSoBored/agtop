package ui

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/HanSoBored/agtop/internal/gpu"
)

var (
	// Catppuccin Mocha Palette
	themeBase     = lipgloss.Color("#1E1E2E")
	themeCrust    = lipgloss.Color("#11111B")
	themeText     = lipgloss.Color("#CDD6F4")
	themeSubtext  = lipgloss.Color("#A6ADC8")
	themeOverlay  = lipgloss.Color("#45475A")
	themeSapphire = lipgloss.Color("#74C7EC")
	themeGreen    = lipgloss.Color("#A6E3A1")
	themeYellow   = lipgloss.Color("#F9E2AF")
	themeRed      = lipgloss.Color("#F38BA8")
	themeMauve    = lipgloss.Color("#CBA6F7")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(themeCrust).
			Background(themeMauve).
			Padding(0, 2)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(themeOverlay).
			Padding(0, 1)

	boxTitleStyle = lipgloss.NewStyle().
			Foreground(themeMauve).
			Bold(true)

	labelStyle = lipgloss.NewStyle().
			Foreground(themeSubtext)

	valueStyle = lipgloss.NewStyle().
			Foreground(themeText).
			Bold(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(themeOverlay).
			Italic(true)

	vendorStyle = lipgloss.NewStyle().
			Foreground(themeSapphire).
			Bold(true)
)

type tickMsg time.Time

type model struct {
	stats    gpu.GPUStats
	history  []int
	width    int
	height   int
	quitting bool
}

func NewModel() model {
	return model{
		stats:   gpu.GetStats(),
		history: make([]int, 0, 120),
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
	case tickMsg:
		m.stats = gpu.GetStats()
		m.history = append(m.history, m.stats.BusyPercentage)
		if len(m.history) > 120 {
			m.history = m.history[1:]
		}
		return m, tick()
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil
	}
	return m, nil
}

func tick() tea.Cmd {
	return tea.Every(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) View() tea.View {
	var v tea.View
	v.AltScreen = true
	v.WindowTitle = "Android GPU Monitor"
	v.BackgroundColor = themeBase
	v.ForegroundColor = themeText

	if m.quitting {
		return v
	}

	if m.width == 0 {
		v.Content = "Initializing..."
		return v
	}

	var layout strings.Builder

	// --- Top Header ---
	vendorBadge := vendorStyle.Render(fmt.Sprintf(" %s ", m.stats.Vendor))
	header := lipgloss.JoinHorizontal(lipgloss.Center,
		titleStyle.Render("GPU MONITOR"),
		"  ",
		vendorBadge,
		"  ",
		helpStyle.Render("Real-time Monitor"),
	)
	layout.WriteString("\n" + header + "\n\n")

	// --- Responsive Grid Layout ---
	var col1, col2, col3 string

	if m.width >= 110 {
		// 3 Columns (Desktop / Wide)
		w := (m.width - 2) / 3
		innerW := w - 4
		if innerW < 10 {
			innerW = 10
		}

		col1 = renderBox(w, m.renderHardwareRows())
		col2 = renderBox(w, m.renderUsageRows(innerW))
		col3 = renderBox(w, m.renderPowerRows())
		layout.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, col1, col2, col3))

	} else if m.width >= 75 {
		// 2 Columns Top, 1 Bottom (Tablet / Mid-size)
		w := m.width / 2
		innerW := w - 4
		if innerW < 10 {
			innerW = 10
		}

		col1 = renderBox(w, m.renderHardwareRows())
		col2 = renderBox(w, m.renderUsageRows(innerW))

		bottomW := m.width
		col3 = renderBox(bottomW, m.renderPowerRows())

		topRow := lipgloss.JoinHorizontal(lipgloss.Top, col1, col2)
		layout.WriteString(lipgloss.JoinVertical(lipgloss.Left, topRow, col3))

	} else {
		// 1 Column (Mobile / Narrow)
		w := m.width
		if w < 30 {
			w = 30
		}
		innerW := w - 4
		if innerW < 10 {
			innerW = 10
		}

		col1 = renderBox(w, m.renderHardwareRows())
		col2 = renderBox(w, m.renderUsageRows(innerW))
		col3 = renderBox(w, m.renderPowerRows())
		layout.WriteString(lipgloss.JoinVertical(lipgloss.Left, col1, col2, col3))
	}

	// --- Footer ---
	layout.WriteString("\n\n" + helpStyle.Render("[q] Quit  [ctrl+c] Exit "))

	v.Content = layout.String()
	return v
}

// renderBox applies the border constraints
func renderBox(width int, content string) string {
	return boxStyle.Copy().Width(width).Render(content)
}

// renderInlineRow returns standard Left-Aligned "Label: Value" format
func renderInlineRow(label, value string) string {
	return labelStyle.Render(label+": ") + valueStyle.Render(value)
}

func (m model) renderHardwareRows() string {
	var s strings.Builder
	s.WriteString(boxTitleStyle.Render("󰍛 HARDWARE") + "\n\n")
	s.WriteString(renderInlineRow("Model", m.stats.Model) + "\n")
	s.WriteString(renderInlineRow("Platform", m.stats.Platform) + "\n")
	s.WriteString(renderInlineRow("Vulkan", m.stats.VulkanDriver) + "\n")
	s.WriteString(renderInlineRow("EGL", m.stats.EGLDriver) + "\n")
	s.WriteString(renderInlineRow("OpenGL ES", m.stats.OpenGLESVersion) + "\n")
	s.WriteString(renderInlineRow("Target Clock", fmt.Sprintf("%d MHz", m.stats.TargetClock)) + "\n")
	s.WriteString(renderInlineRow("IFPC Count", m.stats.IFPCCount) + "\n")
	s.WriteString(renderInlineRow("Preemption", m.stats.PreemptionMode) + "\n")
	s.WriteString(renderInlineRow("FT Policy", m.stats.FTPolicy) + "\n")
	s.WriteString(renderInlineRow("FT Page Fault", m.stats.FTPageFault) + "\n")
	s.WriteString(renderInlineRow("FT Long IB", m.stats.FTLongIB) + "\n")
	s.WriteString(renderInlineRow("FT Hang Intr", m.stats.FTHangIntr) + "\n")
	s.WriteString(renderInlineRow("Idle Timer", m.stats.IdleTimer))
	return s.String()
}

func (m model) renderUsageRows(innerW int) string {
	var s strings.Builder
	s.WriteString(boxTitleStyle.Render("󰾆 USAGE & THERMAL") + "\n\n")

	s.WriteString(renderInlineRow("Effective Load", fmt.Sprintf("%d%%", m.stats.BusyPercentage)) + "\n")
	s.WriteString(renderSmoothBar(m.stats.BusyPercentage, innerW) + "\n\n")

	s.WriteString(renderInlineRow("├─ GPU Busy", fmt.Sprintf("%d%%", m.stats.GpuBusyPercent)) + "\n")
	s.WriteString(renderInlineRow("└─ Devfreq", fmt.Sprintf("%d%%", m.stats.DevfreqLoad)) + "\n\n")

	s.WriteString(labelStyle.Render("Load History:") + "\n")
	s.WriteString(renderSparkline(m.history, innerW) + "\n\n")

	s.WriteString(renderInlineRow("Temperature", formatTemp(m.stats.Temperature)) + "\n")
	s.WriteString(renderInlineRow("Current Clock", fmt.Sprintf("%d MHz", m.stats.CurrentClock)) + "\n")
	s.WriteString(renderInlineRow("Min/Max", fmt.Sprintf("%s/%s MHz", m.stats.MinClock, m.stats.MaxClock)) + "\n")
	s.WriteString(renderInlineRow("Reset Count", m.stats.ResetCount))
	return s.String()
}

func (m model) renderPowerRows() string {
	var s strings.Builder
	s.WriteString(boxTitleStyle.Render(" POWER & LOGIC") + "\n\n")

	s.WriteString(renderInlineRow("Governor", m.stats.Governor) + "\n")
	s.WriteString(renderInlineRow("Throttling", formatBool(m.stats.Throttling, "ACTIVE", "INACTIVE", false)) + "\n")
	s.WriteString(renderInlineRow("Clock Gating", formatBool(m.stats.HWClockGating, "ON", "OFF", true)) + "\n")
	s.WriteString(renderInlineRow("Power Collapse", formatBool(m.stats.IdlePowerCollapse, "ON", "OFF", true)) + "\n\n")

	s.WriteString(renderInlineRow("Thermal Level", m.stats.ThermalPwrlevel) + "\n")
	s.WriteString(renderInlineRow("Power Levels", m.stats.NumPwrLevels) + "\n")
	s.WriteString(renderInlineRow("Range", fmt.Sprintf("%s - %s", m.stats.MinPwrLevel, m.stats.MaxPwrLevel)) + "\n\n")

	s.WriteString(renderInlineRow("Preempt Count", m.stats.PreemptCount) + "\n")
	s.WriteString(renderInlineRow("Reset Count", m.stats.ResetCount) + "\n\n")

	freqStr := strings.Join(m.stats.AvailableFreqs, ", ")
	s.WriteString(renderInlineRow("Freqs", freqStr))

	return s.String()
}

// Btop-style Smooth partial blocks Progress Bar
func renderSmoothBar(percent int, width int) string {
	if percent < 0 {
		percent = 0
	} else if percent > 100 {
		percent = 100
	}

	color := themeSapphire
	if percent > 85 {
		color = themeRed
	} else if percent > 60 {
		color = themeYellow
	} else if percent > 30 {
		color = themeGreen
	}

	totalEighths := (percent * width * 8) / 100
	fullBlocks := totalEighths / 8
	remainder := totalEighths % 8

	blocks := []string{" ", "▏", "▎", "▍", "▌", "▋", "▊", "▉", "█"}
	style := lipgloss.NewStyle().Foreground(color)

	var sb strings.Builder
	sb.WriteString(style.Render(strings.Repeat("█", fullBlocks)))

	if fullBlocks < width {
		sb.WriteString(style.Render(blocks[remainder]))
	}

	emptyBlocks := width - fullBlocks - 1
	if emptyBlocks > 0 {
		sb.WriteString(lipgloss.NewStyle().Foreground(themeOverlay).Render(strings.Repeat("─", emptyBlocks)))
	}

	return sb.String()
}

// Btop-style compact block sparkline
func renderSparkline(history []int, width int) string {
	if len(history) == 0 {
		return strings.Repeat(" ", width)
	}

	blocks := []string{" ", "▂", "▃", "▄", "▅", "▆", "▇", "█"}
	var spark strings.Builder

	start := 0
	if len(history) > width {
		start = len(history) - width
	}

	if len(history) < width {
		spark.WriteString(strings.Repeat(" ", width-len(history)))
	}

	for _, val := range history[start:] {
		idx := (val * (len(blocks) - 1)) / 100
		if idx < 0 {
			idx = 0
		} else if idx >= len(blocks) {
			idx = len(blocks) - 1
		}

		color := themeSapphire
		if val > 85 {
			color = themeRed
		} else if val > 60 {
			color = themeYellow
		} else if val > 20 {
			color = themeGreen
		}

		spark.WriteString(lipgloss.NewStyle().Foreground(color).Render(blocks[idx]))
	}
	return spark.String()
}

func formatBool(b bool, trueStr, falseStr string, trueIsGood bool) string {
	if b {
		if trueIsGood {
			return lipgloss.NewStyle().Foreground(themeGreen).Render(trueStr)
		}
		return lipgloss.NewStyle().Foreground(themeRed).Render(trueStr)
	}
	if trueIsGood {
		return lipgloss.NewStyle().Foreground(themeRed).Render(falseStr)
	}
	return lipgloss.NewStyle().Foreground(themeGreen).Render(falseStr)
}

func formatTemp(t string) string {
	if t == "N/A" {
		return lipgloss.NewStyle().Foreground(themeSubtext).Render(t)
	}
	temp := 0.0
	fmt.Sscanf(t, "%f", &temp)
	style := valueStyle.Copy()

	if temp > 75 {
		style = style.Foreground(themeRed)
	} else if temp > 55 {
		style = style.Foreground(themeYellow)
	} else {
		style = style.Foreground(themeGreen)
	}
	return style.Render(t)
}

func Start() error {
	p := tea.NewProgram(NewModel())
	_, err := p.Run()
	return err
}
