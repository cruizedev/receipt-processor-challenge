package main

import (
	"github.com/cruizedev/receipt-processor-challenge/internal/handler"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func main() {
	e := echo.New()

	r := handler.NewReceipt()
	r.Register(e.Group("/receipts"))

	if err := e.Start(":1378"); err != nil {
		log.Error(err)
	}
}
