package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/camaeel/proxmox-oidc-credential-helper/pkg/browser"
	"github.com/camaeel/proxmox-oidc-credential-helper/pkg/callback"
	"github.com/camaeel/proxmox-oidc-credential-helper/pkg/config"
	"github.com/camaeel/proxmox-oidc-credential-helper/pkg/proxmox"
	"log/slog"
	"os"
	"time"
)

func main() {
	cfg := config.Config{}

	flag.IntVar(&cfg.CallbackPort, "callback-port", 8996, "callback port. Default is 8996")
	flag.StringVar(&cfg.CallbackPath, "callback-path", "/oidc/callback", "callback port. Default is 8996")
	flag.StringVar(&cfg.ProxmoxURL, "proxmox-url", "", "Url to proxmox server with protocol and port, i.e. https://proxmox.example.com:8006")
	flag.StringVar(&cfg.Realm, "realm", "", "Proxmox OIDC realm")
	flag.IntVar(&cfg.TimeoutSeconds, "timeout-url", 180, "Timeout in seconds for whole authentication. By default 180 seconds.")
	flag.BoolVar(&cfg.VerboseLog, "verbose", false, "Verbose logging")
	flag.StringVar(&cfg.OutputFormat, "output", "text", "Output format. One of: text|json. Default is text")
	flag.BoolVar(&cfg.OpenDefaultBrowser, "open-browser", true, "Open default browser. Default is true. Might be useful in situations where default browser is not working well.")

	flag.Parse()

	if cfg.ProxmoxURL == "" {
		slog.Error("Flag proxmox-url is required")
		os.Exit(1)
	}

	if cfg.Realm == "" {
		slog.Error("Flag realm is required")
		os.Exit(1)
	}

	switch cfg.OutputFormat {
	case "text":
		cfg.OutputFormat = "text"
	case "json":
		cfg.OutputFormat = "json"
	default:
		slog.Error("Output format must be one of: text|json")
		os.Exit(1)
	}

	if cfg.VerboseLog {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.TimeoutSeconds)*time.Minute)
	defer cancel()

	result := make(chan callback.CallbackResult)

	cancelHttp := callback.StartHttpServer(cancel, cfg.CallbackPort, cfg.CallbackPath, result)
	defer cancelHttp()

	redirectUrl, err := proxmox.GetOidcURL(cfg)
	if err != nil {
		slog.Error("Unable to obtain oidc URL from proxmox server", slog.String("error", err.Error()))
		os.Exit(1)
	}
	if cfg.OpenDefaultBrowser {
		err = browser.OpenURL(redirectUrl)
		if err != nil {
			slog.Error("Unable to open browser", slog.String("error", err.Error()))
			os.Exit(1)
		}
	} else {
		fmt.Println(browser.PrintableOpenURL(redirectUrl))
	}

	select {
	case <-ctx.Done():
		slog.Error("Timeout for authentication flow exceeded")
		os.Exit(1)
	case resultData := <-result:
		cancelHttp()
		cancel()
		ticket, err := proxmox.ExchangeCallbackResultForTicket(cfg, resultData)
		if err != nil {
			slog.Error("Unable to exchange callback result", slog.String("error", err.Error()))
			os.Exit(1)
		}
		output, err := proxmox.OutputTicket(cfg, ticket)
		if err != nil {
			slog.Error("Unable to obtain output", slog.String("error", err.Error()))
			os.Exit(1)
		}
		fmt.Print(output)
	}

}
