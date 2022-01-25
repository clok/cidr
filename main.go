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
	app.UsageText = `
A CLI tool that is useful for quickly checking and filtering data by IP vs. CIDR blocks.

The 'check'' command allows for a quick check of a list of IPs and Blocks.

	$ cidr check --blocks 172.12.0.0/16,172.10.0.0/16 --ips 172.12.1.56,171.10.123.57,172.10.0.255/32
	172.12.1.56/32 is in CIDR 172.12.0.0/16
	171.10.123.57/32 is NOT in CIDR set
	172.10.0.255/32 is in CIDR 172.10.0.0/16

The 'filter' command is useful for filtering large data sets like access log files.

	$ cidr filter --blocks 10.2.120.0/8,10.2.122.0/8,10.20.128.20/29 --path '/var/log/http/**/access*.log'
	< outputs to STDOUT all lines that contain an IP that is within a CIDR blocks provided >

The 'filter' command can also be used with a pipe.

	$ cidr filter --blocks 10.2.120.0/8,10.2.122.0/8,10.20.128.20/29 < /var/log/http/access-20220120-18.log
	< outputs to STDOUT all lines that contain an IP that is within a CIDR blocks provided >

Finally, the 'filter' command accepts the '--inverse, i' flag which will output all lines that DO NOT contain
an IP within a CIDR block provided. If a line has multiple IP addresses within it, then ALL IPs must not be within
a CIDR block for the line to be output to STDOUT.

	$ cidr filter --blocks 10.2.120.0/8,10.2.122.0/8,10.20.128.20/29 --path '/var/log/http/**/access*.log' --inverse
	< outputs to STDOUT all lines that DO NOT contain an IP that is within a CIDR blocks provided >
`
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
