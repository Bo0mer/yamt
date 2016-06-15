package iostat

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bo0mer/yamt/internal"
)

type DiskStat struct {
	// Major device number.
	Major int
	// Minor device number.
	Minor int
	// Device name.
	Name string

	// Total number of reads completed successfully.
	Reads uint64
	// Reads and writes which are adjacent to each other may be merged for
	// efficiency.  Thus two 4K reads may become one 8K read before it is
	// ultimately handed to the disk, and so it will be counted (and queued)
	// as only one I/O.  This field lets you know how often this was done.
	ReadsMerged uint64
	// Total number of sectors read successfully.
	ReadsSectors uint64
	// Total number of milliseconds spent by all reads.
	ReadsTimeMs uint64
	// Total number of writes completed successfully.
	Writes uint64
	// See ReadsMerged.
	WritesMerged uint64
	// Total number of sectors written successfully.
	WritesSectors uint64
	// Total number of milliseconds spent by all writes.
	WritesTimeMs uint64

	// Number of I/Os currently in progress.
	InFlight uint64
	// Number of milliseconds spent doing I/Os.
	IOTimeMs uint64
	// Weighted number of milliseconds spent doing I/Os.
	// This field is incremented at each I/O start, I/O completion, I/O
	// merge, or read of these stats by the number of I/Os in progress
	// (field 9) times the number of milliseconds spent doing I/O since the
	// last update of this field.  This can provide an easy measure of both
	// I/O completion time and the backlog that may be accumulating.
	WeightedIOTimeMS uint64
}

// ReadDiskStats reads statistics for all available disks.
// It does so by reading from /proc/diskstats.
func ReadDiskStats() ([]DiskStat, error) {
	return readDiskStats("/proc/diskstats")
}

func readDiskStats(path string) ([]DiskStat, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("readdiskstats: error reading from /proc/diskstats: %v", err)
	}
	return parseStats(data)
	return nil, nil
}

func parseStats(data []byte) ([]DiskStat, error) {
	lines := strings.Split(string(data), "\n")
	stats := make([]DiskStat, len(lines)-1)

	for i, line := range lines {
		if line == "" {
			break
		}
		stat, err := parseLine(line)
		if err != nil {
			return nil, fmt.Errorf("readidiskstats: error parsing line %d: %v", i, err)
		}
		stats[i] = stat
	}
	return stats, nil
}

func parseLine(line string) (DiskStat, error) {
	fields := strings.Fields(line)
	p := &internal.ErrParser{}

	stat := DiskStat{}
	stat.Name = fields[2]

	stat.Major = p.ParseInt(fields[0])
	stat.Minor = p.ParseInt(fields[1])
	// fields[2] is disk name, see above
	stat.Reads = p.ParseUint64(fields[3])
	stat.ReadsMerged = p.ParseUint64(fields[4])
	stat.ReadsSectors = p.ParseUint64(fields[5])
	stat.ReadsTimeMs = p.ParseUint64(fields[6])
	stat.Writes = p.ParseUint64(fields[7])
	stat.WritesMerged = p.ParseUint64(fields[8])
	stat.WritesSectors = p.ParseUint64(fields[9])
	stat.WritesTimeMs = p.ParseUint64(fields[10])
	stat.InFlight = p.ParseUint64(fields[11])
	stat.IOTimeMs = p.ParseUint64(fields[12])
	stat.WeightedIOTimeMS = p.ParseUint64(fields[13])

	if err := p.Err(); err != nil {
		return DiskStat{}, fmt.Errorf("readidiskstats: error reading stats for %s: %v", stat.Name, err)
	}

	return stat, nil
}
