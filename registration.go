package metcli

import (
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"strings"
)

func genClient(serv string, magic []byte, magicstring string, magicterm []byte, magictermstr string, regfile string, interval int, delta int, obfseed int, obftext string) Metclient {
	m := Metclient{serv, magic, magicstring, magicterm, magictermstr, regfile, interval, delta, obfseed, obftext}
	return m
}

//get IP of default interface, which the DB will use as hostname
func getIP() string {
	conn, _ := net.Dial("udp", "8.8.8.8:80")
	defer conn.Close()
	ad := conn.LocalAddr().(*net.UDPAddr)
	ipStr := ad.IP.String()
	return ipStr
}

//write the UUID to somewhere on disk
func storeUUID(uid string, m Metclient) {
	obf := obfuscateUUID(uid, m.obfseed, m.obftext)
	f, _ := os.Create(m.regfile)
	io.WriteString(f, obf)
	f.Close()
	return
}

//grab the obfuscated UUID from somewhere on disk
func fetchUUID(m Metclient) string {
	obf, _ := ioutil.ReadFile(m.regfile)
	deobf := deobfuscateUUID(string(obf), m)
	return deobf
}

//simple obfuscation so you cant just search the filesystem for uuid formatted strings
func obfuscateUUID(uid string, seed int, text string) string {
	splituid := strings.Split(uid, "-")
	l1 := strings.Repeat(text, rand.Intn(seed))
	l2 := strings.Repeat(text, rand.Intn(seed))
	l3 := strings.Repeat(text, rand.Intn(seed))
	l4 := strings.Repeat(text, rand.Intn(seed))
	obf := splituid[0] + l1 + splituid[1] + splituid[2] + l2 + splituid[3] + l3 + splituid[4] + l4
	// its' really crappy obfuscation but it's a small deterrent
	return obf
}

//undo UUID obfuscations
func deobfuscateUUID(obf string, m Metclient) string {
	p := strings.Replace(obf, m.obftext, "", -1)
	p = p[:8] + "-" + p[8:]
	p = p[:13] + "-" + p[13:]
	p = p[:18] + "-" + p[18:]
	p = p[:23] + "-" + p[23:]
	return p
}
