package testdata

// PtrAppConfig uses a nil pointer root block (initialized at traverse).
type PtrAppConfig struct {
	Server *ServerBlock `ecfg:"SERVER"`
}
