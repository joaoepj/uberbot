package deletepolicy

import (
	"strings"

	bot "github.com/ubergeek77/uberbot/core"
)

// The reusable strings for this command.
var failTitle = "Failed to set delete policy"
var successTitle = "Delete Policy updated"

// The information about this command
//var deletepolicyInfo = bot.CommandInfo{
//	Trigger:     "deletepolicy",
//	Description: "Enable or disable the deletion of messages containing commands",
//	Group:       bot.Utility,
//	Usage: []string{
//		"true",
//		"yes",
//		"y",
//		"false",
//		"no",
//		"n",
//	},
//	Public: false,
//}.
var deletePolicyInfo = bot.CreateCommandInfo(
	"deletepolicy",
	"Enable or disable the deletion of messages containing commands",
	false,
	bot.Utility,
)

// The function to run for this command.
func deletepolicy(ctx *bot.Context) {
	// Start a new response
	response := bot.NewResponse(ctx, false, true)

	// This command takes exactly 1 argument
	if len(ctx.Args) != 1 {
		response.Send(false, failTitle, "Invalid syntax")
		return
	}

	// Detect the policy the user wants
	var newPolicy bool
	switch strings.ToLower(ctx.Args["choice"].StringValue()) {
	case "y", "yes", "true":
		newPolicy = true
	case "n", "no", "false":
		newPolicy = false
	default:
		response.Send(false, failTitle, "Invalid syntax")
		return
	}

	// Set the policy
	ctx.Guild.SetDeletePolicy(newPolicy)

	// Send a response based on what happened
	if newPolicy == true {
		response.Send(true, successTitle, "Understood! I'll go ahead and delete messages that trigger commands!")
	} else {
		response.Send(true, successTitle, "OK! I'll stop deleting messages that trigger commands!")
	}
}

// When this package is initialized, add this command to the bot.
func init() {
	deletePolicyInfo.AddArg("choice", bot.String, bot.ArgOption, "choice", true, "")
	deletePolicyInfo.AddChoices("choice", []string{
		"yes",
		"no",
	})
	bot.AddSlashCommand(deletePolicyInfo)
	bot.AddCommand(deletePolicyInfo, deletepolicy)
}
