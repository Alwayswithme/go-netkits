// netkit project main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"netkits/util"
)

func main() {
	flag.String("c", "", "command, support[ping|portscan|nbtscan|portscan2]")
	flag.String("r", "", "ip address range to scan, when use nbtscan")
	debugPtr := flag.Bool("v", false, "print error message to stderr")
	toPtr := flag.Int("t", 500, "command timeout")

	flag.Parse()
	subcmd := flag.Lookup("c").Value.String()
	if !*debugPtr {
		log.SetOutput(ioutil.Discard)
	}
	timeout := *toPtr

	switch subcmd {
	case "ping":
		util.Ping(timeout)
	case "portscan":
		util.Portscan(timeout)
	case "nbtscan":
		util.Nbtscan(timeout)
	case "portscan2":
		util.Conntest(timeout)
	default:
		fmt.Println("invalid command -- " + subcmd)
		fmt.Println("(use -h for help)")
	}
}
