package commands

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

var (
	rIPs = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
)
var (
	CommandPipe = &cli.Command{
		Name:  "pipe",
		Usage: "command | cidr pipe",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "blocks",
				Aliases:  []string{"b"},
				Usage:    "CIDR blocks to be checked (csv)",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			kl := k.Extend("pipe")
			klb := kl.Extend("blocks")
			klc := kl.Extend("clear")
			kld := kl.Extend("debug")

			// Get CIDR blocks
			var blocks []*net.IPNet
			for _, cidr := range strings.Split(c.String("blocks"), ",") {
				_, ipv4Net, err := net.ParseCIDR(cidr)
				if err != nil {
					return err
				}
				blocks = append(blocks, ipv4Net)
			}
			klb.Printf("total CIDR Blocks %d", len(blocks))

			info, err := os.Stdin.Stat()
			if err != nil {
				return err
			}

			if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
				return fmt.Errorf("please use this command with a pipe")
			}

			reader := bufio.NewReader(os.Stdin)
			var output []rune
			var lines int64

			for {
				input, _, err := reader.ReadRune()
				if err != nil && err == io.EOF {
					break
				}
				kl.Printf("%c", input)
				output = append(output, input)
				if input == '\n' {
					line := string(output)
					kl.Printf("stringify %s", line)
					matches := rIPs.FindAllStringSubmatch(line, -1)
					for _, v := range matches {
						ip := v[1]
						if !rMask.MatchString(ip) {
							ip = fmt.Sprintf("%s/32", ip)
						}
						_, ipv4Net, err := net.ParseCIDR(ip)
						if err != nil {
							return err
						}
						if ok, cidr := isCiderIn(ipv4Net, blocks); ok {
							fmt.Printf("%s is in CIDR %s\n", v[1], cidr)
						}
					}
					lines++
					output = []rune{}
					klc.Println("------")
				}
			}

			if len(output) > 0 {
				for j := 0; j < len(output); j++ {
					fmt.Printf("%c", output[j])
				}
				lines++
			}
			kld.Printf("%d lines processed", lines)
			return nil
		},
	}
)
