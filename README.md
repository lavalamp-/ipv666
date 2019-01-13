# ipv666

**NOTE** - This software was previously licensed under GPLv3 and has since been updated to BSD.

`ipv666` is a set of tools that enables the discovery of IPv6 addresses both in the global IPv6 address space and in more narrow IPv6 network ranges. These tools are designed to work out of the box with minimal knowledge of their workings.

If you're interested in how these tools work please refer to [this blog post](https://l.avala.mp/?p=285).

The tools included in this codebase are as follows:

* [`discover`](#discover) - Locates live hosts over IPv6 using statistical modeling and ICMP ping scans
* [`alias`](#alias) - Tests a single IPv6 network range to see if the network range is aliased
* [`blgen`](#blgen) - Adds the contents of a file containing IPv6 network ranges to the aliased network blacklist
* [`clean`](#clean) - Cleans the contents of a file containing IPv6 addresses based on an aliased network blacklist
* [`convert`](#convert) - Converts the contents of a file containing IPv6 addresses to another IP address representation

Unless you're doing more complicated IPv6 research it is likely that the [`discover`](#discover) tool is what you're looking for. 

To get started check out the [Installation](#installation) section first and then head to whichever section details the tool you're looking to use.

This software is brought to you free of charge by [@_lavalamp](https://twitter.com/_lavalamp) and [@marcnewlin](https://twitter.com/marcnewlin), and we hope that you find it useful. If you do find it useful and you'd like to support our continued contributions to the codebase, please consider donating via any of the following:

* **PayPal** - [paypal.me/thelavalamp](https://www.paypal.me/thelavalamp)
* **BTC** - 371FzLrE7dzd3cZNjDytSyV5hDhDpLj1Mr
* **ETH** - 0x2A35C6987a7E2515EEdB8fB43a7FA86a9Ea917f4
* **LTC** - MGkiBazpfs17ek7DuVJKwzwcFjVcmwrXur

## Installation

As of `v0.2` `ipv666` should be much easier to install!

First you'll need to have the [Golang environment installed](https://golang.org/doc/install#install). As `ipv666` makes use of Go modules we recommend using Golang `v1.11` and above to ensure a working installation.

Once you have Golang `v1.11` or above installed, do the following:

```$xslt
go get github.com/lavalamp-/ipv666/ipv666
```

Once this command completes you should have the `ipv666` binary on your path. If you don't, double check to make sure that `$GOPATH/bin` is in your `PATH`.

## discover

The `discover` tool is the main workhorse of this toolset. It uses some fairly complicated statistical modeling, analysis, and blacklisting to predict legitimate IPv6 addresses and scan for their presence. More details on how exactly this tool works can be found [in this blog post](https://l.avala.mp/?p=285).

Please note that any networks that you scan with this tool will receive a considerable amount of traffic for a significant variety of IPv6 addresses. In some cases the networking infrastructure that is carrying your traffic will be unhappy and may either fall over and/or block you. We recommend exercising caution when using this tool (especially for targeted network scans) and choosing a `bandwidth` value with care (default is currently 20 Mbps).

### Usage

```$xslt
This utility scans for live hosts over IPv6 based on the network range you specify. If no range is specified, then this utility scans the global IPv6 address space (e.g. 2000::/4). The scanning process generates candidate addresses, scans for them, tests the network ranges where live addresses are found for aliased conditions, and adds legitimate discovered IPv6 addresses to an output list.

Usage:
  ipv666 scan discover [flags]

Flags:
  -h, --help                 help for discover
  -o, --output string        The path to the file where discovered addresses should be written.
  -t, --output-type string   The type of output to write to the output file (txt or bin).

Global Flags:
  -b, --bandwidth string   The maximum bandwidth to use for ping scanning
  -f, --force              Whether or not to force accept all prompts (useful for daemonized scanning).
  -l, --log string         The log level to emit logs at (one of debug, info, success, warn, error).
  -n, --network string     The IPv6 CIDR range to scan.
```

### Examples

Scan the global address space for live hosts over IPv6 with a maximum speed of 20 Mbps and write the results to a file entitled `discovered_addrs.txt`:

```$xslt
ipv666 scan discover
```

Scan the network `2600:6000::/32` for live hosts over IPv6 with a maximum speed of 10 Mbps and write the results to a file entitled `addresses.txt`:
```$xslt
ipv666 scan discover -b 10M -o addresses.txt -n 2600:6000::32
```

## alias

The `alias` tool will test a target network to see if it exhibits traits of being an aliased network (ie: all addresses in the range respond to ICMP pings). If the target network is aliased it will perform a binary search to find the exact network length for how large the aliased network is.

### Usage

```$xslt
A utility for testing whether or not a network range exhibits traits of an aliased network range. Aliased network ranges are ranges in which every host responds to a ping request, thereby making it look like the range is full of IPv6 hosts. Pointing this utility at a network range will let tell you whether or not that network range is aliased and, if it is, the boundary of the network range that is aliased.

Usage:
  ipv666 scan alias [flags]

Flags:
  -h, --help   help for alias

Global Flags:
  -b, --bandwidth string   The maximum bandwidth to use for ping scanning
  -f, --force              Whether or not to force accept all prompts (useful for daemonized scanning).
  -l, --log string         The log level to emit logs at (one of debug, info, success, warn, error).
  -n, --network string     The IPv6 CIDR range to scan.
```

### Examples

Test the network range `2600:9000:2173:6d50:5dca:2d48::/96` to see if it's aliased with a maximum speed of 20 Mbps:

```$xslt
ipv666 scan alias -n 2600:9000:2173:6d50:5dca:2d48::/96
```

Test the network range `2600:9000:2173:6d50:5dca:2d48::/96` to see ifit's aliased with maximum speed of 10 Mbps and show debug logging:
```$xslt
ipv666 scan alias -n 2600:9000:2173:6d50:5dca:2d48::/96 -b 10M -l debug
```

## blgen

The `blgen` tool processes the content of a file containing IPv6 CIDR ranges (new-line delimited) and adds all of the network ranges to either (1) a new blacklist or (2) your existing blacklist. These blacklists are automatically located and loaded from specific file paths during the operation of [`discover`](#discover), [`alias`](#alias), and [`clean`](#clean).

You will be prompted after invocation asking whether you'd like to create a new blacklist or add these new networks to your existing blacklist.

### Usage

```$xslt
This utility takes a list of IPv6 CIDR ranges from a text file (new-line delimited), adds them to the current network blacklist, and sets the new blacklist as the one to use for the 'scan' command.

Usage:
  ipv666 blgen [flags]

Flags:
  -h, --help           help for blgen
  -i, --input string   An input file containing IPv6 network ranges to build a blacklist from.

Global Flags:
  -f, --force        Whether or not to force accept all prompts (useful for daemonized scanning).
  -l, --log string   The log level to emit logs at (one of debug, info, success, warn, error).
```

### Examples

Add the IPv6 CIDR ranges found in the file `/tmp/addrranges` to a blacklist:

```$xslt
ipv666 blgen -i /tmp/addrranges
```

Add the IPv6 CIDR ranges found in the file `/tmp/addrranges` to a blacklist and force accept all prompts:
```$xslt
ipv666 blgen -i /tmp/addrranges -f
```

## clean

The `clean` tool processes the content of a file containing IPv6 addresses (new-line delimited), removes all the addresses that are found within blacklisted networks, and writes the results to an output file. This tool is an easy way to remove addresses in aliased network ranges from a set of IP addresses.

### Usage

```$xslt
This utility will clean the contents of an IPv6 address file (new-line delimited, standard ASCII hex representation) based on the contents of an IPv6 network blacklist file. If no blacklist path is supplied then the utility will use the default blacklist. The cleaned results will then be written to an output file.

Usage:
  ipv666 clean [flags]

Flags:
  -b, --blacklist string   The local file path to the blacklist to use. If not specified, defaults to the most recent blacklist in the configured blacklist directory.
  -h, --help               help for clean
  -i, --input string       An input file containing IPv6 addresses to clean via a blacklist.
  -o, --out string         The file path where the cleaned results should be written to.

Global Flags:
  -f, --force        Whether or not to force accept all prompts (useful for daemonized scanning).
  -l, --log string   The log level to emit logs at (one of debug, info, success, warn, error).
```

### Examples

Process the IPv6 addresses in the file `/tmp/addresses`, remove all addresses found in the most up-to-date blacklist found in the default file path, and write the results to `/tmp/cleanedaddrs`):

```$xslt
ipv666 clean -i /tmp/addresses -o /tmp/cleanedaddrs
```

Process the IPv6 addresses in the file `/tmp/addresses`, remove all addresses found in the blacklist at `/tmp/blacklist`, and write the results to `/tmp/cleanedaddrs`):

```$xslt
ipv666 clean -i /tmp/addresses -o /tmp/cleanedaddrs -b /tmp/blacklist
```

## convert

The `convert` tool is useful for converting a file containing IPv6 addresses to different file formats. It currently supports the three different output types of `txt` (standard ASCII hex IPv6 addresses), `bin` (the raw 16 bytes of all input addresses are written sequentially to a file) and `hex` (the full 32 character ASCII hex representation is written to a file delimited by new lines).

### Usage

```$xslt
This utility will process the contents of a file as containing IPv6 addresses, convert
those addresses to another format, and then write a new file with the same addresses in
the new format. This functionality is (hopefully) intelligent enough to determine how the
addresses are stored in the file without having to specify an input type.

Usage:
  ipv666 convert [flags]

Flags:
  -h, --help           help for convert
  -i, --input string   The file to process IPv6 addresses out of.
  -o, --out string     The file path to write the converted file to.
  -t, --type string    The format to write the IPv6 addresses in (one of 'txt', 'bin', 'hex').

Global Flags:
  -f, --force        Whether or not to force accept all prompts (useful for daemonized scanning).
  -l, --log string   The log level to emit logs at (one of debug, info, success, warn, error).
```

### Examples

Convert the contents of the file at `/tmp/addresses` (in standard text format) to binary format and write the results to `/tmp/out`:

```$xslt
ipv666 convert -i /tmp/addresses -o /tmp/out -t bin
```

Convert the contents of the file at `/tmp/addresses` (in binary format) to fat hex format and write the results to `/tmp/out`:

```$xslt
ipv666 convert -i /tmp/addresses -o /tmp/out -t hex
```

## License

This software is licensed via the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).

We invite people to contribute to the codebase, fork it, do whatever you'd like! The only requirement that we have with this license is that derivative work is similarly open sourced. 

## Thanks

Many thanks to the following people for their contributions, inspiration, and help.

* Vasyl Pihur
* Zakir Durumeric
* David Adrian
* Eric Wustrow
* J. Alex Halderman
* Paul Pearce
* Ariana Mirian
* HD Moore
* Oliver Gasser
* Quirin Scheitle
* Tobias Fiebig
* Matthew Bryant
