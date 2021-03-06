package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/Bo0mer/yamt/internal/flagvar"
	"github.com/Bo0mer/yamt/iostat"
	"github.com/Bo0mer/yamt/metric"
	"github.com/Bo0mer/yamt/metric/riemann"
	"github.com/Bo0mer/yamt/netstat"
)

var (
	host       string
	port       int
	eventHost  string
	interval   int
	tags       flagvar.Array
	attributes flagvar.Map

	net       bool
	ignoreIfs string

	disk          bool
	ignoreDevices string
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
	flag.Var(&tags, "t", "Tag to add to events (shorthand)")
	flag.Var(&tags, "tag", "Tag to add to events")
	flag.Var(&attributes, "a", "Attribute to add to the events (shorthand)")
	flag.Var(&attributes, "attribute", "Attribute to add to the events")

	flag.BoolVar(&net, "net", false, "Report network interface metrics")
	flag.StringVar(&ignoreIfs, "g", "lo", "Interfaces to ignore (shorthand)")
	flag.StringVar(&ignoreIfs, "ignore-interfaces", "lo", "Interfaces to ignore")

	flag.BoolVar(&disk, "disk", false, "Report disk metrics")
	flag.StringVar(&ignoreDevices, "d", "ram|loop", "Devices to exclude")
	flag.StringVar(&ignoreDevices, "ignore-devices", "ram|loop", "Devices to exclude")
}

func main() {
	flag.Parse()

	collectors := make([]metric.Collector, 0)
	if net {
		except, err := regexp.Compile(ignoreIfs)
		if err != nil {
			log.Fatalf("yamt: invalid network interface regexp: %v\n", err)
		}
		netCollector, err := netstat.NewIfStatCollector(netstat.DefaultIfStatReader, except)
		if err != nil {
			log.Fatalf("yamt: error creating interface stats collector: %v\n", err)
		}
		collectors = append(collectors, netCollector)
		log.Printf("yamt: attached network interface stats collector")
	}

	if disk {
		except, err := regexp.Compile(ignoreDevices)
		if err != nil {
			log.Fatalf("yamt: invalid io device regexp: %v\n", err)
		}
		ioCollector, err := iostat.NewDeviceStatCollector(iostat.DefaultDevStatReader, except)
		if err != nil {
			log.Fatalf("yamt: error creating io stats collector: %v\n", err)
		}
		collectors = append(collectors, ioCollector)
		log.Printf("yamt: attached io device stats collector")
	}

	log.Printf("yamt: sticking tags to events: %v\n", tags)
	log.Printf("yamt: sticking attributes to events: %v\n", attributes)
	emitter := riemann.NewEmitter(fmt.Sprintf("%s:%d", host, port),
		riemann.Host(eventHost),
		riemann.Tags(tags),
		riemann.Attributes(attributes))

	d := time.Duration(interval) * time.Second
	reporter := metric.NewReporter(emitter, collectors, metric.Interval(d))
	reporter.Start()
	defer reporter.Close()

	log.Printf("yamt: started emitting metrics\n")
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	sig := <-c
	fmt.Printf("yamt: exiting due to %s\n", sig)
}
