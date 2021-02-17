package web

import (
	"redenvelope/infra"
	"redenvelope/infra/base"
	"redenvelope/services"

	"github.com/kataras/iris/v12"
	"github.com/sirupsen/logrus"
)

//资金账户根路径：/account
//版本号：/v1/account

var groupRouter iris.Party

func init() {
	infra.RegisterApi(&accountAPI{})
}

type accountAPI struct {
	service services.AccountService
}

func (a *accountAPI) Init() {
	a.service = services.GetAccountService()
	groupRouter := base.Iris().Party("/v1/account")
	groupRouter.Post("/create", a.createHandler)
	logrus.Info("account API init")
}

func (a *accountAPI) createHandler(ctx iris.Context) {

	account := services.AccountCreatedDTO{}
	err := ctx.ReadJSON(&account)
	res := base.Res{
		Code: base.ResCodeOK,
	}
	if err != nil {
		res.Code = base.ResCodeRequestParamError
		res.Message = err.Error()
		ctx.JSON(&res)
		return
	}
	//创建账户
	accountDto, err := a.service.CreateAccount(account)
	if err != nil {
		res.Code = base.ResCodeInnerServerError
		res.Message = err.Error()
	}
	res.Data = accountDto
	ctx.JSON(&res)
}

//账户创建接口: /v1/account/create
//转账接口： /v1/account/transfer
