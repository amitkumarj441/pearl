package cmd

import (
	"io"
	"os"
	"time"

	"github.com/inconshreveable/log15"
	"github.com/mmcloughlin/pearl"
	"github.com/mmcloughlin/pearl/check"
	"github.com/mmcloughlin/pearl/log"
	"github.com/mmcloughlin/pearl/telemetry"
	"github.com/mmcloughlin/pearl/telemetry/expvar"
	"github.com/mmcloughlin/pearl/telemetry/logging"
	"github.com/spf13/cobra"
	"github.com/uber-go/tally"
	"github.com/uber-go/tally/multi"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start a relay server",
	RunE: func(cmd *cobra.Command, args []string) error {
		return serve()
	},
}

var (
	logfile       string
	telemetryAddr string
)

func init() {
	serveCmd.Flags().StringVarP(&logfile, "logfile", "l", "pearl.json", "log file")
	serveCmd.Flags().StringVarP(&telemetryAddr, "telemetry", "t", "localhost:7142", "telemetry address")

	Register(serveCmd.Flags(), cfg, authorities)

	rootCmd.AddCommand(serveCmd)
}

func logger(logfile string) (log.Logger, error) {
	base := log15.New()
	fh, err := log15.FileHandler(logfile, log15.JsonFormat())
	if err != nil {
		return nil, err
	}
	base.SetHandler(log15.MultiHandler(
		log15.LvlFilterHandler(log15.LvlInfo,
			log15.StreamHandler(os.Stdout, log15.TerminalFormat()),
		),
		fh,
	))
	return log.NewLog15(base), nil
}

func metrics(l log.Logger) (tally.Scope, io.Closer) {
	return tally.NewRootScope(tally.ScopeOptions{
		Prefix: "pearl",
		Tags:   map[string]string{},
		CachedReporter: multi.NewMultiCachedReporter(
			expvar.NewReporter(),
			logging.NewReporter(l),
		),
	}, 1*time.Second)
}

func serve() error {
	l, err := logger(logfile)
	if err != nil {
		return err
	}

	scope, closer := metrics(l)
	defer check.Close(l, closer)

	config, err := cfg.Config()
	if err != nil {
		return err
	}

	r, err := pearl.NewRouter(config, scope, l)
	if err != nil {
		return err
	}

	// Start telemetry server.
	go telemetry.Serve(telemetryAddr, l)

	// Report runtime metrics
	go telemetry.ReportRuntime(scope, 10*time.Second)

	// Start serving
	go func() {
		if err := r.Serve(); err != nil {
			log.Err(l, err, "router error")
		}
	}()

	// Publish to directory authorities
	p := &pearl.Publisher{
		Router:      r,
		Interval:    16 * time.Hour,
		Authorities: authorities.Addresses(),
		Logger:      l,
	}
	go p.Start()

	select {}
}
