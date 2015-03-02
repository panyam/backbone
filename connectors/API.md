
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

**Endpoints:** POST /users/register/

**Auth Required:** NO

**Parameters:**
- username: Must be unique
- address: An address the user can be sent the verification details to (similar the invite flow above).
- password: Optional password.  If a password is provided than username/password based logins will be allowed otherwise all calls that require an authentication MUST be with access token and secret key.  These can be used to change the password later on.  Also even if a password is not set, with a successful confirmation the access token and secret keys will be returned to the user.
    
**Return:** HTTP Status 200 on success and a registration ID that is valid for 5 minutes.

### Confirm a registration

**Endpoints:** POST /users/&lt;username&gt;/confirm/&lt;registrationid&gt;

**Auth Required:** NO

**Parameters:**
- verification_code: A verification code if it was present.

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
- password: Password associated with the account (if present).
    
**Return:**

HTTP Status 200 along with the sessionid cookie set that can be used in subsequent requests that require authentication.

### Logging out

**Endpoints:** POST /users/logout/
    
**Auth Required:** NO

**Return:**
HTTP Status 200 and the session ID cookies are cleared.

## Teams API

### List Teams

**Endpoints:** GET /users/&lt;userid&gt;/teams/
    
**Auth Required:** YES

**Return:**

A list of teams that the given user is subscribed to.  If the userid is not specified then the currently logged in user is queried, eg:

```
[ {"id": "123", "name": "Dream Team", "organization": "Dream Owner"} ]
```

### Create a team

**Endpoints:** POST /teams/
    
**Auth Required:** YES

**Parameters:**
- organization: Organization the team belongs to (optional)
- name: Name of the team (required and must be unique within the organization).
    
**Return:** HTTP Status 200 on success along with team details, eg:

```
{"id": "123", "name": "Dream Team", "organization": "Dream Owner"}
```

### Get team details

**Endpoints:** GET /teams/&lt;teamid&gt;/
    
**Auth Required:** NO

**Return:** HTTP Status 200 with team details, eg:

```
{"id": "123", "name": "Dream Team", "organization": "Dream Owner"}
```

### Invite user to a team

**Endpoints:** POST /teams/&lt;teamid&gt;/
    
**Auth Required:** YES.  User must also be permitted to invite users to a group (see user registration API)

**Parameters:**
- address: Invitee's address (can be a phone or email address).  
    If address is a phone number then the invitee is sent a verification code that is valid for 5 minutes.  The invitee can join (see Join API) with the given phone number and the verification code.
    If the address is an email then the invitee is sent a verification link that must be accessed to continue the joining process.
    Similar schemes will be used to allow other login methods (such as FB, Google and other OAuth).
    
**Return:** HTTP Status 200 with invitation ID.

```
{"id": "Invite101"}
```

### Accept an invite to join a team

**Endpoints:** GET /teams/&lt;teamid&gt;/join/&lt;invitationid&gt;
    
**Auth Required:** NO

**Parameters:**
- verification_code: If a verification code was sent then this MUST be present and match the invitation.
    
**Return:**
- HTTP Status 200 if invitation was successfully accepted.
- HTTP Status 401 if invitation id was invalid or if it has expired or if verification code did not match.

### List channels in a team

**Endpoints:** GET /teams/&lt;teamid&gt;/channels/
    
**Auth Required:** YES and current user must belong to the team.

**Parameters:**
- participants: Comma seperated list of userids of which atleast one user is a participant.
- status: Filter by channel status
- metadata.&lt;keypath&gt;: Filter by predicates on metadata entries.  See metadata filtering.
- order_by: Order by fields (prefixed by - indicates descending order):
    -   name - order name of the group
    -   created - order by created date
    -   updated - order by last updated
    -   last_messaged - order by last message date

**Return:**

All channels in the team that are visible to the current user.

## Channels API

### List channels for a user

**Endpoints:** GET /channels/

**Auth Required:** YES

**Parameters:**
- is_owner: Return channels user is an owner of
- participants: Comma seperated list of userids of which atleast one user is a participant.
- status: Filter by channel status
- metadata.&lt;keypath&gt;: Filter by predicates on metadata entries.  See metadata filtering.
- order_by: Order by fields (prefixed by - indicates descending order):
    -   name - order name of the group
    -   created - order by created date
    -   updated - order by last updated
    -   last_messaged - order by last message date
**Return:** List of channels that the user belongs to.

### Create new channel

**Endpoints:** POST /channels/

**Auth Required:** YES

**Parameters:**
- public: Whether the channel is public or private (default = true)
- participants: comma seperated list of user IDs.
- metadata: Dictionary of key value pairs.

**Returns:**
- HTTP stauts 200 with the channel details.


### Get channel details

**Endpoints:** GET /channels/&lt;channelid&gt;/

**Auth Required:** NO

**Returns:**
- HTTP Status 403 - if channel is not visible to the current user
- HTTP stauts 200 - If the channel is visible to the current user or is public and the channel details are returned, eg:

```
{
 'id': "channelid", 'name': "Channel Name", 'creator': "creatorid",
 'participants': [ "user1", "user2", "user3" ], 'status': 0
}
```

### Invite users to a channel

**Endpoints:** PUT /channels/&lt;channelid&gt;/invite/

**Auth Required:** YES and user must have permission to invite users.

**Parameters:**
- participants: comma seperated list of user IDs

**Returns:**
- HTTP Status 403 - if user does not have permissions to invite users.
- HTTP stauts 200 - If successful

### Join a channel

**Endpoints:** PUT /channels/&lt;channelid&gt;/join/

**Auth Required:** YES

**Returns:**
- HTTP Status 403 - if channel is private and user has not been sent an invitation
- HTTP stauts 200 - If successful

### Leave a channel

**Endpoints:** PUT /channels/&lt;channelid&gt;/leave/

**Auth Required:** YES and user must be a participant in the channel.

**Returns:**
- HTTP Status 403 - if not allowed.
- HTTP stauts 200 - If successful

### Delete a channel

**Endpoints:** PUT /channels/&lt;channelid&gt;/delete/

**Auth Required:** YES and user must be the owner of the channel.

**Returns:**
- HTTP Status 403 - if not allowed.
- HTTP stauts 200 - If successful

## Messages API

### Send a message in a channel

**Endpoints:** POST /channels/&lt;channelid&gt;/messages/

**Auth Required:** YES and the user must be authorized to send messages.

**Parameters:**
- type: Type of message
- body: Body of the message (as list of message parts).
- metadata: Metadata for the message.
- persist: Whether message is to be persisted or not.

**Return:** List of messages in the channel.

### Get messages in the channel

**Endpoints:** GET /channels/&lt;channelid&gt;/messages/

**Auth Required:** YES

**Parameters:**
- sender: Filter messages by sender
- type: Filter by message type.
- metadata.&lt;keypath&gt;: Filter by predicates on metadata entries.  See metadata filtering.
- text: Message content filtering.

**Return:** List of messages in the channel.

### Get message details

**Endpoints:** GET /messages/&lt;messageid&gt;/

**Auth Required:** YES and user must be able to read the message

**Return:** HTTP Status 200 and message details:

```
{
 'id': "messageid", 'sender': "senderuserid", 'sentAt': "sent at date",
 'type': "messagetype", 'metadata': {...}, 'channel': "channelid",
 'fragments': [
    <Fragment1>,
    <Fragment2>,
    ...
    <FragmentN>
 ]
}
```
