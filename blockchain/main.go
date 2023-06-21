package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"lukehungngo/concave-test/v2/docs"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v0
func main() {

	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server concave test server."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost" + ":" + "8080"
	docs.SwaggerInfo.BasePath = "/api/v0"
	docs.SwaggerInfo.Schemes = []string{"http"}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(fmt.Errorf("recovery: f=%+v", r))
		}
	}()

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
	defer db.Close()

	blockRepository := NewLevelDbBlockRepository(db)

	blockProducer, err := Init(&Account{
		PrivateKey: testKey,
		Address:    testAddr,
	}, dataChan, blockRepository)
	if err != nil {
		panic(fmt.Sprintf("can't create block producer: error=%+v", err))
	}

	done := make(chan bool, 1)
	defer func() {
		close(done)
	}()

	go blockProducer.Run(done)

	engine := gin.Default()
	engine.Use(cors.Default())
	routerV0 := engine.Group("/api/v0")

	NewDataFetcher(routerV0, dataChan)
	NewBlockHandler(routerV0, blockRepository)
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	port := "8080"
	engineError := engine.Run(":" + port)
	if engineError != nil {
		fmt.Println("Start server error", engineError)
	}
}
