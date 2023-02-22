# sio02_host
Host examples and protobuf files for SIO02 "Tracelet"

## Protobuf definitions
The `proto` files can be found [here](./io4edge_api/tracelet/proto/v1/tracelet.proto)
...

## Examples

The `examples` folder contains some simple examples in python:
* [location_server](examples/location_server.py): A TCP server that receives location messages from SIO02 devices and prints the contents.

# Usage
- [SIO02-Device-Simulator](README_DEVSIM.md): A SIO02 simulator

- [io4edge-cli](https://github.com/ci4rail/io4edge-client-go): Contains a command line tool to manage io4edge devices, such as the SIO02.
