package core

type BucketObject []byte

type Bucket struct {
	Name   string
	Object BucketObject
}
