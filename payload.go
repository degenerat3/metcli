package metcli

import (
	"encoding/base64"
	"strings"
)

// EncodePayload will base64 the payload and prepend/append magic chars
func EncodePayload(data string, mode string, aid string, m Metclient) string {
	preEnc := mode + "||" + aid + "||" + data
	encStr := base64.StdEncoding.EncodeToString([]byte(preEnc))
	fin := m.magicstring + encStr + m.magictermstr
	return fin
}

//DecodePayload will decode from MAD into plaintext string
func DecodePayload(payload string, m Metclient) string {
	encodedPayload := strings.Replace(payload, m.magicstring, "", -1) //trim magic chars from payload
	encodedPayload = strings.Replace(encodedPayload, m.magictermstr, "", -1)
	data, err := base64.StdEncoding.DecodeString(encodedPayload)
	if err != nil {
		return ""
	}
	return string(data)
}
