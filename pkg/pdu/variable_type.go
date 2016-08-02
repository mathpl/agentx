package pdu

import "fmt"

const (
	VariableTypeInteger          VariableType = 2
	VariableTypeOctetString      VariableType = 4
	VariableTypeNull             VariableType = 5
	VariableTypeObjectIdentifier VariableType = 6
	VariableTypeIPAddress        VariableType = 64
	VariableTypeCounter32        VariableType = 65
	VariableTypeGauge32          VariableType = 66
	VariableTypeTimeTicks        VariableType = 67
	VariableTypeOpaque           VariableType = 68
	VariableTypeCounter64        VariableType = 70
	VariableTypeNoSuchObject     VariableType = 128
	VariableTypeNoSuchInstance   VariableType = 129
	VariableTypeEndOfMIBView     VariableType = 130
)

type VariableType uint16

func (v VariableType) String() string {
	switch v {
	case VariableTypeInteger:
		return "variableTypeInteger"
	case VariableTypeOctetString:
		return "variableTypeOctetString"
	case VariableTypeNull:
		return "variableTypeNull"
	case VariableTypeObjectIdentifier:
		return "variableTypeObjectIdentifier"
	case VariableTypeIPAddress:
		return "variableTypeIPAddress"
	case VariableTypeCounter32:
		return "variableTypeCounter32"
	case VariableTypeGauge32:
		return "variableTypeGauge32"
	case VariableTypeTimeTicks:
		return "variableTypeTimeTicks"
	case VariableTypeOpaque:
		return "variableTypeOpaque"
	case VariableTypeCounter64:
		return "variableTypeCounter64"
	case VariableTypeNoSuchObject:
		return "variableTypeNoSuchObject"
	case VariableTypeNoSuchInstance:
		return "variableTypeNoSuchInstance"
	case VariableTypeEndOfMIBView:
		return "variableTypeEndOfMIBView"
	default:
		return fmt.Sprintf("variableTypeUnknown (%d)", v)
	}
}
