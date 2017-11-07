package main

import (
	"fmt"
	"net"
	"os"

	owl "github.com/billglover/go-owl"
)

func main() {

	// parse the address
	addr, _ := net.ResolveUDPAddr("udp", owl.MulticastAddress)

	// open a connection
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Printf("unable to listen on address: %s: %v\n", owl.MulticastAddress, err)
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
