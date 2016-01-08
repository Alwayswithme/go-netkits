// netkit project main.go
package main

import (
	"flag"
	"fmt"
	"netkits/util"
)

func main() {
	flag.String("c", "", "command, support[ping|portscan|nbtscan]")
	flag.String("r", "", "address to scan, when use nbtscan")
	toPtr := flag.Int("t", 500, "command timeout")

	flag.Parse()

	subcmd := flag.Lookup("c").Value.String()
	timeout := *toPtr

	switch subcmd {
	case "ping":
		util.Ping(timeout)
	case "portscan":
		util.Portscan(timeout)
	case "nbtscan":
		util.Nbtscan(timeout)
	default:
		fmt.Println("invalid command -- " + subcmd)
		fmt.Println("(use -h for help)")
	}
}
