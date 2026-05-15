package teststruct

type Inner struct {
	Key string `ecfg:"KEY"`
}

type Nested struct {
	Items []Inner `ecfg:"ITEMS"`
}
