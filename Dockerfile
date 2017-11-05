FROM golang:1.9 as builder
WORKDIR /go/src/github.com/billglover/go-owl/ 
COPY main.go    .
RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o go-owl .

FROM scratch 
WORKDIR /root/
COPY --from=builder /go/src/github.com/billglover/go-owl/go-owl .
CMD ["./go-owl"]  