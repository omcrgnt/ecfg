package ecfgtool

import (
	"errors"
	"strings"
	"testing"

	"buf.build/go/protovalidate"
	commonv1 "github.com/omcrgnt/proto/gen/go/common/v1"
)

func TestValidateProto_initFail(t *testing.T) {
	old := newProtovalidateFunc
	t.Cleanup(func() {
		newProtovalidateFunc = old
		resetValidatorForTest()
	})
	newProtovalidateFunc = func(...protovalidate.ValidatorOption) (protovalidate.Validator, error) {
		return nil, errors.New("protovalidate init failed")
	}
	resetValidatorForTest()

	err := validateProto(&commonv1.Port{Value: 1})
	if err == nil || !strings.Contains(err.Error(), "protovalidate") {
		t.Fatalf("got %v", err)
	}
}
