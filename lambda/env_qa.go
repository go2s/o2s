//+build qa

// authors: wangoo
// created: 2018-05-30
// dev env

package main

const (
	LambdaStaging   = "qa"

	rdsAddr     = "sop-redis-for-dev.4qisdv.0001.cnn1.cache.amazonaws.com.cn:6379"
	rdsPassword = ""

	mgoDatabase  = "oauth2"
	mgoUsername  = "oauth2"
	mgoPassword  = "oauth2"
	mgoPoolLimit = 10
)

var mgoAddrs = []string{
	"172.31.0.10:27017",
	"172.31.0.11:27017",
	"172.31.0.12:27017",
}
