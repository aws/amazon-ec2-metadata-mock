# AEMM Usage
This page serves as documentation for AEMM's advanced use cases and behavior.


## Commands
AEMM's supported commands (`spot`, `events`) are viewed using `--help`:
```
$ ec2-metadata-mock --help

...
Available Commands:
  events          Mock EC2 maintenance events
  help            Help about any command
  spot            Mock EC2 Spot interruption notice

...
```
commands are designed as follows:
* Run independently from other commands
  * i.e. when AEMM is started with `events` subcommand, `spot` routes will **NOT** be available and vice-versa 
* Local flag availability so that commands can be configured directly via CLI parameters
    * With validation checks
* Contain additional `--help` documentation
* Default values sourced from code

Metadata categories that are always available, irrespective of the CLI command run are referred to as **static metadata**.

### Spot Interruption
To view the available flags for the Spot Interruption command use `spot --help`:
```
$ ec2-metadata-mock spot --help
Mock EC2 Spot interruption notice

Usage:
  ec2-metadata-mock spot [--action ACTION] [flags]

Aliases:
  spot, spotitn

Examples:
  ec2-metadata-mock spot -h 	spot help
  ec2-metadata-mock spot -d 5 --action terminate		mocks spot interruption only

Flags:
  -h, --help                      help for spot
  -a, --action string             action in the spot interruption notice (default: terminate)
                                  action can be one of the following: terminate,hibernate,stop
  -t, --time string               time specifies the approximate time when the spot instance will receive the shutdown signal in RFC3339 format to execute instance action E.g. 2020-01-07T01:03:47Z (default: request time + 2 minutes in UTC)

Global Flags:
  -c, --config-file string         config file for cli input parameters in json format (default: $HOME/aemm-config.json)
  -n, --hostname string            the HTTP hostname for the mock url (default: 0.0.0.0)
  -I, --imdsv2                     whether to enable IMDSv2 only, requiring a session token when submitting requests (default: false, meaning both IMDS v1 and v2 are enabled)
  -d, --mock-delay-sec int         mock delay in seconds, relative to the application start time (default: 0 seconds)
  -x, --mock-ip-count int          number of IPs in a cluster that can receive a Spot Interrupt Notice and/or Scheduled Event (default 2)
      --mock-trigger-time string   mock trigger time in RFC3339 format. This takes priority over mock-delay-sec (default: none)
  -p, --port string                the HTTP port where the mock runs (default: 1338)
  -s, --save-config-to-file        whether to save processed config from all input sources in .ec2-metadata-mock/.aemm-config-used.json in $HOME or working dir, if homedir is not found (default: false)
```

1.) **Overriding `spot::action` via CLI flag**:

```
$ ec2-metadata-mock spot -a stop
Initiating amazon-ec2-metadata-mock for EC2 Spot interruption notice on port 1338

Flags:
action: stop

Serving the following routes: ... (truncated for readability)
```

Send the request:
```
$ curl localhost:1338/latest/meta-data/spot/instance-action
{
	"action": "stop",
	"time": "2020-04-24T17:11:44Z"
}
```

### Events
Similar to spot, the `events` command, view the local flags using `events --help`:

```
$ ec2-metadata-mock events --help
Mock EC2 Events

Usage:
  ec2-metadata-mock events [--code CODE] [--state STATE] [--not-after] [--not-before-deadline] [flags]

Aliases:
  events, se, scheduledevents

Examples:
  ec2-metadata-mock events -h 	events help
  ec2-metadata-mock events -o instance-stop --state active -d		mocks an active and upcoming scheduled event for instance stop with a deadline for the event start time

Flags:
  -o, --code string                  event code in the scheduled event (default: system-reboot)
                                     event-code can be one of the following: instance-reboot,system-reboot,system-maintenance,instance-retirement,instance-stop
  -h, --help                         help for events
  -a, --not-after string             the latest end time for the scheduled event in RFC3339 format E.g. 2020-01-07T01:03:47Z default: application start time + 7 days in UTC))
  -b, --not-before string            the earliest start time for the scheduled event in RFC3339 format E.g. 2020-01-07T01:03:47Z (default: application start time in UTC)
  -l, --not-before-deadline string   the deadline for starting the event in RFC3339 format E.g. 2020-01-07T01:03:47Z (default: application start time + 9 days in UTC)
  -t, --state string                 state of the scheduled event (default: active)
                                     state can be one of the following: active,completed,canceled

(Truncated Global Flags for readability)
```

1.) **Starting AEMM with `events` invalid flag overrides**: as noted above, all commands have validation logic for overrides via CLI flags. If the user attempts to pass an invalid override value, then AEMM will panic, kill the server, and return an error message with what went wrong:

```
$ ec2-metadata-mock events --code FOO

panic: Fatal error while executing the root command: Invalid CLI input "FOO" for flag code. 
Allowed value(s): instance-reboot,system-reboot,system-maintenance,instance-retirement,instance-stop.
```

## Static Metadata
Additional properties of static metadata:
* delays do **NOT** affect static metadata availability
* values are overridden via config file and/or env variables **only**
  * values **cannot** be overridden using flags


### Static Metadata Overrides & Path Substitutions
Some metadata categories use data unique to the instance as part of the query.
* **Example:**  network/interfaces/macs/*MAC*/interface-id where *MAC* is a placeholder for the instance's mac address
  * [AWS Documentation](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-categories.html)

Applying overrides to these *placeholder values* will automatically update paths containing these values **unless** the path itself is explicitly overridden.

1.) **Starting AEMM with JSON config overrides:** The static metadata will use the overridden values AND paths using overridden instance data will be updated as well.

* config-overrides.json:

```
{
    "metadata": {
        "paths": {
            "mac-device-number": "/latest/meta-data/network/interfaces/macs/BAR/device-number"
        },
        "values": {
            "mac": "FOO"
        }
    }
}

```

* start the server with overrides:

```
$ ec2-metadata-mock -c config-overrides.json
Initiating amazon-ec2-metadata-mock for all mocks on port 1338

Flags:
config-file: config-overrides.json
```

* querying the available routes will show updated placeholder paths **except** for those explicitly overridden in the "paths" blob of *config-overrides.json*:

```
$ curl localhost:1338/latest/meta-data

network/interfaces/macs/BAR/device-number
network/interfaces/macs/FOO/interface-id
network/interfaces/macs/FOO/ipv4-associations/192.0.2.54
network/interfaces/macs/FOO/ipv6s
network/interfaces/macs/FOO/local-hostname
network/interfaces/macs/FOO/local-ipv4s
network/interfaces/macs/FOO/mac
network/interfaces/macs/FOO/owner-id
network/interfaces/macs/FOO/public-hostname
network/interfaces/macs/FOO/public-ipv4s
network/interfaces/macs/FOO/security-group-ids
network/interfaces/macs/FOO/security-groups
network/interfaces/macs/FOO/subnet-id
network/interfaces/macs/FOO/subnet-ipv4-cidr-block
network/interfaces/macs/FOO/subnet-ipv6-cidr-blocks
network/interfaces/macs/FOO/vpc-id
network/interfaces/macs/FOO/vpc-ipv4-cidr-block
network/interfaces/macs/FOO/vpc-ipv4-cidr-blocks
network/interfaces/macs/FOO/vpc-ipv6-cidr-blocks
```

* querying the mac address will reflect overridden value:

```
$ curl http://localhost:1338/latest/meta-data/mac
FOO
```