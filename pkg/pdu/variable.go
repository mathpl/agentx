package pdu

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"

	"github.com/mathpl/typed"
)

type Variable struct {
	Type  VariableType
	Name  ObjectIdentifier
	Value interface{}
}

func (v *Variable) Set(oid OID, t VariableType, value interface{}) {
	v.Name.SetIdentifier(oid)
	v.Type = t
	v.Value = value
}

func (v *Variable) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	wb.WriteUint16Endian(bo, uint16(v.Type))
	wb.WriteUint16(0)
	v.Name.Write(bo, wb)

	switch v.Type {
	case VariableTypeInteger:
		value := v.Value.(int32)
		wb.WriteUint32Endian(bo, uint32(value))
	case VariableTypeOctetString:
		octetString := &OctetString{Text: v.Value.(string)}
		octetString.Write(bo, wb)
	case VariableTypeNull, VariableTypeNoSuchObject, VariableTypeNoSuchInstance, VariableTypeEndOfMIBView:
		break
	case VariableTypeObjectIdentifier:
		targetOID, err := ParseOID(v.Value.(string))
		if err != nil {
			return err
		}

		oi := &ObjectIdentifier{}
		oi.SetIdentifier(targetOID)
		oi.Write(bo, wb)
	case VariableTypeIPAddress:
		ip := v.Value.(net.IP)
		octetString := &OctetString{Text: string(ip)}
		octetString.Write(bo, wb)
	case VariableTypeCounter32, VariableTypeGauge32:
		value := v.Value.(uint32)
		wb.WriteUint32Endian(bo, value)
	case VariableTypeTimeTicks:
		value := uint32(v.Value.(time.Duration).Seconds() * 100)
		wb.WriteUint32Endian(bo, value)
	case VariableTypeOpaque:
		octetString := &OctetString{Text: string(v.Value.([]byte))}
		octetString.Write(bo, wb)
	case VariableTypeCounter64:
		value := v.Value.(uint64)
		wb.WriteUint64Endian(bo, value)
	default:
		return fmt.Errorf("unhandled variable type %s", v.Type)
	}

	return wb.Err()
}

func (v *Variable) Read(f Flags, rb *typed.ReadBuffer) error {
	bo := f.GetByteOrder()

	v.Type = VariableType(rb.ReadUint16Endian(bo))

	// reserved
	rb.ReadUint16()

	v.Name.Read(f, rb)

	switch v.Type {
	case VariableTypeInteger, VariableTypeCounter32, VariableTypeGauge32:
		v.Value = rb.ReadUint32Endian(bo)
	case VariableTypeOctetString:
		octetString := &OctetString{}
		octetString.Read(f, rb)
		v.Value = octetString.Text
	case VariableTypeNull, VariableTypeNoSuchObject, VariableTypeNoSuchInstance, VariableTypeEndOfMIBView:
		v.Value = nil
	case VariableTypeObjectIdentifier:
		oid := &ObjectIdentifier{}
		oid.Read(f, rb)
		v.Value = oid.GetIdentifier()
	case VariableTypeIPAddress:
		octetString := &OctetString{}
		octetString.Read(f, rb)
		v.Value = net.IP(octetString.Text)
	case VariableTypeTimeTicks:
		v.Value = time.Duration(rb.ReadUint32Endian(bo)) * time.Second / 100
	case VariableTypeOpaque:
		octetString := &OctetString{}
		octetString.Read(f, rb)
		v.Value = []byte(octetString.Text)
	case VariableTypeCounter64:
		v.Value = rb.ReadUint64Endian(bo)
	default:
		return fmt.Errorf("unhandled variable type %s", v.Type)
	}

	return rb.Err()
}
