package gae

import (
	"appengine"
	"appengine/datastore"
)

func ChannelKeyFor(context appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(context, "Channel", id, 0, nil)
}

func TeamKeyFor(context appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(context, "Team", id, 0, nil)
}

func UserKeyFor(context appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(context, "User", id, 0, nil)
}

func MessageKeyFor(context appengine.Context, id string) *datastore.Key {
	return datastore.NewKey(context, "Message", id, 0, nil)
}
