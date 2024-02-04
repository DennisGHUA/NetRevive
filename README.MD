# NetRevive

NetRevive is an application for monitoring Ethernet connectivity and ensuring server uptime through automated recovery mechanisms. This program is intended for any Windows server with potentially unstable Ethernet drivers. While it does not fix the root cause of the problem, it does ensure the server will come back online after an Ethernet driver failure.

NetRevive constantly monitors the Ethernet connection to a specified router and does not require a working internet connection. If the connection to the router fails, NetRevive attempts to restore it through various means, including resetting the Ethernet adapter or rebooting the system if all else fails. It also provides options for logging incidents and configuring settings such as the Ethernet adapter name and router IP address.

## Installation

1. Download `NetRevive.exe` from the repository.
2. Place the downloaded `NetRevive.exe` file in a directory of your choice on your system.
3. Open Command Prompt with administrative privileges.
4. Navigate to the directory where `NetRevive.exe` is located using the `cd` command.
5. Execute the program installation by typing `NetRevive.exe install` and pressing Enter.
6. Start the program by typing `NetRevive.exe start` and pressing Enter.

* (Optional) Modify settings in `NetRevive.json` to customize behavior. Refer to the Configuration section below.
* (Optional) You can automate the Windows logon process using Microsoft's Autologon utility. Download it [here](https://learn.microsoft.com/en-us/sysinternals/downloads/autologon).


## Configuration

NetRevive can be configured through a JSON configuration file named `NetRevive.json`. This file contains settings such as the Ethernet adapter name (which can be found at `Control Panel\Network and Internet\Network Connections`), router IP address, and logging preferences. This file gets created after running the program for the first time.

## Uninstallation

Once installed as a Windows service using `NetRevive.exe install`, NetRevive runs in the background and automatically monitors the Ethernet connection. You can uninstall the service using the following command:

`NetRevive.exe uninstall`

To help your server start up faster and avoid getting stuck by requiring user input before it has booted, the following commands were applied:

- `bcdedit /set {current} bootstatuspolicy IgnoreAllFailures`
- `bcdedit /set recoveryenabled NO`

If you ever want to go back to the way things were, just type these commands:

- `bcdedit /set {current} bootstatuspolicy DisplayAllFailures`
- `bcdedit /set recoveryenabled YES`


## Acknowledgements

NetRevive utilizes the following third-party libraries:

- [github.com/kardianos/service](https://github.com/kardianos/service) - Go package for creating and managing Windows services
- [github.com/prometheus-community/pro-bing](https://github.com/prometheus-community/pro-bing) - Go package for ICMP ping operations