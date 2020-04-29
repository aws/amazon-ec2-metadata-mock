# Amazon EC2 Metadata Mock

**Amazon EC2 Metadata Mock (AEMM)** is a tool to simulate [Amazon EC2 instance metadata service](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html) for local testing.

<p>
   <a href="https://golang.org/doc/go1.13">
   <img src="https://img.shields.io/github/go-mod/go-version/aws/amazon-ec2-metadata-mock?color=blueviolet" alt="go-version">
   </a>
   <a href="https://opensource.org/licenses/Apache-2.0">
   <img src="https://img.shields.io/badge/License-Apache%202.0-ff69b4.svg?color=orange" alt="license">
   </a>
   <a href="https://travis-ci.org/aws/amazon-ec2-instance-advisor">
   <img src="https://travis-ci.org/aws/amazon-ec2-instance-advisor.svg?branch=master" alt="build-status">
   </a>
   <a href="https://codecov.io/gh/aws/amazon-ec2-instance-advisor">
   <img src="https://img.shields.io/codecov/c/github/aws/amazon-ec2-instance-advisor" alt="build-status">
   </a>
   <a href="https://hub.docker.com/r/amazon/amazon-ec2-instance-advisor">
   <img src="https://img.shields.io/docker/pulls/amazon/amazon-ec2-instance-advisor" alt="docker-pulls">
   </a>
</p>

# Table of Contents

   * [Project Summary](#project-summary)
   * [Major Features](#major-features)
   * [Supported Metadata Categories](#supported-metadata-categories)
   * [Getting Started](#getting-started)
      * [Installation](#installation)
      * [Starting AEMM](#starting-aemm)
      * [Making a Request](#making-a-request)
   * [Configuration](#configuration)
      * [Defaults](#defaults)
      * [Config Precedence (Highest to Lowest)](#config-precedence-highest-to-lowest)
      * [Configuring AEMM](#configuring-aemm)
   * [Usage](#usage)
      * [Commands](#commands)
         * [Spot Interruption](#spot-interruption)
         * [Scheduled Events](#scheduled-events)
      * [Instance Metadata Service Versions](#instance-metadata-service-versions)
      * [Static Metadata](#static-metadata)
         * [Static Metadata Overrides &amp; Path Substitutions](#static-metadata-overrides--path-substitutions)
   * [Troubleshooting](#troubleshooting)
      * [Warnings and Expected Outcome](#warnings-and-expected-outcome)
   * [Building](#building)
   * [Communication](#communication)
   * [Contributing](#contributing)
   * [License](#license)

# Project Summary
AWS EC2 Instance metadata is data about your instance that you can use to configure or manage the running instance. Instance metadata is divided into categories like hostname, instance id, maintenance events, spot instance action. See the complete list of metadata categories [here](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-categories.html).

The instance metadata can be accessed from within the instance. Some instance metadata is available only when an instance is affected by the event. E.g. A Spot instance's metadata item `spot/instance-action` is available only when AWS decides to interrupt the Spot instance.
These bring forth some challenges like not being able to test one's application in the event of Spot interruption or other such events and requiring an EC2 instance for testing.
This project attempts to bridge these gaps by providing mocks for **most** of these metadata categories. The mock responses are designed to replicate those from the actual instance metadata service for accurate, local testing.

# Major Features
- Emulate Spot Instance Interruption (ITN) events
- Delay mock response from the mock serve start time
- Configure metadata in mock responses via CLI flags, config file, env variables
- IMDSv1 and v2 support (configurable for IMDSv2 support only)
- Save processed configuration to a local file

# Supported Metadata Categories
AEMM supports most [metadata categories](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-categories.html) **except for:**
* ancestor-ami-ids
* elastic-gpus/associations/elastic-gpu-id
* events/maintenance/history
* kernel-id
* ramdisk-id
* Dynamic data categories

# Getting Started
AEMM is simple to get up and running.

## Installation
Download binary from the latest release:
```
curl -Lo amazon-ec2-metadata-mock https://github.com/aws/amazon-ec2-metadata-mock/releases/download/latest/amazon-ec2-metadata-mock
```

## Starting AEMM
Use `amazon-ec2-metadata-mock --help` to view examples and explanations of supported flags and commands:

```
$ amazon-ec2-metadata-mock --help

amazon-ec2-metadata-mock is a tool to mock Amazon EC2 instance metadata.

Usage:
  amazon-ec2-metadata-mock <command> [arguments] [flags]
  amazon-ec2-metadata-mock [command]

Examples:
  amazon-ec2-metadata-mock --mock-delay-sec 10	mocks all metadata paths
  amazon-ec2-metadata-mock spotitn --instance-action terminate	mocks spot ITN only

Available Commands:
  help            Help about any command
  scheduledevents Mock EC2 Scheduled Events
  spotitn         Mock EC2 Spot interruption notice

Flags:
  -c, --config-file string    config file for cli input parameters in json format (default: $HOME/ec2-mock-config.json)
  -h, --help                  help for amazon-ec2-metadata-mock
  -n, --hostname string       the HTTP hostname for the mock url (default: localhost)
  -I, --imdsv2                whether to enable IMDSv2 only requiring a session token when submitting requests (default: false)
  -d, --mock-delay-sec int    mock delay in seconds, relative to the application start time (default: 0 seconds)
  -p, --port string           the HTTP port where the mock runs (default: 1338)
  -s, --save-config-to-file   whether to save processed config from all input sources in .amazon-ec2-metadata-mock/.ec2-mock-config-used.json in $HOME or working dir, if homedir is not found (default: false)

Use "amazon-ec2-metadata-mock [command] --help" for more information about a command.
```

Starting AEMM with default configurations using `amazon-ec2-metadata-mock` will start the server on the default host and port:

```
$ amazon-ec2-metadata-mock

Initiating amazon-ec2-metadata-mock for all mocks on port 1338
Serving the following routes: /latest/meta-data/product-codes, /latest/meta-data/iam/info, /latest/meta-data/instance-type, ...(truncated for readability)
```

## Making a Request
With the server running, send a request to one of the supported routes above using `curl localhost:1338/<route>`. Example request to display all supported routes:

```
$ curl localhost:1338/latest/meta-data

ami-id
ami-launch-index
ami-manifest-path
block-device-mapping/ami
block-device-mapping/ebs0
block-device-mapping/ephemeral0
block-device-mapping/root
block-device-mapping/swap
elastic-inference/associations
elastic-inference/associations/eia-bfa21c7904f64a82a21b9f4540169ce1
events/maintenance/scheduled
hostname
iam/info
iam/security-credentials
iam/security-credentials/baskinc-role
instance-action
instance-id
instance-type
latest/api/token
local-hostname
local-ipv4
mac
network/interfaces/macs/0e:49:61:0f:c3:11/device-number
network/interfaces/macs/0e:49:61:0f:c3:11/interface-id
network/interfaces/macs/0e:49:61:0f:c3:11/ipv4-associations/192.0.2.54
network/interfaces/macs/0e:49:61:0f:c3:11/ipv6s
network/interfaces/macs/0e:49:61:0f:c3:11/local-hostname
network/interfaces/macs/0e:49:61:0f:c3:11/local-ipv4s
network/interfaces/macs/0e:49:61:0f:c3:11/mac
network/interfaces/macs/0e:49:61:0f:c3:11/owner-id
network/interfaces/macs/0e:49:61:0f:c3:11/public-hostname
network/interfaces/macs/0e:49:61:0f:c3:11/public-ipv4s
network/interfaces/macs/0e:49:61:0f:c3:11/security-group-ids
network/interfaces/macs/0e:49:61:0f:c3:11/security-groups
network/interfaces/macs/0e:49:61:0f:c3:11/subnet-id
network/interfaces/macs/0e:49:61:0f:c3:11/subnet-ipv4-cidr-block
network/interfaces/macs/0e:49:61:0f:c3:11/subnet-ipv6-cidr-blocks
network/interfaces/macs/0e:49:61:0f:c3:11/vpc-id
network/interfaces/macs/0e:49:61:0f:c3:11/vpc-ipv4-cidr-block
network/interfaces/macs/0e:49:61:0f:c3:11/vpc-ipv4-cidr-blocks
network/interfaces/macs/0e:49:61:0f:c3:11/vpc-ipv6-cidr-blocks
placement/availability-zone
product-codes
public-hostname
public-ipv4
public-keys/0/openssh-key
reservation-id
security-groups
services/domain
services/partition
spot/instance-action
spot/termination-time
```
        

# Configuration
AEMM's wide-range of configurability ranges from overriding port numbers to enabling IMDSv2-only to updating specific metadata values and paths.
These configurations can be loaded from various sources with a deterministic precedence.

## Defaults
Defaults for AEMM configuration are sourced throughout code. Examples below:
* **CLI flags**
  * [server config defaults](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/pkg/config/server.go#L22) 
* **Metadata mock responses**
  * [aemm-metadata-default-values.json](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/pkg/config/defaults/aemm-metadata-default-values.json)
* **Commands**
  * [scheduledevents](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/pkg/cmd/scheduledevents/scheduledevents.go#L72) 


## Config Precedence (Highest to Lowest)
1. overrides
2. flag
3. env
4. config
5. key/value store
6. default

For example, if values from the following sources were loaded:
```
Defaults in code
 {
 	"port": "1338",
	"mock-delay-sec": 0,
	"save-config-to-file": true,
   	"spot-itn": {
	   "instance-action": "terminate"
	}
 }
Config (at $HOME/.aemm-config-overrides.json)
 {
   "port": "1440",
   "spot-itn": {
	   "instance-action": "stop"
   }
 }
Env (export AEMM_MOCK_DELAY_SEC=12;export AEMM_SPOT_ITN_INSTANCE_ACTION=stop)
 {
   "mock-delay-sec": 12  
   "spot-itn": {
		"instance-action": "hibernate"
   }
 }
CLI Flags
 {
   "mock-delay-sec": 8  
 }
```

The resulting config will have the following values:
```
 {
	"config-file": "$HOME/.aemm-config-overrides.json",
   	"mock-delay-sec": 8,
   	"port": "1440",
	"save-config-to-file": true,
   	"spot-itn": {
		"instance-action": "hibernate"
   }
 }
```

AEMM is built using Viper which is where the precedence is sourced from. More details can be found in their documentation: 

<a href="https://github.com/spf13/viper/blob/master/README.md">
	<img src="https://img.shields.io/badge/viper-documentation-green" alt="">
</a>

## Configuring AEMM
The tool can be configured in various ways:

1. CLI flags

	Use help commands to learn more

2. Env variables
	```
	$ export AEMM_MOCK_DELAY_SEC=12
	$ export AEMM_SPOT_ITN_INSTANCE_ACTION=stop
	$ env | grep AEMM    	// To list the tool's env variables
	```

    > NOTE the translation of config key `spot-itn.instance-action` to `AEMM_SPOT_ITN_INSTANCE_ACTION` env variable.

3. Configuration file in JSON format at `path/to/config-overrides.json`
```
{
  "metadata": {
  	"paths": {
  		"ipv4-associations": "/latest/meta-data/network/interfaces/macs/0e:49:61:0f:c3:77/ipv4-associations/192.0.2.54"
  	},
    "values": {
      "mac": "0e:49:61:0f:c3:77",
      "public-ipv4": "54.92.157.77"
    }
  },
  "spot-itn": {
    "instance-action": "terminate",
    "time": "2020-01-07T01:03:47Z"
  }
}
```

Use the `-c` flag to consume the configuration file and the `-s` flag to save an output of the configurations used by AEMM after precedence has been applied:

```
$ amazon-ec2-metadata-mock -c path/to/config-overrides.json -s
Successfully saved final configuration to local file  /path/to/home/.amazon-ec2-metadata-mock/.ec2-mock-config-used.json


$ cat $HOME/.amazon-ec2-metadata-mock/.ec2-mock-config-used.json
(truncated for readability)

{
  "config-file": "path/to/config-overrides.json",
  "metadata": {
    "paths": {
      "ipv4-associations": "/latest/meta-data/network/interfaces/macs/0e:49:61:0f:c3:77/ipv4-associations/192.0.2.54"
    },
    "values": {
      "mac": "0e:49:61:0f:c3:77",
      "public-ipv4": "54.92.157.77"
    }
  },
  "mock-delay-sec": 12,
  "save-config-to-file": true,
  "server": {
    "hostname": "localhost",
    "port": "1338"
  },
  "spot-itn": {
    "instance-action": "stop",
    "time": "2020-01-07T01:03:47Z"
  }
}

```

# Usage
AEMM is primarily used as a developer tool to help test behavior related to Metadata Service. Popular use cases include: emulating spot instance interrupts after a designated delay, mocking scheduled maintenance events, and testing IMDSv1 to IMDSv2 migrations.

## Commands
AEMM's supported commands (`spotitn`, `scheduledevents`) are viewed using `--help`:
```
$ amazon-ec2-metadata-mock --help

...
Available Commands:
  help            Help about any command
  scheduledevents Mock EC2 Scheduled Events
  spotitn         Mock EC2 Spot interruption notice

...
```
commands are designed as follows:
* Run independently from other commands
  * i.e. when AEMM is started with `scheduledevents` subcommand, `spotitn` routes will **NOT** be available and vice-versa 
* Local flag availability so that commands can be configured directly via CLI parameters
	* With validation checks 
* Contain additional `--help` documentation
* Default values sourced from code

Metadata categories that are always available, irrespective of the CLI command run are referred to as **static metadata**. More details on static metadata can be found [here](#static-metadata)

### Spot Interruption
To view the available local flags for the Spot Interruption command use `spotitn --help`:
```
$ amazon-ec2-metadata-mock spotitn --help
Mock EC2 Spot interruption notice

Usage:
  amazon-ec2-metadata-mock spotitn [--instance-action ACTION] [flags]

Aliases:
  spotitn, spot, spot-itn, spotItn

Examples:
  amazon-ec2-metadata-mock spotitn -h 	spotitn help
  amazon-ec2-metadata-mock spotitn -d 5 --instance-action terminate		mocks spot interruption only

Flags:
  -h, --help                      help for spotitn
  -a, --instance-action string    instance action in the spot interruption notice (default: terminate)
                                  instance-action can be one of the following: terminate,hibernate,stop
  -t, --termination-time string   termination time specifies the approximate time when the spot instance will receive the shutdown signal in RFC3339 format to execute instance action E.g. 2020-01-07T01:03:47Z (default: request time + 2 minutes in UTC)

Global Flags:
  -c, --config-file string    config file for cli input parameters in json format (default: $HOME/ec2-mock-config.json)
  -n, --hostname string       the HTTP hostname for the mock url (default: localhost)
  -I, --imdsv2                whether to enable IMDSv2 only, requiring a session token when submitting requests (default: false, meaning both IMDS v1 and v2 are enabled)
  -d, --mock-delay-sec int    mock delay in seconds, relative to the application start time (default: 0 seconds)
  -p, --port string           the HTTP port where the mock runs (default: 1338)
  -s, --save-config-to-file   whether to save processed config from all input sources in .amazon-ec2-metadata-mock/.ec2-mock-config-used.json in $HOME or working dir, if homedir is not found (default: false)
```

1.) **Starting AEMM with `spotitn`**:  `spotitn` routes available immediately:
```
$ amazon-ec2-metadata-mock spotitn
Initiating amazon-ec2-metadata-mock for EC2 Spot interruption notice on port 1338
Serving the following routes: ... (truncated for readability)
```
Send the request:
```
$ curl localhost:1338/latest/meta-data/spot/instance-action
{
	"instance-action": "terminate",
	"time": "2020-04-24T17:11:44Z"
}
```


2.) **Starting AEMM with `spotitn` overrides after Delay**: Users can override `instance-action` via CLI flag as well as apply a *delay* duration in seconds for when the `spotitn` metadata will become available (i.e. simulating when AWS sends the spot interrupt notice); static metadata availability is unaffected:

```
$ amazon-ec2-metadata-mock spotitn -a stop -d 10
Initiating amazon-ec2-metadata-mock for EC2 Spot interruption notice on port 1338

Flags:
instance-action: stop
mock-delay-sec: 10

Serving the following routes: ... (truncated for readability)
```

Sending a request to `spotitn` paths before the delay has passed will return **404 - Not Found:**
```
$ curl localhost:1338/latest/meta-data/spot/instance-action

<?xml version="1.0" encoding="iso-8859-1"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN"
	"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
 <head>
  <title>404 - Not Found</title>
 </head>
 <body>
  <h1>404 - Not Found</h1>
 </body>
</html>


// Server log
Delaying the response by 10s as requested. The mock response will be avaiable in 2s. Returning `notFoundResponse` for now
```

Once the delay is complete, querying `spotitn` paths return expected results:
```
$ curl localhost:1338/latest/meta-data/spot/instance-action

{
	"instance-action": "stop",
	"time": "2020-04-24T17:19:32Z"
}

```

### Scheduled Events
Similar to spotitn, the `scheduledevents` command, view the local flags using `scheduledevents --help`:

```
$ amazon-ec2-metadata-mock scheduledevents --help
Mock EC2 Scheduled Events

Usage:
  amazon-ec2-metadata-mock scheduledevents [--code CODE] [--state STATE] [--not-after] [--not-before-deadline] [flags]

Aliases:
  scheduledevents, se, scheduled-events, scheduledEvents

Examples:
  amazon-ec2-metadata-mock scheduledevents -h 	scheduledevents help
  amazon-ec2-metadata-mock scheduledevents -o instance-stop --state active -d		mocks an active and upcoming scheduled event for instance stop with a deadline for the event start time

Flags:
  -o, --code string                  event code in the scheduled event (default: system-reboot)
                                     event-code can be one of the following: instance-reboot,system-reboot,system-maintenance,instance-retirement,instance-stop
  -h, --help                         help for scheduledevents
  -a, --not-after string             the latest end time for the scheduled event in RFC3339 format E.g. 2020-01-07T01:03:47Z default: application start time + 7 days in UTC))
  -b, --not-before string            the earliest start time for the scheduled event in RFC3339 format E.g. 2020-01-07T01:03:47Z (default: application start time in UTC)
  -l, --not-before-deadline string   the deadline for starting the event in RFC3339 format E.g. 2020-01-07T01:03:47Z (default: application start time + 9 days in UTC)
  -t, --state string                 state of the scheduled event (default: active)
                                     state can be one of the following: active,completed,canceled

(Truncated Global Flags for readability)
```

1.) **Starting AEMM with `scheduledevents`**: `scheduledevents` route available immediately and `spotitn` routes will no longer be available per the *Note* above:

```
$ amazon-ec2-metadata-mock scheduledevents
Initiating amazon-ec2-metadata-mock for EC2 Scheduled Events on port 1338
Serving the following routes: ... (truncated for readability)

```

Send the request:
```
$ curl localhost:1338/latest/meta-data/events/maintenance/scheduled
{
	"Code": "system-reboot",
	"Description": "The instance is scheduled for system-reboot",
	"State": "active",
	"EventID": "instance-event-1234567890abcdef0",
	"NotBefore": "24 Apr 2020 12:32:00 GMT",
	"NotAfter": "01 May 2020 12:32:00 GMT",
	"NotBeforeDeadline": "03 May 2020 12:32:00 GMT"
}
```

2.) **Starting AEMM with `scheduledevents` invalid flag overrides**: as noted above, all commands have validation logic for overrides via CLI flags. If the user attempts to pass an invalid override value, then AEMM will panic, kill the server, and return an error message with what went wrong:

```
$ amazon-ec2-metadata-mock scheduledevents --code FOO

panic: Fatal error while executing the root command: Invalid CLI input "FOO" for flag code. 
Allowed value(s): instance-reboot,system-reboot,system-maintenance,instance-retirement,instance-stop.
```

## Instance Metadata Service Versions
AEMM supports [both versions](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-instance-metadata-service.html) of Instance Metadata service. By default, AEMM starts with supporting v1 and v2; however, it is also possible to enable **IMDSv2 only** via overrides.

1.) **Starting AEMM with IMDSv1 & IMDSv2:** default behavior where providing session token is completely optional and normal request/response method is supported.

```
$ amazon-ec2-metadata-mock
```

Send a v1 request:
```
$ curl localhost:1338/latest/meta-data/mac
0e:49:61:0f:c3:11
```

Send a v2 request:
```
TOKEN=`curl -X PUT "localhost:1338/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600"` \
&& curl -H "X-aws-ec2-metadata-token: $TOKEN" localhost:1338/latest/meta-data/mac
0e:49:61:0f:c3:11
```

2.) **Starting AEMM with IMDSv2 only:** session tokens are required for all requests; v1 requests will return **401 - Unauthorized:**

```
$ amazon-ec2-metadata-mock --imdsv2
```
Send a v1 request:
```
$ curl localhost:1338/latest/meta-data/mac

<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
   <head>
      <title>401 - Unauthorized</title>
   </head>
   <body>
      <h1>401 - Unauthorized</h1>
   </body>
</html>
```

Sending a v2 request will yield the same results as 1.) above.

Requesting a token outside the TTL bounds (between 1-2600 seconds) will return **400 - Bad Request:**
```
$ curl -X PUT "localhost:1338/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 0"
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
   <head>
      <title>400 - Bad Request</title>
   </head>
   <body>
      <h1>400 - Bad Request</h1>
   </body>
</html>
```

Providing an expired token is synonymous to no token at all resulting in **401 - Unauthorized**.

## Static Metadata
Static metadata is classified as instance-specific metadata that is **always** available regardless of which command is used to start the tool. Some additional properties of static metadata:
* delays do **NOT** affect static metadata availability
* values are overridden via config file and/or env variables **only**
  * values **cannot** be overridden using flags 

Examples of static metadata include *all* [metadata categories](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-categories.html) from the non-dynamic category table (ami-id, instance-id, mac) **except for events and spot categories (classified as commands in AEMM).**

*Note that 'static' naming is used within the context of this tool ONLY*

### Static Metadata Overrides & Path Substitutions
Some metadata categories use data unique to the instance as part of the query.
* **Example:**  network/interfaces/macs/*MAC*/interface-id where *MAC* is a placeholder for the instance's mac address
  * [AWS Documentation](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-categories.html)

Applying overrides to these *placeholder values* will automatically update paths containing these values **unless** the path itself is explicitly overridden.

1.) **Starting AEMM with JSON config overrides:** The static metadata will use the overridden values AND paths using overridden instance data will be updated as well.

* config-overrides.json:

```
{
   "metadata":{
      "paths":{
         "mac-device-number":"/latest/meta-data/network/interfaces/macs/BAR/device-number"
      },
      "values":{
         "mac":"FOO"
      }
   }
}

```

* start the server with overrides:

```
$ amazon-ec2-metadata-mock -c config-overrides.json
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

# Troubleshooting

## Warnings and Expected Outcome
|Warning Displayed|Cause|Effect|
|---|---|---|
|Warning: Config File _file name_ Not Found in _locations_ |input configuration file not found|other input sources are used (See above for details)|
|Warning: Failed to save the final configuration to local file - Failed to create directory for final configuration at _path/to/dir_: _error string_|Failure to create the hidden directory `.amazon-ec2-metadata-mock` to store the final configuration file| configuration used by the tool is NOT saved to a local file. The tool continues with its primary job of mocking metadata paths |
|Warning: Failed to save the final configuration to local file - The destination '_path/to/dir_' for saving the configuration already exists, but is not a directory |Failure to create the hidden directory `.amazon-ec2-metadata-mock` to store the final configuration file, because a resource by that name already exists| configuration used by the tool is NOT saved to a local file. The tool continues with its primary job of mocking metadata paths |
|Warning: Failed to save the final configuration to local file _path/to/local/file_: _error string_  |Failure to save final configuration to a file |configuration used by the tool is NOT saved to a local file. The tool continues with its primary job of mocking metadata paths |
|Warning: Failed to find home directory due to error: _error string_|Failure to get home directory| working directory is used instead|

# Building
For build instructions, please consult [BUILD.md](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/BUILD.md)

# Communication
If you've run into a bug or have a new feature request, please open an [issue](https://github.com/aws/amazon-ec2-metadata-mock/issues/new).

# Contributing
Contributions are welcome! Please read our [guidelines](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/CONTRIBUTING.md) and our [Code of Conduct](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/CODE_OF_CONDUCT.md)

# License
This project is licensed under the Apache-2.0 License.
