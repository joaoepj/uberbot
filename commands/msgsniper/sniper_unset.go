//nolint:dupl
package msgsniper

import bot "github.com/ubergeek77/uberbot/core"

var sniperUnsetInfo = bot.CreateCommandInfo("unset", "unsets the bypass role/channel for the sniper", false, bot.Moderation)

func subCommandUnset(ctx *bot.Context) {
	response := bot.NewResponse(ctx, false, true)
	if ctx.Args["option"].StringValue() == "role" {
		role, _ := ctx.Args["roleid"].RoleValue(bot.Session, ctx.Guild.ID)
		if role != nil {
			ctx.Guild.UnsetSniperRole(role.ID)
			response.Send(true, "Sniper", "Successfully removed"+role.Mention()+"from the bypass list")
			return
		}
		response.Send(false, "Sniper", "ID was not a valid role")
	} else if ctx.Args["option"].StringValue() == "channel" {
		channel, _ := ctx.Args["channelid"].ChannelValue(bot.Session)
		if channel != nil {
			ctx.Guild.UnsetSniperChannel(channel.ID)
			response.Send(true, "Sniper", "Successfully removed"+channel.Mention()+"from the watch list")
			return
		}
		response.Send(false, "Sniper", "ID was not a valid channel")
		return
	}
	response.Send(false, "Sniper", sniperFail)
	return
}

func init() {
	sniperUnsetInfo.SetParent(false, "sniper")
	sniperUnsetInfo.AddArg("roleid", bot.Role, bot.ArgOption, "role to unset", false, "")
	sniperUnsetInfo.AddArg("channelid", bot.Channel, bot.ArgOption, "channel to unset", false, "")
	bot.AddChildCommand(sniperUnsetInfo, subCommandUnset)
}
