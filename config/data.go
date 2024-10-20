package config

import "encoding/json"

type AppConfigData struct {
	MMES_URL      string `json:"mmes_url"`
	LOG_FILE_NAME string `json:"log_file_name"`
}

func (a AppConfigData) String() string {
	b, _ := json.MarshalIndent(a, "", "  ")
	return string(b)
}

type DBConfigData struct {
	MODEL_FILE_POSTFIX string `json:"model_file_postfix"`
	INFO_FILE_POSTFIX  string `json:"info_file_postfix"`
	S3_URL             string `json:"s3_url"`
	S3_ACCESS_KEY      string `json:"s3_access_key"`
	S3_SECRET_KEY      string `json:"s3_secret_key"`
	S3_REGION          string `json:"s3_region"`
	PG_USER            string `json:"pg_user"`
	PG_PASSWORD        string `json:"pg_password"`
	PG_HOST            string `json:"pg_host"`
	PG_PORT            string `json:"pg_port"`
	PG_DBNAME          string `json:"pg_dbname"`
}

func (d DBConfigData) String() string {
	b, _ := json.MarshalIndent(d, "", "  ")
	return string(b)
}
