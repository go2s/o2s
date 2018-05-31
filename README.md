# Golang Oauth2 Server

## Feature:
- support memory,redis,mongodb token&client storage
- support memory,redis session

## Run Ouath2 Server:

```
go get -u -v gopkg.in/oauth2.v3/errors
go get -u -v github.com/satori/go.uuid
go get -u -v github.com/codegangsta/inject

go get -u -v github.com/go2s/o2r
go get -u -v github.com/go2s/o2m

cd http

go run server.go
```

## API Samples

### grant_type=client_credentials

```
# request
curl http://localhost:9096/oauth2/token?grant_type=client_credentials&client_id=000000&client_secret=999999&scope=read`

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
  http://127.0.0.1:9096/token \
  -H 'Authorization: Basic MDAwMDAwOjk5OTk5OQ==' \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'grant_type=refresh_token&refresh_token=EKBGAM00XQYBNA_VE_IYPW'

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
http://127.0.0.1:9096/authorize?redirect_uri=http%3A%2F%2Flocalhost&client_id=000000&response_type=token&state=xyz&scope=read

# response
http://localhost/#access_token=2PMWPTCTOXWXGKTSLH4TNQ&expires_in=3600&scope=read&state=xyz&token_type=Bearer
```

### valid token

```
# request
curl -X GET \
  http://127.0.0.1:9096/valid \
  -H 'Authorization: Bearer 2PMWPTCTOXWXGKTSLH4TNQ'

# response
{"ClientID":"000000","UserID":"000000","RedirectURI":"http://localhost","Scope":"read"
,"Code":"","CodeCreateAt":"0001-01-01T00:00:00Z","CodeExpiresIn":0
,"Access":"2PMWPTCTOXWXGKTSLH4TNQ","AccessCreateAt":"2018-05-29T11:10:17.533296011+08:00","AccessExpiresIn":3600000000000
,"Refresh":"","RefreshCreateAt":"0001-01-01T00:00:00Z","RefreshExpiresIn":0}
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
