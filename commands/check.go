package commands

import (
	"fmt"
	"net"
	"strings"

	"github.com/clok/kemba"
	"github.com/urfave/cli/v2"
)

var (
	k = kemba.New("cidr:commands")
)

func isCiderIn(input *net.IPNet, cidrs []*net.IPNet) (bool, string) {
	for _, cidr := range cidrs {
		if cidr.Contains(input.IP) {
			return true, cidr.String()
		}
	}
	return false, ""
}

var (
	CommandCheck = &cli.Command{
		Name:  "check",
		Usage: "check IP against range of CIDR blocks",
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
			kl := k.Extend("check")

			// Get CIDR blocks
			var blocks []*net.IPNet
			for _, cidr := range strings.Split(c.String("blocks"), ",") {
				_, ipv4Net, err := net.ParseCIDR(cidr)
				if err != nil {
					return err
				}
				blocks = append(blocks, ipv4Net)
			}
			kl.Extend("blocks").Printf("total CIDR Blocks %d", len(blocks))

			// Get IPs
			var ips []*net.IPNet
			for _, ip := range strings.Split(c.String("ips"), ",") {
				_, ipv4Net, err := net.ParseCIDR(ip)
				if err != nil {
					return err
				}
				ips = append(ips, ipv4Net)
			}

			kl.Extend("ips").Printf("total IPs %d", len(ips))

			for _, ip := range ips {
				if ip.String() == "0.0.0.0/0" {
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
