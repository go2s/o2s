//+build prod

// authors: wangoo
// created: 2018-05-30
// prod env

package main

const (
	LambdaStaging = "prod"

	rdsAddr     = "127.0.0.1:6379"
	rdsPassword = ""

	mgoDatabase  = "oauth2"
	mgoUsername  = "oauth2"
	mgoPassword  = "oauth2"
	mgoPoolLimit = 10
)

var mgoAddrs = "mongodb://172.16.101.10:27017,172.16.101.11:27017,172.16.100.11:27017"
