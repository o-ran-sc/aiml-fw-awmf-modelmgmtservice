package main

import (
	"fmt"
	"example.com/mmes/apis"
)

func main() {

	fmt.Println("Starting api..")
	//Start apis declared in apis/mmets_apis.go
	//The mmes_apis will have a structure to hold instance of core/s3_manager
	//mmes_will use this se3 instance to all other se core functing
	//s3_manager would have aninstance of s3.session, which will be
	//confifured with aws credentials

	//router := gin.Default()
	mmes_api := &apis.Mmes_apis{}
	//mmes_api.GetModel()
	//router.GET("/getModel", mmes_api.GetModel())

	fmt.Println(mmes_api)
}
