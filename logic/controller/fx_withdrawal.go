package controller

import (
	"fmt"
	
	"github.com/Sirupsen/logrus"
	"github.com/reechou/real-fx/logic/models"
)

func (daemon *Daemon) CreateWithdrawalRecord(info *models.WithdrawalRecord) error {
	fxAccount := &models.FxAccount{
		UnionId: info.UnionId,
	}
	has, err := models.GetFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("create withdrawal record get fx account error: %v", err)
		return err
	}
	if !has {
		logrus.Errorf("create withdrawal record get fx account error: no this account[%s]", info.UnionId)
		return fmt.Errorf("create withdrawal record get fx account error: no this account[%s]", info.UnionId)
	}
	info.AccountId = fxAccount.ID
	info.Status = WITHDRAWAL_WAITING
	err = models.CreateWithdrawalRecord(info)
	if err != nil {
		logrus.Errorf("create withdrawal record error: %v", err)
		return err
	}
	return nil
}

func (daemon *Daemon) GetWithdrawalRecordListCount(unionId string) (int64, error) {
	return models.GetWithdrawalRecordListCount(unionId)
}

func (daemon *Daemon) GetWithdrawalRecordListCountById(accountId int64) (int64, error) {
	return models.GetWithdrawalRecordListCountById(accountId)
}

func (daemon *Daemon) GetWithdrawalRecordList(unionId string, offset, num int64) ([]models.WithdrawalRecord, error) {
	list, err := models.GetWithdrawalRecordList(unionId, offset, num)
	if err != nil {
		logrus.Errorf("get withdrawal record list error: %v", err)
		return nil, err
	}
	return list, nil
}

func (daemon *Daemon) GetWithdrawalRecordListById(accountId int64, offset, num int64) ([]models.WithdrawalRecord, error) {
	list, err := models.GetWithdrawalRecordListById(accountId, offset, num)
	if err != nil {
		logrus.Errorf("get withdrawal record list error: %v", err)
		return nil, err
	}
	return list, nil
}
