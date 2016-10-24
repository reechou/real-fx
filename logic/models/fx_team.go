package models

import (
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
)

type FxTeam struct {
	ID        int64  `xorm:"pk autoincr"`
	Name      string `xorm:"not null default '' varchar(128) unique"`
	CreatedAt int64  `xorm:"not null default 0 int"`
	UpdatedAt int64  `xorm:"not null default 0 int"`
}

type FxTeamMember struct {
	ID        int64  `xorm:"pk autoincr"`
	TeamId    int64  `xorm:"not null default 0 int unique(uni_fx_team)"`
	UnionId   string `xorm:"not null default '' varchar(128) unique(uni_fx_team)"`
	CreatedAt int64  `xorm:"not null default 0 int"`
	UpdatedAt int64  `xorm:"not null default 0 int"`
}

type FxTeamMemberInfo struct {
	FxTeamMember `xorm:"extends"`
	FxAccount    `xorm:"extends"`
}

func CreateFxTeam(info *FxTeam) (err error) {
	if info.Name == "" {
		return fmt.Errorf("fx team[%s] cannot be nil.", info.Name)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	_, err = x.Insert(info)

	if err == nil {
		logrus.Infof("create fx team[%s] success.", info.Name)
	}

	return
}

func CreateFxTeamMember(info *FxTeamMember) error {
	if info.TeamId == 0 || info.UnionId == "" {
		return fmt.Errorf("fx team_id[%d] union_id[%s] cannot be nil", info.TeamId, info.UnionId)
	}

	now := time.Now().Unix()
	info.CreatedAt = now
	info.UpdatedAt = now
	_, err := x.Insert(info)
	if err != nil {
		logrus.Errorf("create fx team error: %v", err)
		return err
	}
	logrus.Infof("create fx team[%d] member[%s] success.", info.TeamId, info.UnionId)

	return nil
}

func GetFxTeamList() ([]FxTeam, error) {
	var teams []FxTeam
	err := x.Find(&teams)
	if err != nil {
		logrus.Errorf("find all fx team error: %v", err)
		return nil, err
	}
	return teams, nil
}

func GetFxTeamMembers(fxTeamId int64) ([]FxTeamMemberInfo, error) {
	var teamMembers []FxTeamMemberInfo
	err := x.Table("fx_team_member").Select("fx_team_member.*, fx_account.*").
		Join("INNER", "fx_account", "fx_account.union_id = fx_team_member.union_id").
		Where("fx_team_member.team_id = ?", fxTeamId).Find(&teamMembers)
	if err != nil {
		logrus.Errorf("find fx team[%d] members error: %v", fxTeamId, err)
		return nil, err
	}
	return teamMembers, nil
}
