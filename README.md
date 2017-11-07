# go-owl
A local listener for the Owl Intuition energy monitor

[![Go Report Card](https://goreportcard.com/badge/github.com/billglover/go-owl)](https://goreportcard.com/report/github.com/billglover/go-owl) [![travis-ci](https://travis-ci.org/billglover/go-owl.svg?branch=master)](https://travis-ci.org/billglover/go-owl.svg?branch=master)

Package owl reads a slice of bytes as broadcast by the Owl Intuition electricity monitor and decodes them into an ElecReading containing three channels of Power and Energy measurements. It also reports battery level, signal strength and timestamp. Although the Owl Intuition broadcasts weather readings these are ignored as they are of limited use. Errors are returned if the byte slice is not decoded successfully.

Further information on the Owl Intuition multicast and UDP messages formats can be found on the [OWL Intuition support pages](https://theowl.zendesk.com/hc/en-gb/articles/201284603-Multicast-UDP-API-Information).

Full documentation: [godoc.org/github.com/billglover/go-owl](https://godoc.org/github.com/billglover/go-owl)

## Examples

 * [x] *basic* – logs electricity readings to the console
 * [x] *multicast* – similar to `basic` but listens to the multicast address
 * [ ] *prometheus* – exposes electricty readings as metrics for Prometheus

To run the basic example:

```bash
examples/multicast ▸ go build
examples/multicast ▸ ./multicast

2017-11-07 21:50:44 +0000 GMT : electricity reading : power=434.00w
2017-11-07 21:50:56 +0000 GMT : electricity reading : power=418.00w
2017-11-07 21:51:32 +0000 GMT : electricity reading : power=402.00w
```

## Benchmarks

```plain
goos: darwin
goarch: amd64
pkg: github.com/billglover/go-owl
BenchmarkRead-8   	   50000	     24739 ns/op
PASS
ok  	github.com/billglover/go-owl	1.496s
```

**Note:** remember, benchmarks allow for relative comparisons as actual performance is system dependence.

## Contributing

Contributions are welcome. Please leave a comment on an issue if you are going to work on it.
