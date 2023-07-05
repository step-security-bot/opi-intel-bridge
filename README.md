# OPI gRPC to Intel SDK bridge

[![Linters](https://github.com/opiproject/opi-intel-bridge/actions/workflows/linters.yml/badge.svg)](https://github.com/opiproject/opi-intel-bridge/actions/workflows/linters.yml)
[![tests](https://github.com/opiproject/opi-intel-bridge/actions/workflows/go.yml/badge.svg)](https://github.com/opiproject/opi-intel-bridge/actions/workflows/go.yml)
[![Docker](https://github.com/opiproject/opi-intel-bridge/actions/workflows/docker-publish.yml/badge.svg)](https://github.com/opiproject/opi-intel-bridge/actions/workflows/docker-publish.yml)
[![License](https://img.shields.io/github/license/opiproject/opi-intel-bridge?style=flat-square&color=blue&label=License)](https://github.com/opiproject/opi-intel-bridge/blob/master/LICENSE)
[![codecov](https://codecov.io/gh/opiproject/opi-intel-bridge/branch/main/graph/badge.svg)](https://codecov.io/gh/opiproject/opi-intel-bridge)
[![Go Report Card](https://goreportcard.com/badge/github.com/opiproject/opi-intel-bridge)](https://goreportcard.com/report/github.com/opiproject/opi-intel-bridge)
[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/opiproject/opi-intel-bridge)
[![Pulls](https://img.shields.io/docker/pulls/opiproject/opi-intel-bridge.svg?logo=docker&style=flat&label=Pulls)](https://hub.docker.com/r/opiproject/opi-intel-bridge)
[![Last Release](https://img.shields.io/github/v/release/opiproject/opi-intel-bridge?label=Latest&style=flat-square&logo=go)](https://github.com/opiproject/opi-intel-bridge/releases)
[![GitHub stars](https://img.shields.io/github/stars/opiproject/opi-intel-bridge.svg?style=flat-square&label=github%20stars)](https://github.com/opiproject/opi-intel-bridge)
[![GitHub Contributors](https://img.shields.io/github/contributors/opiproject/opi-intel-bridge.svg?style=flat-square)](https://github.com/opiproject/opi-intel-bridge/graphs/contributors)

This is a Intel app (bridge) to OPI APIs for storage, inventory, ipsec and networking (future).

## Getting started

build like this:

```bash
go build -v -o /opi-intel-bridge ./cmd/...
```

import like this:

```go
import "github.com/opiproject/opi-intel-bridge/pkg/frontend"
```

## Using docker

on DPU/IPU (i.e. with IP=10.10.10.1) run

```bash
$ docker run --rm -it -v /var/tmp/:/var/tmp/ -p 50051:50051 ghcr.io/opiproject/opi-intel-bridge:main
2022/11/29 00:03:55 plugin serevr is &{{}}
2022/11/29 00:03:55 server listening at [::]:50051
```

on X86 management VM run

reflection

```bash
$ docker run --network=host --rm -it namely/grpc-cli ls --json_input --json_output localhost:50051
grpc.reflection.v1alpha.ServerReflection
opi_api.inventory.v1.InventorySvc
opi_api.security.v1.IPsec
opi_api.storage.v1.AioControllerService
opi_api.storage.v1.FrontendNvmeService
opi_api.storage.v1.FrontendVirtioBlkService
opi_api.storage.v1.FrontendVirtioScsiService
opi_api.storage.v1.MiddleendEncryptionService
opi_api.storage.v1.MiddleendQosVolumeService
opi_api.storage.v1.NVMfRemoteControllerService
opi_api.storage.v1.NullDebugService
```

full test suite

```bash
docker run --rm -it --network=host docker.io/opiproject/godpu:main get --addr="10.10.10.10:50051"
docker run --rm -it --network=host docker.io/opiproject/godpu:main storagetest --addr="10.10.10.10:50051"
docker run --rm -it --network=host docker.io/opiproject/godpu:main test --addr=10.10.10.10:50151 --pingaddr=8.8.8.1"
```

or manually

```bash
docker run --network=host --rm -it namely/grpc-cli ls   --json_input --json_output 10.10.10.10:50051 -l
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 CreateNvmeSubsystem "{nvme_subsystem : {spec : {nqn: 'nqn.2022-09.io.spdk:opitest2', serial_number: 'myserial2', model_number: 'mymodel2', max_namespaces: 11} }, nvme_subsystem_id : 'subsystem2' }"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 ListNvmeSubsystems "{parent : 'todo'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 GetNvmeSubsystem "{name : '//storage.opiproject.org/volumes/subsystem2'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 CreateNvmeController "{nvme_controller : {spec : {nvme_controller_id: 2, subsystem_id : { value : '//storage.opiproject.org/volumes/subsystem2' }, pcie_id : {physical_function : 0}, max_nsq:5, max_ncq:5 } }, nvme_controller_id : 'controller1'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 ListNvmeControllers "{parent : '//storage.opiproject.org/volumes/subsystem2'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 GetNvmeController "{name : '//storage.opiproject.org/volumes/controller1'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 CreateNvmeNamespace "{nvme_namespace : {spec : {subsystem_id : { value : '//storage.opiproject.org/volumes/subsystem2' }, volume_id : { value : 'Malloc0' }, 'host_nsid' : '10', uuid:{value : '1b4e28ba-2fa1-11d2-883f-b9a761bde3fb'}, nguid: '1b4e28ba-2fa1-11d2-883f-b9a761bde3fb', eui64: 1967554867335598546 } }, nvme_namespace_id: 'namespace1'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 ListNvmeNamespaces "{parent : '//storage.opiproject.org/volumes/subsystem2'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 GetNvmeNamespace "{name : '//storage.opiproject.org/volumes/namespace1'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 NvmeNamespaceStats "{namespace_id : {value : '//storage.opiproject.org/volumes/namespace1'} }"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 CreateNVMfRemoteController "{nv_mf_remote_controller : {multipath: 'NVME_MULTIPATH_MULTIPATH'}, nv_mf_remote_controller_id: 'nvmetcp12'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 ListNVMfRemoteControllers "{}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 GetNVMfRemoteController "{name: '//storage.opiproject.org/volumes/nvmetcp12'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 CreateNVMfPath "{nv_mf_path : {controller_id: {value: '//storage.opiproject.org/volumes/nvmetcp12'}, traddr:'11.11.11.2', subnqn:'nqn.2016-06.com.opi.spdk.target0', trsvcid:'4444', trtype:'NVME_TRANSPORT_TCP', adrfam:'NVMF_ADRFAM_IPV4', hostnqn:'nqn.2014-08.org.nvmexpress:uuid:feb98abe-d51f-40c8-b348-2753f3571d3c'}, nv_mf_path_id: 'nvmetcp12path0'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 ListNVMfPaths "{parent : 'todo'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 GetNVMfPath "{name: '//storage.opiproject.org/volumes/nvmetcp12path0'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 DeleteNVMfPath "{name: '//storage.opiproject.org/volumes/nvmetcp12path0'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 DeleteNVMfRemoteController "{name: '//storage.opiproject.org/volumes/nvmetcp12'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 DeleteNvmeNamespace "{name : '//storage.opiproject.org/volumes/namespace1'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 DeleteNvmeController "{name : '//storage.opiproject.org/volumes/controller1'}"
docker run --network=host --rm -it namely/grpc-cli call --json_input --json_output 10.10.10.10:50051 DeleteNvmeSubsystem "{name : '//storage.opiproject.org/volumes/subsystem2'}"
```

## Recipes

### Prerequisites

#### On Host

TBD

#### On IPU

Allocate all required resources...

```bash
TBD
```

Create required transports:

```bash
/opt/ssa/rpc.py --plugin npi nvmf_create_transport -t NPI --max-queue-depth 4096  --max-io-size 65536 --io-unit-size 4096 --lbads 4096 --log-level ERROR
/opt/ssa/rpc.py nvmf_create_transport -t TCP
```

### Create frontend Nvme controllers

#### PF

```bash
# export BRIDGE_ADDR=10.10.10.10:50051
grpc_cli --json_input --json_output call $BRIDGE_ADDR CreateNvmeSubsystem "{nvme_subsystem : {spec : {nqn: 'nqn.2022-09.io.spdk:opitest-0.0', serial_number: 'mev-opi-serial', model_number: 'mev-opi-model', max_namespaces: 11} }, nvme_subsystem_id : 'subsystem00' }"
grpc_cli --json_input --json_output call $BRIDGE_ADDR CreateNvmeController "{nvme_controller : {spec : {nvme_controller_id: 0, subsystem_id : { value : '//storage.opiproject.org/volumes/subsystem00' }, pcie_id : {physical_function : 0, virtual_function : 0}, max_nsq:5, max_ncq:5} }, nvme_controller_id : 'controller0' }"
```

On Host

```bash
# Bind driver to PF
# export PF_BDF=<bdf address of pf>
modprobe nvme
echo 'nvme' > ./virtfn0/driver_override
echo $PF_BDF > /sys/bus/pci/drivers/nvme/bind

# Allocate resources and prepare to VF creation
cd /sys/bus/pci/devices/$PF_BDF
echo 0 > ./sriov_drivers_autoprobe
echo 4 > ./sriov_numvfs
```

#### VF

```bash
grpc_cli --json_input --json_output call $BRIDGE_ADDR CreateNvmeSubsystem "{nvme_subsystem : {spec : {nqn: 'nqn.2022-09.io.spdk:opitest-0.1', serial_number: 'mev-opi-serial', model_number: 'mev-opi-model', max_namespaces: 11} }, nvme_subsystem_id : 'subsystem01' }"
grpc_cli --json_input --json_output call $BRIDGE_ADDR CreateNvmeController "{nvme_controller : {spec : {nvme_controller_id: 2, subsystem_id : { value : '//storage.opiproject.org/volumes/subsystem01' }, pcie_id : {physical_function : 0, virtual_function : 1}, max_nsq:5, max_ncq:5, max_limit: {rd_iops_kiops: 5}} }, nvme_controller_id : 'controller1' }"
```

On Host

```bash
echo 'nvme' > ./virtfn0/driver_override
# VF_BDF can be found in virtfn<X> where X equals to virtual_function in CreateNvmeController minus 1 e.g.
# virtio_function: 1 -> virtfn0
echo $VF_BDF > /sys/bus/pci/drivers/nvme/bind
```

### Connect to backend remote Nvme/TCP controller

```bash
# export TARGET_IP=200.1.1.11
# export TARGET_PORT=4420
grpc_cli --json_input --json_output call $BRIDGE_ADDR CreateNVMfRemoteController "{nv_mf_remote_controller : {multipath: 'NVME_MULTIPATH_MULTIPATH'}, nv_mf_remote_controller_id: 'nvmetcp12'}"
grpc_cli --json_input --json_output call $BRIDGE_ADDR CreateNVMfPath "{nv_mf_path: {controller_id: {value: '//storage.opiproject.org/volumes/nvmetcp12'}, traddr:'$TARGET_IP', subnqn:'nqn.2016-06.io.spdk:cnode1', trsvcid:'$TARGET_PORT', trtype:'NVME_TRANSPORT_TCP', adrfam:'NVMF_ADRFAM_IPV4', hostnqn:'nqn.2016-06.io.spdk:cnode1'}, nv_mf_path_id: 'nvmetcp12path0'}"
```

### Create middleend QoS volume

```bash
grpc_cli --json_input --json_output call $BRIDGE_ADDR CreateQosVolume "{'qos_volume' : {'volume_id' : { 'value':'nvmetcp12n1'}, 'max_limit' : { 'rw_iops_kiops': 3 } }, 'qos_volume_id' : 'qosvolume0' }"
```

### Create middleend Encrypted volume

```bash
grpc_cli --json_input --json_output call $BRIDGE_ADDR CreateEncryptedVolume "{'encrypted_volume': { 'cipher': 'ENCRYPTION_TYPE_AES_XTS_128', 'volume_id': { 'value': 'nvmetcp12n1'}, 'key': 'MDAwMTAyMDMwNDA1MDYwNzA4MDkwYTBiMGMwZDBlMGY='}, 'encrypted_volume_id': 'encnvmetcp12n1' }"
```

### Create Namespace

middleend volumes should be created before namespace is created

```bash
grpc_cli --json_input --json_output call $BRIDGE_ADDR CreateNvmeNamespace "{nvme_namespace : {spec : {subsystem_id : { value : '//storage.opiproject.org/volumes/subsystem01' }, volume_id : { value : 'nvmetcp12n1' }, 'host_nsid' : '5', uuid:{value : '1b4e28ba-2fa1-11d2-883f-b9a761bde3fc'}, nguid: '1b4e28ba-2fa1-11d2-883f-b9a761bde3fc', eui64: 1967554867335598547 } }, nvme_namespace_id: 'namespace1' }"
```

### Delete Namespace

```bash
grpc_cli --json_input --json_output call $BRIDGE_ADDR DeleteNvmeNamespace "{name : '//storage.opiproject.org/volumes/namespace1'}"
```

### Delete middleend Encrypted volume

```bash
grpc_cli --json_input --json_output call $BRIDGE_ADDR DeleteEncryptedVolume "{'name': '//storage.opiproject.org/volumes/encnvmetcp12n1'}"
```

### Delete middleend QoS volume

```bash
grpc_cli --json_input --json_output call $BRIDGE_ADDR DeleteQosVolume "{name : '//storage.opiproject.org/volumes/qosvolume0'}"
```

### Disconnect from backend remote Nvme/TCP controller

```bash
grpc_cli --json_input --json_output call $BRIDGE_ADDR DeleteNVMfPath "{name: '//storage.opiproject.org/volumes/nvmetcp12path0'}"
grpc_cli --json_input --json_output call $BRIDGE_ADDR DeleteNVMfRemoteController "{name: '//storage.opiproject.org/volumes/nvmetcp12'}"
```

### Delete frontend controllers

On Host

```bash
# unbind driver from controller
echo $BDF > /sys/bus/pci/drivers/nvme/unbind
echo '(null)' > ./virtfn0/driver_override
```

VF

```bash
grpc_cli --json_input --json_output call $BRIDGE_ADDR DeleteNvmeController "{name : '//storage.opiproject.org/volumes/controller1'}"
grpc_cli --json_input --json_output call $BRIDGE_ADDR DeleteNvmeSubsystem "{name : '//storage.opiproject.org/volumes/subsystem01'}"
```

PF

```bash
grpc_cli --json_input --json_output call $BRIDGE_ADDR DeleteNvmeController "{name : '//storage.opiproject.org/volumes/controller0'}"
grpc_cli --json_input --json_output call $BRIDGE_ADDR DeleteNvmeSubsystem "{name : '//storage.opiproject.org/volumes/subsystem00'}"
```

## I Want To Contribute

This project welcomes contributions and suggestions.  We are happy to have the Community involved via submission of **Issues and Pull Requests** (with substantive content or even just fixes). We are hoping for the documents, test framework, etc. to become a community process with active engagement.  PRs can be reviewed by any number of people, and a maintainer may accept.

See [CONTRIBUTING](https://github.com/opiproject/opi/blob/main/CONTRIBUTING.md) and [GitHub Basic Process](https://github.com/opiproject/opi/blob/main/doc-github-rules.md) for more details.
