package main

import (
	"log"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// HostMetrics contains stats exported by the beacon
type HostMetrics struct {
	Hostname   string
	HostID     string
	CPUPercent float64
	MemPercent float64
	Uptime     uint64
}

// Map of target type to send function
var targetFuncMap = map[string]func(*HostMetrics, string){
	"cloudwatch": sendToCloudwatch,
	"log":        sendToLog,
	"url":        sendToURL,
}

// Call a given function at an interval
func doEvery(d time.Duration, f func(*Config), config *Config) {
	go f(config)
	for range time.Tick(d) {
		go f(config)
	}
}

// Read system metrics, and send them to target(s)
func heartbeat(config *Config) {
	// Sampple
	cpuPct, _ := cpu.Percent(1*time.Second, false)
	hostInfo, _ := host.Info()
	memVal, _ := mem.VirtualMemory()

	hostMetrics := HostMetrics{
		Hostname:   hostInfo.Hostname,
		HostID:     hostInfo.HostID,
		CPUPercent: cpuPct[0],
		MemPercent: memVal.UsedPercent,
		Uptime:     hostInfo.Uptime,
	}

	// Send merics
	sendMetrics(&hostMetrics, config)
}

// Send metrics to target(s)
func sendMetrics(hostMetrics *HostMetrics, config *Config) {
	for _, target := range config.Targets {
		if fn, ok := targetFuncMap[target.Type]; ok {
			fn(hostMetrics, target.Arg)
		} else {
			log.Fatalf("ERROR: Unsupported target type: %s", target.Type)
		}
	}
}

func main() {
	// Load the configuration
	var config Config
	config.load()

	// Start the periodic heartbeat
	doEvery(config.Interval, heartbeat, &config)
}
