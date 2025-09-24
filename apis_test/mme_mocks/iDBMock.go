package mme_mocks

import (
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/db"
	"gerrit.o-ran-sc.org/r/aiml-fw/awmf/modelmgmtservice/models"
	"github.com/stretchr/testify/mock"
)

type IDBMock struct {
	mock.Mock
	db.IDB
}

func (i *IDBMock) Create(modelInfo models.ModelRelatedInformation) error {
	args := i.Called(modelInfo)
	return args.Error(0)
}

func (i *IDBMock) GetByID(id string) (*models.ModelRelatedInformation, error) {
	return nil, nil
}

func (i *IDBMock) GetAll() ([]models.ModelRelatedInformation, error) {
	args := i.Called()
	if _, ok := args.Get(1).(error); !ok {
		return args.Get(0).([]models.ModelRelatedInformation), nil
	} else {
		var emptyModelInfo []models.ModelRelatedInformation
		return emptyModelInfo, args.Error(1)
	}
}

func (i *IDBMock) Update(modelInfo models.ModelRelatedInformation) error {
	return nil
}

func (i *IDBMock) Delete(id string) (int64, error) {
	return 1, nil
}

func (i *IDBMock) GetModelInfoByName(modelName string) ([]models.ModelRelatedInformation, error) {
	args := i.Called()
	if _, ok := args.Get(1).(error); !ok {
		return args.Get(0).([]models.ModelRelatedInformation), nil
	} else {
		var emptyModelInfo []models.ModelRelatedInformation
		return emptyModelInfo, args.Error(1)
	}
}

func (i *IDBMock) GetModelInfoByNameAndVer(modelName string, modelVersion string) (*models.ModelRelatedInformation, error) {
	args := i.Called()

	if _, ok := args.Get(1).(error); !ok {
		return args.Get(0).(*models.ModelRelatedInformation), nil
	} else {
		var emptyModelInfo *models.ModelRelatedInformation
		return emptyModelInfo, args.Error(1)
	}
}
