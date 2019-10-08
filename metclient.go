package metcli

// Metclient is a struct that holds all the client data that will be utilized
type Metclient struct {
	serv         string // server it will call back to
	magic        []byte // byte representation of magicstring
	magicstring  string // string representation of magicstring
	magicterm    []byte // byte representation of magicterm
	magictermstr string // string representation of magicterm
	regfile      string // location of bot registration file
	interval     int    // bot callback interval
	delta        int    // variation in callback time
	obfseed      string // obfuscation seed integer
	obftext      string // obfuscation seed text
}
