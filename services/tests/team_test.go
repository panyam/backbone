package services

import (
	. "github.com/panyam/relay/services/core"
	. "gopkg.in/check.v1"
	"log"
	// "code.google.com/p/gomock/gomock"
	// "time"
)

func (s *TestSuite) TestCreateTeam(c *C) {
	svc := s.serviceGroup.TeamService
	team, err := svc.CreateTeam(0, "group", "test")
	c.Assert(err, Equals, nil)
	c.Assert(team, Not(Equals), (*Team)(nil))
	c.Assert(team.Name, Equals, "test")
	team, err = svc.GetTeamByName("group", "test")
	c.Assert(err, Equals, nil)
	c.Assert(team, Not(Equals), (*Team)(nil))
	c.Assert(team.Name, Equals, "test")
}

func (s *TestSuite) TestCreateTeamExistsByName(c *C) {
	svc := s.serviceGroup.TeamService
	_, err := svc.CreateTeam(1, "group", "test")
	team, err := svc.CreateTeam(1, "group2", "test2")
	c.Assert(err, Not(Equals), nil)
	c.Assert(team, Equals, (*Team)(nil))
}

func (s *TestSuite) TestDeleteTeam(c *C) {
	svc := s.serviceGroup.TeamService
	team, err := svc.CreateTeam(0, "group", "test")
	c.Assert(err, Equals, nil)
	c.Assert(team, Not(Equals), (*Team)(nil))
	c.Assert(team.Id, Not(Equals), 0)

	log.Println("TeamID: ", team.Id)
	err = svc.DeleteTeam(team)
	c.Assert(err, Equals, nil)

	team, err = svc.GetTeamByName("group", "test")
	c.Assert(team, Equals, (*Team)(nil))
	c.Assert(err, Not(Equals), nil)
}

/*
func (s *TestSuite) TestJoinTeam(c *C) {
	svc := s.serviceGroup.TeamService
	team, _ := svc.CreateTeam(0, "group", "test")
	svc.JoinTeam(team, "user1")
	c.Assert(svc.TeamContains(team, "user1"), Equals, true)
}

func (s *TestSuite) TestLeaveTeam(c *C) {
	svc := s.serviceGroup.TeamService
	team, _ := svc.CreateTeam("", "group", "test")
	user := NewUser("1", "user1")
	svc.JoinTeam(team, "user1")
	c.Assert(svc.TeamContains(team, "user1"), Equals, true)
	svc.LeaveTeam(team, user)
	c.Assert(svc.TeamContains(team, "user1"), Equals, false)
}
*/
