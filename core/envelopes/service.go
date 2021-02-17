package envelopes

import (
	"context"
	"errors"
	"redenvelope/infra/base"
	"redenvelope/services"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

var once sync.Once

func init() {
	once.Do(func() {
		services.IRedEnvelopeService = new(redEnvelopeService)
		logrus.Info("IRedEnvelopeService init xxxxxxxxxxxxx")
	})
}

type redEnvelopeService struct {
}

func (r *redEnvelopeService) SendOut(
	dto services.RedEnvelopeSendingDTO) (activity *services.RedEnvelopeActivity, err error) {

	err = base.ValidateStruct(&dto) //todo
	if err != nil {
		return activity, err
	}

	account := services.GetAccountService().GetEnvelopeAccountByUserID(dto.UserID)
	if account == nil {
		return nil, errors.New("用户账户不存在:" + dto.UserID)
	}
	goods := dto.ToGoods()
	goods.AccountNo = account.AccountNo
	if goods.Blessing == "" {
		goods.Blessing = services.DefaultBlessing
	}
	if goods.EnvelopeType == services.GeneralEnvelopeType {
		goods.AmountOne = goods.Amount
		goods.Amount = decimal.Decimal{}
	}
	//发红包
	domain := new(goodsDomain)
	activity, err = domain.SendOut(*goods)
	if err != nil {
		logrus.Error(err)
	}
	return activity, err
}

func (r *redEnvelopeService) Receive(
	dto services.RedEnvelopeReceiveDTO) (item *services.RedEnvelopeItemDTO, err error) {
	if err = base.ValidateStruct(&dto); err != nil {
		return item, err
	}
	account := services.GetAccountService().GetEnvelopeAccountByUserID(dto.RecvUserID)
	if account == nil {
		return nil, errors.New("红包资金账户不存在:user_id=" + dto.RecvUserID)
	}
	dto.AccountNo = account.AccountNo
	domain := goodsDomain{}
	item, err = domain.Receive(context.Background(), dto)
	return item, err
}
