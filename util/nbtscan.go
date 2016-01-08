package util

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func query(ip string, laddr *net.UDPAddr, ch chan string, timeout int) {
	addr, err := net.ResolveUDPAddr("udp", ip+":137")
	if err != nil {
		ch <- ""
		return
	}
	conn, err := net.DialUDP("udp", laddr, addr)
	if err != nil {
		ch <- ""
		return
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	req := []byte{0x80, 0xf0, 0x00, 0x10, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x20, 0x43, 0x4b, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x00, 0x00, 0x21,
		0x00, 0x01}
	_, err = conn.Write(req)
	if err != nil {
		ch <- ""
		return
	}
	resp := make([]byte, 1024)
	_, err = conn.Read(resp)
	if err != nil {
		ch <- ""
		return
	}
	idx := 57

	name_bytes := resp[idx : idx+15]
	for i := 0; i < 15; i++ {
		if name_bytes[i] < 31 || name_bytes[i] > 126 {
			name_bytes[i] = 0x2E
		}
	}
	name := string(resp[idx : idx+15])
	ch <- (ip + " " + name)
}

func Nbtscan(timeout int) {
	ch := make(chan string)
	defer close(ch)

	ipRange := flag.Lookup("r").Value.String()
	lastDotIdx := strings.LastIndex(ipRange, ".")
	ipPre := ipRange[:lastDotIdx+1]

	var ip string
	for i := 0; i < 256; i++ {
		ip = ipPre + strconv.Itoa(i)
		go query(ip, nil, ch, timeout)
	}

	size := 256
	laddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err == nil {
		for i := 0; i < 256; i++ {
			ip = ipPre + strconv.Itoa(i)
			go query(ip, laddr, ch, timeout)
		}
		size = 512
	}

	for i := 0; i < size; i++ {
		result, _ := <-ch
		if len(result) > 0 {
			fmt.Printf("%s\n", result)
		}
	}
}
