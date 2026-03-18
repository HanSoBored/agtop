package providers

import (
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
	data := GetProp("ro.hardware.vulkan")
	vulkanDriver := strings.TrimSpace(data)
	return strings.Contains(strings.ToLower(vulkanDriver), "adreno") ||
		strings.Contains(strings.ToLower(vulkanDriver), "freedreno") ||
		FileExists("/sys/class/kgsl/kgsl-3d0")
}

// GetStats retrieves Adreno GPU statistics
func (p *AdrenoProvider) GetStats() (GPUStats, error) {
	stats := GPUStats{
		Vendor: "Adreno",

		// Identity
		Model:           strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/gpu_model")),
		Platform:        strings.TrimSpace(GetProp("ro.board.platform")),
		VulkanDriver:    strings.TrimSpace(GetProp("ro.hardware.vulkan")),
		EGLDriver:       strings.TrimSpace(GetProp("ro.hardware.egl")),
		OpenGLESVersion: strings.TrimSpace(GetProp("ro.opengles.version")),
		MinClock:        strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/min_clock_mhz")),
		MaxClock:        strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/max_clock_mhz")),
		Governor:        strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/devfreq/governor")),
		IdleTimer:       strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/idle_timer")),
		ThermalPwrlevel: strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/thermal_pwrlevel")),
		ResetCount:      strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/reset_count")),
		PreemptCount:    strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/preempt_count")),
		PreemptionMode:  strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/preemption")),
		Temperature:     strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/temp")),
		NumPwrLevels:    strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/num_pwrlevels")),
		MinPwrLevel:     strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/min_pwrlevel")),
		MaxPwrLevel:     strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/max_pwrlevel")),
		IFPCCount:       strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/ifpc_count")),
		FTPageFault:     strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/ft_pagefault_policy")),
		FTPolicy:        strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/ft_policy")),
		FTLongIB:        strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/ft_long_ib_detect")),
		FTHangIntr:      strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/ft_hang_intr_status")),
	}

	// Current Clock
	currClkStr := strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/gpuclk"))
	if currClkStr != "N/A" {
		val, _ := strconv.Atoi(currClkStr)
		stats.CurrentClock = val / 1000000
	}

	// Target Clock
	targetClkStr := strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/devfreq/target_freq"))
	if targetClkStr != "N/A" {
		val, _ := strconv.Atoi(targetClkStr)
		stats.TargetClock = val / 1000000
	}

	// GPU Busy
	busyStr := strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/gpu_busy_percentage"))
	if busyStr != "N/A" {
		val, _ := strconv.Atoi(strings.TrimSuffix(busyStr, "%"))
		stats.GpuBusyPercent = val
	}

	// Devfreq Load
	devfreqStr := strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/devfreq/gpu_load"))
	if devfreqStr != "N/A" {
		val, _ := strconv.Atoi(devfreqStr)
		stats.DevfreqLoad = val
	}

	// Effective BusyPercentage
	if stats.DevfreqLoad > stats.GpuBusyPercent {
		stats.BusyPercentage = stats.DevfreqLoad
	} else {
		stats.BusyPercentage = stats.GpuBusyPercent
	}

	// Boolean flags
	stats.HWClockGating = strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/hwcg")) == "1"
	stats.IdlePowerCollapse = strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/ifpc")) == "1"
	stats.Throttling = strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/throttling")) == "1"

	// Available frequencies
	freqsStr := strings.TrimSpace(ReadFile("/sys/class/kgsl/kgsl-3d0/freq_table_mhz"))
	if freqsStr != "N/A" {
		stats.AvailableFreqs = strings.Fields(freqsStr)
	}

	// Format temperature
	if stats.Temperature != "N/A" {
		temp, err := strconv.ParseFloat(stats.Temperature, 64)
		if err == nil && temp > 1000 {
			stats.Temperature = strconv.FormatFloat(temp/1000.0, 'f', '1', 64) + "°C"
		} else {
			stats.Temperature += "°C"
		}
	}

	return stats, nil
}

// Ensure AdrenoProvider implements GPUProvider
var _ GPUProvider = (*AdrenoProvider)(nil)
