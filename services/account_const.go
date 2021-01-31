package services

//转账状态

type TransferedStatus int8

const (
	TransferedStatusFailure         TransferedStatus = -1
	TransferedStatusSufficientFunds TransferedStatus = 0
	TransferedStatusSuccess         TransferedStatus = 1
)

//转账类型 创建账户 进账 支出
type ChangeType int8

const (
	AccountCreated        ChangeType = 0
	AccountStoreValue     ChangeType = 1
	EnvelopeOutgoing      ChangeType = -2
	EnvelopeIncoming      ChangeType = 2
	EnvelopeExpiredRefund ChangeType = 3
)

//资金交易变化标识
type ChangeFlag int8

const (
	FlagAccountCreated ChangeFlag = 0
	FlagTransferOut    ChangeFlag = -1
	FlagTransferIn     ChangeFlag = 1
)

type AccountType int8

const (
	EnvelopeAccountType       AccountType = 1
	SystemEnvelopeAccountType AccountType = 2
)
