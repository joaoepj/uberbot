package ban

import (
	"fmt"
	"strconv"

	bot "github.com/ubergeek77/uberbot/core"
)

var banFail = "Failed to ban user"
var banSuccess = "Banned user"

var banInfo = bot.CreateCommandInfo(
	"ban",
	"Ban a user from the server. Also allows for deletion of message history up to 7 days",
	false,
	bot.Moderation)

func ban(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// This command requires at least 1 argument
	if len(ctx.Args) < 1 {
		response.Send(false, banFail, "Invalid syntax")
		return
	}

	// Parse Args
	banUser, err := ctx.Args["user"].UserValue(bot.Session)

	if err != nil {
		response.Send(false, banFail, "That user ID seems to be invalid!")
	}

	// Try inferring a number of days to delete from the second argument
	deleteDays := ctx.Args["days"].IntValue()

	// Can't have more than 7 days
	if deleteDays > 7 {
		deleteDays = 7
	}

	// Prepend how many days' worth of messages we're going to delete
	if deleteDays != 0 {
		response.PrependField("Delete user's messages newer than: ", strconv.Itoa(deleteDays)+" Days", false)
	}
	// If we have more arguments than just the user, try getting the number of days to delete and the ban reason
	reason := ctx.Args["reason"].StringValue()

	// Prepend the reason to the embed if it is not blank
	if reason != "" {
		response.PrependField("Ban reason:", reason, false)
	}

	// For the benefit of the audit log, add the person who issued this ban to the reason
	reason = fmt.Sprintf("(via %s#%s) ", ctx.Message.Author.Username, ctx.Message.Author.Discriminator) + reason

	// Make sure the USER exists, because they may not be a member
	if banUser == nil {
		response.Send(false, banFail, "Sorry, but I couldn't find that user!")
		return
	}

	// We have a user, start a response field with that user's information
	response.PrependField("Target user:", banUser.Mention(), false)

	// Can't ban the bot
	if banUser.ID == bot.Session.State.User.ID {
		response.Send(false, banFail, "But why do you want to ban me...? Did I do something wrong...? :(")
		return
	}

	// Can't mute yourself
	if banUser.ID == ctx.Message.Author.ID {
		response.Send(false, banFail, "You want to ban *yourself?* Maybe you should talk to someone about this...")
		return
	}

	// Can't mute a bot admin
	if bot.IsAdmin(banUser.ID) {
		response.Send(false, banFail, "Nice try, but that user is a Bot Administrator.\n\nI'd get in a LOT of trouble if I banned my boss!")
		return
	}

	// Can't mute moderators, *unless* you are a bot admin
	if ctx.Guild.IsMod(banUser.ID) && !bot.IsAdmin(ctx.Message.Author.ID) {
		response.Send(false, banFail, "Only Bot Administrators can ban other Bot Moderators.\n\nIs this a joke, or has a Bot Moderator gone rogue?!")
		return
	}

	//Ban 'em!
	if !ctx.Args["debug"].BoolValue() {
		err = ctx.Guild.Ban(banUser.ID, reason, deleteDays)
	}

	if err != nil {
		response.Send(false, banFail, "Sorry, I had some trouble baning that user. Please try again later.")
		bot.SendErrorReport(ctx.Guild.ID, "", banUser.ID, "Failed to ban user", err)
		return
	}

	// Send a success response
	response.Send(true, banSuccess, "And don't come back!!!")
}

func init() {
	banInfo.AddArg("user", bot.User, bot.ArgOption, "User to ban", true, "")
	banInfo.AddArg("days", bot.Int, bot.ArgOption, "days of msgs  to be deleted", false, "0")
	banInfo.AddArg("reason", bot.String, bot.ArgContent, "reason to be banned", false, "")
	banInfo.AddFlagArg("debug", bot.Boolean, bot.ArgFlag, "debug flag", false, "")
	bot.AddSlashCommand(banInfo)
	bot.AddCommand(banInfo, ban)
}
