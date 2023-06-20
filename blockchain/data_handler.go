package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type DataFetcher struct {
	dataChan chan<- interface{}
}

func NewDataFetcher(router *gin.RouterGroup, dataChan chan interface{}) {
	handler := &DataFetcher{
		dataChan: dataChan,
	}
	token := router.Group("/data")
	token.POST("/push", handler.PushData)
}

func (df *DataFetcher) PushData(c *gin.Context) {
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, NewResponseFail(ERROR_PUSH_BAD_DATA, jsonData))
		return
	}
	df.pushData(jsonData)
	c.JSON(http.StatusOK, NewResponseSuccess("data is received"))
}

func (df *DataFetcher) pushData(data interface{}) {
	df.dataChan <- data
}
