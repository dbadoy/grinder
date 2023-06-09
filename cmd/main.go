package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/dbadoy/grinder/api"
	"github.com/dbadoy/grinder/pkg/checkpoint"
	"github.com/dbadoy/grinder/pkg/database"
	"github.com/dbadoy/grinder/pkg/database/es"
	"github.com/dbadoy/grinder/pkg/database/memdb"
	"github.com/dbadoy/grinder/pkg/ethclient"
	"github.com/dbadoy/grinder/server"
	"github.com/dbadoy/grinder/server/cft"
	"github.com/dbadoy/grinder/server/fetcher"
)

func main() {
	var (
		fetchInterval = flag.Duration("fetch", time.Second, "interval time to fetch block from ethereum")
		ethEndpoint   = flag.String("ethendpoint", "", "ethereum endpoint url (suggest: jsonrpc)")
		cp            = flag.String("checkpoint", "checkpoint", "checkpoint name")
		db            = flag.String("db", "elasticsearch", "database (elasticsearch|memory)")
		dbpath        = flag.String("dbpath", "", "database urls (url1,url2,url3...)")
		cluster       = flag.String("cluster", "", "cluster node list (IP:PORT,IP:PORT,IP:PORT...)")
		http          = flag.Int("http", 0, "http listening port (0 = not support http)")
	)
	flag.Parse()

	eth, err := ethclient.New(*ethEndpoint)
	if err != nil {
		panic(fmt.Errorf("invalid ethereum endpoint: %s (%v)", *ethEndpoint, err))
	}

	if err := ethclient.DefaultHeartbeat(context.Background(), *ethEndpoint); err != nil {
		panic(fmt.Errorf("ethereum endpoint has no response: %s (%v)", *ethEndpoint, err))
	}

	// Databse
	var database database.Database
	switch *db {
	case "elasticsearch":
		database, err = es.New(strings.Split(*dbpath, ","))
	case "memory":
		database = memdb.New()
	default:
		err = errors.New("invalid database")
	}

	if err != nil {
		panic(fmt.Errorf("database connection failed: %v", err))
	}

	if err := database.HealthCheck(); err != nil {
		panic(fmt.Errorf("health check failed; kind: %s, path: %v, reason: %v", *db, *dbpath, err))
	}

	// Checkpoint
	checkpoint := checkpoint.New(checkpoint.DefaultBasePath, *cp)

	// Cluster
	var engine cft.Engine
	switch len(*cluster) {
	case 0:
		// Solo
		engine, err = cft.NewSoloEngine(nil, database, checkpoint)
	default:
		// // Print CFT info
		// engine = &cft.CFT{}
		panic("only support solo mode")
	}

	if err != nil {
		panic(fmt.Errorf("engine creation failed: %v", err))
	}

	// Fetcher
	fetcher := fetcher.New(
		eth,
		checkpoint,
		&fetcher.Config{
			PollInterval: *fetchInterval,
		},
	)

	server, err := server.New(
		eth,
		fetcher,
		engine,
		checkpoint,
		&server.Config{
			AllowProxyContract: true,
		},
	)

	if err != nil {
		panic(err)
	}

	server.Run()

	if *http != 0 {
		api := api.New(uint16(*http), server)
		api.Listen()
	}

	select {}
}
