package core

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/ubergeek77/tinylog"
	"github.com/ubergeek77/uberbot/v2/workers"
)

var (
	VERSION       = "2.0.0"
	ENVIRONMENT   string
	Session       *discordgo.Session
	WorkerManager *workers.WorkerManager
)

func CreateSession(token string) {
	var err error
	Session, err = discordgo.New("Bot " + token)
	if err != nil {
		Log.Fatalf("Unable to create session. %s", err)
		return
	}
	// Set Session variables
	Session.State.MaxMessageCount = messageState
	if IsDevEnv() {
		Session.LogLevel = discordgo.LogInformational
	} else {
		Session.LogLevel = discordgo.LogWarning
	}
	Session.SyncEvents = false
	Session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	// Set a gateway status (is updated via a cronjob)
	Session.Identify.Presence = discordgo.GatewayStatusUpdate{
		Game: discordgo.Activity{
			Name: "things load",
			Type: discordgo.ActivityTypeWatching,
		},
		Status: "dnd",
		AFK:    true,
		Since:  91879201,
	}
	return
}

func CreateManagers(location *time.Location) {
	WorkerManager = workers.InitializeWorkers(location)
}

func InitBot() {
	// Setup environment
	setupEnv()
	// Setup dgoLog introp
	discordgo.Logger = dgoLog

	// setup tinylog level
	if ENVIRONMENT == "prod" {
		Log.LogLevel = tinylog.InfoLevel
	} else {
		Log.LogLevel = tinylog.DebugLevel
	}

	// Load guilds
	// WARN: we are currently using the fs provider, but it will be changeable soon
	currentProvider = initProvider()
	loadGuilds()
	// Current Version
	Log.Infof("Running uberbot version %s. Environment: %s", VERSION, ENVIRONMENT)
	Log.Infof("Connected to database: uberbot")
	// Create Session
	Log.Debug("Creating Session")
	CreateSession(botToken)
	CreateManagers(time.UTC)
}

// Run
// runs the bot.
func Run() {
	// Register the event handlers
	// TODO rewrite handler system
	AddHandler(handleInteraction)
	AddHandler(commandHandler)
	addHandlers()

	// Open up a discordgo session
	err := Session.Open()
	if err != nil {
		Log.Fatalf("Failed to connect to Discord: %s", err)
	}
	// Log that the login succeeded
	Log.Infof("Bot logged in as \"" + Session.State.Ready.User.Username + "#" + Session.State.Ready.User.Discriminator + "\"")

	// Print information about the current bot admins
	numAdmins := 0
	for userId := range botAdmins {
		if user, err := GetUser(userId); err == nil {
			numAdmins += 1
			Log.Infof("Added bot admin: %s#%s", user.Username, user.Discriminator)
		} else {
			Log.Errorf("Unable to lookup bot admin user ID: " + userId)
		}
	}

	if numAdmins == 0 {
		Log.Warning("You have not added any bot admins! Only moderators will be able to run commands, and permissions cannot be changed!")
	}

	// Bot ready
	Log.Info("Initialization complete! The bot is now ready.")
	// -- GRACEFUL TERMINATION -- //

	// Set up a sigterm channel, so we can detect when the application receives a TERM signal
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt, os.Kill)

	// Keep this thread blocked forever, until a TERM signal is received
	<-sigChannel

	Log.Info("Received TERM signal, terminating gracefully.")

	Log.Info("Closing the Discord session...")
	closeErr := Session.Close()
	if closeErr != nil {
		Log.Errorf("An error occurred when closing the Discord session: %s", err)
		return
	}

	Log.Info("Session closed.")
}
