package services

/*
import (
	. "github.com/panyam/backbone/services/core"
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestSaveUserEmptyId_ShouldCreateId(c *C) {
	svc := s.serviceGroup.UserService
	team, _ := s.serviceGroup.TeamService.CreateTeam("1", "org", "team")
	user := User{Username: "user1", Team: team}
	err := svc.SaveUser(&user, false)
	c.Assert(err, Not(Equals), nil)
	c.Assert(user.Id, Not(Equals), "")
	c.Assert(user.Username, Equals, "user1")

	fetched_user, _ := svc.GetUserById("1")
	c.Assert(fetched_user, Equals, user)

	fetched_user, _ = svc.GetUser("user1", team)
	c.Assert(fetched_user, Equals, user)
}

func (s *TestSuite) TestCreateUserService(c *C) {
	svc := s.serviceGroup.UserService
	c.Assert(svc, Not(Equals), nil)
}

func (s *TestSuite) TestGetUserFirstTime(c *C) {
	svc := s.serviceGroup.UserService
	user, _ := svc.GetUserById("1")
	c.Assert(user, Equals, (*User)(nil))
	user, _ = svc.GetUser("user1", nil)
	c.Assert(user, Equals, (*User)(nil))
}
func (s *TestSuite) TestCreateUserDuplicate(c *C) {
	svc := s.serviceGroup.UserService
	_, err := svc.CreateUser("1", "user1")
	c.Assert(err, Equals, nil)
	_, err = svc.CreateUser("1", "user1")
	c.Assert(err, Not(Equals), nil)
}

func (s *TestSuite) TestSaveUserNormal(c *C) {
	svc := s.serviceGroup.UserService
	team := s.serviceGroup.TeamService.CreateTeam("1", "org", "team")
	user := User{"Id": "1", "Username": "user1", "Team": team}
	_, err := svc.CreateUser("1", "user1")
	c.Assert(err, Equals, nil)
}

func (s *TestSuite) TestSaveUserIdOrNameExists(c *C) {
	svc := s.serviceGroup.UserService
	user, err := svc.CreateUser("1", "user1")
	c.Assert(err, Equals, nil)
	user, err = svc.CreateUser("1", "user2")
	c.Assert(err, Not(Equals), nil)
	c.Assert(user, Equals, (*User)(nil))
	user, err = svc.CreateUser("2", "user1")
	c.Assert(err, Not(Equals), nil)
	c.Assert(user, Equals, (*User)(nil))
}
*/
