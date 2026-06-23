package testdata

import (
	"fmt"
	"time"

	"github.com/omcrgnt/builder"
)

// AppResourcesFixture mimics an AppResources root for template and gen tests.
type AppResourcesFixture struct {
	App         *WireApp     `ecfg:"APP"`
	ServiceItem *WireService `ecfg:"SERVICE_ITEM"`
	ServerHTTP  *WireServer  `ecfg:"SERVER_HTTP_ITEM"`
}

type WireApp struct{}

func (*WireApp) BuildConfig() (builder.Builder, error) {
	return appSpec{}, nil
}

type appSpec struct {
	ShutdownTimeout shutdownTimeout
}

func (appSpec) Build() (any, error) { return &WireApp{}, nil }

type shutdownTimeout time.Duration

func (shutdownTimeout) Usage() string { return "Grace period for shutdown" }

func (t shutdownTimeout) Validate() error {
	if t < 0 {
		return fmt.Errorf("shutdown timeout must be >= 0")
	}
	return nil
}

type WireService struct{}

func (*WireService) BuildConfig() (builder.Builder, error) {
	return serviceSpec{}, nil
}

type serviceSpec struct {
	MaxListLen maxListLen
}

func (serviceSpec) Build() (any, error) { return &WireService{}, nil }

type maxListLen int

func (maxListLen) Usage() string { return "Maximum items returned by List" }

func (l maxListLen) Validate() error {
	if l < 0 {
		return fmt.Errorf("max list len must be >= 0")
	}
	return nil
}

type WireServer struct{}

func (*WireServer) BuildConfig() (builder.Builder, error) {
	return serverSpec{}, nil
}

type serverSpec struct {
	Label Label
	Host  hostLeaf
	Port  Port
}

func (serverSpec) Build() (any, error) { return &WireServer{}, nil }

type hostLeaf string

func (hostLeaf) Usage() string { return "HTTP listen host" }

func (hostLeaf) Validate() error { return nil }
