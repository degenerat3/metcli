package metcli

import (
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
)

//write the UUID to somewhere on disk
func storeUUID(uid string) {
	obf := obfuscateUUID(uid, OBFSEED, OBFTEXT)
	f, _ := os.Create(REGFILE)
	io.WriteString(f, obf)
	f.Close()
	return
}

//grab the obfuscated UUID from somewhere on disk
func fetchUUID() string {
	obf, _ := ioutil.ReadFile(REGFILE)
	deobf := deobfuscateUUID(string(obf))
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
func deobfuscateUUID(obf string) string {
	p := strings.Replace(obf, OBFTEXT, "", -1)
	p = p[:8] + "-" + p[8:]
	p = p[:13] + "-" + p[13:]
	p = p[:18] + "-" + p[18:]
	p = p[:23] + "-" + p[23:]
	return p
}
