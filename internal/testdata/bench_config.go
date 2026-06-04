package testdata

// BenchBlock5 is a depth-1 block with five leaf fields (used by all bench configs).
type BenchBlock5 struct {
	F1 Label `ecfg:"F1"`
	F2 Label `ecfg:"F2"`
	F3 Label `ecfg:"F3"`
	F4 Label `ecfg:"F4"`
	F5 Label `ecfg:"F5"`
}

// BenchConfig5 has five root blocks.
type BenchConfig5 struct {
	B01 BenchBlock5 `ecfg:"B01"`
	B02 BenchBlock5 `ecfg:"B02"`
	B03 BenchBlock5 `ecfg:"B03"`
	B04 BenchBlock5 `ecfg:"B04"`
	B05 BenchBlock5 `ecfg:"B05"`
}

// BenchConfig10 has ten root blocks.
type BenchConfig10 struct {
	B01 BenchBlock5 `ecfg:"B01"`
	B02 BenchBlock5 `ecfg:"B02"`
	B03 BenchBlock5 `ecfg:"B03"`
	B04 BenchBlock5 `ecfg:"B04"`
	B05 BenchBlock5 `ecfg:"B05"`
	B06 BenchBlock5 `ecfg:"B06"`
	B07 BenchBlock5 `ecfg:"B07"`
	B08 BenchBlock5 `ecfg:"B08"`
	B09 BenchBlock5 `ecfg:"B09"`
	B10 BenchBlock5 `ecfg:"B10"`
}

// BenchConfig15 has fifteen root blocks.
type BenchConfig15 struct {
	B01 BenchBlock5 `ecfg:"B01"`
	B02 BenchBlock5 `ecfg:"B02"`
	B03 BenchBlock5 `ecfg:"B03"`
	B04 BenchBlock5 `ecfg:"B04"`
	B05 BenchBlock5 `ecfg:"B05"`
	B06 BenchBlock5 `ecfg:"B06"`
	B07 BenchBlock5 `ecfg:"B07"`
	B08 BenchBlock5 `ecfg:"B08"`
	B09 BenchBlock5 `ecfg:"B09"`
	B10 BenchBlock5 `ecfg:"B10"`
	B11 BenchBlock5 `ecfg:"B11"`
	B12 BenchBlock5 `ecfg:"B12"`
	B13 BenchBlock5 `ecfg:"B13"`
	B14 BenchBlock5 `ecfg:"B14"`
	B15 BenchBlock5 `ecfg:"B15"`
}
