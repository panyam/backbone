
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

## Teams API
Handles team management.

## Users API
Handles user management, settings and auths.

## Channels API
Channel management and control

## Messages API
Message sending, recetion, creation, updates, notification etc

