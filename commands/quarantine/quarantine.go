package quarantine

import (
	"github.com/bwmarrin/discordgo"
	tlog "github.com/ubergeek77/tinylog"
	bot "github.com/ubergeek77/uberbot/core"
	"strings"
)

var log = tlog.NewTaggedLogger("BotDetection", tlog.NewColor("38;5;160"))

func messageListener(session *discordgo.Session, message *discordgo.MessageCreate) {
	// Get the author object
	author := message.Author

	if author == nil {
		return
	}

	// Don't do anything if we don't have a state, trying to access .User.ID will crash
	if bot.Session.State == nil {
		return
	}

	// Ignore the bot
	if bot.Session.State.User.ID == author.ID {
		return
	}

	// Get the guild object
	if message.GuildID == "" {
		return
	}

	guild, ok := bot.Guilds[message.GuildID]

	if !ok {
		log.Errorf("Quarantine was not able to find %s", message.GuildID)
		return
	}

	// This only applies to the Otto server
	if guild.ID != "705685190244040766" {
		return
	}

	// Ignore Otto
	if author.ID == "194259969359478784" {
		return
	}

	// Ignore moderator role
	// I heard you liked hardcoding
	if guild.HasRole(author.ID, "750114062934737038") {
		return
	}

	// Ignore Otto role
	if guild.HasRole(author.ID, "706294556508946454") {
		return
	}

	// Ignore sub role
	// Ottomatons
	if guild.HasRole(author.ID, "744129398436790342") {
		return
	}

	// Ignore supporter role
	// HUGE SIMP
	if guild.HasRole(author.ID, "781455132428075019") {
		return
	}

	// Ignore boost role
	// Server Booster
	if guild.HasRole(author.ID, "766312976168255500") {
		return
	}

	// Ignore Beta Tester role
	if guild.HasRole(author.ID, "765018086939557928") {
		return
	}

	// Ignore editor role
	if guild.HasRole(author.ID, "779957237862891540") {
		return
	}

	// Ignore Hackathon role
	if guild.HasRole(author.ID, "832093514077700117") {
		return
	}

	// Ignore CrewLink Admin role
	if guild.HasRole(author.ID, "785755227303182349") {
		return
	}

	// Only apply to the main chat channels
	monitorList := []string{
		"705685190860734476",
		"745795131386363934",
		"774907319048994816",
		"774907269031788544",
		"743720183121838171",
		"743742282137862235",
		"748047687353499678",
		"743686586830684180",
	}

	found := false
	for _, c := range monitorList {
		if message.ChannelID == c {
			found = true
			break
		}
	}

	if !found {
		return
	}

	// Define map of strings to quarantine
	detectWords := []string{
		"@here",
		"@everyone",
		"streancommunity",
		"steancommunity",
		"tradeoffer",
		"discord.gg/",
		"discordapp.com/invite/",
		"streammcomunity",
		".ru/",
	}

	// See if the message hits any of the defined strings
	for _, word := range detectWords {
		if strings.Contains(strings.ToLower(message.Content), word) {
			messageDeleted := "true"
			roleApplied := "true"
			log.Infof("Someone has triggered quarantine detection; user: %s#%s (%s); message (%s): %s", author.Username, author.Discriminator, author.ID, message.ID, message.Content)
			// Delete the message
			err := session.ChannelMessageDelete(message.ChannelID, message.ID)
			if err != nil {
				messageDeleted = "false"
				bot.SendErrorReport(message.GuildID, message.ChannelID, author.ID, "Failed to delete message during quarantine handling: "+message.ID, err)
			}

			// Apply the quarantine role
			// Hardcoded, naturally
			err = session.GuildMemberRoleAdd(guild.ID, author.ID, "857037616649470022")
			if err != nil {
				roleApplied = "false"
				bot.SendErrorReport(message.GuildID, message.ChannelID, author.ID, "Failed to apply quarantine role to member: "+author.ID, err)
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
			// Get the member join date
			joinString := message.Member.JoinedAt.String()
			response.PrependField("Role Applied", roleApplied, false)
			response.PrependField("Message Deleted", messageDeleted, false)
			response.PrependField("Channel", "<#"+message.ChannelID+">", false)
			response.PrependField("Message", message.Content, false)
			response.PrependField("Member Join Date", joinString, false)
			response.PrependField("Member ID", author.ID, false)
			response.PrependField("Member", author.Username+"#"+author.Discriminator, false)

			response.Send(true, "Member Quarantined", "A suspected bot has been placed into quarantine. Please review the content of their message and ban if necessary.")

			// Why, yes, I did just hardcode my own ping, why do you ask?
			_, err = session.ChannelMessageSend(guild.Info.ResponseChannelId, "<@211597673713762305> A member has been quarantined. Please review the above report!")
			return
		}
	}
}

func init() {
	bot.AddHandler(messageListener)
}
