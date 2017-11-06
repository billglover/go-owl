package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/billglover/go-owl/lib"

	"github.com/prometheus/client_golang/prometheus"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	counter := prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "readings",
		Subsystem: "electricity",
		Namespace: "home",
		Help:      "number of meter readings received since restart",
	})
	err := prometheus.Register(counter)
	if err != nil {
		log.Fatalf("unable to register counter: %v", err)
	}

	powerGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "power",
		Subsystem: "electricity",
		Namespace: "home",
		Help:      "instantaneous power consumption",
	})
	err = prometheus.Register(powerGauge)
	if err != nil {
		log.Fatalf("unable to register power gauge: %v", err)
	}

	batteryGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "battery",
		Subsystem: "electricity",
		Namespace: "home",
		Help:      "percentage battery remaining",
	})
	err = prometheus.Register(batteryGauge)
	if err != nil {
		log.Fatalf("unable to register battery gauge: %v", err)
	}

	rssiGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "rssi",
		Subsystem: "electricity",
		Namespace: "home",
		Help:      "received signal strength indicator",
	})
	err = prometheus.Register(rssiGauge)
	if err != nil {
		log.Fatalf("unable to register RSSI gauge: %v", err)
	}

	lqiGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "lqi",
		Subsystem: "electricity",
		Namespace: "home",
		Help:      "link quality indicator",
	})
	err = prometheus.Register(lqiGauge)
	if err != nil {
		log.Fatalf("unable to register LQI gauge: %v", err)
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

	go listen(conn, counter, powerGauge, batteryGauge, rssiGauge, lqiGauge)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func listen(conn *net.UDPConn, counter prometheus.Counter, powerGauge, batteryGauge, rssiGauge, lqiGauge prometheus.Gauge) {

	for {
		buf := make([]byte, 1024)

		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		elec, err := owl.Read(buf[:n])
		if err != nil {
			fmt.Println(err)
		}

		counter.Inc()
		powerGauge.Set(elec.Chan[0].Power)
		batteryGauge.Set(elec.Battery)
		rssiGauge.Set(elec.RSSI)
		lqiGauge.Set(elec.LQI)

		fmt.Printf("%v : electricity reading : power=%.2f%s\n", elec.Timestamp, elec.Chan[0].Power, elec.Chan[0].PowerUnits)
	}
}
