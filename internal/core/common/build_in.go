package common

import (
	"bytes"
	"compress/flate"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/url"
	"strconv"
	"strings"
)

func Rot13char(str string) string {
	result := strings.Map(rot13, str)
	return result
}
func Base64Decode(data string) string {
	var ret string
	if v, err := base64.StdEncoding.DecodeString(data); err == nil {
		ret = string(v)
	}
	return ret
}
func Urldecode(str string) string {
	if ret, err := url.QueryUnescape(str); err == nil {
		return ret
	}
	return ""
}

func GzInflate(data string) (string, error) {
	rawData := []byte(data)
	reader := flate.NewReader(bytes.NewReader(rawData))
	defer reader.Close()

	decompressedData, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("failed to decompress data: %w", err)
	}

	return string(decompressedData), nil
}

func StrReplace(search string, replace string, subject string) string {
	return strings.Replace(search, replace, subject, -1)
}

func rot13(r rune) rune {
	switch {
	case r >= 'a' && r <= 'z':
		return 'a' + (r-'a'+13)%26
	case r >= 'A' && r <= 'Z':
		return 'A' + (r-'A'+13)%26
	default:
		return r
	}
}

func Hex2String(hexString string) (string, error) {
	var decodedString strings.Builder
	i := 0
	for i < len(hexString) {
		if hexString[i] == '\\' && i+1 < len(hexString) {
			switch hexString[i+1] {
			case 'x':
				if i+3 < len(hexString) {
					hexByte, err := strconv.ParseUint(hexString[i+2:i+4], 16, 8)
					if err != nil {
						return "", err
					}
					decodedString.WriteByte(byte(hexByte))
					i += 4
				} else {
					return "", fmt.Errorf("Invalid escape sequence at position %d", i)
				}
			case '0', '1', '2', '3', '4', '5', '6', '7':
				if i+2 < len(hexString) {
					octalByte, err := strconv.ParseUint(hexString[i+1:i+4], 8, 8)
					if err != nil {
						return "", err
					}
					decodedString.WriteByte(byte(octalByte))
					i += 4
				} else {
					return "", fmt.Errorf("Invalid escape sequence at position %d", i)
				}
			default:
				decodedString.WriteByte(hexString[i+1])
				i += 2
			}
		} else {
			decodedString.WriteByte(hexString[i])
			i++
		}
	}
	return decodedString.String(), nil
}

func DecodeEscapeSequence(input string) (string, error) {
	var decoded strings.Builder
	items := strings.Split(input, "\\")

	for _, item := range items {
		if item == "" {
			continue
		}
		if item[0] == 'x' {
			hexCode := strings.Trim(item, "x")
			decodedChar, err := strconv.ParseInt(hexCode, 16, 8)
			if err != nil {
				return "", err
			}
			decoded.WriteByte(byte(decodedChar))
		} else if item[0] <= '9' && item[0] >= '0' {
			decodedChar, err := strconv.ParseInt(item, 10, 8)
			if err != nil {
				return "", err
			}
			decoded.WriteByte(byte(decodedChar))
		}
	}
	return decoded.String(), nil
}
