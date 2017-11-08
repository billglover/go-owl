package owl_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/billglover/go-owl"
)

var elec = []byte(`<electricity id='443719005443'>
	<timestamp>1509950911</timestamp>
	<signal rssi='-68' lqi='48'/>
	<battery level='100%'/>
	<chan id='0'>
		<curr units='w'>305.00</curr>
		<day units='wh'>1863.39</day>
	</chan>
	<chan id='1'>
		<curr units='w'>21.00</curr>
		<day units='wh'>3.01</day>
	</chan>
	<chan id='2'>
		<curr units='w'>270.26</curr>
		<day units='wh'>0.00</day>
	</chan>
</electricity>`)

var weather = []byte(`<weather id='443719005443' code='113'>
		<temperature>9.00</temperature>
		<text>Clear/Sunny</text>
	</weather>`)

var invalid = []byte(`<codequality id='443719005443' code='113'>
		<temperature>9.00</temperature>
		<text>Poor/Sunny</text>
	</codequality>`)

var invalidChan = []byte(`<electricity id='443719005443'>
	<timestamp>1509950911</timestamp>
	<signal rssi='-68' lqi='48'/>
	<battery level='80%'/>
	<chan id='0'>
		<curr units='w'>305.00</curr>
		<day units='wh'>1863.39</day>
	</chan>
</electricity>`)

var garbage = []byte(`asjfd中文可以吗😂`)

func TestReadElec(t *testing.T) {
	reading, err := owl.Read(elec)
	if err != nil {
		t.Fatalf("unexpected error when decoding data: %v", err)
	}
	if reading.ID != "443719005443" {
		t.Fatalf("unexpected identifier: got %s, want %s", reading.ID, "443719005443")
	}
	if reading.Timestamp != time.Unix(1509950911, 0) {
		t.Fatalf("unexpected timestamp: got %v, want %v", reading.Timestamp, time.Unix(1509950911, 0))
	}
	if reading.RSSI != -68.0 {
		t.Fatalf("unexpected RSSI: got %f, want %f", reading.RSSI, -68.0)
	}
	if reading.LQI != 48.0 {
		t.Fatalf("unexpected LQI: got %f, want %f", reading.LQI, 48.0)
	}
	if reading.Battery != 100.0 {
		t.Fatalf("unexpected battery level: got %f, want %f", reading.Battery, 100.0)
	}
}

func TestReadInvalidBattery(t *testing.T) {
	reading, err := owl.Read(invalidBat)
	if err == nil {
		t.Errorf("expected an error when decoding invalid battery data, got: %v", err)
	}
	if reading.Battery != 0 {
		t.Errorf("unexpected battery level: got %f, want %f", reading.Battery, 0.0)
	}
}

func TestReadInvalidChan(t *testing.T) {
	_, err := owl.Read(invalidChan)
	if err == nil {
		t.Errorf("expected an error when decoding invalid number of channels, got: %v", err)
	}
}

func TestReadChannels(t *testing.T) {
	reading, _ := owl.Read(elec)
	// test values in the first channel
	if reading.Chan[0].Power != 305.0 {
		t.Errorf("unexpected power value: got %f, want %f", reading.Chan[0].Power, 305.0)
	}
	if reading.Chan[0].PowerUnits != "w" {
		t.Errorf("unexpected power units: got %s, want %s", reading.Chan[0].PowerUnits, "w")
	}
	if reading.Chan[0].Energy != 1863.39 {
		t.Errorf("unexpected energy value: got %f, want %f", reading.Chan[0].Energy, 1863.39)
	}
	if reading.Chan[0].EnergyUnits != "wh" {
		t.Errorf("unexpected energy units: got %s, want %s", reading.Chan[0].EnergyUnits, "w")
	}

	// test values in the second channel
	if reading.Chan[1].Power != 21.0 {
		t.Errorf("unexpected power value: got %f, want %f", reading.Chan[1].Power, 21.0)
	}
	if reading.Chan[1].PowerUnits != "w" {
		t.Errorf("unexpected power units: got %s, want %s", reading.Chan[1].PowerUnits, "w")
	}
	if reading.Chan[1].Energy != 3.01 {
		t.Errorf("unexpected energy value: got %f, want %f", reading.Chan[1].Energy, 3.01)
	}
	if reading.Chan[1].EnergyUnits != "wh" {
		t.Errorf("unexpected energy units: got %s, want %s", reading.Chan[1].EnergyUnits, "w")
	}

	// test values in the third channel
	if reading.Chan[2].Power != 270.26 {
		t.Errorf("unexpected power value: got %f, want %f", reading.Chan[2].Power, 270.26)
	}
	if reading.Chan[2].PowerUnits != "w" {
		t.Errorf("unexpected power units: got %s, want %s", reading.Chan[2].PowerUnits, "w")
	}
	if reading.Chan[2].Energy != 0.0 {
		t.Errorf("unexpected energy value: got %f, want %f", reading.Chan[2].Energy, 0.0)
	}
	if reading.Chan[2].EnergyUnits != "wh" {
		t.Errorf("unexpected energy units: got %s, want %s", reading.Chan[2].EnergyUnits, "w")
	}
}

func TestReadWeather(t *testing.T) {
	_, err := owl.Read(weather)
	if err != owl.ErrWeatherPacket {
		t.Fatalf("unexpected an error when decoding a weather packet: got %v", err)
	}
}

func TestReadInvalid(t *testing.T) {
	_, err := owl.Read(invalid)
	if err != owl.ErrInvalidPacket {
		t.Fatalf("unexpected an error when decoding an invalid packet: got %v", err)
	}
}

func TestReadGarbage(t *testing.T) {
	_, err := owl.Read(garbage)
	if err != owl.ErrInvalidPacket {
		t.Fatalf("unexpected an error when decoding a garbage packet: got %v", err)
	}
}

func BenchmarkRead(b *testing.B) {
	for n := 0; n < b.N; n++ {
		owl.Read(elec)
	}
}

func ExampleRead() {
	var elec = []byte(`<electricity id='443719005443'>
		<timestamp>1509950911</timestamp>
		<signal rssi='-68' lqi='48'/>
		<battery level='100%'/>
		<chan id='0'>
			<curr units='w'>305.00</curr>
			<day units='wh'>1863.39</day>
		</chan>
		<chan id='1'/>
		<chan id='2'/>
	</electricity>`)

	r, err := owl.Read(elec)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%s power=%.2f energy=%.2f battery=%.2f\n", r.Timestamp.UTC(), r.Chan[0].Power, r.Chan[0].Energy, r.Battery)
	// Output: 2017-11-06 06:48:31 +0000 UTC power=305.00 energy=1863.39 battery=100.00
}
