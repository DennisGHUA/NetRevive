package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/prometheus-community/pro-bing"
)

const (
	maxAttemptsForConnectionReboot = 6 // This one should be higher than maxAttemptsForConnectionResets
	maxAttemptsForConnectionResets = 4 // 12*10sec = 2 minutes
	pingerTimeout                  = 20
)

func CheckIfEthernetIsWorking() {
	LogInfo("Starting ethernet checker module")

	// Timeout before starting to check ethernet connection, this gives the computer time to reconnect after a reboot
	time.Sleep(5 * time.Minute)

	// Enable ethernet driver to prevent crashing
	toggleEthernetAdapterOn(true, true)
	time.Sleep(10 * time.Second)

	for {
		checkEthernet(0, 2)
	}
}

func checkEthernet(failedAttempts int, packetCount int) {

	if packetCount > 8 {
		packetCount = 8
	}

	results := packetPinger(RouterIpAddress, packetCount)

	if failedAttempts >= maxAttemptsForConnectionReboot {
		rebootPCWithEthernetEnabled()
	} else if failedAttempts >= maxAttemptsForConnectionResets {
		restoreEthernetConnection()
	}

	// Not a single package was received
	if results == nil || results.PacketsRecv == 0 {
		failedAttempts++

		LogWarning(fmt.Sprintf("Unable to ping the router try: %d/%v", failedAttempts, maxAttemptsForConnectionReboot))
		//LogInfo(fmt.Sprintf("Trying to send %v packets", packetCount))
		if results != nil {
			LogWarning(fmt.Sprintf("Sent: %v Recv: %v Loss: %v%%", results.PacketsSent, results.PacketsRecv, results.PacketLoss))
		} else {
			LogWarning("Pinger results is null")
		}

		// Recursive multiply package count
		time.Sleep(10 * time.Second) // Everything is NOT OK check again in 10 seconds
		checkEthernet(failedAttempts, packetCount*2)
	} else {
		LogInfo(fmt.Sprintf("Ping OK -> Sent: %v Recv: %v Loss: %v%%", results.PacketsSent, results.PacketsRecv, results.PacketLoss))
		time.Sleep(1 * time.Minute) // Everything is OK check again in a minute
	}

}

func packetPinger(ip string, packetCount int) *probing.Statistics {
	pinger, err := probing.NewPinger(ip)
	if err != nil {
		LogError("Error creating pinger:", err)
		return nil
	}
	defer pinger.Stop()

	pinger.SetPrivileged(true)
	pinger.Count = packetCount
	pinger.Timeout = pingerTimeout * time.Second // 60
	pinger.Interval = 1 * time.Second            // 1

	err = pinger.Run() // Blocks until finished.
	if err != nil {
		LogError("Error running pinger:", err)
		return nil
	}

	return pinger.Statistics() // get send/receive/duplicate/rtt stats
}

func restoreEthernetConnection() {
	LogWarning("Resetting ethernet adapter")
	toggleEthernetAdapterOn(false, false)
	time.Sleep(5 * time.Second) // 30
	toggleEthernetAdapterOn(true, false)
	time.Sleep(20 * time.Second)
	LogWarning("Ethernet adapter has been reset")
}

func toggleEthernetAdapterOn(enable bool, silent bool) {
	action := "disable"
	if enable {
		action = "enable"
	}

	cmd := exec.Command("netsh", "interface", "set", "interface", "name="+EthernetAdapterName, action)
	output, err := cmd.CombinedOutput()
	if err != nil {
		LogError("Error running the command: %v", err)
	}

	if !silent {
		LogWarning(fmt.Sprintf("Ethernet adapter %sd -> Command output: %v", action, strings.TrimSpace(string(output))))
	} else {
		LogInfo(fmt.Sprintf("Ethernet adapter %sd -> Command output: %v", action, strings.TrimSpace(string(output))))
	}
}

func rebootPCWithEthernetEnabled() {
	// Enable the Ethernet adapter
	toggleEthernetAdapterOn(true, false)

	time.Sleep(5 * time.Second)

	// Reboot the Windows computer
	cmd := exec.Command("shutdown", "/r")

	// Run the command
	output, err := cmd.CombinedOutput()
	if err != nil {
		LogFatal("Failed to reboot PC: %v", err)
		return
	}

	// Log the rebooting process
	LogWarning(fmt.Sprintf("Rebooting the system. -> Command output: %s", strings.TrimSpace(string(output))))
	time.Sleep(5 * time.Minute) // Pc should reboot within 1 minute
}
