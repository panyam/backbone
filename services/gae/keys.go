package gae

import (
	"appengine"
	"appengine/datastore"
)

func ChannelKeyFor(context appengine.Context, id int64) *datastore.Key {
	return datastore.NewKey(context, "Channel", "", id, nil)
}

func TeamKeyFor(context appengine.Context, id int64) *datastore.Key {
	return datastore.NewKey(context, "Team", "", id, nil)
}

func UserKeyFor(context appengine.Context, id int64) *datastore.Key {
	return datastore.NewKey(context, "User", "", id, nil)
}

func MessageKeyFor(context appengine.Context, id int64) *datastore.Key {
	return datastore.NewKey(context, "Message", "", id, nil)
}
