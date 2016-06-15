package iostat

import "testing"

func TestReadDiskgot(t *testing.T) {
	got, err := readDiskStats("testdata/procDiskstats")
	if err != nil {
		t.Errorf("unexpected error: %v\n", err)
	}

	want := []DiskStat{
		DiskStat{
			Name:             "sda",
			Major:            8,
			Minor:            0,
			Reads:            70705,
			ReadsMerged:      103,
			ReadsSectors:     2596826,
			ReadsTimeMs:      45308,
			Writes:           37009,
			WritesMerged:     64245,
			WritesSectors:    4307616,
			WritesTimeMs:     95692,
			InFlight:         0,
			IOTimeMs:         24440,
			WeightedIOTimeMS: 140248,
		},
		DiskStat{
			Name:             "sda1",
			Major:            8,
			Minor:            1,
			Reads:            70525,
			ReadsMerged:      103,
			ReadsSectors:     2595386,
			ReadsTimeMs:      45288,
			Writes:           36167,
			WritesMerged:     64245,
			WritesSectors:    4307608,
			WritesTimeMs:     95648,
			InFlight:         0,
			IOTimeMs:         24380,
			WeightedIOTimeMS: 140188,
		},
	}

	if len(want) != len(got) {
		t.Errorf("want %v\n\tgot %v\n", want, got)
	}
	for i, stat := range got {
		if stat != want[i] {
			t.Errorf("want %v\n\tgot %v\n", want, got)
		}
	}
}
