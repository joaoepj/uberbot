package kick

import (
	"fmt"
	bot "github.com/ubergeek77/uberbot/core"
)

var kickFail = "Failed to kick member"
var kickSuccess = "Member kicked"

var kickInfo = bot.CreateCommandInfo(
	"kick",
	"Kicks a member",
	false,
	bot.Moderation)

func kick(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)
	// This command takes between 1 and 2 arguments
	if len(ctx.Args) < 1 {
		response.Send(false, kickFail, "Invalid syntax")
		return
	}

	// Get the reason
	reason := ""
	if len(ctx.Args) >= 2 {
		reason = ctx.Args["reason"].StringValue()

		// Prepend the reason
		response.PrependField("Kick reason:", reason, false)
	}

	// Prepend the reason with who issued this kick
	reason = fmt.Sprintf("(via %s#%s) ", ctx.Message.Author.Username, ctx.Message.Author.Discriminator) + reason

	// Get the member to be kicked
	member, err := ctx.Args["member"].UserValue(bot.Session)
	if err != nil {
		response.Send(false, kickFail, "Sorry, but I couldn't find that member!")
		return
	}

	// We have a member, start a response field with that member's information
	response.PrependField("Target member:", member.Mention(), false)

	// Can't kick the bot
	if member.ID == bot.Session.State.User.ID {
		response.Send(false, kickFail, "What?! Me?? But... what did I do??")
		return
	}

	// Can't mute yourself
	if member.ID == ctx.Message.Author.ID {
		response.Send(false, kickFail, "Wait... You want to kick yourself? But why?")
		return
	}

	// Can't mute a bot admin
	if bot.IsAdmin(member.ID) {
		response.Send(false, kickFail, "That member is a Bot Administrator, and is immune to your shenanigans.")
		return
	}

	// Can't mute moderators, *unless* you are a bot admin
	if ctx.Guild.IsMod(member.ID) && !bot.IsAdmin(ctx.Message.Author.ID) {
		response.Send(false, kickFail, "Only Bot Administrators can kick other Bot Moderators.\n\nEverything okay? Maybe you can talk this one out...")
		return
	}

	// Kick!
	if !ctx.Args["debug"].BoolValue() {
		err = ctx.Guild.Kick(member.ID, reason)
		if err != nil {
			response.Send(false, kickFail, "Sorry, I had some trouble kicking that member. Please try again later.")
			bot.SendErrorReport(ctx.Guild.ID, "", member.ID, "Failed to kick member", err)
			return
		}
	}

	// Send a success response
	response.Send(true, kickSuccess, "See ya later!")
}

func init() {
	kickInfo.AddArg("member", bot.User, bot.ArgOption, "User to kick", true, "")
	kickInfo.AddArg("reason", bot.String, bot.ArgContent, "Reason for kick", false, "")
	kickInfo.AddFlagArg("debug", bot.Boolean, bot.ArgFlag, "debugs the command", false, "")
	kickInfo.AddCmdAlias([]string{
		"ppPoof",
		"yeet",
		"begone",
		"k",
	})
	bot.AddSlashCommand(kickInfo)
	bot.AddCommand(kickInfo, kick)
}
