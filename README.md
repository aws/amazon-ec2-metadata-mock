# Amazon EC2 Metadata Mock

**Amazon EC2 Metadata Mock (AEMM)** is a tool to simulate [Amazon EC2 instance metadata service](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/ec2-instance-metadata.html) for local testing.

<br/>
<p>
   <a href="https://hub.docker.com/r/amazon/amazon-ec2-metadata-mock">
   <img src="https://img.shields.io/github/v/release/aws/amazon-ec2-metadata-mock?color=yellowgreen&label=latest%20release&sort=semver" alt="latest release">
   </a>
   <a href="https://golang.org/doc/go1.14">
   <img src="https://img.shields.io/github/go-mod/go-version/aws/amazon-ec2-metadata-mock?color=blueviolet" alt="go-version">
   </a>
   <a href="https://opensource.org/licenses/Apache-2.0">
   <img src="https://img.shields.io/badge/License-Apache%202.0-ff69b4.svg?color=orange" alt="license">
   </a>
   <a href="https://travis-ci.org/aws/amazon-ec2-metadata-mock">
   <img src="https://travis-ci.org/aws/amazon-ec2-metadata-mock.svg?branch=master" alt="build-status">
   </a>
   <a href="https://hub.docker.com/r/amazon/amazon-ec2-metadata-mock">
   <img src="https://img.shields.io/docker/pulls/amazon/amazon-ec2-metadata-mock" alt="docker-pulls">
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
      * [Overrides](#overrides)
   * [Usage](#usage)
      * [Spot Interruption](#spot-interruption)
      * [Scheduled Events](#events)
      * [Instance Metadata Service Versions](#instance-metadata-service-versions)
      * [Static Metadata](#static-metadata)
   * [Troubleshooting](#troubleshooting)
      * [Warnings and Expected Outcome](#warnings-and-expected-outcome)
   * [Building](#building)
   * [Communication](#communication)
   * [Contributing](#contributing)
   * [License](#license)

# Summary
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

### MacOS/Linux
```
curl -Lo ec2-metadata-mock https://github.com/aws/amazon-ec2-metadata-mock/releases/download/v1.1.1/ec2-metadata-mock-`uname | tr '[:upper:]' '[:lower:]'`-amd64
chmod +x ec2-metadata-mock
```

### ARM Linux
```
curl -Lo ec2-metadata-mock https://github.com/aws/amazon-ec2-metadata-mock/releases/download/v1.1.1/ec2-metadata-mock-linux-arm
```

```
curl -Lo ec2-metadata-mock https://github.com/aws/amazon-ec2-metadata-mock/releases/download/v1.1.1/ec2-metadata-mock-linux-arm64
```

### Windows
```
curl -Lo ec2-metadata-mock https://github.com/aws/amazon-ec2-metadata-mock/releases/download/v1.1.1/ec2-metadata-mock-windows-amd64.exe
```

### Docker
```
docker pull amazon/amazon-ec2-metadata-mock:v1.1.1
docker run -it --rm -p 1338:1338 amazon/amazon-ec2-metadata-mock:v1.1.1
```

### On Kubernetes
#### Supported versions
* Kubernetes >= 1.14

#### Helm
[See Helm README here.](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/helm/amazon-ec2-metadata-mock/README.md)

#### kubectl
kubectl apply -f https://github.com/aws/amazon-ec2-metadata-mock/releases/download/v1.1.1/all-resources.yaml

## Starting AEMM
Use `ec2-metadata-mock --help` to view examples and explanations of supported flags and commands:

```
$ ec2-metadata-mock --help

ec2-metadata-mock is a tool to mock Amazon EC2 instance metadata.

Usage:
  ec2-metadata-mock <command> [arguments] [flags]
  ec2-metadata-mock [command]

Examples:
  ec2-metadata-mock --mock-delay-sec 10	mocks all metadata paths
  ec2-metadata-mock spot --action terminate	mocks spot ITN only

Available Commands:
  events          Mock EC2 maintenance events
  help            Help about any command
  spot            Mock EC2 Spot interruption notice

Flags:
  -c, --config-file string    config file for cli input parameters in json format (default: $HOME/aemm-config.json)
  -h, --help                  help for amazon-ec2-metadata-mock
  -n, --hostname string       the HTTP hostname for the mock url (default: localhost)
  -I, --imdsv2                whether to enable IMDSv2 only requiring a session token when submitting requests (default: false)
  -d, --mock-delay-sec int    mock delay in seconds, relative to the application start time (default: 0 seconds)
  -p, --port string           the HTTP port where the mock runs (default: 1338)
  -s, --save-config-to-file   whether to save processed config from all input sources in .ec2-metadata-mock/.aemm-config-used.json in $HOME or working dir, if homedir is not found (default: false)

Use "ec2-metadata-mock [command] --help" for more information about a command.
```

Starting AEMM with default configurations using `ec2-metadata-mock` will start the server on the default host and port:

```
$ ec2-metadata-mock

Initiating ec2-metadata-mock for all mocks on port 1338
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
  * [events](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/pkg/cmd/events/events.go#L72) 

## Overrides
AEMM supports configuration from various sources including: cli flags, env variables, and config files. Details regarding
configuration steps, behavior, and precedence are outlined [here](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/docs/configuration.md).

# Usage
AEMM is primarily used as a developer tool to help test behavior related to Metadata Service. Popular use cases include: emulating spot instance interrupts after a designated delay, mocking scheduled maintenance events, IMDSv2 migrations, 
and requesting static metadata. This section outlines the common use cases of AEMM; advanced usage and behavior are documented [here](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/docs/usage.md).

## Spot Interruption
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
  -c, --config-file string    config file for cli input parameters in json format (default: $HOME/aemm-config.json)
  -n, --hostname string       the HTTP hostname for the mock url (default: localhost)
  -I, --imdsv2                whether to enable IMDSv2 only, requiring a session token when submitting requests (default: false, meaning both IMDS v1 and v2 are enabled)
  -d, --mock-delay-sec int    mock delay in seconds, relative to the application start time (default: 0 seconds)
  -p, --port string           the HTTP port where the mock runs (default: 1338)
  -s, --save-config-to-file   whether to save processed config from all input sources in .amazon-ec2-metadata-mock/.aemm-config-used.json in $HOME or working dir, if homedir is not found (default: false)
```

1.) **Starting AEMM with `spot`**:  `spot` routes available immediately:
```
$ ec2-metadata-mock spot
Initiating ec2-metadata-mock for EC2 Spot interruption notice on port 1338
Serving the following routes: ... (truncated for readability)
```
Send the request:
```
$ curl localhost:1338/latest/meta-data/spot/instance-action
{
	"action": "terminate",
	"time": "2020-04-24T17:11:44Z"
}
```


2.) **Starting AEMM with `spot` after Delay**: Users can apply a *delay* duration in seconds for when the `spot` metadata will become available:

```
$ ec2-metadata-mock spot -d 10
Initiating ec2-metadata-mock for EC2 Spot interruption notice on port 1338

Flags:
mock-delay-sec: 10

Serving the following routes: ... (truncated for readability)
```

Sending a request to `spot` paths before the delay has passed will return **404 - Not Found:**
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

Once the delay is complete, querying `spot` paths return expected results:
```
$ curl localhost:1338/latest/meta-data/spot/instance-action

{
	"action": "terminate",
	"time": "2020-04-24T17:19:32Z"
}

```

## Events
Similar to spot, the `events` command, view the local flags using `events --help`:

```
$ ec2-metadata-mock events --help
Mock EC2 Scheduled Events

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

1.) **Starting AEMM with `events`**: `events` route available immediately and `spot` routes will no longer be available due to the implementation of Commands [detailed here](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/docs/usage.md):

```
$ ec2-metadata-mock events --code instance-reboot -a 2020-01-07T01:03:47Z  -b 2020-01-01T01:03:47Z -l 2020-01-10T01:03:47Z --state completed
Initiating ec2-metadata-mock for EC2 Events on port 1338
Serving the following routes: ... (truncated for readability)

```

Send the request:
```
$ curl localhost:1338/latest/meta-data/events/maintenance/scheduled
{
	"Code": "instance-reboot",
	"Description": "The instance is scheduled for instance-reboot",
	"State": "completed",
	"EventId": "instance-event-1234567890abcdef0",
	"NotBefore": "1 Jan 2020 01:03:47 GMT",
	"NotAfter": "7 Jan 2020 01:03:47 GMT",
	"NotBeforeDeadline": "10 Jan 2020 01:03:47 GMT"
}
```

## Instance Metadata Service Versions
AEMM supports [both versions](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/configuring-instance-metadata-service.html) of Instance Metadata service. By default, AEMM starts with supporting v1 and v2; however, it is possible to enable **IMDSv2 only** via overrides.

1.) **Starting AEMM with IMDSv2 only:** session tokens are required for all requests; v1 requests will return **401 - Unauthorized:**

```
$ ec2-metadata-mock --imdsv2
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

Send a v2 request:
```
TOKEN=`curl -X PUT "localhost:1338/latest/api/token" -H "X-aws-ec2-metadata-token-ttl-seconds: 21600"` \
&& curl -H "X-aws-ec2-metadata-token: $TOKEN" localhost:1338/latest/meta-data/mac
0e:49:61:0f:c3:11
```

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
Static metadata is classified as instance-specific metadata that is **always** available regardless of which command is used to start the tool. 

Examples of static metadata include *all* [metadata categories](https://docs.aws.amazon.com/AWSEC2/latest/UserGuide/instancedata-data-categories.html) from the non-dynamic category table (ami-id, instance-id, mac) **except for events and spot categories (classified as commands in AEMM).**

*Note that 'static' naming is used within the context of this tool ONLY*

1.) **Requesting static metadata `instance-id`**:
```
$ ec2-metadata-mock
```

Send the request:
```
$ curl localhost:1338/latest/meta-data/instance-id
i-1234567890abcdef0
```

Details on overriding static metadata values and behavior can be found [here](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/docs/usage.md#static-metadata)

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
