//nolint:dupl
package msgsniper

import bot "github.com/ubergeek77/uberbot/core"

var sniperSetInfo = bot.CreateCommandInfo("set", "sets bypass role/channel", false, bot.Moderation)

func subCommandSet(ctx *bot.Context) {
	response := bot.NewResponse(ctx, false, true)
	if ctx.Args["option"].StringValue() == "role" {
		role, _ := ctx.Args["roleid"].RoleValue(bot.Session, ctx.Guild.ID)
		if role != nil {
			ctx.Guild.SetSniperRole(role.ID)
			response.Send(true, "Sniper", "Successfully added"+role.Mention()+"to the bypass list")
			return
		}
		response.Send(false, "Sniper", "ID was not a valid role")
	} else if ctx.Args["option"].StringValue() == "channel" {
		channel, _ := ctx.Args["channelid"].ChannelValue(bot.Session)
		if channel != nil {
			ctx.Guild.SetSniperChannel(channel.ID)
			response.Send(true, "Sniper", "Successfully added"+channel.Mention()+"to the watch list")
			return
		}
		response.Send(false, "Sniper", "ID was not a valid channel")
		return
	}
	response.Send(false, "Sniper", sniperFail)
	return
}

func init() {
	sniperSetInfo.SetParent(false, "sniper")
	sniperSetInfo.AddArg("option", bot.SubCmd, bot.ArgOption, "role or channel to set bypass", true, "")
	sniperSetInfo.AddChoices("option", []string{
		"role",
		"channel",
	})
	sniperSetInfo.AddArg("roleid", bot.Role, bot.ArgOption, "no description", true, "")
	sniperSetInfo.AddArg("channelid", bot.Channel, bot.ArgOption, "no description", true, "")

	bot.AddChildCommand(sniperSetInfo, subCommandSet)
}
