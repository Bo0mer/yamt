package netstat

import "testing"

func TestParseStats(t *testing.T) {
	stats, err := readIfStats("testdata/procNetDev")
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
