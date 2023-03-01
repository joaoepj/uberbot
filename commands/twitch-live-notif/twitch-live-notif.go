package twitchlivenotif

import (
	"github.com/bwmarrin/discordgo"
	"github.com/nicklaw5/helix"
	tlog "github.com/ubergeek77/tinylog"
	bot "github.com/ubergeek77/uberbot/core"
	"strings"
	"time"
)

var client *helix.Client
var log = tlog.NewTaggedLogger("Twitch", tlog.NewColor("38;5;63"))

var activeGuild = "459941804398215168"
var activeChannel = "879218455742656602"
var activeMessage = "879479852074414120"
var activeEmoji = "659654611090669570"
var activeRole = "624108056766185482"
var liveNowChannel = "564631363949690890"

var g *bot.Guild

// Workaround for how the worker system works; setting this to 60 executes the worker immediately on bot boot.
var clock = 60

func twitchInit() {
	clientID := "[redacted]"
	secret := "[redacted]"

	var err error
	client, err = helix.NewClient(&helix.Options{
		ClientID:     clientID,
		ClientSecret: secret,
	})
	if err != nil {
		return
	}
	resp, err := client.RequestAppAccessToken([]string{})
	if err != nil {
		return
	}

	client.SetAppAccessToken(resp.Data.AccessToken)
}

func liveWorker() {
	// Workers run once per second, so to avoid rate limits, keep a custom count and only execute at a count of 60
	if clock >= 60 {
		clock = 0
	} else {
		clock++
		return
	}

	if client == nil {
		log.Error("client is nil")
		time.Sleep(time.Second * 60)
		return
	}

	if g == nil {
		g = bot.GetGuild(activeGuild)
	}

	resp, err := client.GetStreams(&helix.StreamsParams{
		UserIDs: []string{
			"36026656", // Risch
		},
	})

	if err != nil {
		log.Errorf("Unable to get twitch channel information %v", err)
		return
	}

	liveState, liveErr := g.GetString("rischLive")
	if liveErr != nil {
		liveState = "false"
		g.StoreString("rischLive", liveState)
	}

	if len(resp.Data.Streams) > 0 && resp.Data.Streams[0].Type == "live" {
		if liveState == "false" {
			log.Info("Risch is live!")
			liveState = "true"
			g.StoreString("rischLive", liveState)
			embed := bot.CreateEmbed(12924967, resp.Data.Streams[0].Title, "", []*discordgo.MessageEmbedField{
				bot.CreateField("Game:", resp.Data.Streams[0].GameName, false),
			})

			embed.Author = &discordgo.MessageEmbedAuthor{
				Name:    "Risch",
				URL:     "https://twitch.tv/risch",
				IconURL: "https://cdn.discordapp.com/attachments/832860670932418571/879493245749571614/risch_128px.png",
			}

			embed.Thumbnail = &discordgo.MessageEmbedThumbnail{
				URL: "https://cdn.discordapp.com/attachments/832860670932418571/879493245749571614/risch_128px.png",
			}

			embed.URL = "https://twitch.tv/risch"

			thumbnail := resp.Data.Streams[0].ThumbnailURL
			thumbnail = strings.Replace(thumbnail, "{width}", "488", 1)
			thumbnail = strings.Replace(thumbnail, "{height}", "274", 1)

			embed.Image = &discordgo.MessageEmbedImage{
				URL:      thumbnail,
				ProxyURL: "",
				Width:    0,
				Height:   0,
			}

			ms := &discordgo.MessageSend{
				Content:         "<@&" + activeRole + "> Risch is live on Twitch! Get in here!\n\nhttps://twitch.tv/risch",
				Embed:           embed,
				TTS:             false,
				Components:      nil,
				Files:           nil,
				AllowedMentions: nil,
				Reference:       nil,
				File:            nil,
			}

			msg, err := bot.Session.ChannelMessageSendComplex(liveNowChannel, ms)

			if err != nil {
				log.Errorf("Unable to send message %v", err)
				bot.SendErrorReport(msg.GuildID, liveNowChannel, "", "Unable to send go live notification", err)
				return
			}
		}
	} else {
		liveState = "false"
		g.StoreString("rischLive", liveState)
	}
}

func removeRole(session *discordgo.Session, message *discordgo.MessageReactionRemove) {
	if message.GuildID != activeGuild {
		return
	}

	if message.ChannelID != activeChannel {
		return
	}

	if message.MessageID != activeMessage {
		return
	}

	if message.MessageReaction.Emoji.ID != activeEmoji {
		return
	}

	session.GuildMemberRoleRemove(message.GuildID, message.MessageReaction.UserID, activeRole)
}

func applyRole(session *discordgo.Session, message *discordgo.MessageReactionAdd) {
	if message.GuildID != activeGuild {
		return
	}

	if message.ChannelID != activeChannel {
		return
	}

	if message.MessageID != activeMessage {
		return
	}

	if message.MessageReaction.Emoji.ID != activeEmoji {
		return
	}

	session.GuildMemberRoleAdd(message.GuildID, message.MessageReaction.UserID, activeRole)
}

func init() {
	twitchInit()
	bot.AddHandler(applyRole)
	bot.AddHandler(removeRole)
	bot.AddWorker(liveWorker)
}
