package testdata

// BadRootInt has a non-struct root field.
type BadRootInt struct {
	Port int `ecfg:"PORT"`
}

// BadNoTag is missing ecfg on a root field.
type BadNoTag struct {
	Server ServerBlock
}

// BadNested has a struct block nested at depth 1.
type BadNested struct {
	Block BlockNested `ecfg:"BLOCK"`
}

// BlockNested contains another struct (not proto).
type BlockNested struct {
	Inner InnerNested `ecfg:"INNER"`
}

// InnerNested is a nested block.
type InnerNested struct {
	X int
}

// BadSlice uses an unsupported slice container.
type BadSlice struct {
	Tags []string `ecfg:"TAGS"`
}

// BadMap uses an unsupported map container.
type BadMap struct {
	Labels map[string]string `ecfg:"LABELS"`
}

// BadDup maps two leaves to the same env key segment.
type BadDup struct {
	Block BlockDup `ecfg:"BLOCK"`
}

// BlockDup has duplicate ecfg segments.
type BlockDup struct {
	A Label `ecfg:"LABEL"`
	B Label `ecfg:"LABEL"`
}

// BadBare uses a bare string leaf without Usage/Validator.
type BadBare struct {
	Block BlockBare `ecfg:"BLOCK"`
}

// BlockBare holds a bare string field.
type BlockBare struct {
	Name string `ecfg:"NAME"`
}

// NoUsageCfg has a bare int leaf (no Usage/Validator).
type NoUsageCfg struct {
	Block BlockNoUsage `ecfg:"BLOCK"`
}

// BlockNoUsage holds bare int.
type BlockNoUsage struct {
	Count int `ecfg:"COUNT"`
}
