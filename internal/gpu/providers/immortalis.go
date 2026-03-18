package providers

import "strings"

// Immortalis is ARM's premium GPU line with hardware ray tracing
// TODO: Implement full Immortalis GPU support with proper sysfs paths
type ImmortalisProvider struct{}

// Name returns the provider name
func (p *ImmortalisProvider) Name() string {
	return "Immortalis"
}

// Detect checks if Immortalis GPU is present
// TODO: Implement proper detection logic for Immortalis GPUs
func (p *ImmortalisProvider) Detect() bool {
	// Common detection methods for Immortalis:
	// - Check /sys/class/misc/mali0/device/name for "Immortalis"
	// - Check getprop ro.hardware.vulkan for "immortalis"
	// - Immortalis-G15, Immortalis-G715, etc.
	data := strings.TrimSpace(GetProp("ro.hardware.vulkan"))
	return strings.Contains(strings.ToLower(data), "immortalis")
}

// GetStats retrieves Immortalis GPU statistics
// TODO: Implement full Immortalis stats collection
// Note: Immortalis shares similar sysfs structure with Mali
func (p *ImmortalisProvider) GetStats() (GPUStats, error) {
	stats := GPUStats{
		Vendor: "Immortalis",

		// Identity - TODO: Map to Immortalis-specific sysfs paths
		Model:           "Immortalis (TODO)",
		Platform:        strings.TrimSpace(GetProp("ro.board.platform")),
		VulkanDriver:    strings.TrimSpace(GetProp("ro.hardware.vulkan")),
		EGLDriver:       strings.TrimSpace(GetProp("ro.hardware.egl")),
		OpenGLESVersion: strings.TrimSpace(GetProp("ro.opengles.version")),

		// Clocks & Usage - TODO: Map to Immortalis-specific paths
		CurrentClock:   0,
		TargetClock:    0,
		MinClock:       "N/A",
		MaxClock:       "N/A",
		BusyPercentage: 0,
		GpuBusyPercent: 0,
		DevfreqLoad:    0,
		Temperature:    "N/A",

		// Power Management
		Governor:          "N/A",
		HWClockGating:     false,
		IdlePowerCollapse: false,
		IdleTimer:         "N/A",
		Throttling:        false,
		ThermalPwrlevel:   "N/A",
		NumPwrLevels:      "N/A",
		MinPwrLevel:       "N/A",
		MaxPwrLevel:       "N/A",
		IFPCCount:         "N/A",

		// Statistics & Logic
		ResetCount:     "N/A",
		PreemptCount:   "N/A",
		PreemptionMode: "N/A",
		AvailableFreqs: []string{},
		FTPageFault:    "N/A",
		FTPolicy:       "N/A",
		FTLongIB:       "N/A",
		FTHangIntr:     "N/A",
	}

	return stats, nil
}

// Ensure ImmortalisProvider implements GPUProvider
var _ GPUProvider = (*ImmortalisProvider)(nil)
