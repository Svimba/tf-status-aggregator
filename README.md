# Tungsten Fabric status agregator

> NOTE: This project is *pre-release* stage

Agregator TF services takes all "tf-status" nodes from K8s API which are labeled by following label.
```
tungstenfabric: status
```
URL `/json` of each node are requested in parallel and proceed to show requested data.

## Usage

| Link    | Descriprion |
| ------- | ----------- |
| /pod-list     | Returns information about all detected \"tf-status\" pods. TFStatus pod have to have following label: `"tungstenfabric": "status"`. |
| /status/json  | Returns agregated json from all "tf-status" pods.                                                                                   |
| /status       | TEST                                                                                                                                |
| /status/node  | Returns formated output from all detected "tf-status" pods at standart format for each node                                         |
| /status/group | Returns formated output from all detected "tf-status" pod tagregated by service for all nodes which handles the service             |


## ENV

| ENV name    | Default |
| ----------- | ------- |
| SERVER_PORT | 80      |
| NODE_PORT   | 80      |

### Build container

Set `TAG` at make file manually (will be fixed), then run following command.

```
$ make build
```

### Build container

```
$ make push
```