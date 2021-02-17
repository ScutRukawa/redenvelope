package testx

import (
	"redenvelope/core/accounts"
	_ "redenvelope/core/envelopes"

	_ "redenvelope/initfortest"

	"redenvelope/services"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRedenvelopeService_SendOut(t *testing.T) {
	ac := accounts.AccountService{}
	account := services.AccountCreatedDTO{
		UserID:       ksuid.New().Next().String(),
		Username:     "测试用户",
		Amount:       "200",
		AccountName:  "测试账户",
		AccountType:  int(services.EnvelopeAccountType),
		CurrencyCode: "CNY",
	}
	re := services.GetRedEnvelopeService()

	Convey("发红包", t, func() {
		acDTO, err := ac.CreateAccount(account)
		So(err, ShouldBeNil)
		So(acDTO, ShouldNotBeNil)

		goods := services.RedEnvelopeSendingDTO{
			EnvelopeType: services.GeneralEnvelopeType,
			UserID:       account.UserID,
			Username:     account.Username,
			Amount:       decimal.NewFromFloat(8.88),
			Quantity:     10,
			Blessing:     services.DefaultBlessing,
		}

		Convey("发普通红包", func() {
			at, err := re.SendOut(goods)
			So(err, ShouldBeNil)
			So(at, ShouldNotBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)

			//验证每一个属性
			dto := at.RedEnvelopeGoodsDTO
			So(dto.Username, ShouldEqual, goods.Username)
			So(dto.UserID, ShouldEqual, goods.UserID)
			So(dto.Quantity, ShouldEqual, goods.Quantity)
			q := decimal.NewFromFloat(float64(dto.Quantity))
			So(dto.Amount.String(), ShouldEqual, goods.Amount.Mul(q).String())
		})

		goods.EnvelopeType = services.LuckyEnvelopeType
		goods.Amount = decimal.NewFromFloat(88.8)
		Convey("发运气红包", func() {
			at, err := re.SendOut(goods)
			So(err, ShouldBeNil)
			So(at, ShouldNotBeNil)
			So(at.Link, ShouldNotBeEmpty)
			So(at.RedEnvelopeGoodsDTO, ShouldNotBeNil)

			//验证每一个属性
			dto := at.RedEnvelopeGoodsDTO
			So(dto.Username, ShouldEqual, goods.Username)
			So(dto.UserID, ShouldEqual, goods.UserID)
			So(dto.Quantity, ShouldEqual, goods.Quantity)
			So(dto.Amount.String(), ShouldEqual, goods.Amount.String())
		})

	})
}
