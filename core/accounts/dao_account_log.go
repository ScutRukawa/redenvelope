package accounts

import (
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

//AccountLogDao x
type AccountLogDao struct {
	runner *dbx.TxRunner
}

//NewAccountLogDao x
func NewAccountLogDao(runner *dbx.TxRunner) AccountLogDao {
	return AccountLogDao{runner: runner}
}

//GetOne x
func (logdao *AccountLogDao) GetOne(logNo string) *AccountLog {
	a := &AccountLog{LogNo: logNo}
	ok, err := logdao.runner.GetOne(a) //必须传入唯一索引字段
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

//GetByTradeNo x
func (logdao *AccountLogDao) GetByTradeNo(tradeNo string) *AccountLog {
	sql := "select * from account_log " + "where trade_no=? "
	a := &AccountLog{}
	ok, err := logdao.runner.Get(a, sql, tradeNo) //必须传入唯一索引字段
	if err != nil {
		logrus.Error(err)
		return nil
	}
	if !ok {
		return nil
	}
	return a
}

//Insert x
func (logdao *AccountLogDao) Insert(accountLog *AccountLog) (id int64, err error) {

	rs, err := logdao.runner.Insert(accountLog)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}
