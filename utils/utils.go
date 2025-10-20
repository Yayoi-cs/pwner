package utils

import (
	"fmt"
	"os"
	"reflect"
	"runtime"
)
import "log"

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

var logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile)

func Fatal(format string, v ...interface{}) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	logger.Fatalf("["+funcName+"] "+format, v...)
}

func Hexdump(data []byte) {
	fmt.Println()
	const bytesPerLine = 16

	for i := 0; i < len(data); i += bytesPerLine {
		fmt.Printf("  %08x  ", i)

		for j := 0; j < bytesPerLine; j++ {
			if i+j < len(data) {
				fmt.Printf("%02x ", data[i+j])
			} else {
				fmt.Print("   ")
			}

			// Add extra space in the middle
			if j == 7 {
				fmt.Print(" ")
			}
		}

		fmt.Print(" |")
		for j := 0; j < bytesPerLine && i+j < len(data); j++ {
			b := data[i+j]
			if b >= 32 && b <= 126 {
				fmt.Printf("%c", b)
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("|")
		fmt.Println()
	}
}
