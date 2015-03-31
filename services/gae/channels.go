package gae

import (
	"appengine"
	"appengine/datastore"
	"errors"
	"fmt"
	. "github.com/panyam/relay/services/core"
	"log"
)

type ChannelService struct {
	Cls     interface{}
	context appengine.Context
}

func NewChannelService(ctx appengine.Context) *ChannelService {
	svc := ChannelService{}
	svc.Cls = &svc
	svc.context = ctx
	return &svc
}

func (c *ChannelService) SaveChannel(channel *Channel, override bool) error {
	key := ChannelKeyFor(c.context, channel.Id)
	if channel.Id == 0 {
		key, err := datastore.Put(c.context, key, channel)
		if err == nil {
			channel.Id = key.IntID()
		}
		return err
	}

	if override {
		key, err := datastore.Put(c.context, key, channel)
		if err == nil {
			channel.Id = key.IntID()
		}
		return err
	}

	// verify key first
	var err error
	datastore.RunInTransaction(c.context, func(c appengine.Context) error {
		var existing_channel Channel
		err = datastore.Get(c, key, &existing_channel)
		var key *datastore.Key
		if err == nil {
			// entry already exists
			err = errors.New("Channel already exist")
		} else {
			key, err = datastore.Put(c, key, channel)
			channel.Id = key.IntID()
		}
		return err
	}, nil)
	return err
}

/**
 * Retrieve a channel by ID
 */
func (c *ChannelService) GetChannelById(id int64) (*Channel, error) {
	var channel Channel
	key := ChannelKeyFor(c.context, id)
	err := datastore.Get(c.context, key, &channel)
	return &channel, err
}

/**
 * Delete a channel.
 */
func (c *ChannelService) DeleteChannel(channel *Channel) error {
	key := ChannelKeyFor(c.context, channel.Id)
	return datastore.Delete(c.context, key)
}

/**
 * Lets a user to join a channel (if allowed)
 */
func (c *ChannelService) JoinChannel(channel *Channel, user *User) error {
	if channel.ContainsUser(user) {
		return nil
	}
	channel.Participants = append(channel.Participants, user)
	return c.SaveChannel(channel, true)
}

/**
 * Lets a user leave a channel or be kicked out.
 */
func (c *ChannelService) LeaveChannel(channel *Channel, user *User) error {
	for index, value := range channel.Participants {
		if value.Id == user.Id {
			channel.Participants = append(channel.Participants[:index], channel.Participants[index+1:]...)
			return c.SaveChannel(channel, true)
		}
	}
	return nil
}

/**
 * Returns the channels the user belongs to.
 */
func (c *ChannelService) ListChannels(user *User, team *Team) ([]*Channel, error) {
	query := datastore.NewQuery("Channel").Filter("Participants =", user).Filter("Team =", team)
	t := query.Run(c.context)
	out := make([]*Channel, 0, 100)
	for {
		var channel Channel
		_, err := t.Next(&channel)
		if err == datastore.Done {
			break
		} else if err != nil {
			return nil, fmt.Errorf("Error fetching results: %v", err)
		}
		out = append(out, &channel)
	}
	return out, nil
}

/**
 * Removes all entries.
 */
func (svc *ChannelService) RemoveAllChannels() {
	var objs []Channel
	q := datastore.NewQuery("Channel").KeysOnly()
	keys, err := q.GetAll(svc.context, &objs)
	if err != nil {
		log.Println("RemoveAll Error: ", err)
	}
	datastore.DeleteMulti(svc.context, keys)
}
