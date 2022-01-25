package commands

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"regexp"
	"strings"

	"github.com/clok/kemba"
)

var (
	rMask = regexp.MustCompile(`/\d{1,2}$`)
	rIPs  = regexp.MustCompile(`\b((([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])(\.)){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5]))\b`)

	k     = kemba.New("cidr:commands")
	kf    = k.Extend("filter")
	kc    = k.Extend("common")
	kfb   = kc.Extend("blocks")
	kci   = kc.Extend("ips")
	kfp   = kc.Extend("processReader")
	kfpl  = kfp.Extend("lines")
	kfpd  = kfp.Extend("debug")
	kfpld = kc.Extend("processLine:debug")
)

func isCiderIn(input *net.IPNet, cidrs []*net.IPNet) (bool, string) {
	for _, cidr := range cidrs {
		if cidr.Contains(input.IP) {
			return true, cidr.String()
		}
	}
	return false, ""
}

func parseInputCIDRs(csv string) ([]*net.IPNet, error) {
	var blocks []*net.IPNet
	for _, cidr := range strings.Split(csv, ",") {
		_, ipv4Net, err := net.ParseCIDR(cidr)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, ipv4Net)
	}
	kfb.Printf("total CIDR Blocks %d", len(blocks))
	return blocks, nil
}

func parseInputIPs(csv string) ([]*net.IPNet, error) {
	var ips []*net.IPNet
	for _, ip := range strings.Split(csv, ",") {
		if !rMask.MatchString(ip) {
			ip = fmt.Sprintf("%s/32", ip)
		}
		_, ipv4Net, err := net.ParseCIDR(ip)
		if err != nil {
			return nil, err
		}
		ips = append(ips, ipv4Net)
	}
	kci.Printf("total IPs %d", len(ips))
	return ips, nil
}

type processReaderInput struct {
	reader  *bufio.Reader
	blocks  []*net.IPNet
	inverse bool
}

func processReader(opts *processReaderInput) error {
	var output []rune
	var lines int64

	for {
		input, _, err := opts.reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		kfpd.Printf("%c", input)
		output = append(output, input)
		if input == '\n' {
			err := processLine(&processLineInput{
				output:  output,
				blocks:  opts.blocks,
				inverse: opts.inverse,
			})
			if err != nil {
				return err
			}
			lines++
			output = []rune{}
			kfpd.Println("-- RESET OUTPUT --")
		}
	}

	if len(output) > 0 {
		err := processLine(&processLineInput{
			output:  output,
			blocks:  opts.blocks,
			inverse: opts.inverse,
		})
		if err != nil {
			return err
		}
		lines++
	}
	kfpl.Printf("%d lines processed", lines)
	return nil
}

type processLineInput struct {
	output  []rune
	blocks  []*net.IPNet
	inverse bool
}

func processLine(opts *processLineInput) error {
	line := string(opts.output)
	kfpld.Printf("stringify %s", line)
	ips := rIPs.FindAllStringSubmatch(line, -1)
	switch {
	case len(ips) == 0 && opts.inverse:
		fmt.Println(strings.ReplaceAll(strings.ReplaceAll(line, "\r\n", ""), "\n", ""))
	default:
		// track to avoid printing duplicate lines
		var hasIP bool
		for _, found := range ips {
			ip := found[1]
			if !rMask.MatchString(ip) {
				ip = fmt.Sprintf("%s/32", ip)
			}
			_, ipv4Net, err := net.ParseCIDR(ip)
			if err != nil {
				return err
			}

			if ok, cidr := isCiderIn(ipv4Net, opts.blocks); ok {
				kfpld.Printf("line contains %s which is in CIDR %s", ip, cidr)
				hasIP = true
			}
		}

		if (hasIP && !opts.inverse) || (!hasIP && opts.inverse) {
			fmt.Println(strings.ReplaceAll(strings.ReplaceAll(line, "\r\n", ""), "\n", ""))
		}
	}
	return nil
}
