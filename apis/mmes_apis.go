package apis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"example.com/mmes/core"
	"github.com/gin-gonic/gin"
)

type ModelInfo struct {
	ModelName string                 `json:"model-name"`
	RAppId    string                 `json:"rapp-id"`
	Metainfo  map[string]interface{} `json:"meta-info"`
}

func init() {
	fmt.Println("Starting api server...")
	router := gin.Default()

	mmes_api := &Mmes_apis{}
	router.GET("/getModel", mmes_api.GetModel)
	router.POST("/createModel", mmes_api.createModel)
	router.Run(os.Getenv("MMES_URL"))
	fmt.Println("Started api server...")
}

type Mmes_apis struct {
	name       string
	s3_manager *core.S3_manager
}

func (api *Mmes_apis) createModel(cont *gin.Context) {
	fmt.Println("Creating model...")
	bodyBytes, _ := io.ReadAll(cont.Request.Body)

	var modelInfo ModelInfo
	//Need to unmarshal JSON to Struct, to access request
	//data such as model name, rapp id etc
	err := json.Unmarshal(bodyBytes, &modelInfo)
	if err != nil {
		fmt.Println("Error in unmarshalling")
	}
	fmt.Println(modelInfo.ModelName, modelInfo.RAppId, modelInfo.Metainfo)

	//modelInfo.RAppId = "newRappId-1" Update Unmarshalled struct as per need
	//Need to convert struct to json to create a io.ReadSeeker instance
	//to insert in to a bucket as file/body
	modelInfo_json, err := json.Marshal(modelInfo)
	modelinfo_reader := bytes.NewReader(modelInfo_json) //bytes.Reader is type of io.ReadSeeker

	//TODO Create singleton for s3_manager
	s3_manager := core.New_s3manager()
	s3_manager.CreateBucket(modelInfo.ModelName)
	s3_manager.UploadFile(modelinfo_reader, modelInfo.ModelName+"_info.json", modelInfo.ModelName)

	cont.JSON(http.StatusOK, gin.H{
		"code":    http.StatusOK,
		"message": string("Model details stored sucessfully"),
	})
}

func (api *Mmes_apis) GetModel(cont *gin.Context) {
	fmt.Println("Fetching model")
	//s3_manager.GetBucketObject(modelInfo.ModelName+"_info.json", modelInfo.ModelName)
	cont.IndentedJSON(http.StatusOK, "Model stored ")
}

func (api *Mmes_apis) updateModel() {
	fmt.Println("Updating model...")
	return
}

func (api *Mmes_apis) deleteModel() {
	fmt.Println("Deleting model...")
	return
}
