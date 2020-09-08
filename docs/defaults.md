# Default configuration

## AEMM configuration
Key | Value
--- | --- 
hostname | 0.0.0.0
port | 1338
input config file | $HOME/aemm-config.json
output / used config file | {$HOME or working dir}/.ec2-metadata-mock/.aemm-config-used.json
mock delay seconds | 0
IMDS v2 only (http requests require a session token) | false
spot instance action | terminate
spot termination time | request time + 2 minutes in UTC
scheduled events code | systemReboot
scheduled events state | active
scheduled events not-before | application start time in UTC
scheduled events not-after | application start time + 7 days in UTC
scheduled events not-before-deadline | application start time + 9 days in UTC

## Default metadata
Key | Value
--- | --- 
ami-id | ami-0a887e401f7654935
ami-launch-index | 0
ami-manifest-path | (unknown)
block-device-mapping-ami | /dev/xvda
block-device-mapping-ebs | sdb
block-device-mapping-ephemeral | sdb
block-device-mapping-root | /dev/xvda
block-device-mapping-swap | sdcs
elastic-inference-accelerator-id | eia-bfa21c7904f64a82a21b9f4540169ce1
elastic-inference-accelerator-type | eia1.medium
elastic-inference-associations | eia-bfa21c7904f64a82a21b9f4540169ce1
event-id | instance-event-1234567890abcdef0
hostname | ip-172-16-34-43.ec2.internal
iam-info code | Success
iam-info instanceprofilearn | arn:aws:iam::896453262835:instance-profile/baskinc-role
iam-info instanceprofileid | AIPA5BOGHHXZELSK34VU4
iam-info lastupdated | 2020-04-02T18:50:40Z
iam-security-credentials accesskeyid | 12345678901
iam-security-credentials code | Success
iam-security-credentials expiration | 2020-04-02T00:49:51Z
iam-security-credentials lastupdated | 2020-04-02T18:50:40Z
iam-security-credentials secretaccesskey | v/12345678901
iam-security-credentials token | TEST92test48TEST+y6RpoTEST92test48TEST/8oWVAiBqTEsT5Ky7ty2tEStxC1T==
iam-security-credentials-role | baskinc-role
instance-action | none
instance-id | i-1234567890abcdef0
instance-type | m4.xlarge
local-hostname | ip-172-16-34-43.ec2.internal
local-ipv4 | 172.16.34.43
mac | 0e:49:61:0f:c3:11
mac-device-number | 0
mac-ipv4-associations | 192.0.2.54
mac-ipv6-associations | 2001:db8:8:4::2
mac-local-hostname | ip-172-16-34-43.ec2.internal
mac-local-ipv4s | 172.16.34.43
mac-mac | 0e:49:61:0f:c3:11
mac-network-interface-id | eni-0f95d3625f5c521cc
mac-owner-id | 515336597381
mac-public-hostname | ec2-192-0-2-54.compute-1.amazonaws.com
mac-public-ipv4s | 192.0.2.54
mac-security-group-ids | sg-0b07f8f6cb485d4df
mac-security-groups | ura-launch-wizard-harry-1
mac-subnet-id | subnet-0ac62554
mac-subnet-ipv4-cidr-block | 192.0.2.0/24
mac-subnet-ipv6-cidr-blocks | 2001:db8::/32
mac-vpc-id | vpc-d295a6a7
mac-vpc-ipv4-cidr-block | 192.0.2.0/24
mac-vpc-ipv4-cidr-blocks | 192.0.2.0/24
mac-vpc-ipv6-cidr-blocks | 2001:db8::/32
placement-availability-zone | us-east-1a
placement-availability-zone-id | use1-az4
placement-group-name | a-placement-group
placement-host-id | h-0da999999f9999fb9
placement-partition-number | 1
placement-region | us-east-1
product-codes | 3iplms73etrdhxdepv72l6ywj
public-hostname | ec2-192-0-2-54.compute-1.amazonaws.com
public-ipv4 | 192.0.2.54
public-key | ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC/JxGByvHDHgQAU+0nRFWdvMPi22OgNUn9ansrI8QN1ZJGxD1ML8DRnJ3Q3zFKqqjGucfNWW0xpVib+ttkIBp8G9P/EOcX9C3FF63O3SnnIUHJsp5faRAZsTJPx0G5HUbvhBvnAcCtSqQgmr02c1l582vAWx48pOmeXXMkl9qe9V/s7K3utmeZkRLo9DqnbsDlg5GWxLC/rWKYaZR66CnMEyZ7yBy3v3abKaGGRovLkHNAgWjSSgmUTI1nT5/S2OLxxuDnsC7+BiABLPaqlIE70SzcWZ0swx68Bo2AY9T9ymGqeAM/1T4yRtg0sPB98TpT7WrY5A3iia2UVtLO/xcTt test
reservation-id | r-046cb3eca3e201d2f
security-groups | ura-launch-wizard-harry-1
services-domain | amazonaws.com
services-partition | aws
