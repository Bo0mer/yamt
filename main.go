package main

import (
	"flag"
	"fmt"
	"regexp"
	"time"

	"github.com/bo0mer/yamt/metric/riemann"
	"github.com/bo0mer/yamt/netstat"
)

var (
	host      string
	port      int
	eventHost string
)

func init() {
	flag.StringVar(&host, "h", "localhost", "Riemann host (shorthand)")
	flag.StringVar(&host, "host", "localhost", "Riemann host")
	flag.IntVar(&port, "p", 5555, "Riemann port (shorthand)")
	flag.IntVar(&port, "port", 5555, "Riemann port")
	flag.StringVar(&eventHost, "e", "", "Event hostname (shorthand)")
	flag.StringVar(&eventHost, "event-host", "", "Event hostname")
}

func main() {
	flag.Parse()

	emitter := riemann.NewEmitter(fmt.Sprintf("%s:%d", host, port),
		riemann.Host(eventHost))

	re := regexp.MustCompile("lo")
	r := netstat.NewReporter(emitter,
		netstat.Interval(time.Second),
		netstat.Except(re))
	r.Start()
	time.Sleep(time.Second * 30)
	r.Close()
}
