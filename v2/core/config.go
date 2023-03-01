package core

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	tlog "github.com/ubergeek77/tinylog"
)

// messageState
// Tells discordgo the amount of messages to cache.
var messageState = 1000

// Log
// The logger for the core bot.
var Log = tlog.NewTaggedLogger("Core", tlog.NewColor("38;5;111"))

// dlog
// The logger for discordgo.
var dlog = tlog.NewTaggedLogger("DG", tlog.NewColor("38;5;111"))

// BotAdmins
// A list of user IDs that are designated as "Bot Administrators"
// These don't get saved to .json, and must be added programmatically
// They receive some privileges higher than guild moderators
// This is a boolean map, because checking its values is dead simple this way.
var botAdmins = make(map[string]bool)

// BotToken
// A string of the current bot token, usually set by the main method
// Similar to BotAdmins, this isn't saved to .json and is added programmatically.
var botToken = ""

// ColorSuccess
// The color to use for response embeds reporting success.
var ColorSuccess = 0x55F485

// ColorFailure
// The color to use for response embeds reporting failure.
var ColorFailure = 0xF45555

func setupEnv() {
	_ = godotenv.Load()

	// Get the botToken
	token, ok := os.LookupEnv("UBERBOT_TOKEN")
	if !ok {
		Log.Fatalf("You have not specified a bot botToken via the UBERBOT_TOKEN environment variable")
	}
	setToken(token)

	// Get the core admin id(s)
	id, _ := os.LookupEnv("ADMIN_IDS")
	if id != "" {
		for _, admin := range strings.Split(id, ",") {
			// must be a snowflake
			if len(EnsureNumbers(admin)) >= 17 {
				addAdmin(admin)
			}
		}
	}
}

// addAdmin
// A function that allows admins to be added, but not removed.
func addAdmin(adminID string) {
	botAdmins[adminID] = true
}

// setToken
// A function that allows a single token to be added, but not removed.
func setToken(tkn string) {
	botToken = tkn
}

// IsAdmin
// Allow commands to check if a user is an admin or not
// Since botAdmins is a boolean map, if they are not in the map, false is the default.
func IsAdmin(userId string) bool {
	return botAdmins[userId]
}

// dgoLog
// Interop for discordgo to call tinylog.
func dgoLog(msgL, caller int, format string, log ...interface{}) {
	pc, file, line, _ := runtime.Caller(caller)
	files := strings.Split(file, "/")
	file = files[len(files)-1]

	name := runtime.FuncForPC(pc).Name()
	fns := strings.Split(name, ".")
	name = fns[len(fns)-1]
	msg := fmt.Sprintf(format, log...)
	if strings.Contains(msg, "First Packet") {
		return
	}
	switch msgL {
	case discordgo.LogError:
		dlog.Errorf("%s:%d:%s() %s", file, line, name, msg)
	case discordgo.LogWarning:
		dlog.Warningf("%s:%d:%s() %s", file, line, name, msg)
	case discordgo.LogInformational:
		dlog.Infof("%s:%d:%s() %s", file, line, name, msg)
	case discordgo.LogDebug:
		dlog.Debugf("%s:%d:%s() %s", file, line, name, msg)
	}
}
