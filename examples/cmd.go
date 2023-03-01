package examples

import bot "github.com/ubergeek77/uberbot/v2/core"

var exampleCmd = bot.CreateCommandInfo("test", "description", true, bot.Utility)

func test(ctx *bot.CmdContext) {

}
func init() {
	bot.AddCommand(exampleCmd, test)
	bot.AddSlashCommand(exampleCmd)
}