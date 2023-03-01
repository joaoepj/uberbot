package purge

import (
	"strconv"

	bot "github.com/ubergeek77/uberbot/core"
)

var purgeFail = "Failed to purge messages"
var purgeSuccess = "Messages purged"

var purgeInfo = bot.CreateCommandInfo("purge", "Purges Messages", false, bot.Moderation)

func purge(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)
	// Tell the response system not to reply
	response.Reply = false
	// This command expects between 1 and 3 arguments
	if len(ctx.Args) < 1 || len(ctx.Args) > 3 {
		response.Send(false, purgeFail, "Invalid syntax")
		return
	}

	// Parse all the possible items we might have
	deleteNum := ctx.Args["deletenum"].IntValue() // Defaults to 100
	purgeUser, err := ctx.Args["purgeuser"].UserValue(bot.Session)

	// If User is not nil, add it to the response
	if err == nil {
		response.PrependField("Target user:", purgeUser.Mention(), false)
	}
	purgeChannel, err := ctx.Args["purgechannel"].ChannelValue(bot.Session)

	// If Channel is not nil, add it to the response
	if err == nil {
		response.PrependField("Target channel:", purgeChannel.Mention(), false)
	}

	// Create a variable to track how many messages actually got deleted
	var confirmedDeleted = 1

	// Create a variable to track success
	var success bool

	// !purge N
	// (purge last N messages in current channel)
	if deleteNum != -1 && purgeUser.Username == "" && purgeChannel.Name == "" {
		confirmedDeleted, err = ctx.Guild.PurgeChannel(ctx.Message.ChannelID, deleteNum)
		if err != nil {
			bot.SendErrorReport(ctx.Guild.ID, ctx.Message.ChannelID, "", "Purge failed", err)
			response.Send(false, purgeFail, "Sorry, something went wrong... Please try again!")
			return
		}
		success = true
	}

	// !purge N #channel
	// (purge last N messages in #channel)
	if deleteNum != -1 && purgeUser.Username == "" && purgeChannel.Name != "" {
		confirmedDeleted, err = ctx.Guild.PurgeChannel(purgeChannel.ID, deleteNum)
		if err != nil {
			bot.SendErrorReport(ctx.Guild.ID, purgeChannel.ID, "", "Purge failed", err)
			response.Send(false, purgeFail, "Sorry, something went wrong... Please try again!")
			return
		}
		success = true
	}

	// !purge N @user
	// (purge user's last N messages server-wide)
	if deleteNum != -1 && purgeUser.Username != "" && purgeChannel.Name == "" {
		confirmedDeleted, err = ctx.Guild.PurgeUser(purgeUser.ID, deleteNum)
		if err != nil {
			bot.SendErrorReport(ctx.Guild.ID, "", purgeUser.ID, "Purge failed", err)
			response.Send(false, purgeFail, "Sorry, something went wrong... Please try again!")
			return
		}
		success = true
	}

	// !purge N @user #channel
	// (purge user's last N messages server-wide)
	if deleteNum != -1 && purgeUser.Username != "" && purgeChannel.Name != "" {
		confirmedDeleted, err = ctx.Guild.PurgeUserInChannel(purgeUser.ID, purgeChannel.ID, deleteNum)
		if err != nil {
			bot.SendErrorReport(ctx.Guild.ID, purgeChannel.ID, purgeUser.ID, "Purge failed", err)
			response.Send(false, purgeFail, "Sorry, something went wrong... Please try again!")
			return
		}
		success = true
	}

	// !purge @user
	// (purge user's last 100 messages server-wide)
	if deleteNum == -1 && purgeUser.Username != "" && purgeChannel.Name == "" {
		confirmedDeleted, err = ctx.Guild.PurgeUser(purgeUser.ID, 100)
		if err != nil {
			bot.SendErrorReport(ctx.Guild.ID, "", purgeUser.ID, "Purge failed", err)
			response.Send(false, purgeFail, "Sorry, something went wrong... Please try again!")
			return
		}
		success = true
	}
	// Determine success
	if success {
		// Add number of messages deleted to the response
		response.PrependField("Messages deleted:", strconv.Itoa(confirmedDeleted), false)

		// Respond with success
		response.Send(true, purgeSuccess, "Ok! I've tidied things up for you!")
	} else {
		// Assume syntax error
		response.Send(false, purgeFail, "Invalid syntax")
	}
}

func init() {
	purgeInfo.AddArg(
		"deletenum",
		bot.Int,
		bot.ArgOption,
		"Number of messages to delete",
		true,
		"5")
	purgeInfo.AddArg(
		"purgeuser",
		bot.User,
		bot.ArgOption,
		"User to purge msgs",
		false,
		"")
	purgeInfo.AddArg(
		"purgechannel",
		bot.Channel,
		bot.ArgOption,
		"Channel to purge msgs",
		false,
		"")
	purgeInfo.SetTyping(true)
	bot.AddCommand(purgeInfo, purge)
	bot.AddSlashCommand(purgeInfo)
}
