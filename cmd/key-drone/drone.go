package main

import (
	"fmt"
	"log"
	"os"

	"github.com/phensley/go-keymaker"
	"github.com/spf13/cobra"
)

const (
	defaultConfig = `
address: 0.0.0.0:10101
concurrency: 0

# TODO: tls config
`
)

var (
	configPath string

	cmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "RPC service that generates keys on demand",
		Run:   run,
	}
)

func main() {
	cmd.Flags().StringVarP(&configPath, "config", "c", "", "Path to YAML config file")
	cmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	cfg := &keymaker.DroneConfig{}
	err := keymaker.LoadConfig(cfg, []byte(defaultConfig))
	logFail(err, "reading config defaults:")
	if configPath != "" {
		err = keymaker.LoadConfigFile(cfg, configPath)
		logFail(err, "config file %s", configPath)
	}

	drone := keymaker.NewDrone(cfg)
	log.Printf("%s on %s", os.Args[0], cfg.Address)
	err = drone.Start()
	logFail(err, "drone start")
}

func logFail(err error, msg string, args ...interface{}) {
	if err != nil {
		log.Fatalln(fmt.Sprintf(msg, args...), err)
	}
}
