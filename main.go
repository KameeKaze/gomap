package main

import(
	"fmt"
	"net"
	"strconv"
	"time"
	"flag"
)

type Scanner struct {
	Ip   string `required:"true"`
	Port int    `required:"true"`
	
}

// parse flags from command-line
func (s *Scanner) Setup(){
	flag.StringVar(&s.Ip, "i", "", "ip address")
	flag.IntVar(&s.Port, "p", 0, "port to scan")
}

func main(){
	//setup scanner
	scanner := Scanner{}
	scanner.Setup()
	flag.Parse()

	//scan port
	open := ScanPortTCP(scanner.Ip, scanner.Port)
	fmt.Println(open)
}

func ScanPortTCP(ip string, port int,) bool {
	//parse address
	address := ip + ":" + strconv.Itoa(port)

	// check port open
	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}