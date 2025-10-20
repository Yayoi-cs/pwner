package utils

import "reflect"

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorGray   = "\033[90m"
)

func Type2Color(v interface{}) string {
	if v == nil {
		return ColorGray
	}

	val := reflect.ValueOf(v)

	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return ColorRed
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return ColorRed
	case reflect.Float32, reflect.Float64:
		return ColorBlue
	case reflect.String:
		return ColorGreen
	case reflect.Bool:
		return ColorYellow
	case reflect.Slice, reflect.Array:
		return ColorPurple
	case reflect.Map:
		return ColorCyan
	case reflect.Struct:
		return ColorWhite
	case reflect.Ptr:
		return ColorGray
	default:
		return ColorReset
	}
}
