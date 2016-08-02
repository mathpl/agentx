package pdu

import "fmt"

const (
	ReasonOther Reason = iota + 1
	ReasonParseError
	ReasonProtocolError
	ReasonTimeouts
	ReasonShutdown
	ReasonByManager
)

type Reason byte

func (r Reason) String() string {
	switch r {
	case ReasonOther:
		return "reasonOther"
	case ReasonParseError:
		return "reasonParseError"
	case ReasonProtocolError:
		return "reasonProtocolError"
	case ReasonTimeouts:
		return "reasonTimeouts"
	case ReasonShutdown:
		return "reasonShutdown"
	case ReasonByManager:
		return "reasonByManager"
	default:
		return fmt.Sprintf("reasonUnknown (%d)", r)
	}
}
