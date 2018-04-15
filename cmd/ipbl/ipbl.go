package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/coreos/pkg/capnslog"
	"github.com/coreos/pkg/flagutil"
	"github.com/galexrt/ipbl/pkg/db"
	"github.com/galexrt/ipbl/pkg/routes"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/gommon/log"
	"github.com/prometheus/common/version"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type CmdLineOptions struct {
	DatabaseDriver string
	DatabaseDSN    string
}

var (
	opts      CmdLineOptions
	ipblFlags = flag.NewFlagSet("apis", flag.ExitOnError)
	logger    = capnslog.NewPackageLogger("github.com/galexrt/ipbl/cmd", "main")
)

func init() {
	ipblFlags.StringVar(&opts.DatabaseDriver, "database-driver", "mysql", "GORM supported database type")
	ipblFlags.StringVar(&opts.DatabaseDSN, "database-dsn", "", "DSN to database")

	ipblFlags.Parse(os.Args[1:])
}

func copyFlag(name string) {
	ipblFlags.Var(flag.Lookup(name).Value, flag.Lookup(name).Name, flag.Lookup(name).Usage)
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]...\n", os.Args[0])
	ipblFlags.PrintDefaults()
	os.Exit(0)
}

func main() {
	flagutil.SetFlagsFromEnv(ipblFlags, "IPBL")

	logger.Infof("starting %s %s %s", os.Args[0], version.Info(), version.BuildContext())

	e := gin.New()
	e.Use(gin.Recovery())
	p := ginprometheus.NewPrometheus("gin")
	p.Use(e)

	dbCon, err := sqlx.Open(opts.DatabaseDriver, opts.DatabaseDSN)
	if err != nil {
		logger.Fatal(err)
	}
	db.DBCon = dbCon
	defer db.DBCon.Close()

	routes.Register(e)

	srv := &http.Server{
		Addr:    ":1232",
		Handler: e,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Info("shutdown server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown:", err)
	}
	logger.Info("exiting")
}
