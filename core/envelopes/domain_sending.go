package envelopes

import (
	"context"
	"path"
	"redenvelope/core/accounts"
	"redenvelope/infra/base"
	"redenvelope/services"

	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

func (d *goodsDomain) SendOut(
	goods services.RedEnvelopeGoodsDTO) (activity *services.RedEnvelopeActivity, err error) {
	//
	d.Create(goods)
	activity = new(services.RedEnvelopeActivity)
	link := base.GetEnvelopeActivityLink()
	domain := base.GetEnvelopeDomain()
	activity.Link = path.Join(domain, link, d.EnvelopeNo)

	accountDomain := accounts.NewAccountDomain()
	err = base.Tx(func(runner *dbx.TxRunner) error {
		ctx := base.WithValueContext(context.Background(), runner)
		id, err := d.Save(ctx)
		if id < 0 || err != nil {
			logrus.Error(err)
			return err
		}
		//扣减
		body := services.TradeParticipator{
			AccountNo: goods.AccountNo,
			UserID:    goods.UserID,
			Username:  goods.Username,
		}
		systemAccount := base.GetSystemAccount()
		target := services.TradeParticipator{
			AccountNo: systemAccount.AccountNo,
			UserID:    systemAccount.UserID,
			Username:  systemAccount.UserName,
		}
		transfer := services.AccountTransferDTO{
			TradeNo:     d.EnvelopeNo,
			TradeBody:   body,
			TradeTarget: target,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeOutgoing,
			ChangeFlag:  services.FlagTransferOut,
			Decs:        "红包金额支付",
		}
		status, err := accountDomain.TransferWithContextTx(ctx, transfer)
		if status != services.TransferedStatusSuccess {
			return err
		}

		//放入系统账户
		transfer = services.AccountTransferDTO{
			TradeNo:     d.EnvelopeNo,
			TradeBody:   target,
			TradeTarget: body,
			Amount:      d.Amount,
			ChangeType:  services.EnvelopeIncoming,
			ChangeFlag:  services.FlagTransferIn,
			Decs:        "红包金额转入",
		}
		status, err = accountDomain.TransferWithContextTx(ctx, transfer)
		if status != services.TransferedStatusSuccess {
			return err
		}
		return err
	})
	activity.RedEnvelopeGoodsDTO = *d.RedEnvelopeGoods.ToDTO()

	return activity, err
}
