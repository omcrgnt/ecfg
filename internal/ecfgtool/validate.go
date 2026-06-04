package ecfgtool

import (
	"fmt"
	"sync"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/proto"
)

var (
	validatorOnce      sync.Once
	validator          protovalidate.Validator
	validatorErr       error
	newProtovalidateFunc = defaultProtovalidate
)

func defaultProtovalidate(opts ...protovalidate.ValidatorOption) (protovalidate.Validator, error) {
	return protovalidate.New(opts...)
}

func getValidator() (protovalidate.Validator, error) {
	validatorOnce.Do(func() {
		validator, validatorErr = newProtovalidateFunc()
	})
	return validator, validatorErr
}

func resetValidatorForTest() {
	validatorOnce = sync.Once{}
	validator = nil
	validatorErr = nil
}

func validateProto(msg proto.Message) error {
	v, err := getValidator()
	if err != nil {
		return fmt.Errorf("ecfg: protovalidate: %w", err)
	}
	if err := v.Validate(msg); err != nil {
		return err
	}
	return nil
}

func validateLeaf(v any, isProto bool, usageText string, fieldPath, envKey string) error {
	if isProto {
		msg, ok := v.(proto.Message)
		if !ok {
			return fmt.Errorf("%w: %T", ErrIncompatibleLeaf, v)
		}
		if err := validateProto(msg); err != nil {
			return wrapValidateErr(fieldPath, envKey, err, usageText)
		}
		return nil
	}
	valid, ok := v.(Validator)
	if !ok {
		return ErrIncompatibleLeaf
	}
	if err := valid.Validate(); err != nil {
		return wrapValidateErr(fieldPath, envKey, err, usageText)
	}
	return nil
}

func wrapValidateErr(fieldPath, envKey string, err error, usage string) error {
	return fmt.Errorf("ecfg: %s (%s): validate: %w; usage: %s", fieldPath, envKey, err, usage)
}
