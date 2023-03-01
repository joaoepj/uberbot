package disable

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	bot "github.com/ubergeek77/uberbot/core"
)

//todo finish migration of this command

// The reusable strings for these commands.
var enableFail = "Failed to enable command"
var enableSuccess = "Command enabled"
var disableFail = "Failed to disable command"
var disableSuccess = "Command disabled"
var channelenableFail = "Failed to enable command in channel"
var channelenableSuccess = "Command enabled in channel"
var channeldisableFail = "Failed to disable command in channel"
var channeldisableSuccess = "Command disabled in channel"

// A map of protected commands that can't be disabled.
var protected = map[string]bool{
	"perms":          true,
	"enable":         true,
	"disable":        true,
	"channelenable":  true,
	"channeldisable": true,
}

// The information about the enable command.
var enableInfo = bot.CreateCommandInfo(
	"enable",
	"Enable a command that was previously disabled (Globally)",
	false,
	bot.Utility)

// The information about the disable command.
var disableInfo = bot.CreateCommandInfo(
	"disable",
	"Disable a command (Globally)",
	false,
	bot.Utility)

// The information about the enable command.
var channelenableInfo = bot.CreateCommandInfo(
	"channelenable",
	"Enable a command that was previously disabled (in a single channel)",
	false,
	bot.Utility)

//Usage: []string{
//	"command (while in the target channel)",
//	"#channel command",
//},

// The information about the disable command.
var channeldisableInfo = bot.CreateCommandInfo(
	"channeldisable",
	"Disable a command (in a single channel)",
	false,
	bot.Utility)

//Usage: []string{
//	"command (while in the target channel)",
//	"#channel command",
//},

// The function to run for the enable command.
func enable(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// This command expects exactly 1 argument
	if len(ctx.Args) != 1 {
		response.Send(false, enableFail, "Invalid syntax")
		return
	}

	// We don't check if the command exists, as a safety measure for removing old commands that no longer exist
	targetCmd := ctx.Args["command"].StringValue()

	// Prepend a field to the response that reports what command is being worked with
	response.PrependField("Target command:", "```\n"+targetCmd+"\n```", false)

	// Make sure it's not already enabled
	if !ctx.Guild.IsGloballyDisabled(targetCmd) {
		response.Send(false, enableFail, "It looks like that command isn't disabled!")
		return
	}

	// Enable the command globally
	err := ctx.Guild.EnableTriggerGlobally(targetCmd)
	if err != nil {
		// We should never reach this, send an error report
		bot.SendErrorReport(ctx.Guild.ID, ctx.Message.ChannelID, "", "Impossible logic in disable.go; trigger was disabled before check, but not disabled on enable step", err)
		response.Send(false, enableFail, "It looks like that command isn't disabled!")
		return
	}

	outputMsg := fmt.Sprintf("Ok! The command `%s%s` is enabled again!", ctx.Guild.Info.Prefix, ctx.Cmd.Trigger)
	response.Send(true, enableSuccess, outputMsg)
}

// The function to run for the disable command.
func disable(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// This command expects exactly 1 argument
	if len(ctx.Args) != 1 {
		response.Send(false, enableFail, "Invalid syntax")
		return
	}

	// Get the command being worked with, and make sure it exists
	targetCmd := strings.ToLower(ctx.Args["command"].StringValue())
	if !bot.IsCommand(targetCmd) && !ctx.Guild.IsCustomCommand(targetCmd) {
		response.Send(false, disableFail, "I looked everywhere, but I couldn't find a command with that name!")
		return
	}

	// Prepend a field to the response that reports what command is being worked with
	response.PrependField("Target command:", "```\n"+targetCmd+"\n```", false)

	// Make sure the command is not protected
	if protected[targetCmd] {
		response.Send(false, disableFail, "Sorry, but that command is protected! It can't be disabled!")
		return
	}

	// Make sure it's not already enabled
	if ctx.Guild.IsGloballyDisabled(targetCmd) {
		response.Send(false, disableFail, "It looks like that command is already disabled!")
		return
	}

	// Disable the command globally
	err := ctx.Guild.DisableTriggerGlobally(targetCmd)
	if err != nil {
		// We should never reach this, send an error report
		bot.SendErrorReport(ctx.Guild.ID, ctx.Message.ChannelID, "", "Impossible logic in disable.go; trigger was enabled before check, but not enabled on disable step", err)
		response.Send(false, disableFail, "It looks like that command isn't disabled!")
		return
	}

	outputMsg := fmt.Sprintf("Ok! The command `%s%s` is now disabled!", ctx.Guild.Info.Prefix, ctx.Cmd.Trigger)
	response.Send(true, disableSuccess, outputMsg)
}

// The function to run for the enable command.
func channelenable(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// This command expects at least 1 argument, and at most 2 arguments
	if ctx.Args["command"].StringValue() == "" {
		response.Send(false, channelenableFail, "Invalid syntax")
		return
	}

	// We don't check if the command exists, as a safety measure for removing old commands that no longer exist
	targetCmd := ctx.Args["command"].StringValue()

	// Prepend a field to the response that reports what command is being worked with
	response.PrependField("Target command:", "```\n"+targetCmd+"\n```", false)

	var channel *discordgo.Channel
	var err error

	// With just 1 argument, the target channel is the current channel
	if channel.ID == "" {
		channel, err = ctx.Guild.GetChannel(ctx.Message.ChannelID)
		if err != nil {
			bot.SendErrorReport(ctx.Guild.ID, ctx.Message.ChannelID, "", "Channel lookup failed for guaranteed-existing channel", err)
			response.Send(false, channelenableFail, "Something went wrong when collecting information about this channel. Please try again!")
			return
		}
	}

	// With 2 arguments, the target channel is the second argument
	argChannel, err := ctx.Args["channelID"].ChannelValue(bot.Session)
	if err != nil {
		response.Send(false, channelenableFail, "Sorry, I wasn't able to find that channel!\n\nPlease be sure to mention the channel by name. For examples, `#channel`")
		return
	}

	if argChannel.ID != "" {
		if err != nil {
			response.Send(false, channelenableFail, "Sorry, I wasn't able to find that channel!\n\nPlease be sure to mention the channel by name. For examples, `#channel`")
			return
		}
		channel = argChannel
	}

	// We know the channel, append the target channel to the embed
	response.PrependField("Target channel:", "```\n"+channel.Mention()+"\n```", false)

	// Check if the command is enabled already
	if !ctx.Guild.TriggerIsDisabledInChannel(targetCmd, channel.ID) {
		response.Send(false, channelenableFail, "It looks like that command isn't disabled in here!")
		return
	}

	// Enable the command for this channel
	ctx.Guild.EnableTriggerInChannel(targetCmd, channel.ID)
	response.Send(true, channelenableSuccess, "OK! I've enabled the command in the specified channel!")
}

// The function to run for the disable command.
func channeldisable(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// This command expects at least 1 argument, and at most 2 arguments
	if len(ctx.Args) != 1 && len(ctx.Args) != 2 {
		response.Send(false, channeldisableFail, "Invalid syntax")
		return
	}

	// Get the command being worked with, and make sure it exists
	targetCmd := ctx.Args["command"].StringValue()
	if !bot.IsCommand(targetCmd) && !ctx.Guild.IsCustomCommand(targetCmd) {
		response.Send(false, channeldisableFail, "I looked everywhere, but I couldn't find a command with that name!")
		return
	}

	// Prepend a field to the response that reports what command is being worked with
	response.PrependField("Target command:", "```\n"+targetCmd+"\n```", false)

	// Make sure the command is not protected
	if protected[targetCmd] {
		response.Send(false, channeldisableFail, "Sorry, but that command is protected! It can't be disabled!")
		return
	}

	var channel *discordgo.Channel
	var err error

	// With just 1 argument, the target channel is the current channel
	if channel.ID == "" {
		channel, err = ctx.Guild.GetChannel(ctx.Message.ChannelID)
		if err != nil {
			bot.SendErrorReport(ctx.Guild.ID, ctx.Message.ChannelID, "", "Channel lookup failed for guaranteed-existing channel", err)
			response.Send(false, channelenableFail, "Something went wrong when collecting information about this channel. Please try again!")
			return
		}
	}

	// With 2 arguments, the target channel is the second argument
	argChannel, err := ctx.Args["channelID"].ChannelValue(bot.Session)
	if err != nil {
		response.Send(false, channelenableFail, "Sorry, I wasn't able to find that channel!\n\nPlease be sure to mention the channel by name. For examples, `#channel`")
		return
	}

	if argChannel.ID != "" {
		if err != nil {
			response.Send(false, channelenableFail, "Sorry, I wasn't able to find that channel!\n\nPlease be sure to mention the channel by name. For examples, `#channel`")
			return
		}
		channel = argChannel
	}

	// We know the channel, append the target channel to the embed
	response.PrependField("Target channel:", "```\n"+channel.Mention()+"\n```", false)

	// Check if the command is disabled already
	if ctx.Guild.TriggerIsDisabledInChannel(targetCmd, channel.ID) {
		response.Send(false, channeldisableFail, "It looks like that command isn't enabled in here!")
		return
	}

	// Disable the command for this channel
	ctx.Guild.DisableTriggerInChannel(targetCmd, channel.ID)
	response.Send(true, channeldisableSuccess, "No problem! That command is now disabled in the specified channel.")
}

// When this package is initialized, add these commands to the bot.
func init() {
	enableInfo.AddArg("command", bot.String, bot.ArgOption, "command to enable", true, "")
	bot.AddCommand(enableInfo, enable)
	bot.AddCommand(disableInfo, disable)
	bot.AddCommand(channelenableInfo, channelenable)
	bot.AddCommand(channeldisableInfo, channeldisable)
}
