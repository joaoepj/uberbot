package setprefix

import (
	bot "github.com/ubergeek77/uberbot/core"
)

// The reusable strings for this command.
var failTitle = "Failed to set prefix"
var successTitle = "Prefix updated"

// The information about this command.
var setprefixInfo = bot.CreateCommandInfo("setprefix", "Set the server prefix of the bot", false, bot.Utility)

// The function to run for this command.
func setprefix(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)
	// This command takes exactly 1 parameter
	if len(ctx.Args) != 1 {
		response.Send(false, failTitle, "Invalid syntax")
		return
	}

	// Set the prefix
	ctx.Guild.SetPrefix(ctx.Args["prefix"].StringValue())

	// Add a field to the response acknowledging the new prefix
	response.PrependField("New prefix:", "```\n"+ctx.Args["prefix"].StringValue()+"\n```", false)

	outputMsg := "OK! I'll respond to that prefix from now on.\n\nRemember, you can also @ me to run commands!"
	response.Send(true, successTitle, outputMsg)
}

// When this package is initialized, add this command to the bot.
func init() {
	setprefixInfo.AddArg("prefix", bot.String, bot.ArgOption, "Prefix to set", true, "")
	setprefixInfo.AddCmdAlias([]string{
		"sp",
	})
	setprefixInfo.SetTyping(true)
	bot.AddCommand(setprefixInfo, setprefix)
	bot.AddSlashCommand(setprefixInfo)
}
