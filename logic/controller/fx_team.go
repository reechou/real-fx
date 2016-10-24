package controller

import (
	"github.com/reechou/real-fx/logic/models"
)

func (daemon *Daemon) CreateFxTeam(info *models.FxTeam) error {
	return models.CreateFxTeam(info)
}

func (daemon *Daemon) CreateFxTeamMember(info *models.FxTeamMember) error {
	return models.CreateFxTeamMember(info)
}

func (daemon *Daemon) GetFxTeamList() ([]models.FxTeam, error) {
	return models.GetFxTeamList()
}

func (daemon *Daemon) GetFxTeamMembers(fxTeamID int64) ([]models.FxTeamMemberInfo, error) {
	return models.GetFxTeamMembers(fxTeamID)
}
