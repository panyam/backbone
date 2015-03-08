package services

import (
	. "github.com/panyam/backbone/services/core"
	. "gopkg.in/check.v1"
)

func (s *TestSuite) TestCreateUserService(c *C) {
	svc := CreateUserService()
	c.Assert(svc, Not(Equals), nil)
}

func (s *TestSuite) TestGetUserFirstTime(c *C) {
	svc := CreateUserService()
	user, _ := svc.GetUserById("1")
	c.Assert(user, Equals, (*User)(nil))
	user, _ = svc.GetUser("user1")
	c.Assert(user, Equals, (*User)(nil))
}

func (s *TestSuite) TestCreateUser(c *C) {
	svc := CreateUserService()
	user, _ := svc.CreateUser("1", "user1")
	c.Assert(user, Not(Equals), (*User)(nil))
	c.Assert(user.Id, Equals, "1")
	c.Assert(user.Username, Equals, "user1")

	fetched_user, _ := svc.GetUserById("1")
	c.Assert(fetched_user, Equals, user)

	fetched_user, _ = svc.GetUser("user1")
	c.Assert(fetched_user, Equals, user)
}

func (s *TestSuite) TestCreateUserDuplicate(c *C) {
	svc := CreateUserService()
	_, err := svc.CreateUser("1", "user1")
	c.Assert(err, Equals, nil)
	_, err = svc.CreateUser("1", "user1")
	c.Assert(err, Not(Equals), nil)
}

func (s *TestSuite) TestSaveUserNormal(c *C) {
	svc := CreateUserService()
	_, err := svc.CreateUser("1", "user1")
	c.Assert(err, Equals, nil)
}

func (s *TestSuite) TestSaveUserIdOrNameExists(c *C) {
	svc := CreateUserService()
	user, err := svc.CreateUser("1", "user1")
	c.Assert(err, Equals, nil)
	user, err = svc.CreateUser("1", "user2")
	c.Assert(err, Not(Equals), nil)
	c.Assert(user, Equals, (*User)(nil))
	user, err = svc.CreateUser("2", "user1")
	c.Assert(err, Not(Equals), nil)
	c.Assert(user, Equals, (*User)(nil))
}
