package runner

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	fileutil "github.com/projectdiscovery/utils/file"
	folderutil "github.com/projectdiscovery/utils/folder"
	logutil "github.com/projectdiscovery/utils/log"
	"gopkg.in/yaml.v3"
)

var (
	configDir             = folderutil.AppConfigDirOrDefault(".", "pagode")
	defaultConfigLocation = filepath.Join(configDir, "config.yaml")
)

// Options contains the configuration options for tuning
// the google enumeration process.
type Options struct {
	Verbose    bool                // Verbose flag indicates whether to show verbose output or not
	NoColor    bool                // NoColor disables the colored output
	Silent     bool                // Silent suppresses any extra text and only writes results to screen
	Stdin      bool                // Stdin specifies whether stdin input was given to the process
	Version    bool                // Version specifies if we should just show version and exit
	Threads    int                 // Threads controls the number of threads to use for enumerations
	Dork       goflags.StringSlice // Dork is the dork(s) to search with
	DorksFile  string              // DorksFile is the file containing list of dorks to search with
	Domain     string              // Domain to specify in the searches
	Output     io.Writer
	OutputFile string // Output is the file to write results to.
	Proxy      string // proxy URL to use for the requests
	Results    int    // Maximum number of results per search
	Config     string // Config contains the location of the config file
}

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	logutil.DisableDefaultLogger()

	options := &Options{}

	var err error
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`Pagode is a passive google dorking enumeration tool.`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringSliceVarP(&options.Dork, "dork", "g", nil, "dork to search", goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.DorksFile, "list", "gL", "", "file containing a list of dorks to search"),
		flagSet.StringVarP(&options.Domain, "domain", "d", "", `domain to search with "site:" operator`),
	)

	flagSet.CreateGroup("output", "Output",
		flagSet.StringVarP(&options.OutputFile, "output", "o", "", "file to write output to"),
	)

	flagSet.CreateGroup("configuration", "Configuration",
		flagSet.StringVarP(&options.Config, "config", "c", defaultConfigLocation, "config file for API keys"),
		flagSet.StringVarP(&options.Proxy, "proxy", "p", "", "HTTP(s)/SOCKS5 proxy to use with pagode"),
		flagSet.IntVarP(&options.Results, "max", "m", 0, "maximum number of results per search"),
		flagSet.IntVarP(&options.Threads, "threads", "t", 10, "number of concurrent goroutines for querying"),
	)

	flagSet.CreateGroup("debug", "Debug",
		flagSet.BoolVar(&options.Silent, "silent", false, "show only results in output"),
		flagSet.BoolVar(&options.Version, "version", false, "show version of pagode"),
		flagSet.BoolVar(&options.Verbose, "v", false, "show verbose output"),
		flagSet.BoolVarP(&options.NoColor, "no-color", "nc", false, "disable color in output"),
	)

	if err := flagSet.Parse(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Default output is stdout
	options.Output = os.Stdout

	// Check if stdin pipe was given
	options.Stdin = fileutil.HasStdin()

	// Read the inputs and configure the logging
	options.configureOutput()

	if options.Version {
		gologger.Info().Msgf("Current Version: %s\n", version)
		os.Exit(0)
	}

	options.preProcessOptions()

	if !options.Silent {
		showBanner()
	}

	/* We need pagode on pdtm api
	if !options.DisableUpdateCheck {
		latestVersion, err := updateutils.GetToolVersionCallback("pagode", version)()
		if err != nil {
			if options.Verbose {
				gologger.Error().Msgf("pagode version check failed: %v", err.Error())
			}
		} else {
			gologger.Info().Msgf("Current pagode version %v %v", version, updateutils.GetVersionDescription(version, latestVersion))
		}
	}
	*/

	// Validate the options passed by the user and if any
	// invalid options have been used, exit.
	err = options.validateOptions()
	if err != nil {
		gologger.Fatal().Msgf("Program exiting: %s\n", err)
	}

	return options
}

func (options *Options) preProcessOptions() {
	for i, dork := range options.Dork {
		options.Dork[i], _ = sanitize(dork)
	}
}

// UnmarshalFrom reads the marshaled yaml config from disk
func (options *Options) UnmarshalFrom(file string) (map[string][]string, error) {
	reader, err := fileutil.SubstituteConfigFromEnvVars(file)
	if err != nil {
		return nil, err
	}
	sourceApiKeysMap := map[string][]string{}
	err = yaml.NewDecoder(reader).Decode(sourceApiKeysMap)
	if err != nil {
		return nil, err
	}
	return sourceApiKeysMap, nil
}
