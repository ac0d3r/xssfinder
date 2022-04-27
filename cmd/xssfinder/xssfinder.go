package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/Buzz2d0/xssfinder/internal/runner"
	"github.com/Buzz2d0/xssfinder/logger"
	"github.com/sirupsen/logrus"
)

var (
	fdebug             = flag.Bool("debug", false, "Set debug mode")
	fvverbose          = flag.Bool("vv", false, "Set very-verbose mode")
	fmitmAddr          = flag.String("maddr", ":8222", "Set mitm-server listen address")
	fmitmVerbose       = flag.Bool("mverbose", false, "Set mitm-server verbose mode")
	fmitmTargetHosts   = flag.String("mhosts", "", "Set mitm-server target hosts .e.g. foo.com,bar.io")
	fmitmParentProxy   = flag.String("mporxy", "", "Set mitm-server parent proxy")
	flogOutjosn        = flag.Bool("outjson", false, "Set logger output json format")
	flogNoColor        = flag.Bool("nocolor", false, "Set logger no-color mode")
	fbrowserExecPath   = flag.String("bexecpath", "", "Set browser exec path")
	fbrowserNoHeadless = flag.Bool("noheadless", false, "Set browser no-headless mode")
	fbrowserIncognito  = flag.Bool("incognito", false, "Set browser incognito mode")
	fbrowserProxy      = flag.String("bproxy", "", "Set browser proxy addr")
)

func main() {
	banner()
	flag.Parse()

	opt := runner.NewOptions()
	opt.Set(runner.WithDebug(*fdebug),
		runner.WithVeryVerbose(*fvverbose),
		runner.WithMitmAddr(*fmitmAddr),
		runner.WithMitmVerbose(*fmitmVerbose),
		runner.WithMitmTargetHosts(parseMultiHosts(*fmitmTargetHosts)...),
		runner.WithMitmParentProxy(*fmitmParentProxy),
		runner.WithLogOutJson(*flogOutjosn),
		runner.WithLogNoColor(*flogNoColor),
		runner.WithBrowserExecPath(*fbrowserExecPath),
		runner.WithBrowserNoHeadless(*fbrowserNoHeadless),
		runner.WithBrowserIncognito(*fbrowserIncognito),
		runner.WithBrowserProxy(*fbrowserProxy),
	)

	opt.Log.Level = opt.LogLevel()
	logger.Init("xssfinder", opt.Log)

	xssfinder := runner.NewRunner(opt)
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		if err := xssfinder.Run(ctx); err != nil {
			logrus.Errorln(err)
			c <- os.Kill
		}
	}()

	<-c
	cancel()
	xssfinder.Shutdown(ctx)
}

func parseMultiHosts(s string) []string {
	if len(s) == 0 {
		return nil
	}
	ss := strings.Split(strings.TrimSpace(s), ",")
	r := make([]string, 0)
	for _, v := range ss {
		r = append(r, strings.TrimSpace(v))
	}
	return r
}

const (
	banners = `
    ▄    ▄▄▄▄▄    ▄▄▄▄▄   ▄████  ▄█    ▄   ██▄   ▄███▄   █▄▄▄▄ 
▀▄   █  █     ▀▄ █     ▀▄ █▀   ▀ ██     █  █  █  █▀   ▀  █  ▄▀ 
  █ ▀ ▄  ▀▀▀▀▄ ▄  ▀▀▀▀▄   █▀▀    ██ ██   █ █   █ ██▄▄    █▀▀▌  
 ▄ █   ▀▄▄▄▄▀   ▀▄▄▄▄▀    █      ▐█ █ █  █ █  █  █▄   ▄▀ █  █  
█   ▀▄                     █      ▐ █  █ █ ███▀  ▀███▀     █   
 ▀                          ▀       █   ██                ▀`
)

func banner() {
	fmt.Println(banners)
}
