package msgsniper

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	tlog "github.com/ubergeek77/tinylog"
	bot "github.com/ubergeek77/uberbot/core"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
	"strings"
	"time"
	"unicode"
)

var sniperFail = "Failed to change sniper settings"

//I took the sniper from ludwig and put it in otto's channel PepeLaugh

var log = tlog.NewTaggedLogger("MSGSniper", tlog.NewColor("38;5;63"))

func messageListener(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author == nil {
		return
	}
	// Gets the author object
	author := message.Author

	if author == nil {
		return
	}

	// Get the guild object
	if message.GuildID == "" {
		return
	}

	guild, ok := bot.Guilds[message.GuildID]
	if !ok {
		log.Errorf("Sniper was not able to find %s", message.GuildID)
		return
	}
	// Checks to see if the guild has the sniper enabled
	if !guild.IsSniperEnabled() {
		return
	}

	// Only take action if the user is not in certain roles
	if !guild.IsSnipeable(author.ID) {
		return
	}

	// If the member struct is nil or the message flags is crossposted
	// return early
	if message.Member == nil || message.Flags == discordgo.MessageFlagsIsCrossPosted {
		return
	}

	// Ignore members who have been in the server for at least 2 months
	joinSeconds := message.Member.JoinedAt.Unix()
	joinThreshold := joinSeconds + (2629800 * 2) // 2 months

	if joinThreshold <= time.Now().Unix() {
		return
	}

	if !guild.IsSniperChannel(message.ChannelID) {
		return
	}

	// Scrub unicode from the message content
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	scrubbed, _, _ := transform.String(t, message.Content)

	for _, word := range guild.Info.GuildBannedWords {
		if strings.Contains(strings.ToLower(scrubbed), word) {
			// Delete the message
			log.Infof("Deleting message from member %s#%s (%s): %s", message.Author.Username, message.Author.Discriminator, author.ID, message.Content)
			err := session.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				bot.SendErrorReport(message.GuildID, message.ChannelID, author.ID, "Failed to delete message: "+message.ID, err)
			}

			// Mute the user for 3 minutes
			muteErr := guild.Mute(author.ID, 180)
			if muteErr != nil {
				bot.SendErrorReport(message.GuildID, message.ChannelID, author.ID, "Failed to mute member during CrewLink defense: "+author.ID, muteErr)
			} else {
				log.Infof("Muted member for CrewLink reasons %s#%s (%s): %s", message.Author.Username, message.Author.Discriminator, author.ID, message.Content)

				// Spoof a context
				ctx := bot.Context{
					Guild:   guild,
					Cmd:     bot.CommandInfo{},
					Args:    nil,
					Message: nil,
				}

				// Generate a response
				response := bot.NewResponse(&ctx, false, false)
				response.PrependField("Target member:", author.Mention(), false)
				response.Send(true, "Member muted for asking about CrewLink", "This member is new to the server and seems to have not read the rules! They have been muted for 3 minutes.")
			}

			// Get the channel ID of the user to DM
			dmChannel, dmCreateErr := session.UserChannelCreate(author.ID)
			if dmCreateErr != nil {
				bot.SendErrorReport(message.GuildID, message.ChannelID, author.ID, fmt.Sprintf("Failed to send DM to member with CrewLink inquiry: %s#%s (%s)", message.Author.Username, message.Author.Discriminator, author.ID), dmCreateErr)
				// Send a message to general if the offending user is dumb
				msg, _ := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Hello %s! Welcome to the Ottomated Discord server!\n\nYou seem to have DMs disabled, so I was not able to send you my message. To keep it brief, I have deleted your message as I believe you were asking about CrewLink, which has been discontinued. This is covered in our <#744425536046104668>, so please take a moment to read them. Thanks for understanding!\n\n*(This message will self-destruct in 30 seconds)*", message.Author.Mention()))
				bot.SendErrorReport(message.GuildID, message.ChannelID, author.ID, fmt.Sprintf("Failed to send DM to member with CrewLink inquiry: %s#%s (%s)", message.Author.Username, message.Author.Discriminator, author.ID), dmCreateErr)
				time.Sleep(30 * time.Second)
				session.ChannelMessageDelete(msg.ChannelID, msg.ID)
				return
			}

			// Create a generic embed
			responseEmbed := bot.CreateEmbed(bot.ColorSuccess, "Hello! Welcome to the Ottomated Discord server!", "You seem to be asking about CrewLink.\n\nUnfortunately, [CrewLink has reached end of support](https://discord.com/channels/705685190244040766/852690862864466010/852691259272593418). CrewLink will not work on the latest version of Among Us, and no more updates are planned.\n\nThere are unofficial forks of CrewLink that may work with the latest version of Among Us, however we cannot provide you with technical support for those versions. We ask that you seek assistance in those forks' respective communities.\n\nYou are welcome to organize games that use CrewLink in our [#gaming](https://discord.com/channels/705685190244040766/774907319048994816) channel, but please avoid asking for or providing technical support in that channel.\n\nPlease remember that all CrewLink inquiries, like all messages in our Discord, are subject to the [#rules](https://discord.com/channels/705685190244040766/744425536046104668/827256674867478559). Please take a moment to read them if you have not already.\n\nYou have been temporarily paused from speaking in the Discord, to give you an opportunity to review the [#rules](https://discord.com/channels/705685190244040766/744425536046104668/827256674867478559). Don't worry! This will only last a few minutes.\n\nThanks for understanding!", nil)
			_, dmSendErr := session.ChannelMessageSendEmbed(dmChannel.ID, responseEmbed)
			if dmSendErr != nil {
				// Send a message to general if the offending user is dumb
				msg, _ := session.ChannelMessageSend(message.ChannelID, fmt.Sprintf("Hello %s! Welcome to the Ottomated Discord server!\n\nYou seem to have DMs disabled, so I was not able to send you my message. To keep it brief, I have deleted your message as I believe you were asking about CrewLink, which has been discontinued. This is covered in our <#744425536046104668>, so please take a moment to read them. Thanks for understanding!\n\n*(This message will self-destruct in 30 seconds)*", message.Author.Mention()))
				bot.SendErrorReport(message.GuildID, message.ChannelID, author.ID, fmt.Sprintf("Failed to send DM to member with CrewLink inquiry: %s#%s (%s)", message.Author.Username, message.Author.Discriminator, author.ID), dmSendErr)
				time.Sleep(30 * time.Second)
				session.ChannelMessageDelete(msg.ChannelID, msg.ID)
			} else {
				log.Infof("Successfully sent CrewLink info DM: %s#%s (%s): %s", message.Author.Username, message.Author.Discriminator, author.ID, message.Content)
			}
			return
		}
	}
}

// messageEditorListener
// Handler for edited messages.
func messageEditorListener(session *discordgo.Session, e *discordgo.MessageUpdate) {
	guildID := e.GuildID
	if guildID == "" {
		return
	}

	if !bot.Guilds[guildID].IsSniperEnabled() {
		return
	}

	messageListener(session, &discordgo.MessageCreate{
		Message: e.Message,
	})
}

var sniperInfo = bot.CreateCommandInfo("sniper", "Controls the sniper module", false, bot.Moderation)

func sniperCommand(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)
	if len(ctx.Args) < 1 || len(ctx.Args) > 4 {
		response.Send(false, sniperFail, "Invalid syntax, you did not enter a sub command")
		return
	}
	response.Send(false, sniperFail, "Invalid syntax, you did not enter a sub command")
	return
}

func init() {
	sniperInfo.SetParent(true, "")
	sniperInfo.AddArg("subcmdgrp", bot.SubCmd, bot.ArgOption, "Control the sniper!", true, "")
	sniperInfo.AddChoices("subcmdgrp", []string{
		"set",
		"unset",
	})
	sniperInfo.AddArg("subcmd", bot.SubCmd, bot.ArgOption, "Control the sniper!", true, "")
	sniperInfo.AddChoices("subcmd", []string{
		"add",
		"remove",
		"disable",
		"enable",
	})
	bot.AddHandler(messageListener)
	bot.AddHandler(messageEditorListener)
	bot.AddCommand(sniperInfo, sniperCommand)
	//bot.AddSlashCommand(sniperInfo)
}
