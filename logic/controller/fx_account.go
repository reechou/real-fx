package controller

import (
	"fmt"
	
	"github.com/Sirupsen/logrus"
	"github.com/reechou/real-fx/logic/models"
)

func (daemon *Daemon) CreateFxAccount(fxAccount *models.FxAccount, fxAccountFollow *models.FxAccountFollow) error {
	if fxAccount.Superior == "" {
		fxAccount.Superior = GodSalesman
	}

	if err := models.CreateFxAccount(fxAccount); err != nil {
		logrus.Errorf("create fx account error: %v", err)
		return err
	}
	if fxAccount.Superior != "" {
		superFxAccount := &models.FxAccount{
			UnionId: fxAccount.Superior,
		}
		err := models.AddFxAccountMoney(float32(daemon.cfg.Score.FollowScore), superFxAccount)
		if err != nil {
			logrus.Errorf("add super fx account money error: %v", err)
			return err
		}
	}
	if err := models.CreateFxAccountFollow(fxAccountFollow); err != nil {
		logrus.Errorf("create fx account follow error: %v", err)
		return err
	}

	return nil
}

func (daemon *Daemon) CreateSalesman(fxAccount *models.FxAccount) error {
	return models.UpdateFxAccountSalesman(fxAccount)
}

func (daemon *Daemon) UpdateFxAccountBaseInfo(fxAccount *models.FxAccount, fxAccountFollow *models.FxAccountFollow) error {
	if err := models.GetFxAccountFollow(fxAccountFollow); err != nil {
		logrus.Errorf("get fx account follow error: %v", err)
		return err
	}
	fxAccount.UnionId = fxAccountFollow.UnionId

	return models.UpdateFxAccountBaseInfo(fxAccount)
}

func (daemon *Daemon) UpdateFxAccountStatus(fxAccount *models.FxAccount) error {

	return models.UpdateFxAccountStatus(fxAccount)
}

func (daemon *Daemon) GetFxAccount(fxAccount *models.FxAccount) error {
	has, err := models.GetFxAccount(fxAccount)
	if err != nil {
		logrus.Errorf("get fx account error: %v", err)
		return err
	}
	if !has {
		return fmt.Errorf("no this account.")
	}

	return nil
}

func (daemon *Daemon) GetLowerPeopleCount(unionId string) (int64, error) {
	return models.GetLowerPeopleCount(unionId)
}

func (daemon *Daemon) GetLowerPeopleList(unionId string, offset, num int64) ([]models.FxAccount, error) {
	return models.GetLowerPeople(unionId, offset, num)
}

func (daemon *Daemon) CreateFxAccountFollow(info *models.FxAccountFollow) error {
	return models.CreateFxAccountFollow(info)
}

func (daemon *Daemon) UpdateFxAccountFollowStatus(info *models.FxAccountFollow) error {
	return models.UpdateFxAccountFollowStatus(info)
}

func (daemon *Daemon) GetFxAccountFollow(info *models.FxAccountFollow) error {
	return models.GetFxAccountFollow(info)
}
