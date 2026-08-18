package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/localbeam/localbeam/internal/config"
	"github.com/localbeam/localbeam/internal/server"
)

const version = "1.0.0"

const banner = `
 _                     _ ____
| |     ___  ___ __ _ | | __ )  ___  __ _ _ __ ___
| |    / _ \/ __/ _'  | |  _ \ / _ \/ _' | '_ ' _ \
| |___| (_) | (_| (_| | | |_) |  __/ (_| | | | | | |
|_____|\___/ \___\__,_|_|____/ \___|\__,_|_| |_| |_|

  Secure local network file & text transfer  v%s
`

func main() {
	var (
		cfgPath  = flag.String("config", "", "Path to config file (default: ~/.localbeam/config.json)")
		port     = flag.Int("port", 0, "Override server port")
		host     = flag.String("host", "", "Override server host")
		showVer  = flag.Bool("version", false, "Show version")
		initConf = flag.Bool("init-config", false, "Create default config file")
	)
	flag.Parse()

	if *showVer {
		fmt.Printf("LocalBeam v%s\n", version)
		os.Exit(0)
	}

	fmt.Printf(banner, version)

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if *initConf {
		if err := config.Save(cfg, *cfgPath); err != nil {
			log.Fatalf("Failed to save config: %v", err)
		}
		home, _ := os.UserHomeDir()
		fmt.Printf("✅ Config saved to %s/.localbeam/config.json\n", home)
		os.Exit(0)
	}

	// Apply overrides
	if *port != 0 {
		cfg.Server.Port = *port
	}
	if *host != "" {
		cfg.Server.Host = *host
	}

	fmt.Println(
		*cfgPath,
		*port,
		*host,
		*showVer,
		*initConf,
	)

	srv := server.New(cfg)
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
