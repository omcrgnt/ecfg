package ecfg

import (
	"flag"
	"net/http"
	"os"
	"reflect"

	json "github.com/goccy/go-json"
)

var parsed []byte

func Parse(t any, options ...option) error {
	return parseWithFlagSet(flag.CommandLine, os.Args[1:], t, options...)
}

func parseWithFlagSet(flagSet *flag.FlagSet, argumentList []string, t any, options ...option) error {
	if flagSet.Parsed() {
		return errWrap(ErrAlreadyParsed)
	}

	if err := checkInput(t); err != nil {
		return errWrap(err)
	}

	option := newOption(options...)

	if err := parseToStruct(t, flagSet, option, ""); err != nil {
		return errWrap(err)
	}

	if err := flagSet.Parse(argumentList); err != nil {
		return errWrap(err)
	}

	parsed, _ = json.Marshal(t)
	return nil
}

func Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(parsed)
	}
}

func checkInput(t any) error {
	if reflect.ValueOf(t).Kind() != reflect.Pointer || reflect.ValueOf(t).Elem().Kind() != reflect.Struct {
		return ErrInvalidInput
	}
	return nil
}
