% cidr 8
# NAME
cidr - tool for checking IPs against CIDR blocks
# SYNOPSIS
cidr


# COMMAND TREE

- [check](#check)
- [pipe](#pipe)
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

## pipe

command | cidr pipe

**--blocks, -b**="": CIDR blocks to be checked (csv)

## version, v

Print version info
