% cidr 8
# NAME
cidr - tool for checking IPs against CIDR blocks
# SYNOPSIS
cidr

# DESCRIPTION

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



# COMMAND TREE

- [check, c](#check-c)
- [filter, f](#filter-f)
- [version, v](#version-v)

**Usage**:
```
cidr [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# COMMANDS

## check, c

Check IP against range of CIDR blocks

```
IPs that are provided without a mask will be assumed to be /32

CIDR blocks require a mask be provided.
```

**--blocks, -b**="": CIDR blocks to be checked (csv)

**--ips, -i**="": CSV list of IPs with masks (csv)

## filter, f

Filters lines in log files of pipe input

```
Filters lines in log files of pipe input, printing to STDOUT the lines that
contain an IP that is within the provided CIDR blocks.

NOTE: If '--path, -p' is NOT set, then a pipe is assumed to be the input.

The '--inverse, -i' flag will output the lines that do not contain an IP
within the provided CIDR blocks. 
```

**--blocks, -b**="": CIDR blocks to be checked (csv)

**--inverse, -i**: Print out lines that DO NOT match the CIDR check

**--path, -p**="": File path to files to filter, can be a glob. If not set, a pipe is assumed.

## version, v

Print version info

