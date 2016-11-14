package fx

import (
	"github.com/reechou/real-fx/router"
)

type FXRouter struct {
	backend Backend
	routes  []router.Route
}

func NewRouter(b Backend) router.Router {
	r := &FXRouter{
		backend: b,
	}
	r.initRoutes()
	return r
}

func (fxr *FXRouter) Routes() []router.Route {
	return fxr.routes
}

func (fxr *FXRouter) initRoutes() {
	fxr.routes = []router.Route{
		// about fx account
		router.NewPostRoute("/fx/create_fx_account", fxr.createFxAccount),
		router.NewPostRoute("/fx/create_fx_salesman", fxr.createFxSalesman),
		router.NewPostRoute("/fx/update_fx_baseinfo", fxr.updateFxBaseInfo),
		router.NewPostRoute("/fx/update_fx_status", fxr.updateFxStatus),
		router.NewPostRoute("/fx/fx_sign", fxr.updateFxSignTime),
		router.NewPostRoute("/fx/get_fx_accout", fxr.getFxAccount),
		router.NewPostRoute("/fx/get_fx_accout_unionid", fxr.getFxAccountFollow),
		router.NewPostRoute("/fx/get_fx_lower_people_list", fxr.getFxLowerPeopleList),
		// about fx account history
		router.NewPostRoute("/fx/get_fx_history", fxr.getFxAccountHistoryList),
		router.NewPostRoute("/fx/get_fx_history_by_type", fxr.getFxAccountHistoryListByType),
		// about fx team
		router.NewPostRoute("/fx/create_fx_team", fxr.createFxTeam),
		router.NewPostRoute("/fx/create_fx_team_member", fxr.createFxTeamMember),
		router.NewPostRoute("/fx/get_fx_team_list", fxr.getFxTeamList),
		router.NewPostRoute("/fx/get_fx_team_members", fxr.getFxTeamMembers),
		// about fx order
		router.NewPostRoute("/fx/create_fx_order", fxr.createFxOrder),
		router.NewPostRoute("/fx/get_fx_order_list", fxr.getFxOrderList),
		router.NewPostRoute("/fx/get_fx_order_wait_sr_sum", fxr.getFxOrderWaitSettlementSum),
		router.NewPostRoute("/fx/get_fx_order_sr_list", fxr.getFxOrderSettlementRecordList),
		router.NewPostRoute("/fx/get_fx_order_wait_sr_list", fxr.getFxOrderWaitSettlementRecordList),
		// about withdrawal
		router.NewPostRoute("/fx/create_fx_withdrawal_record", fxr.createFxWithdrawalRecord),
		router.NewPostRoute("/fx/get_fx_withdrawal_record", fxr.getFxWithdrawalRecordList),
	}
}
