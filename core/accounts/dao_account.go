package accounts

import (
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountDao struct {
	runner *dbx.TxRunner
}

func NewAccountDao(runner *dbx.TxRunner) *AccountDao {
	return &AccountDao{runner: runner}
}

//GetOne 查询数据库持久化对象的单实例
func (dao *AccountDao) GetOne(accountNo string) *Account {
	a := &Account{AccountNo: accountNo}
	ok, err := dao.runner.GetOne(a) //必须传入唯一索引字段
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

//GetByUserID x
func (dao *AccountDao) GetByUserID(
	userID string, accountType int) *Account {
	a := &Account{}
	sql := "select * from account" +
		" where user_id=? and account_type=?"
	ok, err := dao.runner.Get(a, sql, userID, accountType)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

//Insert 账户数据插入
func (dao *AccountDao) Insert(a *Account) (id int64, err error) {
	rs, err := dao.runner.Insert(a)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

//UpdateBalance 账户余额更新
func (dao *AccountDao) UpdateBalance(
	accountNo string,
	amount decimal.Decimal) (rows int64, err error) {
	sql := "update account " +
		"set balance=balance+CAST(? AS DECIMAL(30,6)) " +
		"where account_no=? and balance>=-1*CAST(? AS DECIMAL(30,6)) "
	rs, err := dao.runner.Exec(sql,
		amount.String(),
		accountNo,
		amount.String())
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

//UpdateStatus 账户状态更新
func (dao *AccountDao) UpdateStatus(
	accountNo string,
	status int) (rows int64, err error) {
	sql := "update account" +
		"set status=?" +
		"where account_no=?"
	rs, err := dao.runner.Exec(sql,
		status,
		accountNo)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()

}
