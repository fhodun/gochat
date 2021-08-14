package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	flag "github.com/spf13/pflag"
)

type FlagsType struct {
	Port    string
	Host    string
}

var (
	Flags FlagsType
)

func main() {
	var cmdServer *cobra.Command = &cobra.Command{
		Use:     "server",
		Aliases: []string{"s"},
		Short:   "Run chat server",
		Run:     RunServer,
	}

	var cmdClient *cobra.Command = &cobra.Command{
		Use:     "client",
		Aliases: []string{"c"},
		Short:   "Run chat client",
		Run:     RunClient,
	}

	var rootCmd *cobra.Command = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdServer, cmdClient)

	flag.StringVarP(&Flags.Port, "port", "p", "2137", "http port")
	flag.StringVar(&Flags.Host, "host", "127.0.0.1", "http host")
	
	flag.Parse()

	log.Fatal(rootCmd.Execute())
}
