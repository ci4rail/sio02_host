# tracelet_host
Host examples and protobuf files for "Tracelets" (SIO02/03/04 devices)

## Protobuf definitions
The `proto` files can be found [here](./io4edge_api/tracelet/proto/v1/tracelet.proto)
...

## Examples

The `examples` folder contains some simple examples in python:
* [location_server](examples/location_server.py): A UDP server that receives location messages from tracelets and prints the contents.

# Usage
- [Device-Simulator](README_DEVSIM.md): A tracelet simulator

- [io4edge-cli](https://github.com/ci4rail/io4edge-client-go): Contains a command line tool to manage io4edge devices, such as the SIO02.
