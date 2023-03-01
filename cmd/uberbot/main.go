package main

import (
	_ "github.com/ubergeek77/uberbot/v2/commands"
	bot "github.com/ubergeek77/uberbot/v2/core"
	_ "github.com/ubergeek77/uberbot/v2/eventhandlers"
	"github.com/ubergeek77/uberbot/v2/providers"
)

func main() {
	bot.SetInitProvider(providers.InitProvider)
	// Load environment variables from file
	bot.InitBot()
	// register stuff that should be run before the bot starts

	// run the bot
	bot.Run()
}
