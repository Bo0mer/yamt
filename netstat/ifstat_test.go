package netstat

import "testing"

const procNetDev = `Inter-|   Receive                                                |  Transmit
 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
   eth0:15017954683 11623018    1    1    0     0          91283         0 14743413932 23122406    288    289    0     0       0          0
       lo:  334946    1394    0    0    0     0          1394         0   334946    1394    0    0    0     0       0          1394`

func TestParseStats(t *testing.T) {
	stats, err := parseStats([]byte(procNetDev))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	want := []IfStat{
		IfStat{
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
		IfStat{
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
