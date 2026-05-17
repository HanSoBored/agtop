package providers

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// GPUStats holds all GPU statistics (shared across all providers)
type GPUStats struct {
	// Vendor identifies the GPU provider (Adreno, Mali, PowerVR, Immortalis)
	Vendor string

	// Identity
	Model           string
	Platform        string
	VulkanDriver    string
	EGLDriver       string
	OpenGLESVersion string

	// Clocks & Usage
	CurrentClock   int
	TargetClock    int
	MinClock       string
	MaxClock       string
	BusyPercentage int
	GpuBusyPercent int
	DevfreqLoad    int
	Temperature    string

	// Power Management
	Governor          string
	HWClockGating     bool
	IdlePowerCollapse bool
	IdleTimer         string
	Throttling        bool
	ThermalPwrlevel   string
	NumPwrLevels      string
	MinPwrLevel       string
	MaxPwrLevel       string
	IFPCCount         string

	// Statistics & Logic
	ResetCount     string
	PreemptCount   string
	PreemptionMode string
	AvailableFreqs []string
	FTPageFault    string
	FTPolicy       string
	FTLongIB       string
	FTHangIntr     string
}

// GPUProvider defines the interface that all GPU providers must implement.
// This allows the UI to work with any GPU vendor (Adreno, Mali, PowerVR, Immortalis)
// without knowing the underlying implementation details.
type GPUProvider interface {
	// Name returns the provider name (e.g., "Adreno", "Mali", "PowerVR", "Immortalis")
	Name() string

	// Detect checks if this provider is available on the current system
	// Returns true if the GPU hardware is detected
	Detect() bool

	// GetStats retrieves current GPU statistics
	// Returns a populated GPUStats struct or an error if retrieval fails
	GetStats() (GPUStats, error)
}

// ProviderRegistry holds all registered GPU providers
var ProviderRegistry []GPUProvider

// RegisterProvider adds a provider to the registry
func RegisterProvider(provider GPUProvider) {
	ProviderRegistry = append(ProviderRegistry, provider)
}

// GetActiveProvider iterates through registered providers and returns the first detected one
// Returns nil if no provider is detected
func GetActiveProvider() GPUProvider {
	for _, provider := range ProviderRegistry {
		if provider.Detect() {
			return provider
		}
	}
	return nil
}

// Helper functions for providers

// ReadFile reads a file and returns its content as string
func ReadFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(data))
}

// GetProp reads an Android system property
func GetProp(prop string) string {
	out, err := exec.Command("getprop", prop).Output()
	if err != nil {
		return "N/A"
	}
	return strings.TrimSpace(string(out))
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// stubStats returns a GPUStats struct with all fields set to "N/A" or 0,
// except for Vendor and identity fields populated via GetProp.
// Used by stub provider implementations (Mali, PowerVR, Immortalis).
func stubStats(vendor string) GPUStats {
	return GPUStats{
		Vendor:            vendor,
		Model:             vendor + " (TODO)",
		Platform:          GetProp("ro.board.platform"),
		VulkanDriver:      GetProp("ro.hardware.vulkan"),
		EGLDriver:         GetProp("ro.hardware.egl"),
		OpenGLESVersion:   GetProp("ro.opengles.version"),
		CurrentClock:      0,
		TargetClock:       0,
		MinClock:          "N/A",
		MaxClock:          "N/A",
		BusyPercentage:    0,
		GpuBusyPercent:    0,
		DevfreqLoad:       0,
		Temperature:       "N/A",
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
		ResetCount:        "N/A",
		PreemptCount:      "N/A",
		PreemptionMode:    "N/A",
		AvailableFreqs:    []string{},
		FTPageFault:       "N/A",
		FTPolicy:          "N/A",
		FTLongIB:          "N/A",
		FTHangIntr:        "N/A",
	}
}

// parseIntOrZero parses a string to int, returning 0 on failure.
func parseIntOrZero(s string) int {
	val, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return 0
	}
	return val
}

// parseIntMHz parses a Hz string and returns MHz (divides by 1,000,000).
func parseIntMHz(s string) int {
	return parseIntOrZero(s) / 1_000_000
}
