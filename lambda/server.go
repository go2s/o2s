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
	"flag"
	"log"
)

var (
	initialized = false
	lambdaMode  = true
	ginLambda   *ginadapter.GinLambda
)

type UriRedirector struct {
}

func (u *UriRedirector) FormatRedirectUri(uri string) string {
	if lambdaMode {
		return "/" + LambdaStaging + uri
	}
	return uri
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, auth, Cache-Control, X-Requested-With")
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

var mgoCfg o2m.MongoConfig
var mgoSession *mgo.Session

func main() {
	mgoCfg = o2m.MongoConfig{
		Addrs:     mgoAddrs,
		Database:  mgoDatabase,
		Username:  mgoUsername,
		Password:  mgoPassword,
		PoolLimit: mgoPoolLimit,
	}

	mgoSession = o2m.NewMongoSession(&mgoCfg)

	ts := o2m.NewTokenStore(mgoSession, mgoDatabase, "token")

	cs := o2m.NewClientStore(mgoSession, mgoDatabase, "client")

	us := o2m.NewUserStore(mgoSession, mgoDatabase, "user")

	o2.InitOauth2Server(cs, ts, us, &UriRedirector{})

	initSession()

	flag.BoolVar(&lambdaMode, "lambda", true, "lambda mode enable")
	flag.Parse()

	if lambdaMode {
		log.Println("lambda mode oauth2 server")
		lambda.Start(handleRequest)
	} else {
		log.Println("gin mode oauth2 server")
		e := engine.NewEngine()
		e.Run(":9096")
	}
}
