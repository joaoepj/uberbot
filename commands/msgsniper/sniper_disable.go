
package msgsniper

import bot "github.com/ubergeek77/uberbot/core"

var sniperDisableInfo = bot.CreateCommandInfo("disable", "disables the sniper", false, bot.Moderation)

func subCommandDisable(ctx *bot.Context) {
	response := bot.NewResponse(ctx, false, true)
	if ctx.Guild.IsSniperEnabled() {
		ctx.Guild.SetSniper(false)
		log.Debug("Sniper mode is off")
		response.Send(true, "Sniper", "Sniper mode has been disabled")
		return
	}
	response.Send(false, "Sniper", "Sniper mode is already disabled")
	return
}

func init() {
	sniperDisableInfo.SetParent(false, "sniper")
	bot.AddChildCommand(sniperDisableInfo, subCommandDisable)
}
