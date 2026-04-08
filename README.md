# goping — A Feature-Rich Ping Tool in Go

`ping` implementation in Go with all standard features plus an **ASCII latency histogram** extra feature.

## Requirements

- Go 1.21+
- `golang.org/x/net` package
- Linux/macOS (raw ICMP sockets require root/sudo on most systems)

## Build

```bash
cd go-ping/
go mod tidy
chmod a+x build.sh
./build.sh
```

## Usage

```
Usage: ping [options] <host>

Options:
  -c <count>      Stop after sending N packets (default: infinite)
  -i <interval>   Interval between packets (default: 1s, e.g. 500ms)
  -W <timeout>    Per-packet response timeout (default: 3s)
  -s <size>       Payload size in bytes (default: 56)
  -t <ttl>        IP TTL / Hop Limit (default: 64)
  -f              Flood mode — send as fast as possible
  -q              Quiet — only print final summary
  -6              Force IPv6
  -w <deadline>   Exit after this duration (e.g. 10s, 1m)
  -H              Show ASCII latency histogram after summary [EXTRA FEATURE]
```

## Examples

```bash
# Basic ping (infinite, Ctrl+C to stop)
sudo ./goping google.com

# Send 5 packets with 500ms interval
sudo ./goping -c 5 -i 500ms google.com

# Quiet mode — only print the summary
sudo ./goping -q -c 10 8.8.8.8

# Flood ping 
sudo ./goping -f -c 1000 192.168.1.1

# Ping with a 15-second deadline regardless of packet count
sudo ./goping -w 15s google.com

# 20 pings + ASCII latency histogram  ← EXTRA FEATURE
sudo ./goping -c 20 -H google.com

# Force IPv6
sudo ./goping -6 -c 5 google.com
```

## Sample Output

```
PING google.com (142.250.70.46): 56 bytes of data
64 bytes from 142.250.70.46: icmp_seq=0 ttl=64 time=12.345 ms
64 bytes from 142.250.70.46: icmp_seq=1 ttl=64 time=11.876 ms
...

--- google.com ping statistics ---
10 packets transmitted, 10 received, 0.0% packet loss, time 9012ms
rtt min/avg/max/stddev = 11.234/12.101/14.567/0.987 ms

--- Latency Histogram ---
  11.23 -  11.54 ms | ########                                  2
  11.54 -  11.85 ms | ################                          4
  11.85 -  12.16 ms | ########################################  10
  12.16 -  12.47 ms | ################                          4
  ...
```

## Features

| Feature | Flag | Notes |
|---|---|---|
| Packet count | `-c N` | Standard |
| Interval | `-i duration` | e.g. `-i 200ms` |
| Per-packet timeout | `-W duration` | Default 3s |
| Payload size | `-s bytes` | Default 56 |
| TTL | `-t N` | Default 64 |
| Flood mode | `-f` | Requires root |
| Quiet mode | `-q` | Summary only |
| IPv6 | `-6` | Force IPv6 |
| Deadline | `-w duration` | Hard stop after N time |
| **Latency histogram** | **`-H`** | **Extra feature — ASCII bar chart of RTT distribution** |

## Notes

- Raw ICMP sockets require elevated privileges. Run with `sudo` or grant the binary `CAP_NET_RAW`:
  ```bash
  sudo setcap cap_net_raw+ep ./goping
  ```
- The latency histogram (`-H`) is most useful with `-c 20` or more samples.
