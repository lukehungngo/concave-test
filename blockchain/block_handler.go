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
	block := router.Group("/block")
	block.GET("/get", handler.GetBlocks)
}

// GetBlocks godoc
// @Summary      Get Many Blocks
// @Description  Get many blocks with default limit 10 and offset 0, block will be return as ascendent order
// @Tags         Block
// @Accept       json
// @Produce      json
// @Param        limit        query      int  false   "Limit"
// @Param        offset   	  query      int  false   "Offset"
// @Success      200          {object}  ResponseSuccess
// @Failure      200          {object}  ResponseFail
// @Router       /block/get [get]
func (bh *BlockHandler) GetBlocks(c *gin.Context) {
	limitStr := c.Query("limit")
	offsetStr := c.Query("offset")
	limit := 10
	offset := 0
	var err error
	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			c.JSON(http.StatusOK, NewResponseFail(ERROR_BAD_LIMIT, limitStr))
			return
		}
	}
	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil {
			c.JSON(http.StatusOK, NewResponseFail(ERROR_BAD_OFFSET, offsetStr))
			return
		}
	}
	blocks, err := bh.blockRepository.GetBlocks(limit, offset*limit)
	if err != nil {
		c.JSON(http.StatusOK, NewResponseFail(ERROR_REPOSITORY, err))
		return
	}
	c.JSON(http.StatusOK, NewResponseSuccess(blocks))
	return
}
