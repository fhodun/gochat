package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/fhodun/gochat/client"
	"github.com/fhodun/gochat/server"
)

func main() {
	cmdServer := &cobra.Command{
		Use:     "server [port]",
		Aliases: []string{"s"},
		Short:   "Run chat server",
		Example: "server 2137",
		Run:     server.RunServer,
	}

	cmdClient := &cobra.Command{
		Use:     "client [host:port]",
		Aliases: []string{"c"},
		Short:   "Run chat client",
		Example: "server 127.0.0.1:2137",
		Run:     client.RunClient,
	}

	rootCmd := &cobra.Command{Use: "gochat"}
	rootCmd.AddCommand(cmdServer, cmdClient)

	log.Fatal(rootCmd.Execute())
}
