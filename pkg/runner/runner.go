package runner

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/projectdiscovery/gologger"
	fileutil "github.com/projectdiscovery/utils/file"
)

// Runner is an instance of the dork enumeration
// client used to orchestrate the whole process.
type Runner struct {
	options      *Options
	dorkingAgent *Agent
}

// NewRunner creates a new runner struct instance by parsing
// the configuration options, configuring sources, reading lists
// and setting up loggers, etc.
func NewRunner(options *Options) (*Runner, error) {
	runner := &Runner{options: options}

	// Check if the application loading with any provider configuration, then take it
	// Otherwise load the default provider config
	if !fileutil.FileExists(options.Config) {
		// Create the default configuration file
		file, err := os.Create(options.Config)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}

	gologger.Info().Msgf("Loading config from %s", options.Config)
	source, err := options.UnmarshalFrom(options.Config)
	if err != nil {
		return nil, err
	}
	apiKeys := []string{}
	searchIds := []string{}
	for _, value := range source["google"] {
		keys := strings.Split(value, ":")
		if len(keys) != 2 {
			gologger.Warning().Msgf("Invalid google dork: %s", value)
			continue
		}
		gologger.Debug().Msgf("Adding API key: %s Search ID: %s", keys[0], keys[1])
		apiKeys = append(apiKeys, keys[0])
		searchIds = append(searchIds, keys[1])
	}
	if len(apiKeys) == 0 || len(searchIds) == 0 {
		return nil, errors.New("no valid google dork found")
	}
	agent := NewAgent(apiKeys, searchIds, options.Proxy)
	runner.dorkingAgent = agent

	return runner, nil
}

// RunEnumeration wraps RunEnumerationWithCtx with an empty context
func (r *Runner) RunEnumeration() error {
	if r.options.Dork != nil {
		dorksReader := strings.NewReader(strings.Join(r.options.Dork, "\n"))
		return r.EnumerateDorks(dorksReader, r.options.Output)
	}
	if r.options.DorksFile != "" {
		f, err := os.Open(r.options.DorksFile)
		if err != nil {
			return err
		}
		err = r.EnumerateDorks(f, r.options.Output)
		f.Close()
		return err
	}
	if r.options.Stdin {
		return r.EnumerateDorks(os.Stdin, r.options.Output)
	}
	return nil
}

func (r *Runner) EnumerateDorks(reader io.Reader, writer io.Writer) error {
	var outputWriter io.Writer = writer

	// If OutputFile is specified, create a MultiWriter to write to both file and stdout
	if r.options.OutputFile != "" {
		f, err := os.Create(r.options.OutputFile)
		if err != nil {
			return err
		}
		defer f.Close()
		outputWriter = io.MultiWriter(writer, f)
	}

	var site string
	if r.options.Domain != "" {
		site = " site:" + r.options.Domain
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	// Create error channel to handle errors from goroutines
	errChan := make(chan error, 1)
	// Create results channel with buffer size equal to threads
	resultsChan := make(chan string, r.options.Threads)
	// Create semaphore for rate limiting
	semaphore := make(chan struct{}, r.options.Threads)

	// Start a goroutine to handle writing results
	go func() {
		for result := range resultsChan {
			if _, err := outputWriter.Write([]byte(result + "\n")); err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}
		}
	}()

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		// Check for any errors from goroutines
		select {
		case err := <-errChan:
			return err
		default:
		}

		wg.Add(1)
		dork := scanner.Text() + site

		// Acquire semaphore
		semaphore <- struct{}{}

		go func(dork string) {
			defer wg.Done()
			defer func() { <-semaphore }() // Release semaphore

			results, err := r.dorkingAgent.Dork(dork)
			if err != nil {
				select {
				case errChan <- err:
				default:
				}
				return
			}

			// Send results to channel
			for _, result := range results {
				select {
				case resultsChan <- result:
				case err := <-errChan:
					// If there's an error from another goroutine, propagate it
					select {
					case errChan <- err:
					default:
					}
					return
				}
			}
		}(dork)
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return err
	}

	// Wait for all goroutines to finish
	wg.Wait()
	// Close the results channel
	close(resultsChan)

	// Check if there were any errors
	select {
	case err := <-errChan:
		return err
	default:
		return nil
	}
}
