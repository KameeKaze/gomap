package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Scanner struct {
	Ip      string
	Ports   string
	Timeout int
}

type Port struct {
	Number int
	Open   bool
}

//Terminal colors
const (
	colorRed   = string("\033[31m")
	colorGreen = string("\033[32m")
	colorBlue  = string("\033[34m")
	colorWhite = string("\033[37m")
)

// parse flags from command-line
func (s *Scanner) Setup() {
	flag.StringVar(&s.Ip, "i", "", "IP address or domain")
	flag.StringVar(&s.Ports, "p", "", "Ports separated with comma \nexample: -p 22,80,443")
	flag.IntVar(&s.Timeout, "t", 500, "Set the timeout in milliseconds")
}

func main() {
	//setup scanner
	scanner := Scanner{}
	scanner.Setup()
	flag.Parse()

	//check if ip is set
	if scanner.Ip == "" {
		fmt.Println("Set an IP address with -i")
		return
	}

	// resolve domain to ip && error if unreachable
	var err error
	scanner.Ip, err = domainToIP(scanner.Ip)
	if err != nil {
		fmt.Printf("%sFailed to resolve '%s'%s\n", colorRed, scanner.Ip, colorWhite)
		return
	}

	//create int slice from ports
	var ports []int
	for _, port := range strings.Split(scanner.Ports, ",") {
		p, err := strconv.Atoi(port)
		if err != nil {
			fmt.Printf("%sInvalid ports%s\n", colorRed, colorWhite)
			return
		}
		//check if port in range
		if p > 65535 || p < 0 {
			fmt.Printf("%sPort %d is out of range%s\n", colorRed, p, colorWhite)
			return
		}
		ports = append(ports, p)
	}

	//receive from a channel wether the port is open
	fmt.Println("Scanning ports for", scanner.Ip)
	fmt.Printf("%sPORT   STATUS%s\n", colorBlue, colorWhite)
	open := scanPortsTCP(scanner.Ip, ports, scanner.Timeout)
	for p := range open {
		// red if closed, green if port is open
		switch p.Open { // format with spaces
		case true:
			fmt.Printf("%d %s%sOpen%s\n", p.Number, strings.Repeat(" ", 6-len(strconv.Itoa(p.Number))), colorGreen, colorWhite)
		case false:
			fmt.Printf("%d %s%sClosed%s\n", p.Number, strings.Repeat(" ", 6-len(strconv.Itoa(p.Number))), colorRed, colorWhite)
		}

	}
}

func scanPortsTCP(ip string, ports []int, timeout int) <-chan Port {
	open := make(chan Port)
	go func() {
		//iterate over ports
		for _, port := range ports {
			//parse address
			address := ip + ":" + strconv.Itoa(port)

			//check if port open
			conn, err := net.DialTimeout("tcp", address, time.Millisecond*time.Duration(timeout))
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

//Convert domain into ip adress
func domainToIP(domain string) (string, error) {
	ip, err := net.LookupIP(domain)
	if err != nil {
		return domain, err
	}
	return ip[0].String(), nil
}
