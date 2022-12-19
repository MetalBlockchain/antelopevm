package types

type TransactionTrace struct {
	Id           TransactionIdType
	BlockNum     uint32
	BlockTime    BlockTimeStamp
	Receipt      TransactionReceiptHeader
	Elapsed      Microseconds
	NetUsage     uint64
	Scheduled    bool
	ActionTraces []ActionTrace
	Except       error
	ErrorCode    uint64
}

type ActionTrace struct {
	ActionOrdinal        int
	CreatorActionOrdinal int
	Receipt              ActionReceipt
	Receiver             ActionName
	Action               Action
	ContextFree          bool
	Elapsed              uint64
	Console              string
	TransactionId        TransactionIdType
	BlockNum             uint32
	BlockTime            BlockTimeStamp
	Except               error
	ErrorCode            uint64
}

func NewActionTrace(trace *TransactionTrace, action Action, receiver AccountName, contextFree bool, actionOrdinal int, creatorActionOrdinal int) *ActionTrace {
	return &ActionTrace{
		ActionOrdinal:        actionOrdinal,
		CreatorActionOrdinal: creatorActionOrdinal,
		Receiver:             receiver,
		Action:               action,
		ContextFree:          contextFree,
		TransactionId:        trace.Id,
		BlockNum:             trace.BlockNum,
		BlockTime:            trace.BlockTime,
	}
}
