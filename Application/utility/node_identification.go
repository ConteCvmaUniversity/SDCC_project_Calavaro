package utility

import "strings"

type NodeType int

const (
	Peer      NodeType = 0
	Register           = 1
	Sequencer          = 2
)

// Struct to send information about peer
type Info struct {
	Type    NodeType
	ID      int
	Address string
	Port    string
}

func ParseLine(s string, sep string) (string, string, string) {
	res := strings.Split(s, sep)
	return res[0], res[1], res[2]
}

func TypeToString(nodeType NodeType) string {
	switch nodeType {
	case Peer:
		return "peer"
	case Register:
		return "register"
	case Sequencer:
		return "sequencer"
	}
	return ""
}

func StringToType(s string) NodeType {
	switch s {
	case "peer":
		return Peer
	case "register":
		return Register
	case "sequencer":
		return Sequencer
	}
	return -1
}
