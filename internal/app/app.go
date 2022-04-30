package app

import (
	"github.com/Buzz2d0/xssfinder/internal/logger"
	"github.com/Buzz2d0/xssfinder/internal/options"
	"github.com/Buzz2d0/xssfinder/internal/runner"
	"github.com/urfave/cli/v2"
)

func New(version string) *cli.App {
	app := &cli.App{
		Name:    "xssfinder",
		Version: version,
		Usage:   "XSS discovery tool",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
				Value:   false,
				Usage:   "enable debug mode",
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"vv"},
				Value:   false,
				Usage:   "enable very-verbose mode",
			},
			&cli.StringFlag{
				Name:  "notifier-yaml",
				Value: "",
				Usage: "set notifier yaml configuration file",
			},
			&cli.BoolFlag{
				Name:  "outjson",
				Value: false,
				Usage: "set logger output json format",
			},
			&cli.StringFlag{
				Name:    "exec",
				Aliases: []string{"e"},
				Value:   "",
				Usage:   "set browser exec path",
			},
			&cli.BoolFlag{
				Name:  "noheadless",
				Value: false,
				Usage: "disable browser headless mode",
			},
			&cli.BoolFlag{
				Name:  "incognito",
				Value: false,
				Usage: "enable browser incognito mode",
			},
			&cli.StringFlag{
				Name:  "proxy",
				Value: "",
				Usage: "set proxy and all traffic will be routed from the proxy server through",
			},
		},
		Commands: []*cli.Command{
			{
				Name:  "mitm",
				Usage: "Passive agent scanning",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "addr",
						Value: "127.0.0.1:8222",
						Usage: "set mitm proxy server listen address",
					},
					&cli.StringSliceFlag{
						Name:  "hosts",
						Value: cli.NewStringSlice(),
						Usage: "set mitm scan host whitelist",
					},
				},
				Action: func(c *cli.Context) error {
					return mitmAction(c)
				},
			},
		},
	}

	return app
}

func mitmAction(c *cli.Context) error {
	opt := options.New().Set(
		options.WithDebug(c.Bool("debug")),
		options.WithVeryVerbose(c.Bool("verbose")),
		options.WithNotifierYaml(c.String("notifier-yaml")),
		options.WithLogOutJson(c.Bool("outjson")),
		options.WithBrowserExecPath(c.String("exec")),
		options.WithBrowserNoHeadless(c.Bool("noheadless")),
		options.WithBrowserIncognito(c.Bool("incognito")),
		options.WithBrowserProxy(c.String("proxy")),

		options.WithMitmVerbose(c.Bool("verbose")),
		options.WithMitmParentProxy(c.String("proxy")),
		options.WithMitmAddr(c.String("addr")),
		options.WithMitmTargetHosts(c.StringSlice("hosts")...),
	)

	opt.Log.Level = opt.LogLevel()
	logger.Init(opt.Log)
	r, err := runner.NewRunner(opt)
	if err != nil {
		return err
	}
	r.Start()
	return nil
}
