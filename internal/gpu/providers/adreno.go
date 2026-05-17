package providers

import (
	"fmt"
	"strconv"
	"strings"
)

// Qualcomm Adreno GPUs
type AdrenoProvider struct{}

func (p *AdrenoProvider) Name() string {
	return "Adreno"
}

// Detect checks if Adreno GPU is present by checking kgsl device path
func (p *AdrenoProvider) Detect() bool {
	vulkanDriver := GetProp("ro.hardware.vulkan")
	return strings.Contains(strings.ToLower(vulkanDriver), "adreno") ||
		strings.Contains(strings.ToLower(vulkanDriver), "freedreno") ||
		FileExists("/sys/class/kgsl/kgsl-3d0")
}

// GetStats retrieves Adreno GPU statistics
func (p *AdrenoProvider) GetStats() (GPUStats, error) {
	// Return error if the critical kgsl device path doesn't exist
	if !FileExists("/sys/class/kgsl/kgsl-3d0") {
		return GPUStats{}, fmt.Errorf("kgsl device path /sys/class/kgsl/kgsl-3d0 not found")
	}

	stats := GPUStats{
		Vendor: "Adreno",

		// Identity
		Model:           ReadFile("/sys/class/kgsl/kgsl-3d0/gpu_model"),
		Platform:        GetProp("ro.board.platform"),
		VulkanDriver:    GetProp("ro.hardware.vulkan"),
		EGLDriver:       GetProp("ro.hardware.egl"),
		OpenGLESVersion: GetProp("ro.opengles.version"),
		MinClock:        ReadFile("/sys/class/kgsl/kgsl-3d0/min_clock_mhz"),
		MaxClock:        ReadFile("/sys/class/kgsl/kgsl-3d0/max_clock_mhz"),
		Governor:        ReadFile("/sys/class/kgsl/kgsl-3d0/devfreq/governor"),
		IdleTimer:       ReadFile("/sys/class/kgsl/kgsl-3d0/idle_timer"),
		ThermalPwrlevel: ReadFile("/sys/class/kgsl/kgsl-3d0/thermal_pwrlevel"),
		ResetCount:      ReadFile("/sys/class/kgsl/kgsl-3d0/reset_count"),
		PreemptCount:    ReadFile("/sys/class/kgsl/kgsl-3d0/preempt_count"),
		PreemptionMode:  ReadFile("/sys/class/kgsl/kgsl-3d0/preemption"),
		Temperature:     ReadFile("/sys/class/kgsl/kgsl-3d0/temp"),
		NumPwrLevels:    ReadFile("/sys/class/kgsl/kgsl-3d0/num_pwrlevels"),
		MinPwrLevel:     ReadFile("/sys/class/kgsl/kgsl-3d0/min_pwrlevel"),
		MaxPwrLevel:     ReadFile("/sys/class/kgsl/kgsl-3d0/max_pwrlevel"),
		IFPCCount:       ReadFile("/sys/class/kgsl/kgsl-3d0/ifpc_count"),
		FTPageFault:     ReadFile("/sys/class/kgsl/kgsl-3d0/ft_pagefault_policy"),
		FTPolicy:        ReadFile("/sys/class/kgsl/kgsl-3d0/ft_policy"),
		FTLongIB:        ReadFile("/sys/class/kgsl/kgsl-3d0/ft_long_ib_detect"),
		FTHangIntr:      ReadFile("/sys/class/kgsl/kgsl-3d0/ft_hang_intr_status"),
	}

	// Current Clock
	stats.CurrentClock = parseIntMHz(ReadFile("/sys/class/kgsl/kgsl-3d0/gpuclk"))

	// Target Clock
	stats.TargetClock = parseIntMHz(ReadFile("/sys/class/kgsl/kgsl-3d0/devfreq/target_freq"))

	// GPU Busy
	busyStr := ReadFile("/sys/class/kgsl/kgsl-3d0/gpu_busy_percentage")
	if busyStr != "N/A" {
		stats.GpuBusyPercent = parseIntOrZero(strings.TrimSuffix(busyStr, "%"))
	}

	// Devfreq Load
	devfreqStr := ReadFile("/sys/class/kgsl/kgsl-3d0/devfreq/gpu_load")
	if devfreqStr != "N/A" {
		stats.DevfreqLoad = parseIntOrZero(devfreqStr)
	}

	// Effective BusyPercentage
	if stats.DevfreqLoad > stats.GpuBusyPercent {
		stats.BusyPercentage = stats.DevfreqLoad
	} else {
		stats.BusyPercentage = stats.GpuBusyPercent
	}

	// Boolean flags
	stats.HWClockGating = ReadFile("/sys/class/kgsl/kgsl-3d0/hwcg") == "1"
	stats.IdlePowerCollapse = ReadFile("/sys/class/kgsl/kgsl-3d0/ifpc") == "1"
	stats.Throttling = ReadFile("/sys/class/kgsl/kgsl-3d0/throttling") == "1"

	// Available frequencies
	freqsStr := ReadFile("/sys/class/kgsl/kgsl-3d0/freq_table_mhz")
	if freqsStr != "N/A" {
		stats.AvailableFreqs = strings.Fields(freqsStr)
	}

	// Format temperature
	if stats.Temperature != "N/A" {
		temp, err := strconv.ParseFloat(stats.Temperature, 64)
		if err == nil && temp > 1000 {
			stats.Temperature = strconv.FormatFloat(temp/1000.0, 'f', 1, 64) + "°C"
		} else {
			stats.Temperature += "°C"
		}
	}

	return stats, nil
}

// Ensure AdrenoProvider implements GPUProvider
var _ GPUProvider = (*AdrenoProvider)(nil)
