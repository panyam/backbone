package gae

import (
	"appengine"
	"appengine/datastore"
	"errors"
	. "github.com/panyam/relay/services/messaging/core"
	"log"
)

type TeamService struct {
	Cls         interface{}
	context     appengine.Context
	teamCounter int64
	teamsById   map[int64]*Team
	teamsByKey  map[string]*Team
	usersById   map[int64]*User
}

func NewTeamService(ctx appengine.Context) *TeamService {
	svc := TeamService{}
	svc.Cls = &svc
	svc.context = ctx
	svc.teamCounter = 1
	svc.teamsById = make(map[int64]*Team)
	svc.teamsByKey = make(map[string]*Team)
	svc.usersById = make(map[int64]*User)
	return &svc
}

/**
 * Lets a user create a team.
 */
func (c *TeamService) CreateTeam(id int64, org string, name string) (*Team, error) {
	key := org + ":" + name
	if _, ok := c.teamsByKey[key]; ok {
		return nil, errors.New("Team already exists with org and name")
	}
	if id == 0 {
		id = c.teamCounter
	} else if _, ok := c.teamsById[id]; ok {
		return nil, errors.New("Team already exists by ID")
	}
	team := NewTeam(id, org, name)
	c.teamsByKey[key] = team
	c.teamsById[id] = team
	c.teamCounter++
	return team, nil
}

/**
 * Retrieve teams in a org
 */
func (c *TeamService) GetTeamsInOrg(org string, offset int, count int) ([]*Team, error) {
	return nil, nil
}

/**
 * Retrieve a team by ID.
 */
func (c *TeamService) GetTeamById(id int64) (*Team, error) {
	return nil, nil
}

/**
 * Retrieve a team by Name.
 */
func (c *TeamService) GetTeamByName(org string, name string) (*Team, error) {
	key := org + ":" + name
	if team, ok := c.teamsByKey[key]; ok {
		return team, nil
	}
	return nil, errors.New("No such team")
}

/**
 * Delete a team.
 */
func (c *TeamService) DeleteTeam(team *Team) error {
	key := team.Organization + ":" + team.Name
	if _, ok := c.teamsByKey[key]; ok {
		delete(c.teamsByKey, key)
		if _, ok := c.teamsById[team.Id]; ok {
			delete(c.teamsById, team.Id)
			return nil
		}
	}
	return errors.New("No such team")
}

/**
 * Lets a user to join a team (if allowed)
 */
func (c *TeamService) JoinTeam(team *Team, username string) (*User, error) {
	return nil, nil
}

/**
 * Tells if a user belongs to a team.
 */
func (c *TeamService) TeamContains(team *Team, username string) bool {
	return false
}

/**
 * Returns the teams the user belongs to.
 */
func (c *TeamService) ListTeams(user *User) ([]*Team, error) {
	return nil, nil
}

/**
 * Lets a user leave a team or be kicked out.
 */
func (c *TeamService) LeaveTeam(team *Team, user *User) error {
	delete(c.usersById, user.Id)
	return nil
}

/**
 * Invite a user to a team.
 */
func (c *TeamService) InviteToTeam(inviter *User, invitee *User, team *Team) error {
	return nil
}

/**
 * Removes all entries.
 */
func (svc *TeamService) RemoveAllTeams() {
	q := datastore.NewQuery("Team").KeysOnly()
	keys, err := q.GetAll(svc.context, nil)
	if err != nil {
		log.Println("RemoveAll Error: ", err)
	}
	datastore.DeleteMulti(svc.context, keys)
}
