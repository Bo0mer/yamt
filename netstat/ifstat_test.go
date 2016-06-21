package netstat_test

import (
	"testing"

	"github.com/bo0mer/yamt/netstat"
)

func TestIfStatReader(t *testing.T) {
	r := netstat.NewIfStatReader("testdata/procNetDev")
	stats, err := r.ReadStats()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want := []netstat.IfStat{
		netstat.IfStat{
			Name:         "eth0",
			RxBytes:      15017954683,
			RxPackets:    11623018,
			RxErrs:       1,
			RxDrop:       1,
			RxCompressed: 91283,
			TxBytes:      14743413932,
			TxPackets:    23122406,
			TxErrs:       288,
			TxDrop:       289,
		},
		netstat.IfStat{
			Name:         "lo",
			RxBytes:      334946,
			RxPackets:    1394,
			RxCompressed: 1394,
			TxBytes:      334946,
			TxPackets:    1394,
			TxCompressed: 1394,
		},
	}

	if len(want) != len(stats) {
		t.Fatalf("expected %v\n\tgot: %v\n", want, stats)
	}
	for i, stat := range stats {
		if stat != want[i] {
			t.Fatalf("expected %v\n\tgot: %v\n", want, stats)
		}
	}
}
