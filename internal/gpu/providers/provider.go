package providers

import "os/exec"

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
	CurrentClock    int
	TargetClock     int
	MinClock        string
	MaxClock        string
	BusyPercentage  int
	GpuBusyPercent  int
	DevfreqLoad     int
	Temperature     string

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

// GetAllProviders returns all registered providers (for testing/debugging)
func GetAllProviders() []GPUProvider {
	return ProviderRegistry
}

// Helper functions for providers

// ReadFile reads a file and returns its content as string
func ReadFile(path string) string {
	data, err := exec.Command("cat", path).Output()
	if err != nil {
		return "N/A"
	}
	return string(data)
}

// GetProp reads an Android system property
func GetProp(prop string) string {
	cmd := exec.Command("getprop", prop)
	out, err := cmd.Output()
	if err != nil {
		return "N/A"
	}
	return string(out)
}

// FileExists checks if a file exists
func FileExists(path string) bool {
	_, err := exec.Command("test", "-e", path).Output()
	return err == nil
}
