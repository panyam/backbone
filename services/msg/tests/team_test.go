package services

import (
	msgcore "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
	"log"
	// "code.google.com/p/gomock/gomock"
	// "time"
)

func (s *TestSuite) TestSaveTeam(c *C) {
	svc := s.serviceGroup.TeamService
	team, err := svc.SaveTeam(msgcore.NewTeam(0, "group", "test"))
	c.Assert(err, IsNil)
	c.Assert(team, Not(IsNil))
	c.Assert(team.Name, Equals, "test")
	team, err = svc.GetTeam(msgcore.NewTeam(0, "group", "test"))
	c.Assert(err, IsNil)
	c.Assert(team, Not(IsNil))
	c.Assert(team.Name, Equals, "test")
}

func (s *TestSuite) TestSaveTeamExistsByName(c *C) {
	svc := s.serviceGroup.TeamService
	_, err := svc.SaveTeam(msgcore.NewTeam(1, "group", "test"))
	team, err := svc.SaveTeam(msgcore.NewTeam(1, "group2", "test2"))
	c.Assert(err, Not(IsNil))
	c.Assert(team, IsNil)
}

func (s *TestSuite) TestDeleteTeam(c *C) {
	svc := s.serviceGroup.TeamService
	team, err := svc.SaveTeam(msgcore.NewTeam(0, "group", "test"))
	c.Assert(err, IsNil)
	c.Assert(team, Not(IsNil))
	c.Assert(team.Id, Not(Equals), 0)

	log.Println("TeamID: ", team.Id)
	err = svc.DeleteTeam(team)
	c.Assert(err, IsNil)

	team, err = svc.GetTeam(msgcore.NewTeam(0, "group", "test"))
	c.Assert(team, IsNil)
	c.Assert(err, Not(IsNil))
}

/*
func (s *TestSuite) TestJoinTeam(c *C) {
	svc := s.serviceGroup.TeamService
	team, _ := svc.SaveTeam(msgcore.NewTeam(0, "group", "test"))
	svc.JoinTeam(team, "user1")
	c.Assert(svc.TeamContains(team, "user1"), Equals, true)
}

func (s *TestSuite) TestLeaveTeam(c *C) {
	svc := s.serviceGroup.TeamService
	team, _ := svc.SaveTeam(msgcore.NewTeam("", "group", "test"))
	user := NewUser("1", "user1")
	svc.JoinTeam(team, "user1")
	c.Assert(svc.TeamContains(team, "user1"), Equals, true)
	svc.LeaveTeam(team, user)
	c.Assert(svc.TeamContains(team, "user1"), Equals, false)
}
*/
