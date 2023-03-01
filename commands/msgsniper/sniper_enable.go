
package msgsniper

import bot "github.com/ubergeek77/uberbot/core"

var sniperEnableInfo = bot.CreateCommandInfo("enable", "enables the sniper, optionally with default settings", false, bot.Moderation)

func subCommandEnable(ctx *bot.Context) {
	response := bot.NewResponse(ctx, false, true)
	if ctx.Args["default"].StringValue() == "defaults" {
		if len(ctx.Guild.Info.GuildBannedWords) == 0 {
			if len(ctx.Guild.Info.BannedWordDetectorChannels) == 0 || len(ctx.Guild.Info.BannedWordDetectorRoles) == 0 {
				response.Send(false, "Sniper", "Please setup the roles and channels for the sniper to snipe!")
				return
			}
			ctx.Guild.BulkAddWords(bot.ReadDefaults("./config/defaultbannedwords.json"))
			if !ctx.Guild.IsSniperEnabled() {
				response.Send(true, "Sniper", "Sniper mode has been enabled and defaults have been set")
				ctx.Guild.SetSniper(true)
				log.Debug("Sniper mode is on")
			} else {
				response.Send(true, "Sniper", "Defaults have been set")
				return
			}
		} else {
			if !ctx.Guild.IsSniperEnabled() {
				response.Send(true, "Sniper", "Sniper mode has been enabled")
				ctx.Guild.SetSniper(true)
				log.Debug("Sniper mode is on")
			} else {
				response.Send(false, "Sniper", "Sniper mode is already enabled")
				return
			}
		}
		return
	}

	if !ctx.Guild.IsSniperEnabled() {
		ctx.Guild.SetSniper(true)
		log.Debug("Sniper mode is on")
		response.Send(true, "Sniper", "Sniper mode has been enabled.")
		return
	}

	response.Send(false, "Sniper", "Sniper mode is already enabled")
	return
}

func init() {
	sniperEnableInfo.SetParent(false, "sniper")
	sniperEnableInfo.AddArg("default", bot.String, bot.ArgOption, "loads defaults", false, "")
	bot.AddChildCommand(sniperEnableInfo, subCommandEnable)
}
