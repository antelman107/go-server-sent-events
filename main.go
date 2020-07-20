package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	e.Static("static", "./static")

	sender := NewSender()

	go func() {
		for i := 0; ; i++ {
			sender.Chan <- DataMessage(fmt.Sprintf("%d", i))

			time.Sleep(time.Second)
		}
	}()

	e.GET("/api/events", func(c echo.Context) error {
		useGzip := strings.Contains(c.Request().Header.Get("Accept-Encoding"), "gzip")
		if err := sender.Send(useGzip, c.Response().Writer, c.Request()); err != nil {
			return err
		}

		return nil
	})

	//e.Logger.Fatal(e.StartTLS(":1323", "server.crt", "server.key"))
	e.Logger.Fatal(e.Start(":1323"))
}
