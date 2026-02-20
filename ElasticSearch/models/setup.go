package models

import (
	"context"
	"crypto/tls"
	"net/http"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var ESClient *elasticsearch.Client

const SearchIndex = "blogs"

func ConnectDatabase() {
	dsn := "gin_elasticsearch:tmp_pwd@tcp(127.0.0.1:3306)/gin_elasticsearch?charset=utf8mb4&parseTime=True&loc=Local"

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB = database
}

func DBMigrate() {
	DB.AutoMigrate(&Blog{})
}

func ESClientConnection() {

	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
		},
		Username: "elastic",
		Password: "s3JhX1=qSaoIoAmLghxB", 
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	ESClient = client
}

func ESCreateIndexIfNotExist() {

	res, err := esapi.IndicesExistsRequest{
		Index: []string{SearchIndex},
	}.Do(context.Background(), ESClient)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	// If index does not exist → create it
	if res.StatusCode == 404 {
		_, err := ESClient.Indices.Create(SearchIndex)
		if err != nil {
			panic(err)
		}
	}
}