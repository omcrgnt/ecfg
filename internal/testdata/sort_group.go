package testdata

// SortGroupCfg has two leaves in one root block for template sort tests.
type SortGroupCfg struct {
	Group TwoLeafBlock `ecfg:"GROUP"`
}

// TwoLeafBlock has two labeled leaves.
type TwoLeafBlock struct {
	Z Label `ecfg:"Z"`
	A Label `ecfg:"A"`
}
