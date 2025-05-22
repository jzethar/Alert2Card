package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"release_youtracker/server"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	GitCommit string
	GitTag    string
	BuildTime string
)

func main() {
	cli.VersionPrinter = func(cCtx *cli.Context) {
		fmt.Printf("Git Tag: %s\n", GitTag)
		fmt.Printf("Git Commit: %s\n", GitCommit)
		fmt.Printf("Build Time: %s\n", BuildTime)
	}
	app := &cli.App{
		Name:            "Alert2Card",
		Version:         GitTag,
		HideHelpCommand: true,
		HideVersion:     false,
		Description:     "Receives alerts from AlertManager about new GitLab release and creates a task in Youtrack",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"conf"},
			},
		},
		Action: runServer,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

func runServer(ctx *cli.Context) error {
	log.Info().Msgf("Start Alert2Card server, tag: %v, commit: %v, build: %v", GitTag, GitCommit, BuildTime)
	var err error
	var server server.Server
	if err = server.Init(ctx.String("config")); err != nil {
		log.Fatal().Msg(err.Error())
	}

	commonMux := http.NewServeMux()

	// On Alert - create a task in Youtrack
	commonMux.HandleFunc("/on_alert", server.AlertHandler)
	commonAddress := fmt.Sprintf("%s:%s", server.Config.Host, server.Config.Ports.Rpc)
	commonServer := &http.Server{
		Addr:    commonAddress,
		Handler: commonMux,
	}

	// Debug part
	debugMux := http.NewServeMux()
	debugMux.HandleFunc("/debug/pprof/", pprof.Index)
	debugMux.HandleFunc("/debug/pprof/heap", pprof.Index)
	debugMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	debugMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	debugMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	debugMux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	debugMux.HandleFunc("/metrics", promhttp.Handler().ServeHTTP)
	debugAddress := fmt.Sprintf("%s:%s", server.Config.Host, server.Config.Ports.Debug)
	debugServer := &http.Server{
		Addr:    debugAddress,
		Handler: debugMux,
	}

	go func() {
		log.Info().Msgf("Starting data server on %s \n", commonAddress)
		if err := commonServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Common server failed: %v", err)
		}
	}()

	go func() {
		log.Info().Msgf("Starting debug and metrics server on %s \n", debugAddress)
		if err := debugServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Msgf("Common server failed: %v", err)
		}
	}()

	select {}
}
