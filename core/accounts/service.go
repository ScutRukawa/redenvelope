package accounts

import (
	"errors"
	"redenvelope/services"
	"sync"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"

	"redenvelope/infra/base"
)

//var _ services.AccountService =

var once sync.Once

func init() {
	once.Do(func() {
		services.IAccountService = new(AccountService)
	})
}

type AccountService struct {
}

func (a *AccountService) CreateAccount(
	dto services.AccountCreatedDTO) (*services.AccountDTO, error) {
	domain := AccountDomain{}
	//验证输入参数
	err := base.Validate().Struct(&dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				logrus.Error(e.Translate(base.Transtate()))
			}
		}
		return nil, err
	}
	//创建账户
	amount, err := decimal.NewFromString(dto.Amount)
	if err != nil {
		return nil, err
	}
	account := services.AccountDTO{
		UserID:       dto.UserID,
		Username:     dto.Username,
		AccountType:  dto.AccountType,
		AccountName:  dto.AccountName,
		CurrencyCode: dto.CurrencyCode,
		Status:       1,
		Balance:      amount,
	}
	rdto, err := domain.Create(account)
	return rdto, err
}
func (a *AccountService) Transfer(
	dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	domain := AccountDomain{}
	//验证输入参数
	err := base.Validate().Struct(&dto)
	if err != nil {
		_, ok := err.(*validator.InvalidValidationError)
		if ok {
			logrus.Error("验证错误", err)
		}
		errs, ok := err.(validator.ValidationErrors)
		if ok {
			for _, e := range errs {
				logrus.Error(e.Translate(base.Transtate()))
			}
		}
		return services.TransferedStatusFailure, err
	}
	//执行转账逻辑
	amount, err := decimal.NewFromString(dto.AmountStr)
	if err != nil {
		return services.TransferedStatusFailure, err
	}
	dto.Amount = amount
	if dto.ChangeFlag == services.FlagTransferOut {
		if dto.ChangeType > 0 {
			return services.TransferedStatusFailure,
				errors.New("如果changeflag 为支出，changeType必须为<0")
		}
	} else {
		if dto.ChangeType < 0 {
			return services.TransferedStatusFailure,
				errors.New("如果changeflag 为支出，changeType必须为>0")

		}
	}

	status, err := domain.Transfer(dto)
	return status, err
}

func (a *AccountService) StoreValue(dto services.AccountTransferDTO) (services.TransferedStatus, error) {
	dto.TradeTarget = dto.TradeBody
	dto.ChangeFlag = services.FlagTransferIn
	dto.ChangeType = services.AccountStoreValue
	return a.Transfer(dto)
}

func (a *AccountService) GetEnvelopeAccountByUserID(userID string) *services.AccountDTO {
	domain := AccountDomain{}
	account := domain.GetEnvelopeAccountByUserID(userID)
	return account
}

func (a *AccountService) GetAccount(accountNo string) *services.AccountDTO {
	domain := AccountDomain{}
	return domain.GetAccount(accountNo)
}
