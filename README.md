# go-owl
A local listener for the Owl Intuition energy monitor

## Running the listener

If you have `make` installed, you can build and run the listener in a Docker container with the following commands.

```bash
make build
make run
```

## Running Prometheus

```bash
docker run -p 9090:9090 prom/prometheus
```