package ecfg

import (
	"reflect"
	"regexp"
	"strings"
)

func checkSourceType(t any) error {
	if reflect.ValueOf(t).Kind() != reflect.Pointer || reflect.ValueOf(t).Elem().Kind() != reflect.Struct {
		return ErrInvalidInput
	}
	return nil
}

func checkStructFieldType(t any) error {
	if reflect.ValueOf(t).Elem().Kind() == reflect.Pointer {
		return ErrInvalidInput
	}
	return nil
}

var (
	// Регулярка для поиска границ слов в CamelCase (LoggerLevel -> Logger_Level)
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

// ToEnvName принимает части пути (напр. "Logger", "Level", "Value")
// и превращает в "LOGGER_LEVEL"
func ToEnvName(parts ...string) string {
	var result []string

	for _, part := range parts {
		// 1. Игнорируем техническое поле "Value", которое часто используется в твоих Proto
		if strings.ToLower(part) == "value" {
			continue
		}

		// 2. Превращаем CamelCase в snake_case (LoggerLevel -> logger_level)
		snake := matchFirstCap.ReplaceAllString(part, "${1}_${2}")
		snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")

		result = append(result, strings.ToUpper(snake))
	}

	// 3. Склеиваем всё через подчеркивание
	return strings.Join(result, "_")
}
