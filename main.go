package main

import (
	"fmt"
	log "logger/logger"
	"time"
)

// func ServeHTTP() {
// 	app := fiber.New()
// 	// port := flag.String("p", "8100", "port to serve on")
// 	// flag.Parse()
// 	app.Use(recover.New())
// 	app.Get("/", func(c *fiber.Ctx) error {
// 		panic("I'm an error")
// 	})
// 	app.Static("/log", "./logs", fiber.Static{
// 		ByteRange:     true,
// 		Browse:        true,
// 		CacheDuration: 1 * time.Second,
// 		MaxAge:        10,
// 	})

// 	app.Listen(":8100")
// }

func main() {
	go log.ServeLogFiles()
	l := log.DefaultLogger()
	for {
		errMessage := fmt.Sprint(`errodds----fdsfsd`)
		l.Error(errMessage)

		time.Sleep(time.Second * 5)
	}

}
