# Amazon EC2 Metadata Mock

Amazon EC2 Metadata Mock(AEMM) Helm chart for Kubernetes. For more information on this project see the project repo at https://github.com/aws/amazon-ec2-metadata-mock.

## Prerequisites

* Kubernetes >= 1.11

## Installing the Chart

The helm chart can be installed from several sources. To install the chart with the release name amazon-ec2-metadata-mock and default configuration, pick a source below:

1. Local chart archive: 
Download the chart archive from the latest release and run 
```sh
helm install amazon-ec2-metadata-mock amazon-ec2-metadata-mock-0.1.0.tgz \
  --namespace default
```

2. Unpacked local chart directory: 
Download the source code or unpack the archive from latest release and run
```sh
helm install amazon-ec2-metadata-mock ./helm/amazon-ec2-metadata-mock \
  --namespace default
```
----
To upgrade an already installed chart named amazon-ec2-metadata-mock:
```sh
helm upgrade amazon-ec2-metadata-mock ./helm/amazon-ec2-metadata-mock \
  --namespace default
```

### Installing the Chart with overridden values for AEMM configuration:

AEMM has an [extensive list of parameters](https://github.com/aws/amazon-ec2-metadata-mock#defaults) that can overridden. For simplicity, a selective list of parameters are configurable using Helm custom `values.yaml` and `--set argument`. To override parameters not listed in `values.yaml` use Kubernetes ConfigMap.    

The [configuration](#configuration) section details the selective list of parameters. Alternatively, to retrieve the same information via helm, run:
```sh
helm show values ./helm/amazon-ec2-metadata-mock
```

* Passing a custom values.yaml to helm
```sh
helm install amazon-ec2-metadata-mock ./helm/amazon-ec2-metadata-mock \
  --namespace default -f path/to/myvalues.yaml 
```

* Passing custom values to helm via CLI 
```sh
helm install amazon-ec2-metadata-mock ./helm/amazon-ec2-metadata-mock \
  --namespace default --set aemm.server.port=1660,aemm.mockDelaySec=120  
```

* Passing a config file to AEMM

 1. Create a Kubernetes ConfigMap from a custom AEMM configuration file:
See [Readme](https://github.com/aws/amazon-ec2-metadata-mock#configuration) to learn more about AEMM configuration. [Here](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/test/e2e/testdata/output/aemm-config-used.json) is a reference config file to create your own `aemm-config.json`
```sh
kubectl create configmap aemm-config-map --from-file path/to/aemm-config.json
```

 2. Create `myvalues.yaml` with overridden value for configMap:
```yaml
configMap: "aemm-config-map"
```

 3. Install AEMM with override:
```sh
helm install amazon-ec2-metadata-mock ./helm/amazon-ec2-metadata-mock \
  --namespace default -f path/to/myvalues.yaml 
```



## Making a HTTP request to the AEMM server running on a pod

1. Get the AEMM pod name:
```sh
kubectl get pods --namespace default
```

2. Set up port-forwarding for the port on which AEMM is running:
```sh
kubectl port-forward pod/<AEMM-pod-name> 1660
```

3. Make the HTTP request
```sh
curl http://localhost:1660/latest/meta-data/spot/instance-action
{
	"instance-action": "terminate",
	"time": "2020-05-04T18:11:37Z"
}
```

## Uninstalling the Chart

To uninstall/delete the `amazon-ec2-metadata-mock` release:
```sh
helm uninstall amazon-ec2-metadata-mock
```
The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following tables lists the configurable parameters of the chart and their default values.

Parameter | Description | Default
--- | --- | --- 
`image.repository` | image repository | `amazon/amazon-ec2-metadata-mock` 
`image.tag` | image tag | `<VERSION>` 
`image.pullPolicy` | image pull policy | `IfNotPresent`
`nameOverride` | override for the name of the Helm Chart (default, if not overridden: `amazon-ec2-metadata-mock`) | `""`
`fullnameOverride` | override for the name of the application (default, if not overridden: `amazon-ec2-metadata-mock`) | `""`
`nodeSelector` | tells the DaemonSet where to place the amazon-ec2-metadata-mock pods. | `{}`, meaning every node will receive a pod
`podAnnotations` | annotations to add to each pod | `{}`
`updateStrategy` | the update strategy for a DaemonSet | `RollingUpdate`
`rbac.pspEnabled` | if `true`, create and use a restricted pod security policy | `false`
`serviceAccount.create` | if `true`, create a new service account | `true`
`serviceAccount.name` | service account to be used | `amazon-ec2-metadata-mock-service-account`
`serviceAccount.annotations` | specifies the annotations for service account | `{}`
`securityContext.runAsUserID` | user ID to run the container | `1000`
`securityContext.runAsGroupID` | group ID to run the container | `1000` 
`namespace` | Kubernetes namespace to use for AEMM pods | `default`
`configMap` | name of the Kubernetes ConfigMap to use to pass a config file for AEMM overrides | `""`
`configMapFileName` | name of the file used to create the Kubernetes ConfigMap | `aemm-config.json`

NOTE: A selective list of AEMM parameters are configurable via Helm CLI and values.yaml file. 
Use the [Kubernetes ConfigMap option](#installing-the-chart-with-overridden-values-for-aemm-configuration) to configure [other AEMM parameters](https://github.com/aws/amazon-ec2-metadata-mock/blob/master/test/e2e/testdata/output/aemm-config-used.json). 

Parameter | Description | Default in values.yaml | Default AEMM configuration
--- | --- | --- | ---
`aemm.server.port` | port to run AEMM on | `""` | `1338`
`aemm.server.hostname` | hostname to run AEMM on | `""` | `localhost`
`aemm.mockDelaySec` | mock delay in seconds, relative to the start time of AEMM | `0` | `0`
`aemm.imdsv2` | if true, IMDSv2 only works | `false` | `false`, meaning both IMDSv1/v2 work 
`aemm.spotItn.instanceAction` | instance action in the spot interruption notice | `""` | `terminate`
`aemm.spotItn.terminationTime` | termination time in the spot interruption notice | `""` | HTTP request time + 2 minutes
`aemm.scheduledEvents.code` | event code in the scheduled event | `""` | `system-reboot`
`aemm.scheduledEvents.notAfter` | the latest end time for the scheduled event | `""` | Start time of AEMM  + 7 days
`aemm.scheduledEvents.notBefore` | the earliest start time for the scheduled event | `""` | Start time of AEMM
`aemm.scheduledEvents.notBeforeDeadline` | the deadline for starting the event | `""` | Start time of AEMM  + 9 days
`aemm.scheduledEvents.state` | state of the scheduled event | `""` | `active`