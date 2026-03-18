package gpu

import (
	"strings"

	"github.com/HanSoBored/agtop/internal/gpu/providers"
)

// GPUStats is an alias for the providers.GPUStats type
// This maintains backward compatibility with existing code
type GPUStats = providers.GPUStats

// activeProvider holds the detected GPU provider
var activeProvider providers.GPUProvider

// init registers all available GPU providers
func init() {
	// Register all providers
	providers.RegisterProvider(&providers.AdrenoProvider{})
	providers.RegisterProvider(&providers.MaliProvider{})
	providers.RegisterProvider(&providers.PowerVRProvider{})
	providers.RegisterProvider(&providers.ImmortalisProvider{})

	// Detect and set active provider
	activeProvider = providers.GetActiveProvider()
}

// GetStats retrieves GPU statistics from the active provider
// Returns a populated GPUStats struct
func GetStats() GPUStats {
	if activeProvider == nil {
		return GPUStats{
			Vendor: "Unknown",
			Model:  "No GPU detected",
		}
	}

	stats, err := activeProvider.GetStats()
	if err != nil {
		return GPUStats{
			Vendor: activeProvider.Name(),
			Model:  "Error reading stats",
		}
	}

	return stats
}

// GetActiveProviderName returns the name of the detected GPU provider
func GetActiveProviderName() string {
	if activeProvider == nil {
		return "Unknown"
	}
	return activeProvider.Name()
}

// GetAllProviders returns all registered provider names (for debugging)
func GetAllProviders() []string {
	allProviders := providers.GetAllProviders()
	names := make([]string, len(allProviders))
	for i, p := range allProviders {
		names[i] = p.Name()
	}
	return names
}

// Helper functions for string trimming (backward compatibility)
func trimSpace(s string) string {
	return strings.TrimSpace(s)
}
