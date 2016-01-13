package util

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type ipAndPort struct {
	ip, port string
}

func scanIpAndPort(addr ipAndPort, timeout int) bool {
	conn, err := net.DialTimeout("tcp", addr.ip+":"+addr.port, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true

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

	jobCount := len(ips) * len(ports)
	jobs := make(chan ipAndPort, jobCount)
	results := make(chan string, len(ips))
	defer close(results)
	for i := 0; i < 500; i++ {
		go worker(jobs, timeout, results)
	}

	for _, ip := range ips {
		for _, port := range ports {

			jobs <- ipAndPort{ip, port}
		}
	}

	close(jobs)

	m := make(map[string][]string)
	for i := 0; i < jobCount; i++ {
		res := <-results
		if len(res) > 0 {
			ip := strings.Split(res, ":")
			m[ip[0]] = append(m[ip[0]], ip[1])
		}

	}

	for k, v := range m {
		fmt.Printf("%s %s\n", k, strings.Join(v, " "))
	}
}

func worker(jobs chan ipAndPort, timeout int, results chan string) {
	for addr := range jobs {
		if scanIpAndPort(addr, timeout) {
			results <- (addr.ip + ":" + addr.port)
		} else {
			results <- ""
		}
	}
}
