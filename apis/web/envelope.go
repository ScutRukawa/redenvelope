package web

import (
	"redenvelope/infra"
	"redenvelope/infra/base"
	"redenvelope/services"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

func init() {
	infra.RegisterApi(&EnvelopeAPI{})
}

type EnvelopeAPI struct {
	service services.RedEnvelopeService
}

func (e *EnvelopeAPI) Init() {
	e.service = services.GetRedEnvelopeService()
	groupRouter := base.Iris().Party("/v1/envelope")
	groupRouter.Post("/sendout", e.SendOutHandler)
	groupRouter.Post("/receive", e.ReceiveHandler)
	logrus.Info("envelope api init")
}

func (e *EnvelopeAPI) SendOutHandler(ctx iris.Context) {
	//ctx.Params().
	dto := services.RedEnvelopeSendingDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOK,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	item, err := e.service.SendOut(dto)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = item
	ctx.JSON(&r)
}

func (e *EnvelopeAPI) ReceiveHandler(ctx iris.Context) {
	dto := services.RedEnvelopeReceiveDTO{}
	err := ctx.ReadJSON(&dto)
	r := base.Res{
		Code: base.ResCodeOK,
	}
	if err != nil {
		r.Code = base.ResCodeRequestParamError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	activity, err := e.service.Receive(dto)
	if err != nil {
		r.Code = base.ResCodeInnerServerError
		r.Message = err.Error()
		ctx.JSON(&r)
		return
	}
	r.Data = activity
	ctx.JSON(&r)
}
