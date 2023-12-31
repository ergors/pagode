package runner

import (
	"fmt"
	"io"
	"os"

	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	fileutil "github.com/projectdiscovery/utils/file"
	logutil "github.com/projectdiscovery/utils/log"
	updateutils "github.com/projectdiscovery/utils/update"
)

// Options contains the configuration options for tuning
// the google enumeration process.
type Options struct {
	Verbose            bool                // Verbose flag indicates whether to show verbose output or not
	NoColor            bool                // NoColor disables the colored output
	JSON               bool                // JSON specifies whether to use json for output format or text file
	Silent             bool                // Silent suppresses any extra text and only writes results to screen
	Stdin              bool                // Stdin specifies whether stdin input was given to the process
	Version            bool                // Version specifies if we should just show version and exit
	GHDB               bool                // Fetch Google Hacking Database
	Threads            int                 // Threads controls the number of threads to use for enumerations
	Dork               goflags.StringSlice // Dork is the dork(s) to search with
	DorksFile          string              // DorksFile is the file containing list of dorks to search with
	Domain             string              // Domain is the domain to specify in the searches
	Output             io.Writer
	OutputFile         string              // Output is the file to write found subdomains to.
	Proxy              goflags.StringSlice // HTTP/HTTPS/SOCKS5 proxy(s)
	ProxyFile          string              // File containing list of proxies to use
	Interval           int                 // Seconds to wait between dork searches
	Results            int                 // Maximum number of results per search
	DisableUpdateCheck bool                // DisableUpdateCheck disable update checking
}

// ParseOptions parses the command line flags provided by a user
func ParseOptions() *Options {
	logutil.DisableDefaultLogger()

	options := &Options{}

	var err error
	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`Pagode is a google dorking enumeration tool.`)

	flagSet.CreateGroup("input", "Input",
		flagSet.StringSliceVarP(&options.Dork, "dork", "g", nil, "dork to search with", goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.DorksFile, "list", "gL", "", "file containing list of dorks to search with"),
		flagSet.StringVarP(&options.Domain, "domain", "d", "", `domain to search with "site:" dork`),
	)

	flagSet.CreateGroup("rate-limit", "Rate-limit",
		flagSet.IntVarP(&options.Interval, "interval", "i", 1, "seconds to wait between searches"),
		flagSet.IntVarP(&options.Results, "results", "r", 100, "maximum number of results per search"),
	)

	flagSet.CreateGroup("update", "Update",
		flagSet.CallbackVarP(GetUpdateCallback(), "update", "up", "update subfinder to latest version"),
		flagSet.BoolVarP(&options.DisableUpdateCheck, "disable-update-check", "duc", false, "disable automatic subfinder update check"),
	)

	flagSet.CreateGroup("output", "Output",
		flagSet.StringVarP(&options.OutputFile, "output", "o", "", "file to write output to"),
		flagSet.BoolVarP(&options.JSON, "json", "oJ", false, "write output in JSONL(ines) format"),
	)

	flagSet.CreateGroup("configuration", "Configuration",
		flagSet.StringSliceVarP(&options.Proxy, "proxy", "p", nil, "HTTP(s)/SOCKS5 proxy to use with pagode", goflags.NormalizedStringSliceOptions),
		flagSet.StringVarP(&options.ProxyFile, "plist", "pL", "", "file containing list of proxies to use"),
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

	// We need pagode on pdtm api
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
