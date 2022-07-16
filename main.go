package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
	"sync"
)

type Scanner struct {
	Ip      string
	Ports   string
	Timeout int
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

//create waitgroup
var wg sync.WaitGroup

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

	//wait for all ports to be scanned
	wg.Add(len(ports))
	for _, port := range ports{
		go scanPortTCP(scanner.Ip, port, scanner.Timeout)
	}
	//wait until all ports scanned
	wg.Wait()

}

func scanPortTCP(ip string, port, timeout int){
	defer wg.Done()
	//parse address
	address := ip + ":" + strconv.Itoa(port)
	//check if port is open
	conn, err := net.DialTimeout("tcp", address, time.Millisecond*time.Duration(timeout))
	switch err {
	case nil: //open
		defer conn.Close()
		fmt.Printf("%d %s%sOpen%s\n", port, strings.Repeat(" ", 6-len(strconv.Itoa(port))), colorGreen, colorWhite)
		return
	default: // closed
		fmt.Printf("%d %s%sClosed%s\n", port, strings.Repeat(" ", 6-len(strconv.Itoa(port))), colorRed, colorWhite)
		return
	}
}

//Convert domain into ip adress
func domainToIP(domain string) (string, error) {
	ip, err := net.LookupIP(domain)
	if err != nil {
		return domain, err
	}
	return ip[0].String(), nil
}
