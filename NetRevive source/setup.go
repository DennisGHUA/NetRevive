package main

import (
	"encoding/json"
	"fmt"
	"github.com/kardianos/service"
	"golang.org/x/sys/windows/svc/mgr"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// To prevent the operating system from getting stuck on a manual input screen after rebooting 3 times,
// this command instructs the system to ignore booting failures and proceed with a reboot.
// Command: bcdedit /set {current} bootstatuspolicy IgnoreAllFailures

// This command disables automatic and manual repair options during boot, enhancing booting progress.
// It may prevent interaction screens requiring user input from appearing.
// Command: bcdedit /set recoveryenabled NO

// To undo the changes made by the previous commands, execute the following commands:
// Re-enable displaying booting failures:
// Command: bcdedit /set {current} bootstatuspolicy DisplayAllFailures
// Re-enable automatic and manual repair options during boot:
// Command: bcdedit /set recoveryenabled YES

const (
	configFileName = "NetRevive.json"
)

func Setup() {
	setBootOptions()
	loadSettings()
}

var (
	EthernetAdapterName string
	RouterIpAddress     string
	LogIncidents        bool
)

// Struct to hold settings
type Settings struct {
	EthernetAdapterName string `json:"ethernet_adapter_name"`
	RouterIpAddress     string `json:"router_ip_address"`
	LogIncidents        bool   `json:"log_incidents"`
}

func loadSettings() {
	// Define default settings
	defaultSettings := Settings{
		EthernetAdapterName: "Ethernet",
		RouterIpAddress:     "192.168.178.1",
		LogIncidents:        true,
	}

	// Get the path of the executable
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}

	// Get the directory of the executable
	exeDir := filepath.Dir(exePath)

	// Construct the config file path relative to the executable directory
	configFilePath := filepath.Join(exeDir, configFileName)

	// Check if config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// Config file does not exist, create it with default settings
		saveSettings(defaultSettings)
		LogInfo(fmt.Sprintf("Created %v", configFileName))
		EthernetAdapterName = defaultSettings.EthernetAdapterName
		RouterIpAddress = defaultSettings.RouterIpAddress
		LogIncidents = defaultSettings.LogIncidents
		LogInfo(fmt.Sprintf("ethernet_adapter_name: %s", EthernetAdapterName))
		LogInfo(fmt.Sprintf("router_ip_address: %s", RouterIpAddress))
		LogInfo(fmt.Sprintf("log_incidents: %v", LogIncidents))
		return
	}

	// Read config file
	file, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		LogFatal("Error reading config file", err)
	}

	// Unmarshal config data into settings struct
	var settings Settings
	err = json.Unmarshal(file, &settings)
	if err != nil {
		LogFatal("Error unmarshalling config data", err)
	}

	// Update global variables with loaded settings
	EthernetAdapterName = settings.EthernetAdapterName
	RouterIpAddress = settings.RouterIpAddress
	LogIncidents = settings.LogIncidents

	LogInfo(fmt.Sprintf("Settings loaded from %v", configFilePath))
	LogInfo(fmt.Sprintf("ethernet_adapter_name: %s", EthernetAdapterName))
	LogInfo(fmt.Sprintf("router_ip_address: %s", RouterIpAddress))
	LogInfo(fmt.Sprintf("log_incidents: %v", LogIncidents))
}

// Function to save settings to config file
func saveSettings(settings Settings) {
	// Get the path of the executable
	exePath, err := os.Executable()
	if err != nil {
		log.Fatalf("Error getting executable path: %v", err)
	}

	// Get the directory of the executable
	exeDir := filepath.Dir(exePath)

	// Construct the config file path relative to the executable directory
	configFilePath := filepath.Join(exeDir, configFileName)

	// Marshal settings struct into JSON
	data, err := json.MarshalIndent(settings, "", "    ")
	if err != nil {
		log.Fatalf("Error marshalling settings: %v", err)
	}

	// Write JSON data to config file
	err = ioutil.WriteFile(configFilePath, data, 0644)
	if err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
}

func isAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func setBootOptions() {
	//LogInfo("Changing two windows settings to prevent popups that require user input from appearing.")

	commands := []string{
		"bcdedit /set {current} bootstatuspolicy IgnoreAllFailures",
		"bcdedit /set recoveryenabled NO",
	}

	for _, cmd := range commands {
		err := runCommand(cmd)
		if err != nil {
			LogFatal("Error setting boot option", err)
		}
	}
}

func runCommand(cmd string) error {
	command := exec.Command("cmd", "/C", cmd)
	err := command.Run()
	return err
}

func isServiceInstalled() bool {
	m, err := mgr.Connect()
	if err != nil {
		return false
	}
	defer m.Disconnect()

	_, err = m.OpenService("NetRevive")
	return err == nil
}

func InstallService() error {
	if isServiceInstalled() {
		return nil
	}

	svcConfig := &service.Config{
		Name:        "NetRevive",
		DisplayName: "NetRevive",
		Description: "NetRevive: Keeping Servers Online with Automated Recovery",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		return err
	}

	err = s.Install()
	if err != nil {
		return err
	}
	return nil
}
