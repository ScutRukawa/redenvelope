package base

//ResCode x
type ResCode int64

const (
	ResCodeOK                ResCode = 1000
	ResCodeValifationError   ResCode = 2000
	ResCodeRequestParamError ResCode = 2100
	ResCodeInnerServerError  ResCode = 5000
	ResCodeBizError          ResCode = 6000
)

//Res isris
type Res struct {
	Code    ResCode     `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
