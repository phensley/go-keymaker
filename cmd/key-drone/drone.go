package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/phensley/go-keymaker"
	"github.com/spf13/cobra"
)

const (
	defaultConfig = `
address: 0.0.0.0:10101
concurrency: 0
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
	cfg := &keymaker.DroneConfig{
		Dir: filepath.Dir(configPath),
	}
	err := keymaker.LoadConfig(cfg, []byte(defaultConfig))
	keymaker.LogFail(err, "reading config defaults:")

	if configPath != "" {
		err = keymaker.LoadConfigFile(cfg, configPath)
		keymaker.LogFail(err, "config file %s", configPath)
	}

	drone, err := keymaker.NewDrone(cfg)
	keymaker.LogFail(err, "NewDrone")
	log.Printf("%s on %s", os.Args[0], cfg.Address)
	err = drone.Start()
	keymaker.LogFail(err, "drone.Start()")
}
