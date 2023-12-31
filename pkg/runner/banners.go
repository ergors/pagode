package runner

import (
	"github.com/projectdiscovery/gologger"
	updateutils "github.com/projectdiscovery/utils/update"
)

// Slant
const banner = `
                                 __   
    ____  ____ _____ _____  ____/ /__ 
   / __ \/ __ ` + "`/ __ `" + `/ __ \/ __  / _ \
  / /_/ / /_/ / /_/ / /_/ / /_/ /  __/
 / .___/\__,_/\__, /\____/\__,_/\___/ 
/_/          /____/                   
`

// Name
const ToolName = `pagode`

// Version is the current version of subfinder
const version = `v0.1.0`

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("\t\tprojectdiscovery.io\n\n")
}

// GetUpdateCallback returns a callback function that updates subfinder
func GetUpdateCallback() func() {
	return func() {
		showBanner()
		updateutils.GetUpdateToolCallback("pagode", version)()
	}
}
