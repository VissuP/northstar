The dummy authentication service handles the following requests and sends back hardcoded responses for the
/oauth2/token, /oauth2/token/info and /api/v2/users/me. It uses gin-gonic to send back dummy responses to the portal
and the NS API.

Requests coming from the portal

/oauth2/token
/oauth2/revoke
/api/v2/models
/api/v2/users/me


Requests coming from NS API
/oauth2/token/info
/south/v2/users

Requests coming from the infrastructure
/management/health

The responses are constructed by dumping the actual responses from a real auth service and
sending back the exact same responses (headers and JSON body) through the dummy authentication service.
