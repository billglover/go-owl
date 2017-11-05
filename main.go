package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Signal struct {
	RSSI int `xml:"rssi,attr"`
	LQI  int `xml:"lqi,attr"`
}

type Chan struct {
	ID     int     `xml:"id,attr"`
	Power  Reading `xml:"curr"`
	Energy Reading `xml:"day"`
}

type Reading struct {
	Units string  `xml:"units,attr"`
	Value float64 `xml:",chardata"`
}

type Battery struct {
	Level string `xml:"level,attr"`
}

type ElecReading struct {
	XMLName  xml.Name `xml:"electricity"`
	ID       string   `xml:"id,attr"`
	Time     int64    `xml:"timestamp"`
	Signal   Signal   `xml:"signal"`
	Battery  Battery  `xml:"battery"`
	Channels []Chan   `xml:"chan"`
}

func main() {
	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "reading",
		Subsystem: "electricity",
		Namespace: "zhujia",
		Help:      "number of meter readings received",
	})
	err := prometheus.Register(counter)
	if err != nil {
		log.Fatalf("unable to register counter: %v", err)
	}

	powerGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "power",
		Subsystem: "electricity",
		Namespace: "zhujia",
		Help:      "instantaneous power consumption",
	})
	err = prometheus.Register(powerGauge)
	if err != nil {
		log.Fatalf("unable to register power gauge: %v", err)
	}

	//addr, err := net.ResolveUDPAddr("udp", "224.192.32.19:22600")
	addr, err := net.ResolveUDPAddr("udp", ":41234")
	if err != nil {
		log.Fatalf("unable to parse multicast address")
	}

	//conn, err := net.ListenMulticastUDP("udp", nil, addr)
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("unable to listen to multicast address")
	}
	defer conn.Close()

	go listen(conn, counter, powerGauge)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listen(conn *net.UDPConn, counter prometheus.Counter, powerGauge prometheus.Gauge) {

	log.Printf("listening to: %s", conn.LocalAddr())

	for {
		buf := make([]byte, 1024)

		n, addr, err := conn.ReadFromUDP(buf)
		//fmt.Println("Received ", string(buf[0:n]), " from ", addr)

		if err != nil {
			fmt.Println("Error: ", err)
		}

		r := ElecReading{}
		err = xml.Unmarshal(buf[:n], &r)
		if err != nil {
			fmt.Println("Error: ", err)
		}
		fmt.Println()
		fmt.Println("ID:", r.ID)
		fmt.Println("Source:", addr)
		fmt.Println("Time:", time.Unix(r.Time, 0))
		fmt.Println("Signal:", r.Signal.RSSI, r.Signal.LQI)
		fmt.Println("Battery:", r.Battery.Level)
		fmt.Println("Channel:", r.Channels[0].ID)
		fmt.Println("Power:", r.Channels[0].Power.Value, r.Channels[0].Power.Units)
		fmt.Println("Energy:", r.Channels[0].Energy.Value, r.Channels[0].Energy.Units)
		counter.Inc()
		powerGauge.Set(r.Channels[0].Power.Value)
	}

}
