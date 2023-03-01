package help

import (
	"fmt"
	"sort"

	"github.com/bwmarrin/discordgo"
	bot "github.com/ubergeek77/uberbot/core"
)

// NOTE
// For now, !help only responds in DMs
// This is a holdover until either pagination, or some proper formatting is added as to not spam the server

// The information about this command.
var helpInfo = bot.CreateCommandInfo(
	"help",
	"Show all available commands and their usage",
	true,
	bot.Utility,
)

// The function to run for this command.
func compileCoreCommands(ctx *bot.Context) *discordgo.MessageEmbed {
	// Get a list of all the core commands the bot knows about
	coreCmds := bot.GetCommands()

	// Get a sorted list of core command triggers
	coreSorted := make([]string, 0, len(coreCmds))
	for trigger := range coreCmds {
		coreSorted = append(coreSorted, trigger)
	}
	sort.Strings(coreSorted)

	// Create fields for each of the commands, only if the current user has access
	var fields []*discordgo.MessageEmbedField
	for _, trigger := range coreSorted {
		// We really don't need to print the help command
		if trigger == "help" {
			continue
		}
		current := coreCmds[trigger]
		if current.Public || (bot.IsAdmin(ctx.Message.Author.ID) || ctx.Guild.IsMod(ctx.Message.Author.ID)) {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("%s%s - %s", ctx.Guild.Info.Prefix, trigger, coreCmds[trigger].Description),
				Value:  ctx.Guild.GetCommandUsage(current),
				Inline: false,
			})
		}
	}

	// Create the embed for core commands if we have any
	var coreEmbed *discordgo.MessageEmbed
	if len(fields) > 0 {
		guild, err := bot.Session.Guild(ctx.Guild.ID)
		guildName := ""
		if err != nil {
			guildName = ctx.Guild.ID
		} else {
			guildName = guild.Name
		}
		coreEmbed = bot.CreateEmbed(bot.ColorSuccess, "Command Usage:", "*Commands for server \""+guildName+"\"*", fields)
	}

	return coreEmbed
}

func compileCustomCommands(ctx *bot.Context) *discordgo.MessageEmbed {
	// Get a list of all the custom commands in this guild
	customCmds := map[string]bot.CustomCommand{}
	for trigger, custom := range ctx.Guild.Info.CustomCommands {
		customCmds[trigger] = custom
	}

	// Get a sorted list of custom command triggers
	customSorted := make([]string, 0, len(customCmds))
	for trigger := range customCmds {
		customSorted = append(customSorted, trigger)
	}
	sort.Strings(customSorted)

	// Create fields for each of the commands, only if the current user has access
	var fields []*discordgo.MessageEmbedField
	for _, trigger := range customSorted {
		current := customCmds[trigger]
		if current.Public || (bot.IsAdmin(ctx.Message.Author.ID) || ctx.Guild.IsMod(ctx.Message.Author.ID)) {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   fmt.Sprintf("```\n%s%s\n```", ctx.Guild.Info.Prefix, trigger),
				Value:  fmt.Sprintf("```\n%s\n```", customCmds[trigger].Content),
				Inline: false,
			})
		}
	}

	// Create the embed for core commands if we have any
	var customEmbed *discordgo.MessageEmbed
	if len(fields) > 0 {
		guild, err := bot.Session.Guild(ctx.Guild.ID)
		guildName := ""
		if err != nil {
			guildName = ctx.Guild.ID
		} else {
			guildName = guild.Name
		}
		customEmbed = bot.CreateEmbed(bot.ColorSuccess, "Custom Commands:", "*Custom commands for server \""+guildName+"\"*", fields)
	}

	return customEmbed
}

func help(ctx *bot.Context) {
	if ctx.Interaction != nil {
		_ = bot.Session.InteractionRespond(ctx.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Flags:   1 << 6,
				Content: "Check your dms for help!",
			},
		})
	}
	coreEmbed := compileCoreCommands(ctx)
	customEmbed := compileCustomCommands(ctx)

	// Get intended recipient and DM channel ID
	recipient := ctx.Message.Author.ID
	dmChannel, dmCreateErr := bot.Session.UserChannelCreate(recipient)
	if dmCreateErr != nil {
		bot.SendErrorReport(ctx.Guild.ID, "", recipient, "Failed to obtain DM channel", dmCreateErr)
		return
	}

	if coreEmbed == nil && customEmbed == nil {
		_, dmSendErr := bot.Session.ChannelMessageSend(dmChannel.ID, "Hi there!\n\n"+
			"You just used the `help` command to get a list of commands I can help you use.\n\n"+
			"Unfortunately... You don't have access to any commands! Sorry about that!")
		if dmSendErr != nil {
			bot.SendErrorReport(ctx.Guild.ID, dmChannel.ID, recipient, "Failed sending DM", dmCreateErr)
			return
		}
	}

	if coreEmbed != nil {
		_, dmSendErr := bot.Session.ChannelMessageSendEmbed(dmChannel.ID, coreEmbed)
		if dmSendErr != nil {
			bot.SendErrorReport(ctx.Guild.ID, dmChannel.ID, recipient, "Failed sending DM", dmCreateErr)
			return
		}
	}

	if customEmbed != nil {
		_, dmSendErr := bot.Session.ChannelMessageSendEmbed(dmChannel.ID, customEmbed)
		if dmSendErr != nil {
			bot.SendErrorReport(ctx.Guild.ID, dmChannel.ID, recipient, "Failed sending DM", dmCreateErr)
			return
		}
	}
}

// When this package is initialized, add this command to the bot.
func init() {
	bot.AddSlashCommand(helpInfo)
	bot.AddCommand(helpInfo, help)
}
