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
	Ip      string
	Ports   string
	Timeout int
}

type Port struct {
	Number int
	Open   bool
}

//Terminal colors
const(
    colorRed   = string("\033[31m")
    colorGreen = string("\033[32m")
	colorBlue  = string("\033[34m")
    colorWhite = string("\033[37m")
)

// parse flags from command-line
func (s *Scanner) Setup() {
	flag.StringVar(&s.Ip, "i", "", "IP address")
	flag.StringVar(&s.Ports, "p", "", "Ports \nexample: -p 22,80,443")
	flag.IntVar(&s.Timeout, "t", 500, "Set the timeout in milliseconds\n")
}

func main() {
	//setup scanner
	scanner := Scanner{}
	scanner.Setup()
	flag.Parse()

	//check if ip is set
	if scanner.Ip == ""{
		fmt.Println("Set an IP address with -i")
		return
	}

	// resolve domain to ip && error if unreachable
	var err error
	scanner.Ip, err = domainToIP(scanner.Ip)
	if err != nil{
		fmt.Printf("Failed to resolve '%s'\n", scanner.Ip)
		return
	}

	//create int slice from ports
	var ports []int
	for _, port := range strings.Split(scanner.Ports, ",") {
		p, err := strconv.Atoi(port)
		if err != nil{
			fmt.Println("Invalid ports")
			return
		}
		//check if port in range
		if p > 65535 || p < 0 {
			fmt.Printf("Port %d is out of range\n",p)
			return
		}
		ports = append(ports, p)
	}


	//receive from a channel wether the port is open
	fmt.Println("Scanning ports for", scanner.Ip)
	fmt.Printf("%sPORT   STATUS%s\n",colorBlue,colorWhite)
	open := ScanPortsTCP(scanner.Ip, ports, scanner.Timeout)
	for p := range open {
		// red if closed, green if port is open
		switch p.Open {									  // format with spaces
		case true: fmt.Printf("%d %s%sOpen%s\n",p.Number, strings.Repeat(" ", 6-len(strconv.Itoa(p.Number))), colorGreen, colorWhite)
		
		case false: fmt.Printf("%d %s%sClosed%s\n",p.Number, strings.Repeat(" ", 6-len(strconv.Itoa(p.Number))), colorRed, colorWhite)
		}
		
	}
}

func ScanPortsTCP(ip string, ports []int, timeout int) <-chan Port {
	open := make(chan Port)

	go func() {
		//iterate over ports
		for _, port := range ports {
			//parse address
			address := ip + ":" + strconv.Itoa(port)

			//check if port open
			conn, err := net.DialTimeout("tcp", address, time.Millisecond * time.Duration(timeout))
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

	//Convert domain into ip adress
func domainToIP(domain string) (string, error) {
	ip, err := net.LookupIP(domain)
	if err != nil{
		return domain, err
	}
	return ip[0].String(), nil
}