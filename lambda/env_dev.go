//+build dev

// authors: wangoo
// created: 2018-05-30
// dev env

package main

const (
	LambdaStaging = "dev"

	rdsAddr     = "127.0.0.1:6379"
	rdsPassword = ""

	mgoDatabase  = "oauth2"
	mgoUsername  = "oauth2"
	mgoPassword  = "oauth2"
	mgoPoolLimit = 10
)

var mgoAddrs = "mongodb://127.0.0.1:27017"
