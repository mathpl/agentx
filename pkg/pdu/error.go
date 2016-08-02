package pdu

import "fmt"

type Error uint16

const (
	ErrorNone       Error = 0
	ErrorOpenFailed Error = iota + 256
	ErrorNotOpen
	ErrorIndexWrongType
	ErrorIndexAlreadyAllocated
	ErrorIndexNoneAvailable
	ErrorIndexNotAllocated
	ErrorUnsupportedContext
	ErrorDuplicateRegistration
	ErrorUnknownRegistration
	ErrorUnknownAgentCaps
	ErrorParse
	ErrorRequestDenied
	ErrorProcessing
)

func (e Error) String() string {
	switch e {
	case ErrorNone:
		return "errorNone"
	case ErrorOpenFailed:
		return "errorOpenFailed"
	case ErrorNotOpen:
		return "errorNotOpen"
	case ErrorIndexWrongType:
		return "errorIndexWrongType"
	case ErrorIndexAlreadyAllocated:
		return "errorIndexAlreadyAllocated"
	case ErrorIndexNoneAvailable:
		return "errorIndexNoneAvailable"
	case ErrorIndexNotAllocated:
		return "errorIndexNotAllocated"
	case ErrorUnsupportedContext:
		return "errorUnsupportedContext"
	case ErrorDuplicateRegistration:
		return "errorDuplicateRegistration"
	case ErrorUnknownRegistration:
		return "errorUnknownRegistration"
	case ErrorUnknownAgentCaps:
		return "errorUnknownAgentCaps"
	case ErrorParse:
		return "errorParse"
	case ErrorRequestDenied:
		return "errorRequestDenied"
	case ErrorProcessing:
		return "errorProcessing"
	default:
		return fmt.Sprintf("errorUnknown (%d)", e)
	}
}
