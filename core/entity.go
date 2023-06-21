package core

const (
	AccountType              uint8 = 0
	PermissionType           uint8 = 1
	PermissionLinkType       uint8 = 2
	TableType                uint8 = 3
	KeyValueType             uint8 = 4
	TransactionType          uint8 = 5
	BlockType                uint8 = 6
	IndexObjectType          uint8 = 7
	ResourceUsageType        uint8 = 8
	ResourceLimitType        uint8 = 9
	GlobalPropertyObjectType uint8 = 10
	TransactionObjectType    uint8 = 11
)

type EntityIndex struct {
	Name   string
	Fields []string
}

type Entity interface {
	GetId() []byte
	GetIndexes() map[string]EntityIndex
	GetObjectType() uint8
	MarshalMsg(b []byte) (o []byte, err error)
	UnmarshalMsg(bts []byte) (o []byte, err error)
}
