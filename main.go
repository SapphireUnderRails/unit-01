package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Creating a struct to hold the two tokens.
type Tokens struct {
	DiscordToken string
	GPT3Token    string
}

// Creating a struct that will hold all the GPT3 parameters.
type Parameters struct {
	Chance float64
	Length int64
}

// Globalizing the structs that hold this important data.
var tokens Tokens
var parameters Parameters

// Main functions.
func main() {

	// Retrieve the tokens from the tokens.json file.
	tokensFile, err := os.ReadFile("tokens.json")
	if err != nil {
		log.Fatal("COULD NOT READ 'tokens.json' FILE: ", err)
	}

	// Unmarshal the tokens from tokensFile.
	json.Unmarshal(tokensFile, &tokens)

	// Retrieve the parameters from the GPT3Parameters.json file.
	parametersFile, err := os.ReadFile("parameters.json")
	if err != nil {
		log.Fatal("COULD NOT READ 'parameters.json' FILE: ", err)
	}

	// Unmarshal the tokens from the gp3ParametersFile.
	json.Unmarshal(parametersFile, &parameters)

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + tokens.DiscordToken)
	if err != nil {
		log.Fatal("ERROR CREATING DISCORD SESSION:", err)
	}

	// Identify that we want all intents.
	session.Identify.Intents = discordgo.IntentsAll

	// Now we open a websocket connection to Discord and begin listening.
	err = session.Open()
	if err != nil {
		log.Fatal("ERROR OPENING WEBSOCKET:", err)
	}

	// Making a map of registered commands.
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	// Looping through the commands array and registering them.
	// https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.ApplicationCommandCreate
	for i, command := range commands {
		registered_command, err := session.ApplicationCommandCreate(session.State.User.ID, "1001077854936760352", command)
		if err != nil {
			log.Panicf("CANNOT CREATE '%v' COMMAND: %v", command.Name, err)
		}
		registeredCommands[i] = registered_command
	}

	// Looping through the array of handlers and adding them to the session.
	session.AddHandler(func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		if handler, ok := commandHandlers[interaction.ApplicationCommandData().Name]; ok {
			handler(session, interaction)
		}
	})

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Lopping through the registeredCommands array and deleting all the commands.
	for _, v := range registeredCommands {
		err := session.ApplicationCommandDelete(session.State.User.ID, "1001077854936760352", v.ID)
		if err != nil {
			log.Panicf("CANNOT DELETE '%v' COMMAND: %v", v.Name, err)
		}
	}

	// Cleanly close down the Discord session.
	session.Close()
}

// Decalaring default member permission.
var defaultMemberPermissions int64 = discordgo.PermissionManageServer

// Declaring min and max values of the chance command option.
var minChanceValue float64 = 0
var maxChanceValue float64 = 100

// Declaring the max value allowed for a response.
var minLengthValue float64 = 60
var maxLengthValue float64 = 512

// Creating an array of commands to register.
//https://pkg.go.dev/github.com/bwmarrin/discordgo#ApplicationCommand
var commands = []*discordgo.ApplicationCommand{
	{
		Name:                     "test",
		Description:              "This is just a test command!",
		DefaultMemberPermissions: &defaultMemberPermissions,
	},
	{
		Name:                     "get_chance",
		Description:              "This returns the value of the chance that Shem-Ha will respond to a message.",
		DefaultMemberPermissions: &defaultMemberPermissions,
	},
	{
		Name:                     "get_length",
		Description:              "This returns the maximum length of a response from Shem-Ha in tokens. A token is about 4 characters.",
		DefaultMemberPermissions: &defaultMemberPermissions,
	},
	{
		Name:                     "set_chance",
		Description:              "This sets the value of the chance that Shem-Ha will respond to a message.",
		DefaultMemberPermissions: &defaultMemberPermissions,
		// Registering the option available for this command.
		// https://pkg.go.dev/github.com/bwmarrin/discordgo#ApplicationCommandOption
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionNumber,
				Name:        "percentage",
				Description: "This value is the chance that Shem-Ha will respond to a message, must be between 0 and 100.",
				Required:    true,
				MinValue:    &minChanceValue,
				MaxValue:    maxChanceValue,
			},
		},
	},
	{
		Name:                     "set_length",
		Description:              "This sets the maximum length of a response from Shem-Ha in tokens. A token is about 4 characters.",
		DefaultMemberPermissions: &defaultMemberPermissions,
		// Registering the option available for this command.
		// https://pkg.go.dev/github.com/bwmarrin/discordgo#ApplicationCommandOption
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "percentage",
				Description: "This is the maximum response length in tokens. A token is about 4 characters.",
				Required:    true,
				MinValue:    &minLengthValue,
				MaxValue:    maxLengthValue,
			},
		},
	},
}

// Creating a map of event handlers to respond to application commands.
// https://pkg.go.dev/github.com/bwmarrin/discordgo#EventHandler
var commandHandlers = map[string]func(session *discordgo.Session, interaction *discordgo.InteractionCreate){
	"test": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Responding to the interaction.
		//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Congrats on using the test command!",
			},
		})
	},
	"get_chance": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Responding to the interaction.
		//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The current response chance is %v percent.", parameters.Chance),
			},
		})
	},
	"get_length": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		// Responding to the interaction.
		//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("The current response length is %v tokens.", parameters.Length),
			},
		})
	},
	"set_chance": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		parameters.Chance = interaction.ApplicationCommandData().Options[0].FloatValue()

		// Marshall the new parameters to save.
		jsonBytes, err := json.Marshal(parameters)
		if err != nil {
			log.Panicln("ERROR MARSHALING JSON: ", err)

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("FAILED TO UPDATE CHANCE: %v", err),
				},
			})

			return
		}

		// Save updated parameters into parameters.json.
		err = os.WriteFile("parameters.json", jsonBytes, 0644)
		if err != nil {
			log.Panicln("ERROR SAVING JSON: ", err)

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("FAILED TO UPDATE CHANCE: %v", err),
				},
			})

			return
		}

		// Responding to the interaction.
		//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Successfully updated the response chance. The reponse chance is now %v percent.", parameters.Chance),
			},
		})
	},
	"set_length": func(session *discordgo.Session, interaction *discordgo.InteractionCreate) {
		parameters.Length = interaction.ApplicationCommandData().Options[0].IntValue()

		// Marshall the new parameters to save.
		jsonBytes, err := json.Marshal(parameters)
		if err != nil {
			log.Panicln("ERROR MARSHALING JSON: ", err)

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("FAILED TO UPDATE LENGTH: %v", err),
				},
			})

			return
		}

		// Save updated parameters into parameters.json.
		err = os.WriteFile("parameters.json", jsonBytes, 0644)
		if err != nil {
			log.Panicln("ERROR SAVING JSON: ", err)

			session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("FAILED TO UPDATE LENGTH: %v", err),
				},
			})

			return
		}

		// Responding to the interaction.
		//https://pkg.go.dev/github.com/bwmarrin/discordgo#Session.InteractionRespond
		session.InteractionRespond(interaction.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Successfully updated the response length. The reponse length is now %v tokens.", parameters.Length),
			},
		})
	},
}
