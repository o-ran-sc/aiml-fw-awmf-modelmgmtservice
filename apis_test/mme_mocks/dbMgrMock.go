package mme_mocks

import (
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/core"
	"github.com/stretchr/testify/mock"
)

type DbMgrMock struct {
	mock.Mock
	core.DBMgr
}

func (d *DbMgrMock) CreateBucket(bucketName string) (err error) {
	args := d.Called(bucketName)
	return args.Error(0)
}

func (d *DbMgrMock) UploadFile(dataBytes []byte, file_name string, bucketName string) error {
	args := d.Called()
	// If error is passed, return the error
	if _, ok := args.Get(0).(error); ok {
		return args.Get(0).(error)
	}

	return nil
}

func (d *DbMgrMock) ListBucket(bucketObjPostfix string) ([]core.Bucket, error) {
	args := d.Called()
	return args.Get(0).([]core.Bucket), args.Error(1)
}
