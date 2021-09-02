package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DisgoOrg/disgo/core"
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgo/gateway"
	"github.com/DisgoOrg/disgo/info"
	"github.com/DisgoOrg/log"
)

var (
	token     = os.Getenv("disgo_token")
	publicKey = os.Getenv("disgo_public_key")
	guildID   = discord.Snowflake(os.Getenv("disgo_guild_id"))

	commands = []discord.ApplicationCommandCreate{
		{
			Type:              discord.ApplicationCommandTypeSlash,
			Name:              "say",
			Description:       "says what you say",
			DefaultPermission: true,
			Options: []discord.ApplicationCommandOption{
				{
					Type:        discord.CommandOptionTypeString,
					Name:        "message",
					Description: "What to say",
					Required:    true,
				},
			},
		},
	}
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.Info("starting example...")
	log.Infof("disgo version: %s", info.Version)

	disgo, err := core.NewBuilder(token).
		SetGatewayConfig(gateway.Config{
			GatewayIntents: discord.GatewayIntentsNone,
		}).
		AddEventListeners(&events.ListenerAdapter{
			OnSlashCommand: commandListener,
		}).
		Build()

	if err != nil {
		log.Fatal("error while building disgo instance: ", err)
		return
	}

	defer disgo.Close()

	_, err = disgo.SetGuildCommands(guildID, commands)
	if err != nil {
		log.Fatalf("error while registering commands: %s", err)
	}

	err = disgo.Connect()
	if err != nil {
		log.Fatalf("error while connecting to gateway: %s", err)
	}

	log.Infof("example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-s
}

func commandListener(event *events.SlashCommandEvent) {
	if event.CommandName() == "say" {
		_ = event.Create(core.NewMessageCreateBuilder().
			SetContent(event.Option("message").String()).
			Build(),
		)
	}
}
