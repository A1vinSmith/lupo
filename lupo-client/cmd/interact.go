package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/InjectionSoftwareandSecurityLLC/lupo/lupo-client/core"
	"github.com/desertbit/grumble"
)

// activeSession - Active session that is being interacted with by the user
//
// This data is supplied as a parameter when switching sessions with either the "interact" command or "session" sub-shell
var activeSession int

// init - Initializes the primary "interact" grumble command
//
// "interact" accepts an argument of "id" that is used to generate a new SessionApp with the SessionAppConfig
//
//  "interact" subcommands include:
//
//  	"show" - Shows all registered sessions. Accepts andargument of "id" that can be used to show a specific session based on the id.
//
//  	"kill" - Accepts an argument of "id" that is used to de-register a session.
//
//  	"clean" - De-registers all sessions marked as "DEAD" based on a pre-determined "Check-In" update interval.

func init() {

	interactCmd := &grumble.Command{
		Name:     "interact",
		Help:     "interact with a session",
		LongHelp: "Interact with an available session by specifying the Session ID",
		Args: func(a *grumble.Args) {
			a.Int("id", "Session ID to interact with")
		},
		Run: func(c *grumble.Context) error {

			activeSession = c.Args.Int("id")

			// Exec interact with server goes here to switch sessions

			reqString := "&command="
			commandString := "interact " + strconv.Itoa(activeSession)

			reqString = core.AuthURL + reqString + url.QueryEscape(commandString)

			resp, err := core.WolfPackHTTP.Get(reqString)

			if err != nil {
				fmt.Println(err)
				return nil
			}

			defer resp.Body.Close()

			jsonData, err := ioutil.ReadAll(resp.Body)

			if err != nil {
				fmt.Println(err)
				return nil
			}

			// Parse the JSON response
			// We are expecting a JSON string with the key "response" by default, the value is just a raw string response that can be printed to the output
			var coreResponse map[string]interface{}
			err = json.Unmarshal(jsonData, &coreResponse)

			fmt.Println(coreResponse)
			if err != nil {
				//fmt.Println(err)
				return nil
			}

			if coreResponse["response"].(string) == "true" {
				App = grumble.New(SessionAppConfig)
				App.SetPrompt("lupo session " + strconv.Itoa(activeSession) + " ☾ ")
				InitializeSessionCLI(App, activeSession)

				grumble.Main(App)

			} else {

				errorMessage := "Session " + strconv.Itoa(activeSession) + " does not exist"

				return errors.New(errorMessage)

			}
			return nil

		},
	}
	App.AddCommand(interactCmd)

	showCmd := &grumble.Command{
		Name:     "show",
		Help:     "show all sessions",
		LongHelp: "Show all available session information",
		Args: func(a *grumble.Args) {
			a.Int("id", "Filter on session id", grumble.Default(-1))
		},
		Run: func(c *grumble.Context) error {

			//filterID := c.Args.Int("id")

			// Exec interact with server goes here to get a list of current sessions

			table := tabwriter.NewWriter(os.Stdout, 0, 2, 2, ' ', 0)
			fmt.Fprintf(table, "ID\tRemote Host\tArch\tProtocol\tLast Check In\tUpdate Interval\tStatus\t\n")
			fmt.Fprintf(table, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t\n",
				strings.Repeat("=", len("ID")),
				strings.Repeat("=", len("Remote Host")),
				strings.Repeat("=", len("Arch")),
				strings.Repeat("=", len("Protocol")),
				strings.Repeat("=", len("Last Check In")),
				strings.Repeat("=", len("Update Interval")),
				strings.Repeat("=", len("Status")))

			// Populate this based on returned data
			/*
				if filterID != -1 {

					_, sessionExists := core.Sessions[filterID]

					if !sessionExists {

						errorMessage := "cannot filter show on session " + strconv.Itoa(activeSession) + " because the session does not exist"

						return errors.New(errorMessage)
					}

					updateInterval := core.Sessions[filterID].Implant.Update
					lastCheckIn := core.Sessions[filterID].RawCheckin

					status, err := calculateSessionStatus(updateInterval, lastCheckIn)

					var textStatus string

					if err != nil {
						textStatus = "UNKNOWN"
						core.SessionStatusUpdate(filterID, "UNKNOWN")
					} else if status {
						textStatus = core.GreenColorIns("ALIVE")
						core.SessionStatusUpdate(filterID, "ALIVE")
					} else if !status {
						textStatus = core.RedColorIns("DEAD")
						core.SessionStatusUpdate(filterID, "DEAD")
					} else {
						textStatus = core.ErrorColorBoldIns("ERROR")
						core.SessionStatusUpdate(filterID, "ERROR")
					}

					fmt.Fprintf(table, "%s\t%s\t%s\t%s\t%s\t%f\t%s\t\n",
						strconv.Itoa(core.Sessions[filterID].ID),
						core.Sessions[filterID].Rhost,
						core.Sessions[filterID].Implant.Arch,
						core.Sessions[filterID].Protocol,
						core.Sessions[filterID].Checkin,
						core.Sessions[filterID].Implant.Update,
						textStatus)

				} else {
					core.LogData(operator + " executed: interact show")

					for i := range core.Sessions {

						updateInterval := core.Sessions[i].Implant.Update
						lastCheckIn := core.Sessions[i].RawCheckin

						status, err := calculateSessionStatus(updateInterval, lastCheckIn)

						var textStatus string

						if err != nil {
							textStatus = "UNKNOWN"
							core.SessionStatusUpdate(i, "UNKNOWN")
						} else if status {
							textStatus = core.GreenColorIns("ALIVE")
							core.SessionStatusUpdate(i, "ALIVE")
						} else if !status {
							textStatus = core.RedColorIns("DEAD")
							core.SessionStatusUpdate(i, "DEAD")
						} else {
							textStatus = core.ErrorColorBoldIns("ERROR")
							core.SessionStatusUpdate(i, "ERROR")
						}

						fmt.Fprintf(table, "%s\t%s\t%s\t%s\t%s\t%f\t%s\t\n",
							strconv.Itoa(core.Sessions[i].ID),
							core.Sessions[i].Rhost,
							core.Sessions[i].Implant.Arch,
							core.Sessions[i].Protocol,
							core.Sessions[i].Checkin,
							core.Sessions[i].Implant.Update,
							textStatus)
					}
				}

			*/

			table.Flush()

			return nil
		},
	}
	interactCmd.AddCommand(showCmd)

	killCmd := &grumble.Command{
		Name:     "kill",
		Help:     "kills a specified session",
		LongHelp: "Kills a session with a specified ID",
		Args: func(a *grumble.Args) {
			a.Int("id", "Session ID to kill")
		},
		Run: func(c *grumble.Context) error {

			//id := c.Args.Int("id")

			// Exec command on server to return sessions
			/*
				delete(core.Sessions, id)

				core.WarningColorBold.Println("Session " + strconv.Itoa(id) + " has been terminated...")
			*/

			return nil
		},
	}
	interactCmd.AddCommand(killCmd)

	cleanCmd := &grumble.Command{
		Name:     "clean",
		Help:     "cleans all sessions marked as DEAD",
		LongHelp: "Kills all sessions marked as DEAD to clear up the session list.",
		Run: func(c *grumble.Context) error {

			// Exec to get sessions

			/*

				for i := range core.Sessions {

					sessionStatus := core.Sessions[i].Status

					if sessionStatus == "DEAD" {
						delete(core.Sessions, i)
						core.WarningColorBold.Println("Session " + strconv.Itoa(i) + " has been terminated...")
					}

				}
			*/

			return nil
		},
	}

	interactCmd.AddCommand(cleanCmd)

}
