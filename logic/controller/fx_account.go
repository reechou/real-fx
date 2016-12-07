package controller

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/real-fx/logic/ext"
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
	if fxAccount.Superior != "" && fxAccount.Superior != GodSalesman {
		superFxAccount := &models.FxAccount{
			UnionId: fxAccount.Superior,
		}
		err := models.AddFxAccountMoney(float32(daemon.cfg.Score.FollowScore), superFxAccount)
		if err != nil {
			logrus.Errorf("add super fx account money error: %v", err)
			return err
		}
		h := models.FxAccountHistory{
			UnionId:    fxAccount.Superior,
			Score:      float32(daemon.cfg.Score.FollowScore),
			ChangeType: int64(FX_HISTORY_TYPE_INVITE),
			ChangeDesc: FxHistoryDescs[FX_HISTORY_TYPE_INVITE],
			CreatedAt:  time.Now().Unix(),
		}
		models.CreateFxAccountHistoryList([]models.FxAccountHistory{h})

		// send msg to wechat
		has, err := models.GetFxAccount(superFxAccount)
		if err != nil {
			logrus.Errorf("get super fx account info error: %v", err)
		} else {
			if has {
				superFxAccountFollow := &models.FxAccountFollow{
					UnionId:   superFxAccount.UnionId,
					WXAccount: WX_WGLS_ACCOUNT,
				}
				err = models.GetFxAccountFollowFromUnionId(superFxAccountFollow)
				if err == nil {
					wxSendInfo := &ext.WeixinMsgSendReq{
						OpenId:    superFxAccountFollow.OpenId,
						Score:     float32(daemon.cfg.Score.FollowScore),
						LeftScore: superFxAccount.CanWithdrawals,
						Reason:    FxHistoryDescs[FX_HISTORY_TYPE_INVITE],
						UserName:  superFxAccount.Name,
					}
					daemon.we.AsyncWxSendMsg(wxSendInfo)
				}
			}
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

func (daemon *Daemon) UpdateFxAccountSignTime(fxAccount *models.FxAccount) (int64, error) {
	affected, err := models.UpdateFxAccountSignTime(float32(daemon.cfg.Score.SignScore), fxAccount)
	if err != nil {
		return 0, err
	}
	if affected > 0 {
		h := models.FxAccountHistory{
			UnionId:    fxAccount.UnionId,
			Score:      float32(daemon.cfg.Score.SignScore),
			ChangeType: int64(FX_HISTORY_TYPE_SIGN),
			ChangeDesc: FxHistoryDescs[FX_HISTORY_TYPE_SIGN],
			CreatedAt:  time.Now().Unix(),
		}
		models.CreateFxAccountHistoryList([]models.FxAccountHistory{h})

		// send msg to wechat
		has, err := models.GetFxAccount(fxAccount)
		if err != nil {
			logrus.Errorf("get fx account info error: %v", err)
		} else {
			if has {
				fxAccountFollow := &models.FxAccountFollow{
					UnionId:   fxAccount.UnionId,
					WXAccount: WX_WGLS_ACCOUNT,
				}
				err = models.GetFxAccountFollowFromUnionId(fxAccountFollow)
				if err == nil {
					wxSendInfo := &ext.WeixinMsgSendReq{
						OpenId:    fxAccountFollow.OpenId,
						Score:     float32(daemon.cfg.Score.SignScore),
						LeftScore: fxAccount.CanWithdrawals,
						Reason:    FxHistoryDescs[FX_HISTORY_TYPE_SIGN],
						UserName:  fxAccount.Name,
					}
					daemon.we.AsyncWxSendMsg(wxSendInfo)
				}
			}
		}
	}

	return affected, nil

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

func (daemon *Daemon) GetFxAccountRank(offset, num int64) ([]models.FxAccount, error) {
	return models.GetFxAccountRank(offset, num)
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
