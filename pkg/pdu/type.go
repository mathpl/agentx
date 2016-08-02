package pdu

const (
	TypeOpen            Type = 1
	TypeClose           Type = 2
	TypeRegister        Type = 3
	TypeUnregister      Type = 4
	TypeGet             Type = 5
	TypeGetNext         Type = 6
	TypeGetBulk         Type = 7
	TypeTestSet         Type = 8
	TypeCommitSet       Type = 9
	TypeUndoSet         Type = 10
	TypeCleanupSet      Type = 11
	TypeNotify          Type = 12
	TypePing            Type = 13
	TypeIndexAllocate   Type = 14
	TypeIndexDeallocate Type = 15
	TypeAddAgentCaps    Type = 16
	TypeRemoveAgentCaps Type = 17
	TypeResponse        Type = 18
)

type Type byte

type TypeOwner interface {
	Type() Type
}

func (t Type) String() string {
	switch t {
	case TypeOpen:
		return "typeOpen"
	case TypeClose:
		return "typeClose"
	case TypeRegister:
		return "typeRegister"
	case TypeUnregister:
		return "typeUnregister"
	case TypeGet:
		return "typeGet"
	case TypeGetNext:
		return "typeGetNext"
	case TypeGetBulk:
		return "typeGetBulk"
	case TypeTestSet:
		return "typeTestSet"
	case TypeCommitSet:
		return "typeCommitSet"
	case TypeUndoSet:
		return "typeUndoSet"
	case TypeCleanupSet:
		return "typeCleanupSet"
	case TypeNotify:
		return "typeNotify"
	case TypePing:
		return "typePing"
	case TypeIndexAllocate:
		return "typeIndexAllocate"
	case TypeIndexDeallocate:
		return "typeIndexDeallocate"
	case TypeAddAgentCaps:
		return "typeAddAgentCaps"
	case TypeRemoveAgentCaps:
		return "typeRemoveAgentCaps"
	case TypeResponse:
		return "typeResponse"
	default:
		return "typeUnknown"
	}
}
