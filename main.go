package main

import (
	"flag"
	"fmt"
	"github.com/open-falcon/alarm/cron"
	"github.com/open-falcon/alarm/g"
	"github.com/open-falcon/alarm/db"
	"github.com/open-falcon/alarm/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg := flag.String("c", "cfg.json", "configuration file")
	version := flag.Bool("v", false, "show version")
	help := flag.Bool("h", false, "help")
	flag.Parse()

	if *version {
		fmt.Println(g.VERSION)
		os.Exit(0)
	}

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	g.ParseConfig(*cfg)
	g.InitRedisConnPool()
	db.Init()

	go http.Start()

	go cron.ReadHighEvent()
	go cron.ReadLowEvent()
	go cron.CombineSms()
	go cron.CombineMail()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Println()
		g.RedisConnPool.Close()
		os.Exit(0)
	}()

	select {}
}
