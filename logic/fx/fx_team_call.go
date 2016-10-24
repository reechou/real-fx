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

func (fxr *FXRouter) createFxTeam(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &createFxTeamReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxTeam := &models.FxTeam{
		Name: req.Name,
	}
	if err := fxr.backend.CreateFxTeam(fxTeam); err != nil {
		logrus.Errorf("Error create fx team: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error create fx team: %v", err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) createFxTeamMember(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &createFxTeamMemberReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	fxTeam := &models.FxTeamMember{
		TeamId:  req.TeamId,
		UnionId: req.UnionId,
	}
	if err := fxr.backend.CreateFxTeamMember(fxTeam); err != nil {
		logrus.Errorf("Error create fx team member: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error create fx team member: %v", err)
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxTeamList(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	list, err := fxr.backend.GetFxTeamList()
	if err != nil {
		logrus.Errorf("Error get fx team list: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx team list: %v", err)
	} else {
		rsp.Data = list
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}

func (fxr *FXRouter) getFxTeamMembers(ctx context.Context, w http.ResponseWriter, r *http.Request, vars map[string]string) error {
	if err := utils.ParseForm(r); err != nil {
		return err
	}

	req := &getFxTeamMembersReq{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	rsp := &FxResponse{Code: RspCodeOK}

	list, err := fxr.backend.GetFxTeamMembers(req.FxTeamId)
	if err != nil {
		logrus.Errorf("Error get fx team members: %v", err)
		rsp.Code = RspCodeErr
		rsp.Msg = fmt.Sprintf("Error get fx team members: %v", err)
	} else {
		rsp.Data = list
	}

	return utils.WriteJSON(w, http.StatusOK, rsp)
}
