package iptable

import "strings"

func (b *iptableBackend) parseActionREJECT(input string) interface{} {
	parts := strings.Fields(input)
	action := ActionReject{}

	for i := 2; i < len(parts); i++ {
		switch parts[i] {
		case "--reject-with":
			if i+1 < len(parts) {
				action.RejectWith = parts[i+1]
				i++
			}
		}
	}

	return action
}

func (b *iptableBackend) parseActionLOG(input string) interface{} {
	action := ActionLog{}

	return action
}

func (b *iptableBackend) parseActionNFLOG(input string) interface{} {
	action := ActionNFLOG{}

	return action
}

func (b *iptableBackend) parseActionNAT(input string) interface{} {
	action := ActionNat{}

	return action
}

func (b *iptableBackend) parseActionMARK(input string) interface{} {
	action := ActionMark{}

	return action
}

func (b *iptableBackend) parseActionCONNMARK(input string) interface{} {
	action := ActionConnMark{}

	return action
}

func (b *iptableBackend) parseActionTEE(input string) interface{} {
	action := ActionTEE{}

	return action
}
