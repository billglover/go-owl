# go-owl
A local listener for the Owl Intuition energy monitor

[![Go Report Card](https://goreportcard.com/badge/github.com/billglover/go-owl)](https://goreportcard.com/report/github.com/billglover/go-owl)

Package owl reads a slice of bytes as broadcast by the Owl Intuition electricity monitor and decodes them into an ElecReading containing three channels of Power and Energy measurements. It also reports battery level, signal strength and timestamp. Although the Owl Intuition broadcasts weather readings these are ignored as they are of limited use. Errors are returned if the byte slice is not decoded successfully.

Further information on the Owl Intuition multicast and UDP messages formats can be found on the [OWL Intuition support pages](https://theowl.zendesk.com/hc/en-gb/articles/201284603-Multicast-UDP-API-Information).

Full documentation: [godoc.org/github.com/billglover/go-owl](https://godoc.org/github.com/billglover/go-owl)

## Examples

 [ ] *basic* – logs electricity readings to the console
 [ ] *multicast* – similar to `basic` but listens to the multicast address
 [x] *prometheus* – exposes electricty readings as metrics for Prometheus

## Contributing

Contributions are welcome. Please leave a comment on an issue if you are going to work on it.
