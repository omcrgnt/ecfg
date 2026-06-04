package testdata

import "google.golang.org/protobuf/types/known/wrapperspb"

// ProtoNoUsageCfg uses a proto wrapper without options.v1.usage on value.
type ProtoNoUsageCfg struct {
	Server ProtoNoUsageServer `ecfg:"SERVER"`
}

// ProtoNoUsageServer holds a standard wrapper proto (no usage extension).
type ProtoNoUsageServer struct {
	Port *wrapperspb.UInt32Value `ecfg:"PORT"`
}
