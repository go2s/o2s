// authors: wangoo
// created: 2018-05-30
// lambda oauth2 server

package main

import (
	"flag"
	"github.com/mongodb/mongo-go-driver/mongo"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/go2s/o2s/captcha"
	"github.com/go2s/o2s/engine"
	"github.com/go2s/o2s/jwtex"
	"github.com/go2s/o2s/o2"
	"github.com/go2s/o2s/o2m"
	"github.com/go2s/o2s/o2r"
	"github.com/golang/glog"
)

var (
	initialized = false
	lambdaMode  = true
	ginLambda   *ginadapter.GinLambda
)

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !initialized {
		ginLambda = ginadapter.New(engine.GetGinEngine())
		initialized = true
	}
	// If no name is provided in the HTTP request body, throw an error
	response, err := ginLambda.Proxy(request)

	// Enable CORS
	response.Headers["Access-Control-Allow-Origin"] = "*"
	return response, err
}

var mgoCfg o2m.MongoConfig
var mgoClient *mongo.Client

func main() {
	flag.BoolVar(&lambdaMode, "lambda", false, "lambda mode enable")
	flag.Parse()
	flag.Set("logtostderr", "true") // Log to stderr only, instead of file.

	mgoCfg = o2m.MongoConfig{
		Address:   mgoAddrs,
		Database:  mgoDatabase,
		Username:  mgoUsername,
		Password:  mgoPassword,
		PoolLimit: mgoPoolLimit,
	}

	mgoClient = o2m.NewMongoClient(&mgoCfg)

	cfg := o2.DefaultServerConfig()
	cfg.ServerName = "Lambda Oauth2 Server"
	cfg.JWTSupport = true
	cfg.JWT = jwtex.JWTConfig{
		SignedKey:     []byte("go2s"),
		SigningMethod: jwt.SigningMethodHS512,
	}
	ts := o2m.NewTokenStore(mgoClient, mgoDatabase, "token")
	cs := o2m.NewClientStore(mgoClient, mgoDatabase, "client")
	us := o2m.NewUserStore(mgoClient, mgoDatabase, "user", o2m.DefaultMgoUserCfg())
	as := o2m.NewAuthStore(mgoClient, mgoDatabase, "auth")

	svr := o2.InitOauth2Server(cs, ts, us, as, cfg, engine.GinMap)
	redisOptions := &redis.Options{
		Addr: "127.0.0.1:6379",
	}
	mcs, err := o2r.NewRedisCaptchaStore(redisOptions, time.Minute*5)
	if err != nil {
		panic(err)
	}
	captcha.EnableCaptchaAuth(svr, mcs, captcha.CaptchaLogSender)

	initSession()

	if lambdaMode {
		glog.Info("lambda mode oauth2 server")
		cfg.URIPrefix = "/" + LambdaStaging
		lambda.Start(handleRequest)
	} else {
		glog.Info("gin mode oauth2 server")
		e := engine.GetGinEngine()
		e.Run(":9096")
	}
}
