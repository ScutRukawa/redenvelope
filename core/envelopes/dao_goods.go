package envelopes

import (
	"redenvelope/services"
	"time"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

//RedEnvelopeGoodsDao x
type RedEnvelopeGoodsDao struct {
	runner *dbx.TxRunner
}

//Insert 红包
func (dao *RedEnvelopeGoodsDao) Insert(po *RedEnvelopeGoods) (int64, error) {
	rs, err := dao.runner.Insert(po)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}

//UpdateBalance 更新余额 数量，
func (dao *RedEnvelopeGoodsDao) UpdateBalance(
	envelopeNo string, amount decimal.Decimal) (int64, error) {
	sql := "update red_envelope_goods set remain_quantity=remain_quantity-1, " +
		"remain_amount=remain_amount-CAST(? AS DECIMAL(30,6)) " +
		"where envelope_no=? " +
		"and remain_quantity>0 " +
		"and remain_amount>=CAST(? AS DECIMAL(30,6))"
	rs, err := dao.runner.Exec(sql, amount.String(), envelopeNo, amount.String())
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

//UpdateOrderStatus 更新订单状态
func (dao *RedEnvelopeGoodsDao) UpdateOrderStatus(
	envelopeNo string, status services.OrderStatus) (int64, error) {
	sql := "update red_envelope_goods " +
		"set order_status=? " +
		" where envelopeNo=? "
	rs, err := dao.runner.Exec(sql, status, envelopeNo)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

//GetOne 查询
func (dao *RedEnvelopeGoodsDao) GetOne(envelopeNo string) *RedEnvelopeGoods {
	po := &RedEnvelopeGoods{EnvelopeNo: envelopeNo}
	ok, err := dao.runner.GetOne(po)
	if err != nil || !ok {
		logrus.Error(err)
		return po
	}
	return po
}

//过期，分页查询 limit offset size
func (dao *RedEnvelopeGoodsDao) FindExpired(
	offset, size int) []RedEnvelopeGoods {
	var goods []RedEnvelopeGoods
	now := time.Now()
	sql := "select * from red_envelope_goods where expired_at>? limit ?,?"
	err := dao.runner.Find(&goods, sql, now, offset, size)
	//dao.runner.Query()
	if err != nil {
		logrus.Error(err)
	}
	return goods
}
