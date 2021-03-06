package main

import (
	"flag"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	maxConcurrency     = 1
	maxErrors      int = 1e3
	maxVisited     int = 1e3
	defaultWait        = 3 * time.Second
	defaultDelay       = 1 * time.Second
)

var (
	concurrency = flag.Int("concurrency-level", maxConcurrency, "")
	// flags
	// internal
	// version
	whatVersion = flag.Bool("version", false, "")
	v           = flag.Bool("v", false, "")
	// license
	whatLicense = flag.Bool("license", false, "")
	l           = flag.Bool("l", false, "")
	debug       = flag.Bool("debug", false, "")
	d           = flag.Bool("d", false, "")

	// general
	watchHREF      = flag.Bool("watch-href", false, "")
	watchSRC       = flag.Bool("watch-src", false, "")
	rePattern      = flag.String("watch-pattern", "", "")
	spanHosts      = flag.Bool("span-hosts", false, "")
	spanSubdomains = flag.Bool("span-subdomains", false, "")
	outJSON        = flag.Bool("json", false, "")
	j              = flag.Bool("j", false, "")
	// what things you care about
	check5xx = flag.Bool("check-server-errors", false, "")
	check4xx = flag.Bool("check-client-errors", false, "")
	check3xx = flag.Bool("check-redirection", false, "")
	// set limits
	maxErrsCount    = flag.Int("max-errors-count", maxErrors, "")
	timeWait        = flag.Duration("time-wait", defaultWait, "")
	timeDelay       = flag.Duration("time-delay", defaultDelay, "")
	maxVisitedCount = flag.Int("max-visited-count", maxVisited, "")
	//
	website,
	hostName string
	re *regexp.Regexp
)

func init() {
	if isROOT() {
		exit(2, "Running as ROOT isn't your worst mistake, is it!!")
	}

	flag.Usage = usage
	flag.Parse()

	setupCmd()
}

func setupCmd() {
	args := flag.Args()

	if *l || *whatLicense {
		showLicense()
	}

	if *v || *whatVersion {
		showVersion()
	}

	if len(args) > 1 {
		exit(1, "Error: URL cannot be given more than once, for usage see: -h | --help")
	} else if len(args) == 0 {
		exit(1, "Error: URL to be checked is missing, for usage see: -h | --help")
	} else {
		parsed, err := url.Parse(args[0])
		if err != nil {
			exit(1, "Error: URL given is not parsable, for usage see: -h | --help")
		}

		if parsed.Host == "" {
			exit(1, "Error: URL host is missing, for usage see: -h | --help")
		}

		if parsed.Scheme == "" {
			parsed.Scheme = "http"
		}
		website = parsed.String()
		hostName = parsed.Host
	}

	if !*watchHREF && !*watchSRC {
		exit(1, "Error: Nothing to 'watch', for usage see: -h | --help")
	}

	if !*check5xx && !*check4xx && !*check3xx {
		exit(1, "Error: Nothing to 'check', for usage see: -h | --help")
	}

	if *rePattern != "" {
		if strings.Contains(*rePattern, `/`) {
			exit(1, "Error: regexp file pattern should not contain a slash '/'")
		}
		r, err := regexp.Compile(*rePattern)
		if err != nil {
			exit(1, "Error: Nasty regexp pattern failed to compile", err)
		}
		re = r
	}

	if *concurrency <= 0 {
		*concurrency = maxConcurrency
	}

	if *maxErrsCount <= 0 {
		*maxErrsCount = maxErrors
	}

	if *d || *debug {
		*debug = true
	}

	if *j || *outJSON {
		*outJSON = true
	}

}

func usage() {
	const tpl = `

Usage: %s [-v | --version] [-h | --help] [-l | --license] [-vvv | --verbose] [--watch-href] [--watch-src] [--watch-pattern regexp] [--span-hosts] [--span-subdomains] [-j | --json] [--check-server-errors] [--check-client-errors] [--check-redirection] [--max-errors-count NUM] [--max-visited-count NUM] [--time-wait DURATION] [--time-delay DURATION]   URL


FLAGS:
 -v | --version
    Show version and exit.
 -l | --license
    Show License and exit.
 -h | --help
    Show help and exit.
 -d | --debug
    Display whats going on.

 --check-server-errors
    Check for HTTP 5xx servers errors (default: false)
 --check-client-errors
    Check for HTTP 4xx servers errors (default: false)
 --check-redirection
    Check for HTTP 3xx servers responses, redirection. (default: false)

 --watch-href
    Watch 'href' attributes URL (default: false)
 --watch-src
    Watch 'src' attributes URL (default: false)
 --span-hosts
    Follow links hosted on other websites, including subdomains (default: false)
 --span-subdomains
    Follow links hosted on other subdomains (default: false)

 -j | --json
    Display check results as JSON (default: false)


OPTIONS:
 --watch-pattern regexp
    Regular expression pattern to match filename against, if URL doesn't match fetching will be skipped (default: '')
 --concurrency-level num
    Number of concurrent requests to be performed at once (default: %d)
 --max-errors-count NUM
    Max errors count before exit (default: %d)
 --max-visited-count NUM
    Max visited URLs remembered, higher values likely will prevent checking already checked URLs yet requires more memory (RAM) (default: %d)
 --time-wait DURATION
    Duration of time of inactivity before the program exits (default: %v)
 --time-delay DURATION
    Duration of time to delay the next URL fetching (default: %v)



AURGUMENTS:
 URL
    The website's URL to be checked

`
	exitf(1,
		tpl,
		os.Args[0],
		maxConcurrency,
		maxErrors,
		maxVisited,
		defaultWait,
		defaultDelay,
	)
}
