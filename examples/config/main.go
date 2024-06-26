package main

import (
	"log"

	"github.com/wyy-go/wzo/core/config"
)

func main() {
	config.OnChange(func(c config.Config) {
		log.Println("on change ...")
	})

	log.Println(config.GetStringMap("app"))
	log.Printf("app.name=%s", config.GetString("app.name"))

	select {}
}
