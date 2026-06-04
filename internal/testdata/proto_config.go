package testdata

import commonv1 "github.com/omcrgnt/proto/gen/go/common/v1"

// ProtoCfg uses a proto wrapper leaf.
type ProtoCfg struct {
	Server ProtoServer `ecfg:"SERVER"`
}

// ProtoServer holds a Port proto message.
type ProtoServer struct {
	Port *commonv1.Port `ecfg:"PORT"`
}

// BadProtoWrapper uses an invalid non-proto struct as leaf shape.
type BadProtoWrapper struct {
	Server BadProtoServer `ecfg:"SERVER"`
}

// BadProtoServer holds invalid wrapper.
type BadProtoServer struct {
	Wrap BadWrap `ecfg:"WRAP"`
}

// BadWrap is not a valid proto wrapper.
type BadWrap struct {
	Value int
	Extra int
}
