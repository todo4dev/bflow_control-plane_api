package env

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	envMutex sync.RWMutex
	envStore = map[string]string{}
)

func loadEnvFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		key, value, ok := parseEnvLine(line)
		if !ok {
			continue
		}

		expandedValue := expandValue(value)

		envMutex.Lock()
		envStore[key] = expandedValue
		envMutex.Unlock()

		_ = os.Setenv(key, expandedValue)
	}

	return scanner.Err()
}

func parseEnvLine(rawLine string) (string, string, bool) {
	line := strings.TrimSpace(rawLine)
	if line == "" {
		return "", "", false
	}

	if strings.HasPrefix(line, "#") ||
		strings.HasPrefix(line, "//") {
		return "", "", false
	}

	if strings.HasPrefix(line, "export ") {
		line = strings.TrimSpace(line[len("export "):])
	}

	equalIndex := strings.Index(line, "=")
	if equalIndex < 0 {
		return "", "", false
	}

	keyPart := strings.TrimSpace(line[:equalIndex])
	if keyPart == "" {
		return "", "", false
	}

	valuePart := strings.TrimSpace(line[equalIndex+1:])
	valuePart = stripInlineComments(valuePart)
	valuePart = strings.TrimSpace(valuePart)

	if len(valuePart) >= 2 {
		firstChar := valuePart[0]
		lastChar := valuePart[len(valuePart)-1]
		if (firstChar == '"' && lastChar == '"') || (firstChar == '\'' && lastChar == '\'') {
			valuePart = valuePart[1 : len(valuePart)-1]
		}
	}

	if keyPart == "" {
		return "", "", false
	}

	return keyPart, valuePart, true
}

func stripInlineComments(valueWithComments string) string {
	inQuotes := false
	var quoteChar rune

	for index, char := range valueWithComments {
		if char == '"' || char == '\'' {
			if !inQuotes {
				inQuotes = true
				quoteChar = char
			} else if quoteChar == char {
				inQuotes = false
			}
		}

		if inQuotes {
			continue
		}

		if char == '#' {
			return strings.TrimSpace(valueWithComments[:index])
		}
		if char == '/' && index+1 < len(valueWithComments) && valueWithComments[index+1] == '/' {
			return strings.TrimSpace(valueWithComments[:index])
		}
	}

	return strings.TrimSpace(valueWithComments)
}

func expandValue(raw string) string {
	mappingFunc := func(key string) string {
		envMutex.RLock()
		value, ok := envStore[key]
		envMutex.RUnlock()

		if ok {
			return value
		}

		return os.Getenv(key)
	}

	return os.Expand(raw, mappingFunc)
}

func lookupEnv(key string) (string, bool) {
	envMutex.RLock()
	value, ok := envStore[key]
	envMutex.RUnlock()

	if ok {
		return value, true
	}

	return os.LookupEnv(key)
}

func convertStringToType[T any](raw string) (T, error) {
	var zero T

	targetType := reflect.TypeOf(zero)
	if targetType == nil {
		var anyValue any = raw
		return anyValue.(T), nil
	}

	if targetType == reflect.TypeFor[json.RawMessage]() {
		bytes := []byte(raw)
		if !json.Valid(bytes) {
			return zero, fmt.Errorf("env: invalid JSON for json.RawMessage: %q", raw)
		}
		var anyValue any = json.RawMessage(bytes)
		return anyValue.(T), nil
	}

	if targetType == reflect.TypeFor[time.Time]() {
		parsedTime, err := parseTime(raw)
		if err != nil {
			return zero, err
		}
		var anyValue any = parsedTime
		return anyValue.(T), nil
	}

	if targetType == reflect.TypeFor[time.Duration]() {
		duration, err := time.ParseDuration(raw)
		if err != nil {
			return zero, err
		}
		var anyValue any = duration
		return anyValue.(T), nil
	}

	switch targetType.Kind() {
	case reflect.String:
		var anyValue any = raw
		return anyValue.(T), nil

	case reflect.Bool:
		parsedBool, err := strconv.ParseBool(raw)
		if err != nil {
			return zero, err
		}
		value := reflect.ValueOf(parsedBool).Convert(targetType)
		return value.Interface().(T), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsedInt, err := strconv.ParseInt(raw, 10, targetType.Bits())
		if err != nil {
			return zero, err
		}
		value := reflect.New(targetType).Elem()
		value.SetInt(parsedInt)
		return value.Interface().(T), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		parsedUint, err := strconv.ParseUint(raw, 10, targetType.Bits())
		if err != nil {
			return zero, err
		}
		value := reflect.New(targetType).Elem()
		value.SetUint(parsedUint)
		return value.Interface().(T), nil

	case reflect.Float32, reflect.Float64:
		parsedFloat, err := strconv.ParseFloat(raw, targetType.Bits())
		if err != nil {
			return zero, err
		}
		value := reflect.New(targetType).Elem()
		value.SetFloat(parsedFloat)
		return value.Interface().(T), nil

	default:
		return zero, fmt.Errorf("env: unsupported target type %v", targetType)
	}
}

func parseTime(raw string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	var lastErr error

	for _, layout := range layouts {
		parsed, err := time.Parse(layout, raw)
		if err == nil {
			return parsed, nil
		}
		lastErr = err
	}

	if lastErr == nil {
		lastErr = fmt.Errorf("env: could not parse time: %q", raw)
	}

	return time.Time{}, lastErr
}

func Load(filePaths ...string) {
	for _, filePath := range filePaths {
		if filePath == "" {
			continue
		}

		fileInfo, err := os.Stat(filePath)
		if err != nil || fileInfo.IsDir() {
			continue
		}

		if err := loadEnvFile(filePath); err == nil {
			return
		}
	}
}

func Get[T any](envVar string, optionalDefault ...T) T {
	var zero T
	rawValue, found := lookupEnv(envVar)
	if !found || strings.TrimSpace(rawValue) == "" {
		if len(optionalDefault) > 0 {
			return optionalDefault[0]
		}
		return zero
	}

	converted, err := convertStringToType[T](rawValue)
	if err != nil {
		if len(optionalDefault) > 0 {
			return optionalDefault[0]
		}
		return zero
	}

	return converted
}
