package util

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func scanIpAndPort(ip string, port string, ipPortCh chan string, timeout int) {
	conn, err := net.DialTimeout("tcp", ip+":"+port, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		ipPortCh <- ""
		return
	}
	defer conn.Close()
	ipPortCh <- port
}

func scanIp(ip string, ports []string, ipCh chan string, timeout int) {
	ipPortCh := make(chan string, len(ports))
	defer close(ipPortCh)
	for _, port := range ports {
		go scanIpAndPort(ip, port, ipPortCh, timeout)
	}
	output := []string{ip}
	for range ports {
		result, _ := <-ipPortCh
		if len(result) > 0 {
			output = append(output, result)
		}
	}

	if len(output) > 1 {
		ipCh <- strings.Join(output, " ")
	} else {
		ipCh <- ""
	}
}

func Portscan(timeout int) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	ports := strings.Split(scanner.Text(), " ")

	ips := make([]string, 0)

	for scanner.Scan() {
		input := scanner.Text()
		_, err := net.ResolveIPAddr("ip", input)
		if err != nil {
			break
		} else {
			ips = append(ips, input)
		}
	}

	ipCh := make(chan string, len(ips))
	defer close(ipCh)

	for _, ip := range ips {
		go scanIp(ip, ports, ipCh, timeout)
	}

	for range ips {
		result, _ := <-ipCh
		if len(result) > 0 {
			fmt.Printf("%s\n", result)
		}
	}
}
