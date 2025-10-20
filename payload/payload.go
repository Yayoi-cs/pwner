package payload

import (
	"log"
	"pwner/utils"
)

func Pay(args ...interface{}) []byte {
	var payload []byte

	for _, arg := range args {
		switch v := arg.(type) {
		case []byte:
			payload = append(payload, v...)
		case string:
			payload = append(payload, []byte(v)...)
		case uint64:
			payload = append(payload, utils.P64(v)...)
		case uint32:
			payload = append(payload, utils.P32(v)...)
		case uint16:
			payload = append(payload, utils.P16(v)...)
		case uint8:
			payload = append(payload, utils.P8(v)...)
		case int:
			payload = append(payload, utils.P64(uint64(v))...)
		case int64:
			payload = append(payload, utils.P64(uint64(v))...)
		case int32:
			payload = append(payload, utils.P32(uint32(v))...)
		default:
			log.Fatalf("unsupported type: %T", v)
		}
	}

	return payload
}
