package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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
	block.GET("/:number", handler.GetBlockBytesByNumber)
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

// GetBlockBytesByNumber godoc
// @Summary      Get Block as by Byte array representation given a block number
// @Description  Get Block as by Byte array representation given a block number
// @Tags         Block
// @Produce      json
// @Param        number        path      int  true   "Number"
// @Success      200          {object}  ResponseSuccess
// @Failure      200          {object}  ResponseFail
// @Router       /block/{number} [get]
func (bh *BlockHandler) GetBlockBytesByNumber(c *gin.Context) {
	numberStr := c.Param("number")
	number := uint64(0)
	var err error
	fmt.Println(numberStr)
	if numberStr != "" {
		number, err = strconv.ParseUint(numberStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, NewResponseFail(ERROR_BAD_LIMIT, numberStr))
			return
		}
	}
	block, has, err := bh.blockRepository.GetBlockByNumber(number)
	if err != nil {
		c.JSON(http.StatusOK, NewResponseFail(ERROR_REPOSITORY, err))
		return
	}
	if !has {
		c.JSON(http.StatusOK, NewResponseFail(ERROR_NOT_FOUND, fmt.Sprintf("block not found: block=%d", number)))
		return
	}
	m := make(map[string]interface{})
	m["block"] = block
	blockBytes, err := block.toBytes()
	if err != nil {
		c.JSON(http.StatusOK, NewResponseFail(ERROR_BAD_DATA, fmt.Sprintf("bad block: block=%d - error=%+v", number, err)))
		return
	}
	m["blockBytes"] = fmt.Sprintf("0x%v", common.Bytes2Hex(blockBytes))
	c.JSON(http.StatusOK, NewResponseSuccess(m))
	return
}
