package ecfgtool

import (
	"fmt"
	"go/types"
	"reflect"
	"strings"

	"github.com/omcrgnt/ecfg/internal/usage"
	"github.com/omcrgnt/ecfg/pkg/walk"
	optionsv1 "github.com/omcrgnt/proto/gen/go/options/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type usageInput struct {
	IsProto      bool
	ReflectType  reflect.Type
	ProtoMsgType reflect.Type
	TypesType    types.Type
}

func resolveUsage(f walk.FieldDesc, isProto bool) (string, error) {
	in := usageInput{IsProto: isProto}
	if f.TypesType != nil {
		in.TypesType = f.TypesType
	}
	if f.ReflectType != nil {
		in.ReflectType = f.ReflectType
		if isProto {
			t := f.ReflectType
			for t.Kind() == reflect.Ptr {
				t = t.Elem()
			}
			in.ProtoMsgType = t
		}
	}
	return resolveUsageInput(in)
}

func resolveUsageInput(in usageInput) (string, error) {
	if in.IsProto {
		if in.ProtoMsgType != nil {
			return protoUsage(in.ProtoMsgType)
		}
		if named := typesNamed(in.TypesType); named != nil {
			return protoUsageFromNamed(named)
		}
		return "", ErrMissingUsage
	}
	if in.ReflectType != nil {
		return goUsage(in.ReflectType)
	}
	if named, ok := in.TypesType.(*types.Named); ok {
		return goUsageFromNamed(named)
	}
	return "", ErrMissingUsage
}

func protoUsageFromNamed(named *types.Named) (string, error) {
	if named.Obj() == nil || named.Obj().Pkg() == nil {
		return "", ErrMissingUsage
	}
	mt, err := protoMessageFromGoNamed(named)
	if err != nil {
		return "", ErrMissingUsage
	}
	return protoUsageFromMessageType(mt)
}

func protoMessageFromGoNamed(named *types.Named) (protoreflect.MessageType, error) {
	if full, ok := protoFullNameFromGoNamed(named); ok {
		if mt, err := protoregistry.GlobalTypes.FindMessageByName(full); err == nil {
			return mt, nil
		}
	}
	full := protoreflect.FullName(named.Obj().Pkg().Path() + "." + named.Obj().Name())
	return protoregistry.GlobalTypes.FindMessageByName(full)
}

// protoFullNameFromGoNamed maps org generated Go import paths to protobuf message names.
// e.g. github.com/omcrgnt/proto/gen/go/common/v1 + Label → common.v1.Label
func protoFullNameFromGoNamed(named *types.Named) (protoreflect.FullName, bool) {
	const marker = "/gen/go/"
	path := named.Obj().Pkg().Path()
	i := strings.Index(path, marker)
	if i < 0 {
		return "", false
	}
	suffix := strings.ReplaceAll(path[i+len(marker):], "/", ".")
	return protoreflect.FullName(suffix + "." + named.Obj().Name()), true
}

func protoUsageFromMessageType(mt protoreflect.MessageType) (string, error) {
	fields := mt.Descriptor().Fields()
	if fields.Len() != 1 || fields.Get(0).Name() != "value" {
		return "", ErrInvalidProtoWrapper
	}
	opts := fields.Get(0).Options()
	if !proto.HasExtension(opts, optionsv1.E_Usage) {
		return "", ErrMissingUsage
	}
	ext := proto.GetExtension(opts, optionsv1.E_Usage)
	text, ok := ext.(string)
	if !ok || strings.TrimSpace(text) == "" {
		return "", ErrEmptyUsage
	}
	return text, nil
}

func goUsageFromNamed(named *types.Named) (string, error) {
	if named.Obj() == nil || named.Obj().Pkg() == nil {
		return "", ErrMissingUsage
	}
	text, err := usage.GoUsageFromAST(named.Obj().Pkg().Path(), named.Obj().Name())
	if err != nil {
		return "", ErrMissingUsage
	}
	if strings.TrimSpace(text) == "" {
		return "", ErrEmptyUsage
	}
	return text, nil
}

func goUsage(typ reflect.Type) (string, error) {
	if typ == nil {
		return "", ErrMissingUsage
	}
	var u Usage
	if typ.Implements(reflect.TypeOf((*Usage)(nil)).Elem()) {
		v := reflect.New(typ).Elem()
		u = v.Interface().(Usage)
	} else if reflect.PointerTo(typ).Implements(reflect.TypeOf((*Usage)(nil)).Elem()) {
		u = reflect.New(typ).Interface().(Usage)
	} else {
		return "", ErrIncompatibleLeaf
	}
	text := strings.TrimSpace(u.Usage())
	if text == "" {
		return "", ErrEmptyUsage
	}
	return text, nil
}

func protoUsage(typ reflect.Type) (string, error) {
	if typ == nil {
		return "", ErrMissingUsage
	}
	for typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	msg, ok := reflect.New(typ).Interface().(proto.Message)
	if !ok {
		return "", ErrIncompatibleLeaf
	}
	desc := msg.ProtoReflect().Descriptor()
	fields := desc.Fields()
	if fields.Len() != 1 {
		return "", fmt.Errorf("%w: proto message field count %d", ErrInvalidProtoWrapper, fields.Len())
	}
	field := fields.Get(0)
	if field.Name() != "value" {
		return "", ErrInvalidProtoWrapper
	}
	opts := field.Options()
	if !proto.HasExtension(opts, optionsv1.E_Usage) {
		return "", ErrMissingUsage
	}
	ext := proto.GetExtension(opts, optionsv1.E_Usage)
	text, ok := ext.(string)
	if !ok || strings.TrimSpace(text) == "" {
		return "", ErrEmptyUsage
	}
	return text, nil
}
