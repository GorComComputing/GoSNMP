package main

import (
    //"fmt"
    "os"
)

	
type Cmd struct {
	addr 	func([]string) string
	descr string     
}

// Command list for interpretator
var cmd =  map[string]Cmd{ 
	//"tst": Cmd{addr: cmd_tst, descr: "Test command"},
	//"ls": Cmd{addr: cmd_ls, descr: "Test command: print all file names from catalog (ls)"},

	".quit": Cmd{addr: cmd_quit, descr: "Exit from this program"},
	".help": Cmd{addr: cmd_help, descr: "Print this Help"},
	
	"get": Cmd{addr: cmd_get, descr: "Get"},
	"get_param": Cmd{addr: cmd_get_param, descr: "Get with set params (IP address, port)"},
	"get_v3": Cmd{addr: cmd_get_v3, descr: "Get by v3"},
	"set": Cmd{addr: cmd_set, descr: "Set variable"},
	"get_hex": Cmd{addr: cmd_get_hex, descr: "Get hex string"},
	"walk": Cmd{addr: cmd_walk, descr: "Walk"},
	"trap_v1": Cmd{addr: cmd_trap_v1, descr: "Send trap v1"},
	"trap_v2": Cmd{addr: cmd_trap_v2, descr: "Send trap v2"},
	"trap_v3": Cmd{addr: cmd_trap_v3, descr: "Send trap v3"},
	"trap_srv": Cmd{addr: cmd_trap_srv, descr: "Trap server"},
	
}



// Interpretator 
func interpretator(words []string) string {
	if _, ok := cmd[words[0]]; ok {
		return cmd[words[0]].addr(words)
	} else{
		return "Unknown command: " + words[0] + "\n"
	}
}


// HELP - Print command list
var cmd_print = make(map[string]Cmd)
func cmd_help(words []string) string {
	var output string
	for key, val := range cmd_print {
		output += key 
		for i := len(key); i < 10; i++ {
			output += " "
		} 
		output += " - " + val.descr + "\n"
	}
	return output
}


// Exit from this program
func cmd_quit(words []string) string {
	os.Exit(0)
	return ""
}


