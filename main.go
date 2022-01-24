package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/clok/cdocs"
	"github.com/clok/cidr/commands"
	"github.com/clok/kemba"
	"github.com/urfave/cli/v2"
)

var (
	version string
	k       = kemba.New("cidr")
)

func main() {
	k.Println("executing")

	im, err := cdocs.InstallManpageCommand(&cdocs.InstallManpageCommandInput{
		AppName: "cidr",
		Hidden:  true,
	})
	if err != nil {
		log.Fatal(err)
	}

	app := cli.NewApp()
	app.Name = "cidr"
	app.Copyright = "(c) 2022 Derek Smith"
	app.Authors = []*cli.Author{
		{
			Name:  "Derek Smith",
			Email: "derek@clokwork.net",
		},
	}
	app.Version = version
	app.Usage = "tool for checking IPs against CIDR blocks"
	app.Commands = []*cli.Command{
		// Check checks input flag with provided CIDR blocks
		commands.CommandCheck,
		// Filter accepts a file or pipe and filters rows based on CIDR blocks
		commands.CommandFilter,
		im,
		{
			Name:    "version",
			Aliases: []string{"v"},
			Usage:   "Print version info",
			Action: func(c *cli.Context) error {
				fmt.Printf("%s %s (%s/%s)\n", "cidr", version, runtime.GOOS, runtime.GOARCH)
				return nil
			},
		},
	}

	if os.Getenv("DOCS_MD") != "" {
		docs, err := cdocs.ToMarkdown(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	if os.Getenv("DOCS_MAN") != "" {
		docs, err := cdocs.ToMan(app)
		if err != nil {
			panic(err)
		}
		fmt.Println(docs)
		return
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
