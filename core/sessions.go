package core

import (
	"strconv"

	"github.com/InjectionSoftwareandSecurityLLC/lupo/server"
	"github.com/desertbit/grumble"
	"github.com/fatih/color"
)

type session struct {
	id      int
	implant implant
	rhost   string
	checkin string
	status  string
}

var activeSession = -1

// SessionAppConfig - Primary session nested grumble CLI config construction
var SessionAppConfig = &grumble.Config{
	Name:                  "session",
	Description:           "Interactive Session CLI",
	HistoryFile:           "/tmp/lupo.log",
	Prompt:                "lupo session " + strconv.Itoa(activeSession) + " ☾ ",
	PromptColor:           color.New(color.FgGreen, color.Bold),
	HelpHeadlineColor:     color.New(color.FgWhite),
	HelpHeadlineUnderline: true,
	HelpSubCommands:       true,
}

// InitializeSession - Initialize the nested session CLI arguments
func InitializeSession(sessionApp *grumble.App, activeSession int) {

	backCmd := &grumble.Command{
		Name:     "back",
		Help:     "go back to core lupo cli (or use the exit command)",
		LongHelp: "Exit interactive session cli and return to lupo cli (The 'exit' command is an optional built-in for to go back as well) ",
		Run: func(c *grumble.Context) error {
			activeSession = -1

			server.CMD = ""

			sessionApp.Close()

			return nil
		},
	}
	sessionApp.AddCommand(backCmd)

	sessionSwitchCmd := &grumble.Command{
		Name:     "session",
		Help:     "switch to session id",
		LongHelp: "Interact with a different available session by specifying the Session ID",
		Args: func(a *grumble.Args) {
			a.Int("id", "Session ID to interact with")
		},
		Run: func(c *grumble.Context) error {
			activeSession = c.Args.Int("id")
			sessionApp.SetPrompt("lupo session " + strconv.Itoa(activeSession) + " ☾ ")

			return nil
		},
	}

	sessionApp.AddCommand(sessionSwitchCmd)

	sessionCMDCmd := &grumble.Command{
		Name:     "cmd",
		Help:     "execute command on session",
		LongHelp: "Executes a standard OS command that the implant for the current session will execute.",
		Args: func(a *grumble.Args) {
			a.String("cmd", "OS Command to be executed by the target session")
		},
		Run: func(c *grumble.Context) error {

			server.CMD = c.Args.String("cmd")

			return nil
		},
	}

	sessionApp.AddCommand(sessionCMDCmd)
}
