package shell

import (
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf16"
)

func parseParcel(text string) (string, error) {
	if !strings.HasPrefix(text, "Result: Parcel(") {
		return "", fmt.Errorf("%s", text)
	}
	re := regexp.MustCompile(`0x[0-9a-fA-F]+: ([0-9a-fA-F ]{35})`)
	matches := re.FindAllStringSubmatch(text, -1)

	var u16Data []uint16

	for _, match := range matches {
		hexStr := strings.ReplaceAll(match[1], " ", "")

		for i := 0; i < len(hexStr); i += 8 {
			chunk := hexStr[i : i+8]
			val1 := hexToUint16(chunk[4:8])
			val2 := hexToUint16(chunk[0:4])

			u16Data = append(u16Data, val1, val2)
		}
	}

	runes := utf16.Decode(u16Data)
	var sb strings.Builder
	for _, r := range runes {
		if r >= 32 && r <= 126 {
			sb.WriteRune(r)
		} else if r > 128 {
			sb.WriteRune(r)
		}
	}

	return sb.String(), nil
}

func hexToUint16(s string) uint16 {
	b, _ := hex.DecodeString(s)
	if len(b) < 2 {
		return 0
	}
	return uint16(b[0])<<8 | uint16(b[1])
}
