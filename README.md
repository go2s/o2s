# Golang Oauth2 Server

## Feature:
- storage support: memory,redis,mongodb 
- session support: memory,redis 

## Run Oauth2 Server:

```
go run http/server.go
```

## API Samples

### grant_type=client_credentials

```
# request
curl http://localhost:9096/oauth2/token?grant_type=client_credentials&scope=read \
  -H 'Authorization: Basic MDAwMDAwOjk5OTk5OQ=='

# response
{"access_token":"FZGYOSWDMQMX23BM5BHCWQ","expires_in":7200,"scope":"read","token_type":"Bearer"}
```

### grant_type=password

```
# request
curl -X POST \
  http://localhost:9096/oauth2/token \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -H 'Authorization: Basic MDAwMDAwOjk5OTk5OQ==' \
  -d 'grant_type=password&username=u1&password=123456'

# response
{
    "access_token": "WH3S8VDTOSGWJOWN5NCZPG",
    "expires_in": 7200,
    "refresh_token": "EKBGAM00XQYBNA_VE_IYPW",
    "token_type": "Bearer"
}
```

### grant_type=refresh_token

```
# request
curl -X POST \
  http://127.0.0.1:9096/oauth2/token \
  -H 'Authorization: Basic MDAwMDAwOjk5OTk5OQ==' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'grant_type=refresh_token&refresh_token=LX5J_I57WPOW8ZATJRDLYQ'

# response
{
    "access_token": "BXP5MQ1-NDO2WVNVJ1KT4Q",
    "expires_in": 7200,
    "refresh_token": "PTETKR3MX_6I63V-KDD2ZQ",
    "token_type": "Bearer"
}
```

### grant_type=implicit

```
# request
http://127.0.0.1:9096/oauth2/authorize?redirect_uri=http%3A%2F%2Flocalhost&client_id=000000&response_type=token&state=xyz&scope=read

# response
http://localhost/#access_token=2PMWPTCTOXWXGKTSLH4TNQ&expires_in=3600&scope=read&state=xyz&token_type=Bearer
```

### valid token

```
# request
curl -X GET \
  http://127.0.0.1:9096/oauth2/valid \
  -H 'Authorization: Bearer 48ZVGMI3PTWUYYCKMVKPAQ'

# response
{"client_id":"000000","user_id":"u1"}
```


## Mongodb

create mongodb user:

```
use oauth2
db.createUser(
   {
     user: "oauth2",
     pwd: "oauth2",
     roles: [ "readWrite", "dbAdmin" ]
   }
)
```
