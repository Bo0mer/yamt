# yamt
Yet another metric tool

## Installation
Just go get it.
```bash
go get github.com/bo0mer/yamt
```
## Usage
Example usage - send network and disk statistics to local riemann instance:
```
yamt -net -disk
```
You can configure the interval between different metric reports:
```
yamt -net -disk -i 20 # send report every 20 seconds
```

Following is a list of all supported command line arguments.
```
Usage of yamt:
  -d string
    	Devices to exclude (default "ram|loop")
  -disk
    	Report disk metrics
  -e string
    	Event hostname (shorthand)
  -event-host string
    	Event hostname
  -g string
    	Interfaces to ignore (shorthand) (default "lo")
  -h string
    	Riemann host (shorthand) (default "localhost")
  -host string
    	Riemann host (default "localhost")
  -i int
    	Seconds between updates (shorthand) (default 5)
  -ignore-devices string
    	Devices to exclude (default "ram|loop")
  -ignore-interfaces string
    	Interfaces to ignore (default "lo")
  -interval int
    	Seconds between updates (default 5)
  -net
    	Report network interface metrics
  -p int
    	Riemann port (shorthand) (default 5555)
  -port int
    	Riemann port (default 5555)
```

## Development

### Testing
yamt uses counterfeiter to create fakes. For more information see 
https://github.com/maxbrunsfeld/counterfeiter.

ATM no test framework is used.
```
go test $(go list ./... | grep -v vendor) # do not run tests for vendored deps
```

