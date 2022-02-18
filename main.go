package main

import (
	"log"

	"github.com/charmbracelet/keygen"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	lm "github.com/charmbracelet/wish/logging"
	"github.com/charmbracelet/wishlist"
	"github.com/gliderlabs/ssh"
	"github.com/jon4hz/wishbox/config"
	"github.com/jon4hz/wishbox/netbox"
)

func main() {
	k, err := keygen.New(".wishlist", "server", nil, keygen.Ed25519)
	if err != nil {
		log.Fatalln(err)
	}
	if !k.IsKeyPairExists() {
		if err := k.WriteKeys(); err != nil {
			log.Fatalln(err)
		}
	}

	cfg, err := config.Get()
	if err != nil {
		log.Fatalln(err)
	}

	endpoints, err := netbox.GetInventory(cfg.Netbox)
	if err != nil {
		log.Fatal(err)
	}

	// wishlist config
	wcfg := &wishlist.Config{
		Listen: cfg.Listen,
		Port:   cfg.Port,
		Users:  cfg.Users,
		Factory: func(e wishlist.Endpoint) (*ssh.Server, error) {
			return wish.NewServer(
				wish.WithAddress(e.Address),
				wish.WithHostKeyPEM(k.PrivateKeyPEM),
				wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
					return true
				}),
				wish.WithMiddleware(
					append(
						e.Middlewares, // this is the important bit: the middlewares from the endpoint
						lm.Middleware(),
						activeterm.Middleware(),
					)...,
				),
			)
		},
		Endpoints: endpoints,
	}

	if err := wishlist.Serve(wcfg); err != nil {
		log.Fatalln(err)
	}
}
