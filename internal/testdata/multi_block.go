package testdata

// MultiBlock has two root blocks for template grouping tests.
type MultiBlock struct {
	Server ServerBlock `ecfg:"SERVER"`
	Worker WorkerBlock `ecfg:"WORKER"`
}

// WorkerBlock is a second root block.
type WorkerBlock struct {
	Label Label `ecfg:"LABEL"`
}
