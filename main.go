package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"time"
	"strings"
)

type Scanner struct {
	Ip    string
	Ports string
}

type Port struct {
	Number int
	Open   bool
}

// parse flags from command-line
func (s *Scanner) Setup() {
	flag.StringVar(&s.Ip, "i", "", "IP address")
	flag.StringVar(&s.Ports, "p", "", "Ports \nexample: -p 22,80,443")
}


func main() {
	//setup scanner
	scanner := Scanner{}
	scanner.Setup()
	flag.Parse()

	// validate ip address
	if validateIP(scanner.Ip) {
		fmt.Println("Invalid ip address")
		return
	}

	//some testing ports
	var ports []int
	for _, port := range strings.Split(scanner.Ports, ",") {
		p, err := strconv.Atoi(port)
		if err != nil{
			fmt.Println("Invalid ports")
			return
		}
		ports = append(ports, p)
	}

	//receive from a channel wether the port is open
	open := ScanPortsTCP(scanner.Ip, ports)
	for p := range open {
		fmt.Println(p.Number, p.Open)
	}
}

func ScanPortsTCP(ip string, ports []int) <-chan Port {
	open := make(chan Port)

	go func() {
		//iterate over ports
		for _, port := range ports {
			//parse address
			address := ip + ":" + strconv.Itoa(port)

			//check if port open
			conn, err := net.DialTimeout("tcp", address, time.Second*2)
			switch err {
			case nil: //open
				defer conn.Close()
				open <- Port{Number: port, Open: true}
			default: // closed
				open <- Port{Number: port, Open: false}
			}
		}
		//close channel
		close(open)
	}()
	return open
}

func validateIP(ip string) bool {
	return net.ParseIP(ip) == nil
}
