package main

import (
	"dtmgin/routes"
	//"github.com/go-resty/resty/v2"
)

func main() {
	app := routes.GetApp()

	go routes.Run(app)
	//err := app.Run(fmt.Sprintf(":%d", common.BusiPort))
	//go func() {
	//
	//}()

	select {}
	//dtmimp.FatalIfError(err)
}
