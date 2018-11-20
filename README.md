# ipv666

`ipv666` is a set of tools that enables the discovery of IPv6 addresses both in the global IPv6 address space and in more narrow IPv6 network ranges. These tools are designed to work out of the box with minimal knowledge of their workings.

If you're interested in how these tools work please refer to *this blog post*.

The tools included in this codebase are as follows:

* [`666scan`](#666scan) - Locates live hosts over IPv6 using statistical modeling and ICMP ping scans
* [`666alias`](#666alias) - Tests a single IPv6 network range to see if the network range is aliased
* [`666blgen`](#666blgen) - Adds the contents of a file containing IPv6 network ranges to the aliased network blacklist
* [`666clean`](#666clean) - Cleans the contents of a file containing IPv6 addresses based on an aliased network blacklist

Unless you're doing more complicated IPv6 research it is likely that the [`666scan`](#666scan) tool is what you're looking for. 

To get started check out the [Installation](#installation) section first and then head to whichever section details the tool you're looking to use.

This software is brought to you free of charge by [@_lavalamp](https://twitter.com/_lavalamp) and [@marcnewlin](https://twitter.com/marcnewlin), and we hope that you find it useful. If you do find it useful and you'd like to support our continued contributions to the codebase, please consider donating via any of the following:

* **PayPal** - [paypal.me/thelavalamp](https://www.paypal.me/thelavalamp)
* **BTC** - 371FzLrE7dzd3cZNjDytSyV5hDhDpLj1Mr
* **ETH** - 0x2A35C6987a7E2515EEdB8fB43a7FA86a9Ea917f4
* **LTC** - MGkiBazpfs17ek7DuVJKwzwcFjVcmwrXur

## Installation

`ipv666` is largely self-contained but it does rely on the IPv6 ZMap port authored by a research group out of the TUM Department of Informatics at the Technical University of Munich. We encourage users to [check out their research](https://www.net.in.tum.de/projects/gino/ipv6-hitlist.html).

On a host that has at least one IPv6 interface available pull down the `ZMapv6` code from [GitHub](https://github.com/tumi8/zmap):

```$xslt
git clone https://github.com/tumi8/zmap.git
```

[Follow the instructions found here](https://github.com/tumi8/zmap/blob/master/INSTALL.md) to install `ZMapv6` from source. The exact commands will differ depending on your OS (we've done the following and everything works pretty easily on Ubuntu 16.04 x64).

```$xslt
sudo apt-get install -y build-essential cmake libgmp3-dev gengetopt libpcap-dev flex byacc libjson-c-dev pkg-config libunistring-dev
cd zmap
cmake .
make -j4
sudo make install
```

The `ZMapv6` tool should now be working and found in your `$PATH`:

```$xslt
ubuntu@scanner:~/zmap# which zmap
/usr/local/sbin/zmap
```

We have noticed in some cases that `ZMapv6` fails silently, indicating that no ICMP ping responses are ever received. To ensure that it's working first get the IPv6 address of the interface you'd like to use for scanning:

```$xslt
ubuntu@scanner:~/zmap# ifconfig | grep inet6 | grep Global | awk '{print $3}'
2600:ffff:bbbb:cccc::1234/64
```

Create a file with a test IPv6 address that we know will respond to a ping (in this case we're using an IPv6 address associated with Google):

```$xslt
ubuntu@scanner:~/zmap# dig aaaa +short ipv6.google.com | grep ":" > scantarget
ubuntu@scanner:~/zmap# cat scantarget
2607:f8b0:4005:807::200e
```

We can now use `ZMapv6` to scan this address and confirm that an ICMP ping response is received:

```$xslt
ubuntu@scanner:~/zmap# zmap --bandwidth=20m --output-file=scantest --ipv6-target-file=scantarget --ipv6-source-ip=2600:ffff:bbbb:cccc::1234 --probe-module=icmp6_echoscan
ubuntu@scanner:~/zmap# cat scantest
2607:f8b0:4006:81a::200e
```

If the `scantest` file has the Google IPv6 address in it, then `ZMapv6` is working and you can continue with installing our software. If it doesn't work then... turn it off and on again? Honestly our approach has been to destroy the host and provision a new one and try again.

In order to build `ipv666` you'll need to have the [Golang environment installed](https://golang.org/doc/install). An alternative to pulling the code down and building it on the box where you're scanning from is to pull it down locally, build it locally, and then push the code to a scanning box. While this is possible, please note that there are a number of directories and files that come with the codebase that are required to be on the scanning box. If you do pull/build/push from local, we recommend zipping up the entire `ipv666` directory (after being built) and pushing it to the scanning box. We realize this is not optimal and are open to recommendations for how to better support portability.

Pull down the `ipv666` code using `go get`:

```$xslt
go get github.com/lavalamp-/ipv666
```

We can now test the code to make sure that everything looks good.

```$xslt
cd $GOPATH/github.com/lavalamp-/ipv666
make get-deps
make test
```

If all of the tests pass then do the following to build the tools:

```$xslt
make build-all
```

The binaries should now be present in the `build` directory:

```$xslt
ubuntu@scanner:~/gocode/src/github.com/lavalamp-/ipv666# ls -l build/
total 29996
-rwxr-xr-x 1 ubuntu ubuntu  3469112 Nov 20 19:37 666alias
-rwxr-xr-x 1 ubuntu ubuntu  7232230 Nov 20 19:37 666blgen
-rwxr-xr-x 1 ubuntu ubuntu  7218927 Nov 20 19:37 666clean
-rwxr-xr-x 1 ubuntu ubuntu 12787808 Nov 20 19:37 666scan
```

The last thing we need to do is fill out some of the fields in the `config.json` file. There are many values in here and most of them are irrelevant for most users. The values that either **must** be changed or are likely of interest to most users are below:

```$xslt
{
  ...
  "GenerateAddressCount": 10000000,         // The number of IPv6 addresses generated and tested for each loop of the scanner's state machine
  "ZmapExecPath": "/usr/local/sbin/zmap",   // The path to the ZMap executable 
  "ZmapBandwidth": "20M",                   // The maximum bandwidth for ZMap scans
  "ZmapSourceAddress": "[REPLACE]",         // The IPv6 address to scan from (identified above)
  ...
}
```

Note that you **must** put the correct IPv6 address in the `ZmapSourceAddress` field above. The rest can optionally be updated.

Once the configuration file is updated you should be good to go with any of the commands listed below.

## 666scan

The `666scan` tool is the main workhorse of this toolset. It uses some fairly complicated statistical modeling, analysis, and blacklisting to predict legitimate IPv6 addresses and scan for their presence. More details on how exactly this tool works can be found *in this blog post*.

Please note that any networks that you scan with this tool will receive a considerable amount of traffic for a significant variety of IPv6 addresses. In some cases the networking infrastructure that is carrying your traffic will be unhappy and may either fall over and/or block you. We recommend exercising caution when using this tool (especially for targeted network scans) and turning down the `ZmapBandwidth` value in the configuration file as needed.

### Usage

```$xslt
Usage of ./build/666scan:
  -config string
    	Local file path to the configuration file to use. (default "config.json")
  -force
    	Whether or not to force accept all prompts (useful for daemonized scanning).
  -input string
    	An input file containing IPv6 addresses to initiate scanning from.
  -input-type string
    	The type of file pointed to by the 'input' argument (bin or txt). (default "txt")
  -network string
    	The target IPv6 network range to scan in. If empty, defaults to 2000::/4
  -output string
    	The path to the file where discovered addresses should be written.
  -output-type string
    	The type of output to write to the output file (txt or bin). (default "txt")
```

### Examples

Scan the global IPv6 address space using the configuration values in the file `/foo/bar/config.json`:

```$xslt
./build/666scan -config /foo/bar/config.json
```

Scan the network `2600:6000::/32` using the configuration values in a file named `config.json` residing in the present working directory and write the results in hexadecimal format to the file `/tmp/2600_3000__32_results.txt`:

```$xslt
./build/666scan -network 2600:6000::/32 -output /tmp/2600_3000__32_results
```

## 666alias

The `666alias` tool will test a target network to see if it exhibits traits of being an aliased network (ie: all addresses in the range respond to ICMP pings). If the target network is aliased it will perform a binary search to find the exact network length for how large the aliased network is.

### Usage

```$xslt
Usage of ./build/666alias:
  -config string
    	Local file path to the configuration file to use. (default "config.json")
  -net string
    	An IPv6 CIDR range to test as an aliased network.
```

### Examples

Check if the network at `2600:9000:2173:6d50:5dca:2d48::/96` is aliased using the configuration values in the `config.json` file in the present working directory:

```$xslt
./build/666alias -net 2600:9000:2173:6d50:5dca:2d48::/96
``` 

Check if the network at `2600:9000:2173:6d50:5dca:2d48::/96` is aliased using the configuration values in the file at `/tmp/config.json`:

```$xslt
./build/666alias -net 2600:9000:2173:6d50:5dca:2d48::/96 -config /tmp/config.json
``` 

## 666blgen

The `666blgen` tool processes the content of a file containing IPv6 CIDR ranges (new-line delimited) and adds all of the network ranges to either (1) a new blacklist or (2) your existing blacklist. These blacklists are automatically located and loaded from specific file paths during the operation of [`666scan`](#666scan), [`666alias`](#666alias), and [`666clean`](#666clean).

You will be prompted after invocation asking whether you'd like to create a new blacklist or add these new networks to your existing blacklist.

### Usage

```$xslt
Usage of ./build/666blgen:
  -config string
    	Local file path to the configuration file to use. (default "config.json")
  -input string
    	An input file containing IPv6 network ranges to build a blacklist from.
```

### Examples

Add the IPv6 CIDR ranges found in the file `/tmp/addrranges` to a blacklist based on the configuration values found in a `config.json` file in your present working directory:

```$xslt
./build/666blgen -input /tmp/addrranges
```

Add the IPv6 CIDR ranges found in the file `/tmp/addrranges` to a blacklist based on the configuration values found in the file `/tmp/config.json`:

```$xslt
./build/666blgen -input /tmp/addrranges -config /tmp/config.json
```

## 666clean

The `666clean` tool processes the content of a file containing IPv6 addresses (new-line delimited), removes all the addresses that are found within blacklisted networks, and writes the results to an output file. This tool is an easy way to remove addresses in aliased network ranges from a set of IP addresses.

### Usage

```$xslt
Usage of ./build/666clean:
  -blacklist string
    	The local file path to the blacklist to use. If not specified, defaults to the most recent blacklist in the configured blacklist directory.
  -config string
    	Local file path to the configuration file to use. (default "config.json")
  -input string
    	An input file containing IPv6 addresses to clean via a blacklist.
  -out string
    	The file path where the cleaned results should be written to.
```

### Examples

Process the IPv6 addresses in the file `/tmp/addresses`, remove all addresses found in the most up-to-date blacklist found in the default file path, and write the results to `/tmp/cleanedaddrs` (using configuration values in the `config.json` file in the present working directory):

```$xslt
./build/666clean -input /tmp/addresses -out /tmp/cleanedaddrs
```

Process the IPv6 addresses in the file `/tmp/addresses`, remove all addresses found in the blacklist in the file at `/tmp/blacklist`, and write the results to `/tmp/cleanedaddrs` using configuration values from the file `/tmp/config.json`:

```$xslt
./build/666clean -input /tmp/addresses -blacklist /tmp/blacklist -out /tmp/cleanedaddrs -config /tmp/config.json
```

## License

This software is licensed via the [GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html).

We invite people to contribute to the codebase, fork it, do whatever you'd like! The only requirement that we have with this license is that derivative work is similarly open sourced. 

## Thanks

Many thanks to the following people for their contributions, inspirations, and help.

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
