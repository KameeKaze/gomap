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

const(

    colorRed   = "\033[31m"
    colorGreen = "\033[32m"
    colorWhite = "\033[37m"
)
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

	// convert domain into ip address
	var err error
	scanner.Ip, err = domainToIP(scanner.Ip)
	if err != nil{
		fmt.Println("Invalid ip address")
		return
	}
	
	// validate ip address
	if validateIP(scanner.Ip) {
		fmt.Println("Invalid ip address")
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
		ports = append(ports, p)
	}

	//receive from a channel wether the port is open
	fmt.Println("Scanning ports for", scanner.Ip)
	open := ScanPortsTCP(scanner.Ip, ports)
	for p := range open {
		// red if closed, green if port is open
		switch p.Open {
		case true: fmt.Println(p.Number, string(colorGreen), "Open",string(colorWhite))
		case false: fmt.Println(p.Number, string(colorRed), "Closed",string(colorWhite))
		}
		
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

	//Convert domain into ip adress
func domainToIP(domain string) (string, error) {
	ip, err := net.LookupIP(domain)
	return ip[0].String(), err
}