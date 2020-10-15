package state

import (
	"crypto/sha256"
	"fmt"
	"net"
	"strings"

	"github.com/neghmurken/mdr/pkg/util"
)

const BUFFER_SIZE = 1024

type Peer struct {
	Name string
	Conn net.Conn
	Authenticated bool
	PublicAuthority string
	Debug bool
}

func NewPeer(conn net.Conn, name string, debug bool) *Peer {
	return &Peer {
		Name: name,
		Conn: conn,
		Authenticated: false,
		Debug: debug,
	}
}

func (this *Peer) Authority() string {
	return this.Conn.RemoteAddr().String()
}

func (this *Peer) Hash() string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(this.Authority())))
}

func (this *Peer) Authenticate(name string, authority string) {
	this.Name = name
	this.Authenticated = true
	this.PublicAuthority = authority
}

func (this *Peer) Send(data string) {
	if this.Debug {
		fmt.Printf("[DEBUG]{%s} > %s\n", this.Authority(), data)
	}
	fmt.Fprintf(this.Conn, data + "\n")
}

func (this *Peer) SendInit(name string, authority string) {
	this.Send(fmt.Sprintf("KIKOO \"%s\" %s ASV", name, authority))
}

func (this *Peer) SendIdent(name string, authority string, otherAuthorities []string) {
	others := ""
	for _, other := range otherAuthorities {
		if other != this.PublicAuthority {
			others = others + " / " + other
		}
	}

	this.Send(fmt.Sprintf("OKLM \"%s\" %s%s", name, authority, others))
}

func (this *Peer) SendMsg(message string) {
	this.Send(fmt.Sprintf("TAVU \"%s\"", strings.ReplaceAll(message, "\"", "\"\"")))
}

func (this *Peer) SendAck() {
	this.Send("LOL")
}

func (this *Peer) SendQuit() {
	this.Send("JPP")
}

func (this *Peer) SendErr(explanation string) {
	this.Send(fmt.Sprintf("WTF \"%s\"", strings.ReplaceAll(explanation, "\"", "\"\"")))
}

func (this *Peer) Listen(callback func(string, int) bool) {
	for {
		var buff = make([]byte, BUFFER_SIZE)
        byteRead, err := this.Conn.Read(buff)

        if err != nil {
            break
		}
		
		data := util.CleanData(string(buff))

		if this.Debug {
			fmt.Printf("[DEBUG]{%s} < %s\n", this.Authority(), data)
		}
		if !callback(data, byteRead) {
			break
		}
    }
}