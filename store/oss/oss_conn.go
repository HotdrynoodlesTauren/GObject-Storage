package oss

import (
	"fmt"
	cfg "gobject-storage/config"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

var ossCli *oss.Client

// Client : create oss client object
func Client() *oss.Client {
	if ossCli != nil {
		return ossCli
	}
	ossCli, err := oss.New(cfg.OSSEndpoint,
		cfg.OSSAccesskeyID, cfg.OSSAccessKeySecret)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return ossCli
}

// Bucket : get the storage space of bucket
func Bucket() *oss.Bucket {
	cli := Client()
	if cli != nil {
		bucket, err := cli.Bucket(cfg.OSSBucket)
		if err != nil {
			fmt.Println(err.Error())
			return nil
		}
		return bucket
	}
	return nil
}

// DownloadURL : get the temporary url for downloading
func DownloadURL(objName string) string {
	signedURL, err := Bucket().SignURL(objName, oss.HTTPGet, 3600)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	return signedURL
}

// BuildLifecycleRule : Set bucket lifecycle rule
func BuildLifecycleRule(bucketName string) {
	// Set the expiration after 1 month for files with test as the prefix
	ruleTest1 := oss.BuildLifecycleRuleByDays("rule1", "test/", true, 30)
	rules := []oss.LifecycleRule{ruleTest1}

	Client().SetBucketLifecycle(bucketName, rules)
}
