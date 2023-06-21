package chain

type DeadlineException struct{}

func (e DeadlineException) Code() int64 { return 3080006 }

type BlockCpuUsageExceededException struct{}

func (e BlockCpuUsageExceededException) Code() int64 { return 3080005 }

type TxCpuUsageExceededException struct{}

func (e TxCpuUsageExceededException) Code() int64 { return 3080004 }

type LeewayDeadlineException struct{}

func (e LeewayDeadlineException) Code() int64 { return 3081001 }
