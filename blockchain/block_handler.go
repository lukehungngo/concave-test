package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BlockHandler struct {
	blockRepository BlockRepository
}

func NewBlockHandler(router *gin.RouterGroup, blockRepository BlockRepository) {
	handler := &BlockHandler{
		blockRepository: blockRepository,
	}
	token := router.Group("/data")
	token.POST("/push", handler.GetBlocks)
}

func (bh *BlockHandler) GetBlocks(c *gin.Context) {
	limitStr := c.Param("limit")
	offsetStr := c.Param("offset")
	limit := 10
	offset := 0
	var err error
	if limitStr == "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusOK, NewResponseFail(ERROR_BAD_LIMIT, limitStr))
			return
		}
	}
	if offsetStr == "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusOK, NewResponseFail(ERROR_BAD_OFFSET, offsetStr))
			return
		}
	}

	blocks, err := bh.blockRepository.GetBlocks(limit, offset)
	if err != nil {
		c.JSON(http.StatusOK, NewResponseFail(ERROR_REPOSITORY, err))
		return
	}
	c.JSON(http.StatusOK, NewResponseSuccess(blocks))
	return
}
