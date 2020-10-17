package state

import (

)

type State struct {
	Name string
	Authority string
	Peers map[string]*Peer
	Debug bool
}

func NewState(name string, authority string, debug bool) State {
	return State {
		Name: name,
		Authority: authority,
		Peers: make(map[string]*Peer),
		Debug: debug,
	}
}

func (this *State) Add(peer *Peer) {
	if !this.Has(peer) {
		this.Peers[peer.Hash()] = peer
	}
}

func (this *State) Has(peer *Peer) bool {
	_, ok := this.Peers[peer.Hash()]

	return ok
}

func (this *State) Remove(peer *Peer) {
	if this.Has(peer) {
		delete(this.Peers, peer.Hash())
	}
}

func (this *State) Count() int {
	return len(this.Peers)
}

func (this *State) Broadcast(data string) {
	for _, peer := range this.Peers {
		if peer.Authenticated {
			peer.SendMsg(data)
		}
	}
}

func (this *State) Quit() {
	for _, peer := range this.Peers {
		if peer.Authenticated {
			peer.SendQuit()
		}
	}
}

func (this *State) AllAuthorities() []string {
	var authorities []string
	for _, peer := range this.Peers {
		if peer.Authenticated {
			authorities = append(authorities, peer.PublicAuthority)
		}
	}

	return authorities
}

func (this *State) Knows(authority string) bool {
	if authority == this.Authority {
		return true
	}
	
	for _, peer := range this.Peers {
		if peer.PublicAuthority == authority {
			return true
		}
	}

	return false
}