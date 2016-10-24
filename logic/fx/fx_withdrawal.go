package fx

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Sirupsen/logrus"
	"github.com/reechou/real-fx/logic/models"
	"github.com/reechou/real-fx/utils"
	"golang.org/x/net/context"
)

func (fxr *FXRouter) createFxWithdrawalRecord(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &withdrawalMoneyReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}
	
	wInfo := &models.WithdrawalRecord{
		UnionId: req.UnionId,
		WithdrawalMoney: req.Money,
	}
	err := fxr.backend.CreateWithdrawalRecord(wInfo)
	if err != nil {
		logrus.Errorf("create withdrawal record[%v] error: %v", wInfo, err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("create withdrawal record[%v] error: %v", wInfo, err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxWithdrawalRecordList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getWithdrawalListReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	type FxWithdrawalRecordList struct {
		Count int64                     `json:"count"`
		List  []models.WithdrawalRecord `json:"list"`
	}
	count, err := fxr.backend.GetWithdrawalRecordListCount(req.UnionId)
	if err != nil {
		logrus.Errorf("Error get fx withdrawal record list count: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx withdrawal record list count: %v", err)
	} else {
		list, err := fxr.backend.GetWithdrawalRecordList(req.UnionId, req.Offset, req.Num)
		if err != nil {
			logrus.Errorf("Error get fx withdrawal record list: %v", err)
			rsp.Code = RspCodeErr
			rsp.Msg = fmt.Sprintf("Error get fx withdrawal record list: %v", err)
		} else {
			var listInfo FxWithdrawalRecordList
			listInfo.Count = count
			listInfo.List = list
			rsp.Data = listInfo
		}
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}
