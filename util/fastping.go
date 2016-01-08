package util

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	fastping "github.com/tatsushid/go-fastping"
)

func Ping(timeout int) {
	scanner := bufio.NewScanner(os.Stdin)

	ips := make([]string, 0)
	var input string
	// single ip
	for scanner.Scan() {
		input = scanner.Text()
		_, err := net.ResolveIPAddr("ip", input)
		if err != nil {
			break
		} else {
			ips = append(ips, input)
		}
	}
	// whole lan
	if strings.Contains(input, "/") {
		lastDotIdx := strings.LastIndex(input, ".")
		ipPre := input[:lastDotIdx+1]
		for i := 0; i < 256; i++ {
			ip := ipPre + strconv.Itoa(i)
			ips = append(ips, ip)
		}
	}

	runFastPing(ips, "", timeout)
	runFastPing(ips, "127.0.0.1", timeout)

}

func runFastPing(ips []string, src string, timeout int) {
	p := fastping.NewPinger()
	p.Network("udp")
	if len(src) > 0 {
		p.Source(src)
	}

	for _, ip := range ips {
		p.AddIP(ip)
	}

	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("%s %d\n", addr.String(), rtt.Nanoseconds()/1000000)
	}
	p.MaxRTT = time.Duration(timeout) * time.Millisecond
	err := p.Run()
	if err != nil {
		fmt.Println(err)
	}
}
