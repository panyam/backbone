
# API Specs

All API calls are done will the following conventions described below.

Input
=====

* API calls are made with the endpoint prefix at - /api/
* When appropriate Results can be paginated with the parameters offset and count (defaulting to 0 and -1, respectively if not present)
* Parameters specified in the POST body will override those specified as URL query params.

Output
======
* Where possible http status codes will be used to denote return statuses.
* Body of the message will be in JSON unless the output format is specified
(either or via Accept header or via a "format" query parameter or POST
parameter in that order of precedence)

API is divided into following sections:

## Users API

### Register a User

**Endpoints:** 
    POST /users/register/
    
**Auth Required:** NO

**Parameters:**
    username: Must be unique
    address: An address the user can be sent the verification details to (similar the invite flow above).
    password: Optional password.  If a password is provided than username/password based logins will be allowed otherwise all calls that require an authentication MUST be with access token and secret key.  These can be used to change the password later on.  Also even if a password is not set, with a successful confirmation the access token and secret keys will be returned to the user.
    
**Return:**
    HTTP Status 200 on success and a registration ID that is valid for 5 minutes.

### Confirm a registration

**Endpoints:** 
    POST /users/&lt;username&gt;/confirm/&lt;registrationid&gt;
    
**Auth Required:** NO

**Parameters:**
    verification_code: A verification code if it was present.

**Return:**
    HTTP Status 200 on success along with user details:
    ```
    {'id': "userid", 'username': "username", 'token': "api_access_token", 'secret': "api_secret_key"}
    ```

### Logging in

**Endpoints:** 
    POST /users/&lt;username&gt;/login/
    
**Auth Required:** NO

**Parameters:**
    password: Password associated with the account (if present).
    
**Return:**
    HTTP Status 200 along with the sessionid cookie set that can be used in subsequent requests that require authentication.

### Logging out

**Endpoints:** 
    POST /users/logout/
    
**Auth Required:** NO

**Return:**
    HTTP Status 200 and the session ID cookies are cleared.

## Teams API

### List Teams
**Endpoints:** 
    GET /users/&lt;userid&gt;/teams/
    
**Auth Required:** YES

**Return:**

A list of teams that the given user is subscribed to.  If the userid is not specified then the currently logged in user is queried, eg:

```
[ {"id": "123", "name": "Dream Team", "organization": "Dream Owner"} ]
```

### Create a team

**Endpoints:** 
    POST /teams/
    
**Auth Required:** YES

**Parameters:**
    organization: Organization the team belongs to (optional)
    name: Name of the team (required and must be unique within the organization).
    
**Return:**
    HTTP Status 200 on success along with team details, eg:
```
{"id": "123", "name": "Dream Team", "organization": "Dream Owner"}
```

### Get team details

**Endpoints:** 
    GET /teams/&lt;teamid&gt;/
    
**Auth Required:** NO

**Return:**
    HTTP Status 200 with team details, eg:
```
{"id": "123", "name": "Dream Team", "organization": "Dream Owner"}
```

### Invite user to a team

**Endpoints:** 
    POST /teams/&lt;teamid&gt;/
    
**Auth Required:** YES.  User must also be permitted to invite users to a group (see user registration API)

**Parameters:**
    address: Invitee's address (can be a phone or email address).  
    If address is a phone number then the invitee is sent a verification code that is valid for 5 minutes.  The invitee can join (see Join API) with the given phone number and the verification code.
    If the address is an email then the invitee is sent a verification link that must be accessed to continue the joining process.
    Similar schemes will be used to allow other login methods (such as FB, Google and other OAuth).
    
**Return:**
    HTTP Status 200 with invitation ID.
```
{"id": "Invite101"}
```

### Accept an invite to join a team

**Endpoints:** 
    GET /teams/&lt;teamid&gt;/join/&lt;invitationid&gt;
    
**Auth Required:** NO

**Parameters:**
    verification_code: If a verification code was sent then this MUST be present and match the invitation.
    
**Return:**
    HTTP Status 200 if invitation was successfully accepted.
    HTTP Status 401 if invitation id was invalid or if it has expired or if verification code did not match.


### List channels in a team

**Endpoints:** 
    GET /teams/&lt;teamid&gt;/channels/
    
**Auth Required:** YES and current user must belong to the team.

**Return:**

All channels in the team that are visible to the current user.

## Channels API
Channel management and control

## Messages API
Message sending, recetion, creation, updates, notification etc

