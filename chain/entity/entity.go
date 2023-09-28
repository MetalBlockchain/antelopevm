package entity

const (
	AccountType uint8 = iota
	PermissionType
	PermissionLinkType
	TableType
	KeyValueType
	BlockType
	IndexObjectType
	ResourceUsageType
	ResourceLimitType
	GlobalPropertyObjectType
	TransactionObjectType
	AccountMetaDataObjectType
	AccountRamCorrectionObjectType
	CodeObjectType
)

type EntityIndex struct {
	Name   string
	Fields []string
}

type Entity interface {
	GetId() []byte
	GetIndexes() map[string]EntityIndex
	GetObjectType() uint8
}
