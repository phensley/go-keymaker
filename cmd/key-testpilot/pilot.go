package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/phensley/go-keymaker"
	"github.com/spf13/cobra"
)

// testpilot - demonstrates generating keys from a cluster of drones

var (
	verbose    bool
	droneAddrs []string
	keyTypes   []string

	cmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "Generates keys and displays stats",
		Run:   run,
	}
)

func main() {
	cmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose mode")
	cmd.Flags().StringSliceVarP(&droneAddrs, "drone", "d", nil, "Drone address")
	cmd.Flags().StringSliceVarP(&keyTypes, "type", "t", nil, "Key type")
	cmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	switch {
	case droneAddrs == nil:
		log.Fatalln("no drone addresses specified")
	case keyTypes == nil:
		log.Fatalln("no key types specified")
	}

	client := keymaker.NewClient(&keymaker.ClientConfig{
		BufferSize: 8,
		Addresses:  droneAddrs,
	})

	state := &keymaker.State{}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	go func() {
		<-sigCh
		state.Touch()
		client.Stop()
	}()

	for _, k := range keyTypes {
		go generate(state, client, k)
	}
	for state.OK() {
		<-time.After(time.Second)
	}
}

func generate(state *keymaker.State, client *keymaker.Client, keyType string) {
	ch := client.Generate(keyType)
	count := int32(0)
	go func() {
		for state.OK() {
			<-time.After(time.Second)
			log.Println(keyType, atomic.LoadInt32(&count), "generated")
		}
	}()
	for key := range ch {
		count++
		if verbose {
			fmt.Println(string(key))
		}
	}
	fmt.Println(keyType, "complete")
}
