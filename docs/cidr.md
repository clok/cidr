% cidr 8
# NAME
cidr - tool for checking IPs against CIDR blocks
# SYNOPSIS
cidr


# COMMAND TREE

- [check](#check)
- [filter](#filter)
- [version, v](#version-v)

**Usage**:
```
cidr [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# COMMANDS

## check

check IP against range of CIDR blocks

**--blocks, -b**="": CIDR blocks to be checked (csv)

**--ips, -i**="": CSV list of IPs with masks (csv)

## filter

Filters lines in log files of pipe input

```
Filters lines in log files of pipe input, printing to STDOUT the lines that
contain an IP that is within the provided CIDR blocks.

NOTE: If '--path, -p' is NOT set, then a pipe is assumed to be the input.

The '--inverse, -i' flag will output the lines that do not contain an IP
within the provided CIDR blocks. 
```

**--blocks, -b**="": CIDR blocks to be checked (csv)

**--inverse, -i**: print out lines that DO NOT match the CIDR check

**--path, -p**="": File path to files to filter, can be a glob. If not set, a pipe is assumed.

## version, v

Print version info

