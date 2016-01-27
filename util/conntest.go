package util

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func Conntest(timeout int) {
	log.SetPrefix("[Conntest] ")
	scanner := bufio.NewScanner(os.Stdin)

	results := make(chan string, 10)
	defer close(results)

	count := 0
	for scanner.Scan() {
		split := strings.Split(scanner.Text(), " ")
		go test(split[0], split[1], timeout, results)
		count = count + 1
	}

	for i := 0; i < count; i++ {
		res := <-results
		if len(res) > 0 {
			fmt.Printf("%s\n", res)
		}
	}
}

func test(host, port string, timeout int, results chan string) {
	start := time.Now()
	conn, err := net.DialTimeout("tcp", host+":"+port, time.Duration(timeout)*time.Millisecond)
	if err != nil {
		log.Printf("error occur: %v\n", err)
		results <- ""
		return
	}
	defer conn.Close()
	end := time.Now()
	ms := end.Sub(start).Seconds() * 1000
	results <- host + " " + fmt.Sprintf("%.2f", ms)
}
