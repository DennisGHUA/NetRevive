package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kardianos/service"
)

type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}

func (p *program) run() {
	LogInfo("Running NetRevive application")
	verifyAdminRights()
	runningNotAsService()
	Setup()
	go CheckIfEthernetIsWorking()

	// Block the main thread until a signal is received
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	defer LogInfo("Exiting NetRevive")
	<-quit
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	return nil
}

func verifyAdminRights() {
	if isAdmin() {
		//LogInfo("This program is correctly running with admin permissions.")
	} else {
		LogError("This program is NOT running with admin permissions. Exiting in 1 minute...", nil)
		time.Sleep(60 * time.Second)
		LogFatal("Please run the program with admin permissions.", nil)
	}
}

func runningNotAsService() {
	// Prints instructions if ran without args
	if len(os.Args) <= 1 {
		fmt.Println("")
		LogWarning("The program isn't running as a service, which is highly recommended.")
		LogWarning("To install as a Windows service in the current location, run 'NetRevive.exe install'.")
		LogWarning("To uninstall as a Windows service, run 'NetRevive.exe uninstall'.")
		fmt.Println("")
	}

}

func main() {

	svcConfig := &service.Config{
		Name:        "NetRevive",
		DisplayName: "NetRevive",
		Description: "NetRevive application for monitoring ethernet.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		LogFatal("Error running service", err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "install":
			// Install the service if it's not already installed
			// err = s.Install()
			if err := InstallService(); err != nil {
				LogError("Error installing service: ", err)
				return
			}
			LogInfo("Service installed successfully.")
		case "uninstall":
			err = s.Uninstall()
			if err != nil {
				LogError("Error uninstalling service:", err)
				return
			}
			LogInfo("Service uninstalled successfully.")
			LogWarning("\nTo prevent the operating system from getting stuck on a manual input screen after rebooting 3 times,\nthese actions were done: command instructs the system to ignore booting failures and proceed with a reboot.\nCommand: bcdedit /set {current} bootstatuspolicy IgnoreAllFailures\n\nThis command disables automatic and manual repair options during boot, enhancing booting progress.\nIt may prevent interaction screens requiring user input from appearing.\nCommand: bcdedit /set recoveryenabled NO\n\nTo undo the changes made by the previous commands, execute the following commands:\nRe-enable displaying booting failures:\nCommand: bcdedit /set {current} bootstatuspolicy DisplayAllFailures\nRe-enable automatic and manual repair options during boot:\nCommand: bcdedit /set recoveryenabled YES")
		case "start":
			err = s.Start()
			if err != nil {
				LogError("Error starting service:", err)
				return
			}
			LogInfo("Service started successfully.")
		case "stop":
			err = s.Stop()
			if err != nil {
				LogError("Error stopping service:", err)
				return
			}
			LogInfo("Service stopped successfully.")
		case "restart":
			err = s.Restart()
			if err != nil {
				LogError("Error restarting service:", err)
				return
			}
			LogInfo("Service restarted successfully.")
		default:
			LogInfo("Unknown command")
		}
		return
	}

	err = s.Run()
	if err != nil {
		LogFatal("Error running service", err)
	}
}
