package main

import(
	"fmt"
	"net"
	"strconv"
	"time"
)


func main(){
	open := ScanPortTCP("127.0.0.1", 3306)
	fmt.Println(open)
}

func ScanPortTCP(ip string, port int,) bool {
	address := ip + ":" + strconv.Itoa(port)

	conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}