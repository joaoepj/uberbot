package sendmessage

import (
	"github.com/bwmarrin/discordgo"
	bot "github.com/ubergeek77/uberbot/core"
	"strings"
)

var sendMessageInfo = bot.CreateCommandInfo("sendmsg", "sends a message, only works in dms", false, bot.Utility)

func sendMessage(ctx *bot.Context) {
	if ch, err := bot.Session.State.Channel(ctx.Message.ChannelID); err == nil {
		if ch.Type != discordgo.ChannelTypeDM {
			return
		}
		if ctx.Args["channelID"].StringValue() == "" && ctx.Args["url"].StringValue() == "" {
			return
		}
	}

	if ctx.Args["url"].StringValue() != "" {
		urlStr := strings.TrimPrefix(ctx.Args["url"].StringValue(), "https://canary.discord.com/channels/")
		splitStr := strings.Split(urlStr, "/")
		if len(splitStr) < 3 {
			return
		}
		_, err := bot.Session.ChannelMessageSendReply(splitStr[1], ctx.Args["msg"].StringValue(), &discordgo.MessageReference{
			MessageID: splitStr[2],
			ChannelID: splitStr[1],
			GuildID:   splitStr[0],
		})
		if err != nil {
			return
		}
		return
	}
	_, err := bot.Session.ChannelMessageSend(ctx.Args["channelID"].StringValue(), ctx.Args["msg"].StringValue())
	if err != nil {
		return
	}
}

func init() {
	sendMessageInfo.AddArg("msg", bot.String, bot.ArgContent, "yep", true, "")
	sendMessageInfo.AddFlagArg("channelID", bot.Channel, bot.ArgOption, "yep", false, "")
	sendMessageInfo.AddFlagArg("url", bot.Message, bot.ArgOption, "yep", false, "")
	bot.AddCommand(sendMessageInfo, sendMessage)
}
