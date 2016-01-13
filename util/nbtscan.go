package util

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

func query(raddr *net.UDPAddr, laddr *net.UDPAddr, ch chan<- string, timeout int, req []byte) {

	name := dial(raddr, laddr, timeout, req)
	ch <- name
	name = dial(raddr, nil, timeout, req)
	ch <- name

}
func dial(raddr *net.UDPAddr, laddr *net.UDPAddr, timeout int, req []byte) string {
	conn, err := net.DialUDP("udp", laddr, raddr)
	if err != nil {
		return ""
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Millisecond))
	_, err = conn.Write(req)
	if err != nil {
		return ""
	}
	resp := make([]byte, 1024)
	_, err = conn.Read(resp)
	if err != nil {
		return ""
	}
	idx := 57

	name_bytes := resp[idx : idx+15]
	for i, b := range name_bytes {
		if b < 31 || b > 126 {
			name_bytes[i] = 0x2E
		}
	}
	name := string(name_bytes)
	return (fmt.Sprintf("%v ", raddr.IP) + name)
}

func Nbtscan(timeout int) {
	size := 256
	ch := make(chan string, size)
	defer close(ch)

	ipRange := flag.Lookup("r").Value.String()
	lastDotIdx := strings.LastIndex(ipRange, ".")
	ipPre := ipRange[:lastDotIdx+1]

	laddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:0")
	if err != nil {
		return
	}

	req := []byte{0x80, 0xf0, 0x00, 0x10, 0x00, 0x01, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x20, 0x43, 0x4b, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41, 0x41,
		0x41, 0x41, 0x41, 0x41, 0x41, 0x00, 0x00, 0x21,
		0x00, 0x01}

	for i := 0; i < size; i++ {
		ip := ipPre + strconv.Itoa(i)
		raddr, err := net.ResolveUDPAddr("udp", ip+":137")

		if err != nil {
			continue
		}
		go query(raddr, laddr, ch, timeout, req)
	}

	for i := 0; i < size*2; i++ {
		result, _ := <-ch
		if len(result) > 0 {
			fmt.Println(result)
		}
	}
}
