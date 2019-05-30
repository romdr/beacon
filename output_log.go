package main

import "log"

// Write to log
func sendToLog(hostMetrics *HostMetrics, arg string) {
	log.Printf("%+v\n", *hostMetrics)
}
