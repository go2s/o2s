// authors: wangoo
// created: 2018-05-30
// lambda oauth2 server

package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/gin-gonic/gin"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/go2s/o2s/engine"
	"github.com/go2s/o2s/o2"
	"github.com/go2s/o2m"
	"gopkg.in/mgo.v2"
)

var (
	initialized = false
	ginLambda   *ginadapter.GinLambda
)

type UriRedirector struct {
}

func (u *UriRedirector) FormatRedirectUri(uri string) string {
	return "/" + LambdaStaging + uri
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func handleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if !initialized {
		e := engine.NewEngine()
		ginLambda = ginadapter.New(e)
		initialized = true
	}
	// If no name is provided in the HTTP request body, throw an error
	response, err := ginLambda.Proxy(request)

	// Enable CORS
	response.Headers["Access-Control-Allow-Origin"] = "*"
	return response, err
}

var mgoCfg mongo.MongoConfig
var mgoSession *mgo.Session

func main() {
	mgoCfg = mongo.MongoConfig{
		Addrs:     mgoAddrs,
		Database:  mgoDatabase,
		Username:  mgoUsername,
		Password:  mgoPassword,
		PoolLimit: mgoPoolLimit,
	}

	mgoSession = mongo.NewMongoSession(&mgoCfg)

	ts := mongo.NewTokenStore(mgoSession)

	cs := mongo.NewClientStore(mgoSession, "oauth2", "client")

	us := mongo.NewUserStore(mgoSession, "oauth2", "user")

	o2.InitOauth2Server(cs, ts, us, &UriRedirector{})

	initSession()

	lambda.Start(handleRequest)
}
