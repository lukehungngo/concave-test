package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	var (
		// testKey is a private key to use for funding a tester account.
		testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
		// testAddr is the address of the tester account.
		testAddr = crypto.PubkeyToAddress(testKey.PublicKey)
		dataChan = make(chan interface{}, 10)
	)
	db, err := NewLevelDb("./db", 256, 0)
	if err != nil {
		panic(fmt.Sprintf("can't open database connection: error=%+v", err))
	}
	blockRepository := NewLevelDbBlockRepository(db)

	Init(&Account{
		PrivateKey: testKey,
		Address:    testAddr,
	}, dataChan, blockRepository)

	go Run()

	engine := gin.Default()
	engine.Use(cors.Default())
	routerV0 := engine.Group("/api/v0")

	NewDataFetcher(routerV0, dataChan)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	port := 8080
	if err := engine.Run(":" + fmt.Sprintf("%d", port)); err != nil {
		panic(err)
	} else {
		fmt.Printf("Application is running at port %d\n", port)
	}
}
