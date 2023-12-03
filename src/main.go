package main

import (
	"fmt"
	"telegram-bot/src/app"
)

func main() {
	telegramer, err := app.NewTelegramer()
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	err = telegramer.Start()
	fmt.Printf("%v", err)
}
