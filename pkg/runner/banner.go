package runner

import (
	"github.com/projectdiscovery/gologger"
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

// Version is the current version of pagode
const version = `v1.0.0`

// showBanner is used to show the banner to the user
func showBanner() {
	gologger.Print().Msgf("%s\n", banner)
	gologger.Print().Msgf("\t\tprojectdiscovery.io\n\n")
}
