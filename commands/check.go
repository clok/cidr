package commands

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	CommandCheck = &cli.Command{
		Name:  "check",
		Usage: "check IP against range of CIDR blocks",
		UsageText: `
IPs that are provided without a mask will be assumed to be /32

CIDR blocks require a mask be provided.
`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "ips",
				Aliases:  []string{"i"},
				Usage:    "CSV list of IPs with masks (csv)",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "blocks",
				Aliases:  []string{"b"},
				Usage:    "CIDR blocks to be checked (csv)",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			// Get CIDR blocks
			blocks, err := parseInputCIDRs(c.String("blocks"))
			if err != nil {
				return err
			}

			// Get IPs
			ips, err := parseInputIPs(c.String("ips"))
			if err != nil {
				return err
			}

			for _, ip := range ips {
				fmt.Println(ip.IP.String(), ip.String(), ip.Mask.String(), ip.IP.To16())
				if strings.HasPrefix(ip.String(), "0.0.0.0") {
					fmt.Printf("%s is in ALL CIDR sets\n", ip.String())
				} else if ok, cidr := isCiderIn(ip, blocks); ok {
					fmt.Printf("%s is in CIDR %s\n", ip.String(), cidr)
				} else {
					fmt.Printf("%s is NOT in CIDR set\n", ip.String())
				}
			}

			return nil
		},
	}
)
