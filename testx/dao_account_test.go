package testx

import (
	"database/sql"
	"redenvelope/core/accounts"
	"redenvelope/infra/base"
	"testing"

	"github.com/segmentio/ksuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/tietang/dbx"
)

func TestAccountDao_GetOne(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := accounts.NewAccountDao(runner)

		Convey("通过编号查询账户数据", t, func() {
			a := &accounts.Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "test account2",
				UserID:      ksuid.New().Next().String(),
				Username:    sql.NullString{String: "test account2", Valid: true},
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
			na := dao.GetOne(a.AccountNo)
			So(na, ShouldNotBeNil)
			So(na.Balance.String(), ShouldEqual, a.Balance.String())
			So(na.CreatedAt, ShouldNotBeNil)
			So(na.UpdatedAt, ShouldNotBeNil)
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_GetByUserID(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := accounts.NewAccountDao(runner)
		userID := ksuid.New().Next().String()
		Convey("通过编号查询账户数据", t, func() {
			a := &accounts.Account{
				Balance:     decimal.NewFromFloat(100),
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "test account",
				UserID:      userID,
				Username:    sql.NullString{String: "test account", Valid: true},
				AccountType: 2,
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)
			na := dao.GetByUserID(userID, 2)
			So(na, ShouldNotBeNil)
			So(na.Balance.String(), ShouldEqual, a.Balance.String())
			So(na.CreatedAt, ShouldNotBeNil)
			So(na.UpdatedAt, ShouldNotBeNil)
		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}

func TestAccountDao_UpdateBalance(t *testing.T) {
	err := base.Tx(func(runner *dbx.TxRunner) error {
		dao := accounts.NewAccountDao(runner)
		balance := decimal.NewFromFloat(100)
		Convey("更新账户余额", t, func() {
			a := &accounts.Account{
				Balance:     balance,
				Status:      1,
				AccountNo:   ksuid.New().Next().String(),
				AccountName: "测试资金账户",
				UserID:      ksuid.New().Next().String(),
				Username:    sql.NullString{String: "测试用户", Valid: true},
			}
			id, err := dao.Insert(a)
			So(err, ShouldBeNil)
			So(id, ShouldBeGreaterThan, 0)

			//1.增加余额
			Convey("增加余额", func() {
				amount := decimal.NewFromFloat(10)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 1)
				na := dao.GetOne(a.AccountNo)
				newBalance := balance.Add(amount)
				So(na, ShouldNotBeNil)
				So(na.Balance.String(), ShouldEqual, newBalance.String())
			})

			//2.扣减余额，余额足够
			//把updateBalance余额扣减的第二个测试用例
			//作为作业留给同学们来编写

			//3.扣减余额，余额不够
			Convey("扣减余额，余额不够", func() {
				a1 := dao.GetOne(a.AccountNo)
				So(a1, ShouldNotBeNil)
				amount := decimal.NewFromFloat(-300)
				rows, err := dao.UpdateBalance(a.AccountNo, amount)
				So(err, ShouldBeNil)
				So(rows, ShouldEqual, 0)
				a2 := dao.GetOne(a.AccountNo)
				So(a2, ShouldNotBeNil)
				So(a1.Balance.String(), ShouldEqual, a2.Balance.String())
			})

		})
		return nil
	})
	if err != nil {
		logrus.Error(err)
	}
}
