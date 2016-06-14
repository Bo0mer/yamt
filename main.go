package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/bo0mer/yamt/metric/riemann"
	"github.com/bo0mer/yamt/netstat"
)

var (
	host      string
	port      int
	eventHost string
	interval  int
	ignoreIfs string
)

func init() {
	flag.StringVar(&host, "h", "localhost", "Riemann host (shorthand)")
	flag.StringVar(&host, "host", "localhost", "Riemann host")
	flag.IntVar(&port, "p", 5555, "Riemann port (shorthand)")
	flag.IntVar(&port, "port", 5555, "Riemann port")
	flag.StringVar(&eventHost, "e", "", "Event hostname (shorthand)")
	flag.StringVar(&eventHost, "event-host", "", "Event hostname")
	flag.IntVar(&interval, "i", 5, "Seconds between updates (shorthand)")
	flag.IntVar(&interval, "interval", 5, "Seconds between updates")
	flag.StringVar(&ignoreIfs, "g", "lo", "Interfaces to ignore (shorthand)")
	flag.StringVar(&ignoreIfs, "ignore-interfaces", "lo", "Interfaces to ignore")
}

func main() {
	flag.Parse()

	emitter := riemann.NewEmitter(fmt.Sprintf("%s:%d", host, port),
		riemann.Host(eventHost))

	re, err := regexp.Compile(ignoreIfs)
	if err != nil {
		log.Fatalf("yamt: invalid regular expression: %s\n", err)
	}
	r := netstat.NewReporter(emitter,
		netstat.Interval(time.Second*time.Duration(interval)),
		netstat.Except(re))
	r.Start()
	select {}
	defer r.Close()
}
