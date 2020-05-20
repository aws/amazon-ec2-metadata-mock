# AEMM Configuration
This page serves as documentation for AEMM's configuration details.

## Configuring AEMM
The tool can be configured in various ways:

1. CLI flags

    Use help commands to learn more

2. Env variables
    ```
    $ export AEMM_MOCK_DELAY_SEC=12
    $ export AEMM_SPOT_ITN_INSTANCE_ACTION=stop
    $ env | grep AEMM     // To list the tool's env variables
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
$ ec2-metadata-mock -c path/to/config-overrides.json -s
Successfully saved final configuration to local file  /path/to/home/.ec2-metadata-mock/.aemm-config-used.json


$ cat $HOME/.ec2-metadata-mock/.aemm-config-used.json
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

## Precedence (Highest to Lowest)
1. overrides
2. flag
3. env
4. config
5. key/value store
6. default

For example, if values from the following sources were loaded:
```
Defaults in code:
{
    "config-file": "$HOME/aemm-config.json", # by default, AEMM looks here for a config file
    "server": {
         "port": "1338"
    },
    "mock-delay-sec": 0,
    "save-config-to-file": false,
    "spot-itn": {
       "instance-action": "terminate"
    }
}

Env variables:
export AEMM_MOCK_DELAY_SEC=12
export AEMM_SPOT_ITN_INSTANCE_ACTION=hibernate
export AEMM_CONFIG_FILE=/path/to/my-custom-aemm-config.json

Config File (at /path/to/my-custom-aemm-config.json):
{
    "imdsv2": true,
    "server": {
         "port": "1550"
    },
    "spot-itn": {
       "instance-action": "stop"
    }
}

CLI Flags:
 {
   "mock-delay-sec": 8
 }
```

The resulting config will have the following values (non-overriden values are truncated for readability):
```
 {
    "mock-delay-sec": 8,                                        # from CLI flag
    "config-file": "/path/to/my-custom-aemm-config.json",       # from env
    "spot-itn": {
        "instance-action": "hibernate"                          # from env
   }
    "imdsv2": true,                                             # from custom config file at /path/to/my-custom-aemm-config.json
     "server": {
        "port": "1550"                                          # from custom config file at /path/to/my-custom-aemm-config.json
    },
    "save-config-to-file": false,                               # from defaults in code
 }
```

AEMM is built using Viper which is where the precedence is sourced from. More details can be found in their documentation: 

<a href="https://github.com/spf13/viper/blob/master/README.md">
    <img src="https://img.shields.io/badge/viper-documentation-green" alt="">
</a>
