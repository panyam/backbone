package services

import (
	. "github.com/panyam/relay/services/msg/core"
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestCreateUserService(c *C) {
	svc := s.serviceGroup.UserService
	c.Assert(svc, Not(IsNil))
}

func (s *TestSuite) TestSaveUserEmptyId_ShouldCreateId(c *C) {
	svc := s.serviceGroup.UserService
	team, err := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := NewUser(0, "user1", team)
	err = svc.SaveUser(user, false)
	c.Assert(err, IsNil)
	c.Assert(user.Id, Not(Equals), 0)
	c.Assert(user.Username, Equals, "user1")

	fetched_user, err := svc.GetUserById(user.Id)
	c.Assert(fetched_user.Id, Equals, user.Id)
	c.Assert(fetched_user.Team, Not(IsNil))
	c.Assert(fetched_user.Team.Id, Equals, user.Team.Id)

	fetched_user, _ = svc.GetUser("user1", team)
	c.Assert(fetched_user.Id, Equals, user.Id)
	c.Assert(fetched_user.Team, Not(IsNil))
	c.Assert(fetched_user.Team.Id, Equals, user.Team.Id)
}

func (s *TestSuite) TestGetUserFirstTime(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user, _ := svc.GetUserById(1)
	c.Assert(user, IsNil)
	_, err := svc.GetUser("user1", team)
	c.Assert(err, Not(IsNil))
}

func (s *TestSuite) TestSaveUserNormal(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := NewUser(0, "user1", team)
	user.Object = Object{Id: 1}
	err := svc.SaveUser(user, false)
	c.Assert(err, IsNil)
}

func (s *TestSuite) TestCreateUserDuplicate(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.CreateTeam(1, "org", "team")
	user := User{Username: "user1", Team: team}

	err := svc.SaveUser(&user, false)
	c.Assert(err, IsNil)

	// err = svc.SaveUser(&user, false)
	// c.Assert(err, Not(IsNil))
}

/*
func (s *TestSuite) TestSaveUserIdOrNameExists(c *C) {
	svc := s.serviceGroup.UserService
	user, err := svc.CreateUser("1", "user1")
	c.Assert(err, IsNil)
	user, err = svc.CreateUser("1", "user2")
	c.Assert(err, Not(IsNil))
	c.Assert(user, IsNil)
	user, err = svc.CreateUser("2", "user1")
	c.Assert(err, Not(IsNil))
	c.Assert(user, IsNil)
}
*/
