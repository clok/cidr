package commands

import (
	"bufio"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"github.com/yargevad/filepathx"
)

var (
	CommandFilter = &cli.Command{
		Name:    "filter",
		Aliases: []string{"f"},
		Usage:   "Filters lines in log files of pipe input",
		UsageText: `
Filters lines in log files of pipe input, printing to STDOUT the lines that
contain an IP that is within the provided CIDR blocks.

NOTE: If '--path, -p' is NOT set, then a pipe is assumed to be the input.

The '--inverse, -i' flag will output the lines that do not contain an IP
within the provided CIDR blocks. 
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "blocks",
				Aliases:  []string{"b"},
				Usage:    "CIDR blocks to be checked (csv)",
				Required: true,
			},
			&cli.StringFlag{
				Name:    "path",
				Aliases: []string{"p"},
				Usage:   "File path to files to filter, can be a glob. If not set, a pipe is assumed.",
			},
			&cli.BoolFlag{
				Name:    "inverse",
				Aliases: []string{"i"},
				Usage:   "Print out lines that DO NOT match the CIDR check",
			},
		},
		Action: func(c *cli.Context) error {
			// Get CIDR blocks
			blocks, err := parseInputCIDRs(c.String("blocks"))
			if err != nil {
				return err
			}

			// Verify inputs
			if c.String("path") != "" {
				kf.Printf("globbing files with pattern: %s", c.String("path"))
				files, err := filepathx.Glob(c.String("path"))
				if err != nil {
					return err
				}
				kf.Printf("found %d files", len(files))
				kf.Log(files)

				// filter files
				// For each file, create reader, pass in reader
				for _, fPath := range files {
					kf.Printf("processing file: %s", fPath)
					file, err := os.Open(fPath)
					if err != nil {
						return err
					}

					reader := bufio.NewReader(file)
					err = processReader(&processReaderInput{
						reader:  reader,
						blocks:  blocks,
						inverse: c.Bool("inverse"),
					})
					if err != nil {
						return err
					}
				}

				return nil
			}

			// run in pipe pass blocks
			info, err := os.Stdin.Stat()
			if err != nil {
				return err
			}

			noNamedPipe := info.Mode()&os.ModeNamedPipe == 0
			noUnixPipe := info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0
			k.Printf("noNamedPipe: %t noUnixPipe: %t\n", noNamedPipe, noUnixPipe)

			if noNamedPipe && noUnixPipe {
				// if neither, throw error
				return fmt.Errorf("please use this command with a pipe or the --path flag set")
			}

			reader := bufio.NewReader(os.Stdin)
			err = processReader(&processReaderInput{
				reader:  reader,
				blocks:  blocks,
				inverse: c.Bool("inverse"),
			})
			if err != nil {
				return err
			}

			return nil
		},
	}
)
