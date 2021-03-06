package netstat

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/Bo0mer/yamt/internal"
)

//go:generate counterfeiter . InterfaceStatReader

// IfStat represents statistics about a network interface.
type IfStat struct {
	Name string

	RxBytes      uint64
	RxPackets    uint64
	RxErrs       uint64
	RxDrop       uint64
	RxFIFO       uint64
	RxFrame      uint64
	RxCompressed uint64
	RxMulticast  uint64

	TxBytes      uint64
	TxPackets    uint64
	TxErrs       uint64
	TxDrop       uint64
	TxFIFO       uint64
	TxColls      uint64
	TxCarrier    uint64
	TxCompressed uint64
}

// InterfaceStatReader should read statistics for all available network
// interfaces.
type InterfaceStatReader interface {
	ReadStats() ([]IfStat, error)
}

// IfStatReader reads statistics for network interfaces.
type IfStatReader struct {
	path string
}

// NewIfStatReader creates IfStatReader that reads from the specified path.
// It expects well defined format and may cause panics if it is not present.
func NewIfStatReader(path string) *IfStatReader {
	return &IfStatReader{
		path: path,
	}
}

// DefaultIfStatReader is the default implementation of InterfaceStatReader.
// It reads interface statistics from /proc/net/dev.
var DefaultIfStatReader InterfaceStatReader = NewIfStatReader("/proc/net/dev")

// ReadIfStats is shorthand for DefaultIfStatReader.ReadIfStats.
func ReadIfStats() ([]IfStat, error) {
	return DefaultIfStatReader.ReadStats()
}

func (r *IfStatReader) ReadStats() ([]IfStat, error) {
	data, err := ioutil.ReadFile(r.path)
	if err != nil {
		return nil, fmt.Errorf("readifstats: error reading from /proc/net/dev: %s", err)
	}
	return r.parseStats(data)
}

func (r *IfStatReader) parseStats(data []byte) ([]IfStat, error) {
	lines := strings.Split(string(data), "\n")
	lines = lines[2:] // reamove header
	stats := make([]IfStat, len(lines)-1)

	for i, line := range lines {
		if line == "" {
			break
		}
		stat, err := r.parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("readifstats: error parsing line %d: %s", i, err)
		}
		stats[i] = stat
	}
	return stats, nil
}

func (r *IfStatReader) parseLine(line string) (IfStat, error) {
	colon := strings.Index(line, ":")
	if colon <= 0 {
		return IfStat{}, fmt.Errorf("readifstats: unsupported format: %q", line)
	}

	stat := IfStat{}
	stat.Name = strings.Replace(line[:colon], " ", "", -1)

	fields := strings.Fields(line[colon+1:])
	p := &internal.ErrParser{}
	stat.RxBytes = p.ParseUint64(fields[0])
	stat.RxPackets = p.ParseUint64(fields[1])
	stat.RxErrs = p.ParseUint64(fields[2])
	stat.RxDrop = p.ParseUint64(fields[3])
	stat.RxFIFO = p.ParseUint64(fields[4])
	stat.RxFrame = p.ParseUint64(fields[5])
	stat.RxCompressed = p.ParseUint64(fields[6])
	stat.RxMulticast = p.ParseUint64(fields[7])

	stat.TxBytes = p.ParseUint64(fields[8])
	stat.TxPackets = p.ParseUint64(fields[9])
	stat.TxErrs = p.ParseUint64(fields[10])
	stat.TxDrop = p.ParseUint64(fields[11])
	stat.TxFIFO = p.ParseUint64(fields[12])
	stat.TxColls = p.ParseUint64(fields[13])
	stat.TxCarrier = p.ParseUint64(fields[14])
	stat.TxCompressed = p.ParseUint64(fields[15])

	if err := p.Err(); err != nil {
		return IfStat{}, fmt.Errorf("readifstats: error reading stats for %s: %s", stat.Name, err)
	}

	return stat, nil
}
