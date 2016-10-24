package fx

import (
	"github.com/reechou/real-fx/logic/models"
)

type fxAccountBackend interface {
	CreateFxAccount(fxAccount *models.FxAccount, fxAccountFollow *models.FxAccountFollow) error
	CreateSalesman(fxAccount *models.FxAccount) error
	UpdateFxAccountBaseInfo(fxAccount *models.FxAccount, fxAccountFollow *models.FxAccountFollow) error
	UpdateFxAccountStatus(fxAccount *models.FxAccount) error
	GetFxAccount(fxAccount *models.FxAccount) error
	CreateFxAccountFollow(info *models.FxAccountFollow) error
	UpdateFxAccountFollowStatus(info *models.FxAccountFollow) error
	GetFxAccountFollow(info *models.FxAccountFollow) error
}

type fxTeamBackend interface {
	CreateFxTeam(info *models.FxTeam) error
	CreateFxTeamMember(info *models.FxTeamMember) error
	GetFxTeamList() ([]models.FxTeam, error)
	GetFxTeamMembers(fxTeamID int64) ([]models.FxTeamMemberInfo, error)
}

type fxOrderBackend interface {
	CreateFxOrder(info *models.FxOrder) error
	GetFxOrderListCount(unionId string) (int64, error)
	GetFxOrderListCountById(accountId int64) (int64, error)
	GetFxOrderList(unionId string, offset, num, status int64) ([]models.FxOrder, error)
	GetFxOrderListById(accountId int64, offset, num, status int64) ([]models.FxOrder, error)
	GetFxOrderSettlementRecordListCount(unionId string) (int64, error)
	GetFxOrderSettlementRecordListCountById(accountId int64) (int64, error)
	GetFxOrderSettlementRecordList(unionId string, offset, num int64) ([]models.FxOrderSettlementRecord, error)
	GetFxOrderSettlementRecordListById(accountId int64, offset, num int64) ([]models.FxOrderSettlementRecord, error)
	GetFxOrderWaitSettlementRecordListCountById(accountId int64) (int64, error)
	GetFxOrderWaitSettlementRecordListById(accountId int64, offset, num int64) ([]models.FxOrderWaitSettlementRecord, error)
	GetFxOrderWaitSettlementRecordSum(accountId int64) (float32, error)
}

type fxWithdrawalBackend interface {
	CreateWithdrawalRecord(info *models.WithdrawalRecord) error
	GetWithdrawalRecordListCount(unionId string) (int64, error)
	GetWithdrawalRecordListCountById(accountId int64) (int64, error)
	GetWithdrawalRecordList(unionId string, offset, num int64) ([]models.WithdrawalRecord, error)
	GetWithdrawalRecordListById(accountId int64, offset, num int64) ([]models.WithdrawalRecord, error)
}

type Backend interface {
	fxAccountBackend
	fxTeamBackend
	fxOrderBackend
	fxWithdrawalBackend
}
