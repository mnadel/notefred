package main

import (
	"log"

	"github.com/mnadel/notefred/cmd/search"
	"github.com/mnadel/notefred/cmd/version"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "notefred",
		Short: "A CLI for an Alfred+Notes.app integration",
		Long:  "Search note titles",
	}

	cmd.AddCommand(search.New())
	cmd.AddCommand(version.New())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
