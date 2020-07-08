// Package owl reads a slice of bytes as broadcast by the Owl Intuition electricity
// monitor and decodes them into an ElecReading containing three channels of Power and
// Energy measurements. It also reports battery level, signal strength and timestamp.
// Although the Owl Intuition broadcasts weather readings these are ignored as they
// are of limited use. Errors are returned if the byte slice is not decoded successfully.
//
// Further information on the Owl Intuition multicast and UDP messages formats can be
// found on the OWL Intuition support pages.
//
// https://theowl.zendesk.com/hc/en-gb/articles/201284603-Multicast-UDP-API-Information
package owl

import (
	"encoding/xml"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// MulticastAddress is the default address for the Owl Intuition.
	MulticastAddress string = "224.192.32.19:22600"
)

var (
	// ErrWeatherPacket indicates we have received a weather packet.
	ErrWeatherPacket = errors.New("weather packets are not decoded")

	// ErrInvalidPacket indicates we were unable to decode a valid packet.
	ErrInvalidPacket = errors.New("unable to decode packet")
)

// ElecReading represents a single electricity reading from the Owl Intuition.
type ElecReading struct {
	ID        string
	Timestamp time.Time
	RSSI      float64
	LQI       float64
	Battery   float64
	Chan      [3]ElecChan
}

// ElecChan represents a single channel electricity reading containing both
// instantaneous Power and Energy used throughout the day.
type ElecChan struct {
	Energy      float64
	EnergyUnits string
	Power       float64
	PowerUnits  string
}

// Read takes a byte slice and returns an ElecReading containing three channels
// of data. It returns an empty reading and an error if decoding the byte slice was
// unsuccessful.
func Read(b []byte) (ElecReading, error) {
	elec := ElecReading{}

	p := packet{}
	err := xml.Unmarshal(b, &p)
	if err != nil {
		return elec, ErrInvalidPacket
	}

	switch p.XMLName.Local {
	case "weather":
		return elec, ErrWeatherPacket
	case "electricity":
		elec, err := parseElectric(p, elec)
		return elec, err
	default:
		return elec, ErrInvalidPacket
	}
}

// packet represents a single data packet from the Owl Intuition.
type packet struct {
	XMLName xml.Name `xml:""`
	ID      string   `xml:"id,attr"`
	elecPacket
}

// elecPacket represents a single packet of electricity data.
type elecPacket struct {
	Time     int64     `xml:"timestamp"`
	Signal   signal    `xml:"signal"`
	Battery  battery   `xml:"battery"`
	Channels []channel `xml:"chan"`
}

// signal represents the signal strength at the the Owl Intuition.
// receiver
type signal struct {
	RSSI float64 `xml:"rssi,attr"`
	LQI  float64 `xml:"lqi,attr"`
}

// battery represents the battery level in the Owl Intuition.
// transmitter
type battery struct {
	Level string `xml:"level,attr"`
}

// channel represents a single channel electricity reading.
type channel struct {
	ID     int     `xml:"id,attr"`
	Power  reading `xml:"curr"`
	Energy reading `xml:"day"`
}

// reading represents a single value read from the Owl Intuition.
type reading struct {
	Units string  `xml:"units,attr"`
	Value float64 `xml:",chardata"`
}

// parseElectric populates an ElecReading struct with data from a packet.
func parseElectric(p packet, elec ElecReading) (ElecReading, error) {
	elec.ID = p.ID
	elec.Timestamp = time.Unix(p.Time, 0)
	batStr := strings.Replace(p.Battery.Level, "%", "", -1)
	bat, err := strconv.ParseFloat(batStr, 64)
	if err != nil {
		return elec, fmt.Errorf("unexpected value for battery level: got %s, want <float>%%", p.Battery.Level)
	}
	elec.Battery = bat

	elec.RSSI = p.Signal.RSSI
	elec.LQI = p.Signal.LQI

	if len(p.Channels) != 3 {
		return elec, fmt.Errorf("expected 3 channels, received %d", len(p.Channels))
	}

	elec.Chan[0] = ElecChan{
		Energy:      p.Channels[0].Energy.Value,
		EnergyUnits: p.Channels[0].Energy.Units,
		Power:       p.Channels[0].Power.Value,
		PowerUnits:  p.Channels[0].Power.Units,
	}
	elec.Chan[1] = ElecChan{
		Energy:      p.Channels[1].Energy.Value,
		EnergyUnits: p.Channels[1].Energy.Units,
		Power:       p.Channels[1].Power.Value,
		PowerUnits:  p.Channels[1].Power.Units,
	}
	elec.Chan[2] = ElecChan{
		Energy:      p.Channels[2].Energy.Value,
		EnergyUnits: p.Channels[2].Energy.Units,
		Power:       p.Channels[2].Power.Value,
		PowerUnits:  p.Channels[2].Power.Units,
	}

	return elec, nil
}
