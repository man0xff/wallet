package main

import (
	golog "log"
	"os"

	"github.com/kelseyhightower/envconfig"

	"example/core"
	"example/http"
)

func main() {
	log := golog.New(os.Stderr, "", 0)

	config := config{}
	err := envconfig.Process("app", &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	core := core.New(&core.Options{
		DBUser:     config.DB.User,
		DBPassword: config.DB.Password,
		DBAddress:  config.DB.Addr,
		DBName:     config.DB.Name,
		Log:        log,
	})
	server := http.New(&http.Options{
		Addr: config.HTTP.Addr,
		Core: core,
		Log:  log,
	})
	server.Serve()
}
