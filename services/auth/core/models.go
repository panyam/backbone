package core

import (
	. "github.com/panyam/relay/services/msg/core"
	msgcore "github.com/panyam/relay/services/msg/core"
	. "github.com/panyam/relay/utils"
	"time"
)

type Registration struct {
	Id               int64
	Team             *msgcore.Team
	Username         string
	Created          time.Time
	ExpiresAt        time.Time
	AddressType      string
	Address          string
	VerificationData string
	Status           int
}

type UserLogin struct {
	Id int64

	/**
	 * Type of login for this user.
	 */
	LoginType string

	/**
	 * The login token for the user - could be phone number, FB ID, access token
	 * etc.
	 */
	LoginToken string

	/**
	 * Credentials to validate the login - eg FB access token + secret key
	 * or password or secret key for access token etc.
	 */
	Credentials map[string]interface{}

	User    *User
	Created time.Time
	Status  int
}

func RegistrationFromDict(data map[string]interface{}) (*Registration, error) {
	registration := Registration{}
	registration.Username = data["Username"].(string)
	registration.Address = data["Address"].(string)
	registration.AddressType = data["AddressType"].(string)
	registration.Id = String2ID(data["Id"].(string))
	return &registration, nil
}

func (r *Registration) ToDict() map[string]interface{} {
	out := map[string]interface{}{
		"Username":    r.Username,
		"Address":     r.Address,
		"AddressType": r.AddressType,
		"Id":          ID2String(r.Id),
	}
	return out
}
