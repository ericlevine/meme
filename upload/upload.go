package upload

import (
  "launchpad.net/goamz/aws"
  "launchpad.net/goamz/s3"
)

var bucket *s3.Bucket
var bucketName string

func init() {
  bucketName = "ericlevine"
  auth, err := aws.EnvAuth()
  if err != nil { panic(err.Error()) }
  s := s3.New(auth, aws.USEast)
  bucket = s.Bucket(bucketName)
}

func Write(data []byte, name, contentType string) (string, error) {
  err := bucket.Put(name, data, contentType, s3.PublicRead)
  if err != nil { return "", err }
  return "https://s3.amazonaws.com/ericlevine/" + name, nil
}
