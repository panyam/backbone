package gae

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	. "github.com/panyam/relay/services/messaging/core"
	"log"
)

type UserService struct {
	Cls     interface{}
	context appengine.Context
}

func NewUserService(ctx appengine.Context) *UserService {
	svc := UserService{}
	svc.Cls = &svc
	svc.context = ctx
	return &svc
}

/**
 * Get user info by ID
 */
func (s *UserService) GetUserById(id int64) (*User, error) {
	var user User
	key := ChannelKeyFor(s.context, id)
	err := datastore.Get(s.context, key, &user)
	if err != nil {
		return nil, err
	}
	return &user, err
}

/**
* Get a user by username.
 */
func (s *UserService) GetUser(username string, team *Team) (*User, error) {
	query := datastore.NewQuery("User").Filter("Username =", username)
	t := query.Run(s.context)
	out := User{}
	_, err := t.Next(&out)
	if err != nil {
		return nil, fmt.Errorf("Error fetching results: %v", err)
	}
	return &out, nil
}

/**
 * Saves a user.
 * 	If the ID param is empty:
 * 		If username/team does not already exist a new one is created.
 * 		otherwise, it is updated and returned if override=true otherwise
 * 		false is returned.
 * 	Otherwise:
 * 		If username/team does not exist then it is written as is (Create or Update)
 * 		otherwise if IDs of curr and existing are different errow is thrown,
 * 		otherwise object is updated.
 */
func (s *UserService) SaveUser(user *User, override bool) error {
	// see if username/team exists
	return datastore.RunInTransaction(s.context, func(c appengine.Context) error {
		query := datastore.NewQuery("User").
			Filter("Username =", user.Username).
			Filter("Team =", user.Team)
		t := query.Run(c)
		var existing_user User
		key, err := t.Next(&existing_user)

		if err == datastore.Done {
			// username/team does NOT exist
			err = nil
			if user.Id == 0 {
				// create a new one
				key = datastore.NewIncompleteKey(c, "User", nil)
			}
			key, err = datastore.Put(c, key, user)
			if err == nil {
				user.Id = key.IntID()
			} else {
				log.Println(" =========== Daaaaaaaaaamn: ", err)
			}
		} else if err == nil {
			// found existing username/team
			if user.Id == 0 {
				if override {
					key, err = datastore.Put(c, key, user)
					if err == nil {
						user.Id = key.IntID()
					}
				} else {
					err = fmt.Errorf("Username/team already exists")
				}
			} else {
				if user.Id == key.IntID() {
					// then update
					key, err = datastore.Put(c, key, user)
				} else {
					err = fmt.Errorf("Username/team already exists")
				}
			}
		}
		return err
	}, nil)
}

/**
 * Removes all entries.
 */
func (svc *UserService) RemoveAllUsers() {
	q := datastore.NewQuery("User").KeysOnly()
	keys, err := q.GetAll(svc.context, nil)
	if err != nil {
		log.Println("RemoveAll Error: ", err)
	}
	datastore.DeleteMulti(svc.context, keys)
}
