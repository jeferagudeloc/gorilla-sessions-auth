package main

import (
	"os"
	"time"

	"github.com/jeferagudeloc/gorilla-sessions-auth/infrastructure"
	"github.com/jeferagudeloc/gorilla-sessions-auth/infrastructure/database"
	"github.com/jeferagudeloc/gorilla-sessions-auth/infrastructure/http"
	"github.com/jeferagudeloc/gorilla-sessions-auth/infrastructure/log"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	var app = infrastructure.NewConfig().
		Name(os.Getenv("APP_NAME")).
		ContextTimeout(10 * time.Second).
		Logger(log.InstanceLogrusLogger).
		SqlSetup(database.InstanceMysql).
		WebServerPort(os.Getenv("APP_PORT")).
		WebServer(http.InstanceGorillaMux)
	app.StartServers()
}
