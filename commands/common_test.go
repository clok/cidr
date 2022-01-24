package commands

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_isCiderIn(t *testing.T) {
	is := assert.New(t)

	blocks, _ := parseInputCIDRs("1.1.0.0/16,1.2.0.0/16,1.3.0.0/16")

	t.Run("ip is in a provided CIDR", func(t *testing.T) {
		_, ipv4Net, _ := net.ParseCIDR("1.2.3.4/32")
		check, cidr := isCiderIn(ipv4Net, blocks)
		is.True(check)
		is.Equal("1.2.0.0/16", cidr)
	})

	t.Run("ip is NOT in a provided CIDR", func(t *testing.T) {
		_, ipv4Net, _ := net.ParseCIDR("123.2.3.4/32")
		check, cidr := isCiderIn(ipv4Net, blocks)
		is.False(check)
		is.Equal("", cidr)
	})
}

func Test_parseInputCIDRs(t *testing.T) {
	is := assert.New(t)

	t.Run("successfully parses", func(t *testing.T) {
		var blocks []*net.IPNet
		var err error
		blocks, err = parseInputCIDRs("1.1.0.0/16")
		is.NoError(err)
		is.Equal("1.1.0.0/16", blocks[0].String())

		blocks, err = parseInputCIDRs("1.1.0.0/16,1.2.0.0/16,1.3.0.0/16")
		is.NoError(err)
		is.Equal("1.1.0.0/16", blocks[0].String())
		is.Equal("1.2.0.0/16", blocks[1].String())
		is.Equal("1.3.0.0/16", blocks[2].String())
	})

	t.Run("error while parsing", func(t *testing.T) {
		var blocks []*net.IPNet
		var err error
		blocks, err = parseInputCIDRs("0.0.0.0")
		is.Error(err)
		is.Nil(blocks)

		blocks, err = parseInputCIDRs("")
		is.Error(err)
		is.Nil(blocks)
	})
}

func Test_processLine(t *testing.T) {
	is := assert.New(t)

	blocks, _ := parseInputCIDRs("1.1.0.0/16,1.2.0.0/16,1.3.0.0/16")

	t.Run("print out line that match CIDR blocks [single]", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := "Test 1.2.3.4 Test"
		err := processLine(&processLineInput{
			output:  []rune(input),
			blocks:  blocks,
			inverse: false,
		})
		is.NoError(err)

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		is.Equal(fmt.Sprintf("%s\n", input), string(out))
	})

	t.Run("print out line that match CIDR blocks [multi]", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := "Test1 123.123.123.123 1.2.3.4 Test1"
		err := processLine(&processLineInput{
			output:  []rune(input),
			blocks:  blocks,
			inverse: false,
		})
		is.NoError(err)

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		is.Equal(fmt.Sprintf("%s\n", input), string(out))
	})

	t.Run("print out line that match CIDR blocks [newline]", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := "Test1 1.2.3.4 Test1\n"
		err := processLine(&processLineInput{
			output:  []rune(input),
			blocks:  blocks,
			inverse: false,
		})
		is.NoError(err)

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		is.Equal(fmt.Sprintf("%s\n", "Test1 1.2.3.4 Test1"), string(out))
	})

	t.Run("print out line that match CIDR blocks [carriage return]", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := "Test1 1.2.3.4 Test1\r\n"
		err := processLine(&processLineInput{
			output:  []rune(input),
			blocks:  blocks,
			inverse: false,
		})
		is.NoError(err)

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		is.Equal(fmt.Sprintf("%s\n", "Test1 1.2.3.4 Test1"), string(out))
	})

	t.Run("no matching line [passthru]", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := "Test1"
		err := processLine(&processLineInput{
			output:  []rune(input),
			blocks:  blocks,
			inverse: false,
		})
		is.NoError(err)

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		is.Equal("", string(out))
	})

	t.Run("skip lines with invalid IPs", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := "Error 1.2.3.777 Error"
		err := processLine(&processLineInput{
			output:  []rune(input),
			blocks:  blocks,
			inverse: false,
		})
		is.NoError(err)

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		is.Equal("", string(out))
	})

	t.Run("print lines with invalid IPs when inverse is set", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := "Error 1.2.3.777 Error"
		err := processLine(&processLineInput{
			output:  []rune(input),
			blocks:  blocks,
			inverse: true,
		})
		is.NoError(err)

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		is.Equal(fmt.Sprintf("%s\n", input), string(out))
	})

	t.Run("print out line that does NOT match CIDR blocks [inverse]", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		input := "Test 1.5.3.4 Test"
		err := processLine(&processLineInput{
			output:  []rune(input),
			blocks:  blocks,
			inverse: true,
		})
		is.NoError(err)

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		is.Equal(fmt.Sprintf("%s\n", input), string(out))
	})
}

func Test_processReader(t *testing.T) {
	is := assert.New(t)

	blocks, _ := parseInputCIDRs("1.2.0.0/16,1.3.0.0/16,123.1.4.0/8")

	sample := `test [255.255.255.255 123.1.4.127 256.256.256.256 123.1.4.1] test
256.256.256.256
8.8.8.8
1.2.4.5
test 1.2.3.256 test
0.0.0.0.0.0.0.0
1.2.3.777`

	t.Run("processes many lines and outputs matches", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := processReader(&processReaderInput{
			reader:  bufio.NewReader(strings.NewReader(sample)),
			blocks:  blocks,
			inverse: false,
		})
		is.NoError(err)

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		is.Equal("test [255.255.255.255 123.1.4.127 256.256.256.256 123.1.4.1] test\n1.2.4.5\n", string(out))
	})

	t.Run("processes many lines and outputs non-matches [inverse]", func(t *testing.T) {
		rescueStdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		err := processReader(&processReaderInput{
			reader:  bufio.NewReader(strings.NewReader(sample)),
			blocks:  blocks,
			inverse: true,
		})
		is.NoError(err)

		_ = w.Close()
		out, _ := ioutil.ReadAll(r)
		os.Stdout = rescueStdout
		is.Equal("256.256.256.256\n8.8.8.8\ntest 1.2.3.256 test\n0.0.0.0.0.0.0.0\n1.2.3.777\n", string(out))
	})
}
