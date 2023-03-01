package mute

import (
	"time"

	bot "github.com/ubergeek77/uberbot/core"
)

//// The logger for the mute command
// currently unused
//var log = tlog.NewTaggedLogger("MuteCmd", tlog.NewColor("38;5;203"))

// The reusable strings for this command.
var muteFail = "Failed to mute member"
var muteSuccess = "Member muted"
var unmuteFail = "Failed to unmute member"
var unmuteSuccess = "Member unmuted"
var muteroleFail = "Failed to set mute role"
var muteroleSuccess = "Mute role updated"

var muteInfo = bot.CreateCommandInfo("mute", "Mute a user for a specified duration, or indefinitely", false, bot.Moderation)

var unmuteInfo = bot.CreateCommandInfo("unmute", "Unmute a muted user", false, bot.Moderation)

var muteroleInfo = bot.CreateCommandInfo("muterole", "Set the role to apply when a user is muted", false, bot.Utility)

// As a worker, this command will be called once every second
// This worker will check the list of mutes, unmuting and re-muting as needed.
func muteWorker() {
	for _, guild := range bot.Guilds {
		for userId, muteExpiry := range guild.Info.MutedUsers {
			// Make sure the user exists
			member, err := guild.GetMember(userId)
			if err != nil {
				continue
			}

			// Check if the mute has expired
			muteActive := (muteExpiry == 0) || time.Now().Unix() < muteExpiry

			// If the mute should be active, but they are not muted, they may be evading their mute
			if muteActive && !guild.HasRole(member.User.ID, guild.Info.MuteRoleId) {
				// Try re-muting the member
				err := guild.Mute(member.User.ID, muteExpiry-time.Now().Unix())
				if err != nil {
					bot.SendErrorReport(guild.ID, "", member.User.ID, "Failed re-muting potential mute evader", err)
					continue
				}

				// Don't bother sending a response if the guild doesn't have a configured response channel
				if guild.Info.ResponseChannelId == "" {
					continue
				}

				// Start a new response
				// Use a spoofed context so logs show this was started by the bot
				response := bot.NewResponse(&bot.Context{
					Guild:       guild,
					Cmd:         bot.CommandInfo{},
					Args:        nil,
					Message:     nil,
					Interaction: nil,
				}, false, false)

				// Prepend the target member to the response
				response.PrependField("Target member:", member.Mention(), false)

				// Send a response logging what happened
				response.Send(true, "Member re-muted", "A member may have been trying to evade their mute, and has been re-muted")

				continue
			}

			// If the mute has expired, unmute the member
			if !muteActive {
				// Unmute the user
				err = guild.UnMute(userId)
				if err != nil {
					bot.SendErrorReport(guild.ID, "", member.User.ID, "Failed to unmute member", err)
					continue
				}

				// Don't bother sending a response if the guild doesn't have a configured response channel
				if guild.Info.ResponseChannelId == "" {
					continue
				}

				// Start a new response
				// Use a spoofed context so logs show this was started by the bot
				response := bot.NewResponse(&bot.Context{
					Guild:       guild,
					Cmd:         bot.CommandInfo{},
					Args:        nil,
					Message:     nil,
					Interaction: nil,
				}, false, false)

				// Prepend the target member to the response
				response.PrependField("Target member:", member.Mention(), false)

				// Send a response logging what happened
				response.Send(true, "Member unmuted", "This member has served their mute and is now unmuted")

				continue
			}
		}
	}
}

func mute(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// This command takes between 1 and 2 arguments
	if ctx.Args["member"].StringValue() == "" {
		response.Send(false, muteFail, "Invalid syntax")
		return
	}

	// Make sure the mute role is set
	if ctx.Guild.Info.MuteRoleId == "" {
		response.Send(false, muteFail, "The mute role is not set! Set it with the `muterole` command.")
		return
	}

	// Get the member to be muted
	member, err := ctx.Args["member"].MemberValue(bot.Session, ctx.Guild.ID)
	if err != nil {
		response.Send(false, muteFail, "Sorry, but I couldn't find that member!")
		return
	}

	// We have a member, start a response field with that member's information
	response.PrependField("Target member:", member.Mention(), false)

	// Can't mute the bot
	if member.User.ID == bot.Session.State.User.ID {
		response.Send(false, muteFail, "You're trying to mute ***me?*** That's not very nice...")
		return
	}

	// Can't mute yourself
	if member.User.ID == ctx.Message.Author.ID {
		response.Send(false, muteFail, "Sorry, I can't let you mute yourself!")
		return
	}

	// Can't mute a bot admin
	if bot.IsAdmin(member.User.ID) {
		response.Send(false, muteFail, "Sorry, but bot administrators are immune to being muted!")
		return
	}

	// Can't mute moderators, *unless* you are a bot admin
	if ctx.Guild.IsMod(member.User.ID) && !bot.IsAdmin(ctx.Message.Author.ID) {
		response.Send(false, muteFail, "Sorry, but only Bot Administrators can mute Bot Moderators!")
		return
	}

	// Set some default values for indefinite mutes
	duration := 0
	displayDuration := "Indefinite"
	//muteMultiplier := 1

	// If it's not an indefinite mute, calculate the duration in seconds
	if ctx.Args["time"].StringValue() != "" {
		// A duration has been specified; calculate how long the mute is for
		muteDuration := ctx.Args["time"].StringValue()
		/*displayDuration = strconv.Itoa(muteMultiplier) + " " + displayDuration*/
		duration, displayDuration = bot.ParseTime(muteDuration)
	}

	// Add the duration field
	response.PrependField("Mute duration:", displayDuration, false)

	// Check if we already had a record of this mute before
	wasMutedBefore := ctx.Guild.HasMuteRecord(member.User.ID)

	// Mute the target for the specified duration
	if !ctx.Args["debug"].BoolValue() {
		err = ctx.Guild.Mute(member.User.ID, int64(duration))
	}
	if err != nil {
		bot.SendErrorReport(ctx.Guild.ID, "", member.User.ID, "Failed to mute member", err)
		response.Send(false, muteFail, "An unexpected error occurred when muting that member. Please try again.")
		return
	}

	// Send a response depending on if this was a mute update or not
	if wasMutedBefore {
		response.Send(true, "Mute updated", member.Mention()+"'s mute has been updated!")
	} else {
		response.Send(true, muteSuccess, "OK! Looks like "+member.Mention()+" won't be talking for a bit!")
	}
}

func unmute(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// This command takes exactly 1 argument
	if len(ctx.Args) != 1 {
		response.Send(false, unmuteFail, "Invalid syntax")
		return
	}

	// Make sure the mute role is set
	if ctx.Guild.Info.MuteRoleId == "" {
		response.Send(false, unmuteFail, "The mute role is not set! Set it with the `muterole` command.")
		return
	}

	// Get the member to be unmuted
	member, err := ctx.Args["member"].MemberValue(bot.Session, ctx.Guild.ID)
	if err != nil {
		response.Send(false, unmuteFail, "Sorry, but I couldn't find that member!")
		return
	}

	// We have a member, start a response field with that member's information
	response.PrependField("Target member:", member.Mention(), false)

	// Unmute them
	err = ctx.Guild.UnMute(member.User.ID)
	if err != nil {
		response.Send(false, unmuteFail, "Something went wrong when trying to unmute "+member.Mention()+". If they're still muted, please try again.")
		return
	}

	response.Send(true, unmuteSuccess, "You got it! "+member.Mention()+" is no longer muted.")
}

func muterole(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// This command expects exactly 1 parameter
	if len(ctx.Args) != 1 {
		response.Send(false, muteroleFail, "Invalid syntax")
		return
	}

	// Get the new role
	role, err := ctx.Args["muterole"].RoleValue(bot.Session, ctx.Guild.ID)
	if err != nil {
		response.Send(false, muteroleFail, "I had some trouble finding that role... Are you sure that role exists?")
		return
	}

	// We have a role, start a response field with that role's information
	response.PrependField("New mute role:", role.Mention(), false)

	// Set the mute role
	err = ctx.Guild.SetMuteRole(role.ID)
	if err != nil {
		bot.SendErrorReport(ctx.Guild.ID, "", "", "Error setting mute role to guaranteed-existing role: "+role.ID, err)
		response.Send(false, muteroleFail, "I had some trouble finding that role... Are you sure that role exists?")
		return
	}

	response.Send(true, muteroleSuccess, "Ok! I'll apply this role when muting users!")
}

func init() {
	// Add Args
	muteInfo.AddArg(
		"member",
		bot.User,
		bot.ArgOption,
		"member to be muted",
		true,
		"")
	muteInfo.AddArg(
		"time",
		bot.Time,
		bot.ArgOption,
		"amount of time to be muted",
		true,
		"")
	muteInfo.AddFlagArg("debug",
		bot.Boolean,
		bot.ArgFlag,
		"debug command",
		false,
		"")
	unmuteInfo.AddArg(
		"user",
		bot.User,
		bot.ArgOption,
		"User to unmute",
		true,
		"")
	muteroleInfo.AddArg(
		"muterole",
		bot.Role,
		bot.ArgOption,
		"role to set when muting",
		true,
		"")
	muteInfo.SetTyping(true)
	muteroleInfo.SetTyping(true)

	// Add commands
	bot.AddCommand(muteInfo, mute)
	bot.AddCommand(unmuteInfo, unmute)
	bot.AddCommand(muteroleInfo, muterole)
	//bot.AddSlashCommand(muteInfo)
	bot.AddSlashCommand(unmuteInfo)
	bot.AddSlashCommand(muteroleInfo)
	// Add Worker
	bot.AddWorker(muteWorker)
}
