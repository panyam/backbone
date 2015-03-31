package core

import ()

/**
 * Service to handle user registrations, authentications and basic profiles.
 */
type IAuthService interface {
	/**
	 * Registers a user in a team.  If the username already exists in the team,
	 * then an error is returned.  Otherwise a registration object is returned.
	 */
	SaveRegistration(registration *Registration) error

	/**
	 * Gets registration by ID.
	 */
	GetRegistrationById(id int64) (*Registration, error)

	/**
	 * Delete a particular registration.
	 */
	DeleteRegistration(registration *Registration) error

	/**
	 * Logs in a user given by username and his/her credentials.
	 */
	LoginUser(username string, credentials map[string]interface{}) (*UserLogin, error)

	/**
	 * Removes all entries.
	 */
	RemoveAllRegistrations()

	/**
	 * Removes all entries.
	 */
	RemoveAllLogins()
}
