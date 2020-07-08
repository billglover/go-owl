package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	owl "github.com/billglover/go-owl"
	"github.com/rcrowley/go-metrics"
	"github.com/wavefronthq/go-metrics-wavefront/reporting"
	"github.com/wavefronthq/wavefront-sdk-go/application"
	"github.com/wavefronthq/wavefront-sdk-go/senders"
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

	token := os.Getenv("WF_TOKEN")
	if token == "" {
		fmt.Println("Environment variable WF_TOKEN must be set")
		os.Exit(1)
	}

	cfg := &senders.DirectConfiguration{
		Server: "https://surf.wavefront.com",
		Token:  token,
	}

	sender, err := senders.NewDirectSender(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	reporter := reporting.NewReporter(
		sender,
		application.New("owl", "electricity"),
		reporting.Source("owl.internal.glvr.io"),
		reporting.Prefix("owl.monitor"),
		reporting.LogErrors(true),
		reporting.RuntimeMetric(true),
	)

	tags := map[string]string{
		"type": "electricity",
	}

	mReadings := metrics.NewCounter()
	err = reporter.RegisterMetric("readings", mReadings, tags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mPower := metrics.NewGaugeFloat64()
	err = reporter.RegisterMetric("power", mPower, tags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mBat := metrics.NewGaugeFloat64()
	err = reporter.RegisterMetric("bat", mBat, tags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mRSSI := metrics.NewGaugeFloat64()
	err = reporter.RegisterMetric("rssi", mRSSI, tags)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mLQI := metrics.NewGaugeFloat64()
	err = reporter.RegisterMetric("lqi", mLQI, tags)
	if err != nil {
		fmt.Println(err)
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
			continue
		}

		mReadings.Inc(1)
		mPower.Update(elec.Chan[0].Power)
		mBat.Update(elec.Battery)
		mRSSI.Update(elec.RSSI)
		mLQI.Update(elec.LQI)
	}
}
