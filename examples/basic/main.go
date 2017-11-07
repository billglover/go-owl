package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	owl "github.com/billglover/go-owl"
)

func main() {

	// allow users to set the address and port on which to listen
	bindAddr := flag.String("addr", ":41234", "the address and port on which to listen for readings")
	flag.Parse()

	// parse the address
	addr, err := net.ResolveUDPAddr("udp", *bindAddr)
	if err != nil {
		fmt.Printf("unable to parse address: %s\n", *bindAddr)
		os.Exit(1)
	}

	// open a connection
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Printf("unable to listen on address: %s: %v\n", *bindAddr, err)
		os.Exit(1)
	}
	defer conn.Close()

	for {
		// read from the network
		buf := make([]byte, 1024)
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
			continue
		}

		// decode the electricity reading
		elec, err := owl.Read(buf[:n])
		if err != nil {
			fmt.Println(err)
		}

		// print a log line
		fmt.Printf("%v : electricity reading : power=%.2f%s\n", elec.Timestamp, elec.Chan[0].Power, elec.Chan[0].PowerUnits)
	}
}
