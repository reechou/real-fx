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

func (fxr *FXRouter) createFxAccount(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &CreateFxAccountReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId:  req.UnionId,
		Superior: req.Superior,
		Name:     req.Name,
	}
	fxAccountFollow := &models.FxAccountFollow{
		UnionId:   req.UnionId,
		WXAccount: req.WXAccount,
		OpenId:    req.OpenId,
	}
	if err := fxr.backend.CreateFxAccount(fxAccount, fxAccountFollow); err != nil {
		logrus.Errorf("Error create fx account: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error create fx account: %v", err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) createFxSalesman(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &CreateSalesmanReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
		Ticket:  req.Ticket,
		Phone:   req.Phone,
	}
	if err := fxr.backend.CreateSalesman(fxAccount); err != nil {
		logrus.Errorf("Error create fx salesman: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error create fx salesman: %v", err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) updateFxBaseInfo(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &updateFxBaseInfoReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		Name:  req.Name,
		Phone: req.Phone,
	}
	fxAccountFollow := &models.FxAccountFollow{
		WXAccount: req.WXAccount,
		OpenId:    req.OpenId,
	}
	if err := fxr.backend.UpdateFxAccountBaseInfo(fxAccount, fxAccountFollow); err != nil {
		logrus.Errorf("Error update fx base info: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error update fx base info: %v", err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) updateFxStatus(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &updateFxStatusReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
		Status:  req.Status,
	}
	if err := fxr.backend.UpdateFxAccountStatus(fxAccount); err != nil {
		logrus.Errorf("Error update fx status: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error update fx status: %v", err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxAccount(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxAccountReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccount := &models.FxAccount{
		UnionId: req.UnionId,
	}
	if err := fxr.backend.GetFxAccount(fxAccount); err != nil {
		logrus.Errorf("Error get fx account: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx account: %v", err)
	} else {
		rsp.Data = fxAccount
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxAccountFollow(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxAccountFollowReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxAccountFollow := &models.FxAccountFollow{
		WXAccount: req.WXAccount,
		OpenId:    req.OpenId,
	}
	if err := fxr.backend.GetFxAccountFollow(fxAccountFollow); err != nil {
		logrus.Errorf("Error get fx account follow: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx account follow: %v", err)
	} else {
		rsp.Data = fxAccountFollow
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}
