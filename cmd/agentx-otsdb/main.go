package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/context"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/mathpl/agentx/pkg/master"
	"github.com/mathpl/go-tsdmetrics"
	"github.com/rcrowley/go-metrics"
)

func readConf(filename string) (*Conf, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	conf := &Conf{
		Freq: 5,
	}
	_, err = toml.DecodeReader(f, conf)
	if err != nil {
		return nil, err
	}

	return conf, nil
}

func main() {
	app := cli.NewApp()
	app.Name = "agentx-otsdb"
	app.Usage = "Agentx master to OpenTSDB metrics."
	//app.Version = fmt.Sprintf("%s+%s", Version, GitCommit)

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:   "verbosity",
			Value:  4,
			Usage:  "Verbosity (0-5)",
			EnvVar: "VERBOSITY",
		},
		cli.StringFlag{
			Name:   "metrics-report-url, r",
			Usage:  "Where to send OpenTSDB metrics.",
			EnvVar: "METRICS_REPORT_URL",
		},
		cli.StringFlag{
			Name:   "config",
			Value:  "/etc/agentx-otsdb.toml",
			Usage:  "Config file",
			EnvVar: "CONFIG",
		},
		cli.StringFlag{
			Name:   "socket-path",
			Value:  "/var/agentx/master",
			Usage:  "Path to socket to manage",
			EnvVar: "SOCKET_PATH",
		},
	}

	rand.Seed(time.Now().Unix())

	ctx, cancel := context.WithCancel(context.Background())

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info("Received interupt, shutting down...")
		cancel()
		return
	}()

	wg := &sync.WaitGroup{}

	app.Action = func(c *cli.Context) error {
		log.SetLevel(log.Level(c.Int("verbosity")))

		conf, err := readConf(c.String("config"))
		if err != nil {
			return fmt.Errorf("Unable to parse config file: %s", err)
		}

		os.Remove(c.String("socket-path"))

		m, err := master.NewMasterAgent("unix", c.String("socket-path"))
		if err != nil {
			return fmt.Errorf("Unable to start master agent: %s", err)
		}

		go func() {
			wg.Add(1)
			defer wg.Done()

			err = m.Run(ctx)
			if err != nil {
				log.Errorf("Master agent runtime error: %s", err)
			}
			cancel()
		}()

		collectMibs := func(r tsdmetrics.TaggedRegistry) {
			var oids []string
			oidToMetric := make(map[string]MIBMetric)
			for _, mib := range conf.MIBS {
				for _, m := range mib.Metrics {
					oid := mib.BaseOid + m.Oid
					oids = append(oids, oid)
					oidToMetric[oid] = m
				}
			}

			variables, err := m.Get(ctx, oids)
			if err != nil {
				log.Error(err)
				return
			}

			for _, variable := range variables {
				oid := variable.Name.String()

				m, found := oidToMetric[oid]
				if !found {
					log.Errorf("Unexpected oid returned: %s", oid)
					continue
				}

				tags, err := tsdmetrics.TagsFromString(m.Tags)
				if err != nil {
					log.Error(err)
					continue
				}

				switch v := variable.Value.(type) {
				case uint32:
					g := metrics.NewGauge()
					r.Add(m.Metric, tags, g)
					g.Update(int64(v))
				case time.Duration:
					g := metrics.NewGauge()
					r.Add(m.Metric, tags, g)
					g.Update(int64(v / time.Second))
				default:
					log.Errorf("Unhandled type: %+v", variable)
				}
			}
		}

		rootRegistry := tsdmetrics.NewSegmentedTaggedRegistry("", conf.Tags, nil)

		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					time.Sleep(time.Duration(conf.Freq) * time.Second)
					collectMibs(rootRegistry)
				}
			}
		}()

		// Add Go metrics
		tsdmetrics.RegisterTaggedRuntimeMemStats(tsdmetrics.NewSegmentedTaggedRegistry("", tsdmetrics.Tags{"j_app": "agentx-otsdb"}, rootRegistry))

		metricsTsdb := tsdmetrics.TaggedOpenTSDB{Addr: c.String("metrics-report-url"),
			Registry:      rootRegistry,
			FlushInterval: time.Duration(conf.Freq) * time.Second,
			DurationUnit:  time.Millisecond,
			Format:        tsdmetrics.Json,
			Logger:        log.StandardLogger()}

		//collectFn := []func(tsdmetrics.TaggedRegistry){collectMibs}
		collectFn := tsdmetrics.RuntimeCaptureFn
		collectFn = append(collectFn, collectMibs)

		metricsTsdb.RunWithPreprocessing(ctx, collectFn)

		return nil
	}

	err := app.Run(os.Args)
	wg.Wait()
	if err != nil {
		log.Error(err)
		os.Exit(-1)
	}
}
