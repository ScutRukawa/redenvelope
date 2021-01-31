package accounts

import (
	"errors"
	"redenvelope/infra/base"
	"redenvelope/services"

	"github.com/shopspring/decimal"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"github.com/tietang/dbx"
)

type AccountDomain struct {
	account    Account
	accountLog AccountLog
}

func (domain *AccountDomain) GetAccount(accountNo string) *services.AccountDTO {
	accountDao := AccountDao{}
	var account *Account

	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetOne(accountNo)
		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()
}

//创建流水的记录
func (domain *AccountDomain) createAccountLogNo() {
	//全局唯一ID ，后面换成分布式ID
	domain.accountLog.LogNo = ksuid.New().Next().String()
}

//
func (domain *AccountDomain) createAccountNo() {
	domain.account.AccountNo = ksuid.New().Next().String()
}

func (domain *AccountDomain) createAccountLog() {
	//先创建账户逻辑
	domain.accountLog = AccountLog{}
	domain.createAccountLogNo()
	domain.accountLog.TradeNo = domain.accountLog.LogNo
	//流水中的交易主体信息
	domain.accountLog.AccountNo = domain.account.AccountNo
	domain.accountLog.UserID = domain.account.UserID
	domain.accountLog.Username = domain.account.Username.String
	//交易对象
	domain.accountLog.TargetAccountNo = domain.account.AccountNo
	domain.accountLog.TargetUserID = domain.account.UserID
	domain.accountLog.TargetUsername = domain.account.Username.String
	//交易金额
	domain.accountLog.Amount = domain.account.Balance //初始化值，第一次创建时
	domain.accountLog.Balance = domain.account.Balance
	//交易变化属性
	domain.accountLog.Decs = "创建账户"
	domain.accountLog.ChangeType = services.AccountCreated
	domain.accountLog.ChangeFlag = services.FlagAccountCreated
}

func (domain *AccountDomain) Create(
	dto services.AccountDTO) (*services.AccountDTO, error) {
	//创建账户持久化对象
	domain.account = Account{}
	domain.account.FromDTO(&dto)
	domain.createAccountNo()
	domain.account.Username.Valid = true

	//流水持久化对象
	domain.createAccountLog()
	accountDao := AccountDao{}
	accountLogDao := AccountLogDao{}

	var rdto *services.AccountDTO
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		id, err := accountDao.Insert(&domain.account)
		if err != nil {
			logrus.Error(err)
			return err
		}
		if id <= 0 {
			return errors.New("账户创建失败")
		}
		accountLogDao.runner = runner
		_, err = accountLogDao.Insert(&domain.accountLog)
		if err != nil {
			logrus.Error(err)
			return err
		}
		if id <= 0 {
			return errors.New("账户创建失败")
		}
		domain.account = *accountDao.GetOne(domain.account.AccountNo)
		return nil
	})
	rdto = domain.account.ToDTO()
	return rdto, err
}

//Transfer x
func (domain *AccountDomain) Transfer(
	dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	//支出类型
	amount := dto.Amount
	if dto.ChangeFlag == services.FlagTransferOut {
		amount = dto.Amount.Mul(decimal.NewFromFloat(-1))
	}
	//创建流水
	var status services.TransferedStatus
	domain.account = Account{}
	domain.accountLog.FromTransferDTO(&dto)
	domain.createAccountLogNo()
	//检查余额是否足够和更新余额，乐观锁
	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao := AccountDao{runner: runner}
		accountLogDao := AccountLogDao{runner: runner}
		rows, err := accountDao.UpdateBalance(domain.accountLog.AccountNo, amount)
		if err != nil {
			status = services.TransferedStatusFailure
			return err
		}
		if rows <= 0 && dto.ChangeFlag == services.FlagTransferOut {
			status = services.TransferedStatusSufficientFunds
			return errors.New("余额不足")
		}
		account := accountDao.GetOne(dto.TradeBody.AccountNo)
		if account == nil {
			return errors.New("账户出错")
		}
		domain.account = *account
		domain.accountLog.Balance = domain.account.Balance
		id, err := accountLogDao.Insert(&domain.accountLog)
		if err != nil || id <= 0 {
			status = services.TransferedStatusFailure
			return errors.New("账户流水创建失败")
		}
		return nil
	})
	if err != nil {
		logrus.Error(err)
		return status, err
	}
	status = services.TransferedStatusSuccess
	return status, nil
}
func (domain *AccountDomain) GetEnvelopeAccountByUserID(userID string) *services.AccountDTO {
	accountDao := AccountDao{}
	var account *Account

	err := base.Tx(func(runner *dbx.TxRunner) error {
		accountDao.runner = runner
		account = accountDao.GetByUserID(userID, int(services.EnvelopeAccountType))
		return nil
	})
	if err != nil {
		return nil
	}
	if account == nil {
		return nil
	}
	return account.ToDTO()

}
