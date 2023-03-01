package perms

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	bot "github.com/ubergeek77/uberbot/core"
)

// The reusable strings for this command.
var failTitle = "Failed to update permissions"
var successTitle = "Permissions updated"

var permsMeta = bot.CreateCommandInfo(
	"perms",
	"Grant or revoke bot permissions for a user or role",
	false,
	bot.Utility)

var permsSlashInfo = []*discordgo.ApplicationCommandOption{
	{
		Type:        discordgo.ApplicationCommandOptionString,
		Name:        "option",
		Description: "Grants or revokes perms",
		Required:    true,
		Choices: []*discordgo.ApplicationCommandOptionChoice{
			{
				Name:  "Grant",
				Value: "grant",
			},
			{
				Name:  "Revoke",
				Value: "revoke",
			},
		},
	},
	{
		Type:        discordgo.ApplicationCommandOptionUser,
		Name:        "user",
		Description: "User to grant/revoke perms",
		Required:    false,
	},
	{
		Type:        discordgo.ApplicationCommandOptionRole,
		Name:        "role",
		Description: "Role to grant/revoke perms",
		Required:    false,
	},
}

func grant(ctx *bot.Context, targetID string, output string, mention string) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// Only grant permissions if they are not already granted
	if !ctx.Guild.IsMod(targetID) {
		err := ctx.Guild.AddMod(targetID)
		if err != nil {
			response.Send(false, failTitle, "Sorry, I couldn't find that member or role!")
			return
		}

		// Add a new field
		response.PrependField("Granted to:", mention, false)

		// Send the response
		response.Send(true, successTitle, output)
	} else {
		response.Send(false, failTitle, "That target already has bot permissions!")
	}
}

func revoke(ctx *bot.Context, targetID string, output string, mention string) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// Only revoke permissions if they are already granted
	if ctx.Guild.IsMod(targetID) {
		ctx.Guild.RemoveMod(targetID)

		// Add a new field
		response.PrependField("Revoked from:", mention, false)

		// Send the response
		response.Send(true, successTitle, output)
	} else {
		response.Send(false, failTitle, "That target doesn't have bot permissions!")
	}
}

// The function to run for this command.
func perms(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// Perms is a special command; only bot administrators can run it
	if !bot.IsAdmin(ctx.Message.Author.ID) {
		response.Send(false, failTitle, "Sorry, only Bot Administrators can manage permissions!")
		return
	}

	// This command expects exactly 3 parameters
	if len(ctx.Args) > 3 {
		response.Send(false, failTitle, "Invalid syntax")
		return
	}

	// Detect the subcommand
	action := strings.ToLower(ctx.Args["type"].StringValue())
	if action != "grant" && action != "revoke" {
		response.Send(false, failTitle, "Invalid syntax")
		return
	}

	// Get the ID being worked with
	// Also generate the success messages in this block
	grantOutput := ""
	revokeOutput := ""
	targetID := ""
	mention := ""
	if member, err := ctx.Args["user"].MemberValue(bot.Session, ctx.Guild.ID); err == nil {
		targetID = member.User.ID
		mention = member.Mention()
		grantOutput = "You got it! " + member.Mention() + " can now run protected commands!"
		revokeOutput = "Understood! I'll stop allowing " + member.Mention() + " to run protected commands."
	} else if role, err := ctx.Args["role"].RoleValue(bot.Session, ctx.Guild.ID); err == nil {
		targetID = role.ID
		mention = role.Mention()
		grantOutput = "You got it! Members from the \"" + role.Name + "\" role can now run protected commands!"
		revokeOutput = "You got it! I'll stop allowing members from the \"" + role.Name + "\" role from running protected commands!"
	} else {
		response.Send(false, failTitle, "Sorry, I can't find a user or role with that ID!")
		return
	}

	// Make sure the bot's own ID is not trying to be messed with
	if targetID == bot.Session.State.User.ID {
		response.Send(false, failTitle, "Hey, that's me! I can't grant or revoke my own permissions!")
		return
	}

	// Make sure the user can't manage their own permissions
	if targetID == ctx.Message.Author.ID {
		response.Send(false, failTitle, "Wait a second! You can't manage your own permissions!")
		return
	}

	// Pass the commands off based on the subcommand
	if action == "grant" {
		grant(ctx, targetID, grantOutput, mention)
	} else if action == "revoke" {
		revoke(ctx, targetID, revokeOutput, mention)
	}
}

// When this package is initialized, add this command to the bot.
func init() {
	permsMeta.AddArg("type", bot.String, bot.ArgOption, "grant or revoke", true, "")
	permsMeta.AddChoices("type", []string{
		"grant",
		"revoke",
	})
	permsMeta.AddArg("user", bot.User, bot.ArgOption, "user to (grant/revoke) perms to", false, "")
	permsMeta.AddArg("role", bot.Role, bot.ArgOption, "role to (grant/revoke) perms to", false, "")
	bot.AddSlashCommand(permsMeta)
	bot.AddCommand(permsMeta, perms)
}
