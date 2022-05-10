package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/charmbracelet/keygen"
	"github.com/charmbracelet/wish"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/charmbracelet/wishlist"
	"github.com/gliderlabs/ssh"
	"github.com/jon4hz/wishbox/internal/config"
	"github.com/jon4hz/wishbox/internal/netbox"
	"github.com/jon4hz/wishbox/internal/version"
	"github.com/muesli/coral"
	mcoral "github.com/muesli/mango-coral"
	"github.com/muesli/roff"
)

var refreshInterval = time.Minute * 5

var rootCmd = &coral.Command{
	Use:     "wishbox",
	Version: version.Version,
	Short:   "wishlist using netbox as inventory source",
	RunE:    root,
}

func root(cmd *coral.Command, args []string) error {
	k, err := keygen.New(".wishlist/server", nil, keygen.Ed25519)
	if err != nil {
		return err
	}
	if !k.KeyPairExists() {
		if err := k.WriteKeys(); err != nil {
			return err
		}
	}

	cfg, err := config.Get()
	if err != nil {
		return err
	}

	endpoints, err := netbox.GetInventory(cfg.Netbox)
	if err != nil {
		return err
	}

	// wishlist config
	wcfg := &wishlist.Config{
		Listen: cfg.Listen,
		Port:   cfg.Port,
		Users:  cfg.Users,
		Factory: func(e wishlist.Endpoint) (*ssh.Server, error) {
			return wish.NewServer(
				wish.WithAddress(e.Address),
				wish.WithHostKeyPEM(k.PrivateKeyPEM()),
				wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
					return true
				}),
				wish.WithMiddleware(
					append(
						e.Middlewares, // this is the important bit: the middlewares from the endpoint
						lm.Middleware(),
					)...,
				),
			)
		},
		Endpoints:    endpoints,
		EndpointChan: make(chan []*wishlist.Endpoint),
	}

	go func() {
		ticker := time.NewTicker(refreshInterval)
		for range ticker.C {
			endpoints, err := netbox.GetInventory(cfg.Netbox)
			if err != nil {
				log.Println("error getting inventory:", err)
				continue
			}
			log.Printf("updated %d endpoints\n", len(endpoints))
			wcfg.EndpointChan <- endpoints
		}
	}()

	if err := wishlist.Serve(wcfg); err != nil {
		log.Fatalln(err)
	}
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {
	rootCmd.AddCommand(manCmd, versionCmd)
}

var manCmd = &coral.Command{
	Use:                   "man",
	Short:                 "generates the manpages",
	SilenceUsage:          true,
	DisableFlagsInUseLine: true,
	Hidden:                true,
	Args:                  coral.NoArgs,
	RunE: func(cmd *coral.Command, args []string) error {
		manPage, err := mcoral.NewManPage(1, rootCmd)
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
		return err
	},
}

var versionCmd = &coral.Command{
	Use:   "version",
	Short: "Print the version info",
	Run: func(cmd *coral.Command, args []string) {
		fmt.Printf("Version: %s\n", version.Version)
		fmt.Printf("Commit: %s\n", version.Commit)
		fmt.Printf("Date: %s\n", version.Date)
		fmt.Printf("Build by: %s\n", version.BuiltBy)
	},
}
