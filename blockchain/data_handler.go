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
	data := router.Group("/data")
	data.POST("/push", handler.PushData)
}

// PushData godoc
// @Summary      Push data for new block
// @Description  Push data for new block, new block will be sent to other go-routine for block creation process
// @Tags         Data
// @Param        request        body      string  true   "Request Data"
// @Accept       json
// @Produce      json
// @Success      200          {object}  ResponseSuccess
// @Failure      200          {object}  ResponseFail
// @Router       /data/push [post]
func (df *DataFetcher) PushData(c *gin.Context) {
	jsonData, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusOK, NewResponseFail(ERROR_PUSH_BAD_DATA, jsonData))
		return
	}
	go df.pushData(jsonData)
	c.JSON(http.StatusOK, NewResponseSuccess("data is received"))
	return
}

func (df *DataFetcher) pushData(data interface{}) {
	df.dataChan <- data
}
