package pdu

import (
	"encoding/binary"

	"github.com/mathpl/typed"
)

type Variables []Variable

func (v *Variables) Add(oid OID, t VariableType, value interface{}) {
	variable := Variable{}
	variable.Set(oid, t, value)
	*v = append(*v, variable)
}

func (v *Variables) Read(f Flags, rb *typed.ReadBuffer) error {
	for rb.BytesRemaining() > 0 {
		variable := Variable{}
		if err := variable.Read(f, rb); err != nil {
			return err
		}
		*v = append(*v, variable)
	}
	return rb.Err()
}

func (v *Variables) Write(bo binary.ByteOrder, wb *typed.WriteBuffer) error {
	for _, variable := range *v {
		variable.Write(bo, wb)
	}
	return wb.Err()
}
