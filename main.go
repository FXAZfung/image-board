package main

import (
	"github.com/FXAZfung/image-board/cmd"
	_ "github.com/FXAZfung/image-board/docs"
)

// @title           Image Board API
// @version         1.0
// @description     This is a image sharing platform API documentation

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:4536
// @BasePath  /
// @schemes   http https

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        Authorization
// @description                 "Type 'YOUR_TOKEN'"

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	cmd.Execute()
}
