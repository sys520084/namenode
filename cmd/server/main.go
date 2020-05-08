package main

import (
	"flag"
	"fmt"
	"os"

	api "github.com/sys520084/namenode/api"
	"github.com/sys520084/namenode/internal/config"
	"github.com/sys520084/namenode/internal/log"
)

func myUsage() {
	fmt.Printf("Usage: %s [OPTIONS] argument ...\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = myUsage

	cfg := flag.String("cfg", "default", "configuration file exp. default -> default.toml、default.yaml、default.json")
	flag.Parse()

	config.Init(*cfg)
	log.InitLogger()
	r := api.SetupRouter()
	r.Run(":8080")
}
