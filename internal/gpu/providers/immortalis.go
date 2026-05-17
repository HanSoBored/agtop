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
	data := GetProp("ro.hardware.vulkan")
	return strings.Contains(strings.ToLower(data), "immortalis")
}

// GetStats retrieves Immortalis GPU statistics
// TODO: Implement full Immortalis stats collection
// Note: Immortalis shares similar sysfs structure with Mali
func (p *ImmortalisProvider) GetStats() (GPUStats, error) {
	return stubStats("Immortalis"), nil
}

// Ensure ImmortalisProvider implements GPUProvider
var _ GPUProvider = (*ImmortalisProvider)(nil)
