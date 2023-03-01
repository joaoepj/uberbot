package output

import (
	bot "github.com/ubergeek77/uberbot/core"
	"strings"
)

// The reusable strings for this command.
var failTitle = "Failed to set output channel"
var successTitle = "Output channel updated"

// The information about this command.
var outputInfo = bot.CreateCommandInfo(
	"output",
	"Set the channel the bot should use for output",
	false,
	bot.Utility)

// The unset subcommand.
func unset(ctx *bot.Context) {
	// Unset the output channel
	ctx.Guild.SetResponseChannel("")

	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// Send acknowledgement
	response.Send(true, successTitle, "Alright! I'll now send my responses to commands in whatever channel they're called from!")
}

// The set subcommand.
func set(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// Make sure the provided ID is a valid channel first
	if channel, err := ctx.Args["channel"].ChannelValue(bot.Session); err == nil {
		// Add a field to the response
		response.PrependField("New output channel:", channel.Mention(), false)

		// Send the response
		response.Send(true, successTitle, "Alright! I'll be sure to use that channel for all future responses!")

		// Set the output channel *after* sending the response so all future responses will go in the output channel
		ctx.Guild.SetResponseChannel(channel.ID)
		return
	}

	// If we haven't returned, then the given argument was not a channel
	response.Send(false, failTitle, "Sorry, but I couldn't find that channel!")
}

// The function to run for this command.
func output(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// Make sure there is at least 1 argument
	if len(ctx.Args) != 0 {
		// Detect the subcommand, and pass on to the correct subcommand
		switch strings.ToLower(ctx.Args["type"].StringValue()) {
		case "unset":
			// Only 1 argument is expected
			if len(ctx.Args) == 1 {
				unset(ctx)
				return
			}
		case "set":
			// 2 arguments are expected
			if len(ctx.Args) == 2 {
				set(ctx)
				return
			}
		}
	}

	// If we haven't returned by now, we can assume a syntax error occurred
	response.Send(false, failTitle, "Invalid syntax")
}

// When this package is initialized, add this command to the bot.
func init() {
	// todo refactor to use subcmds
	outputInfo.AddArg("type", bot.String, bot.ArgOption, "set or unset", true, "")
	outputInfo.AddChoices("type", []string{
		"set",
		"unset",
	})
	outputInfo.AddArg("channel", bot.Channel, bot.ArgOption, "channel to set.", false, "")
	outputInfo.SetTyping(true)
	bot.AddCommand(outputInfo, output)
	bot.AddSlashCommand(outputInfo)
}
