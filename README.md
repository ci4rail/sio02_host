# sio02_host
Host examples and protobuf files for SIO02 device

## Protobuf definitions
The `proto` folder contains
...

## Examples

The `examples` folder contains some simple examples in python:
* [location_server](examples/location_server.py): A TCP server that receives location messages from SIO02 devices and prints the location.

### Usage

Prerequisites:
* Linux machine
* Python >= 3.9 installed

```bash
cd examples
pip3 install -r requirements.txt
```

#### Run location server:
```bash
./location_server.py
```
On the SIO02, configure parameter `loc-srv` to the IP address of the machine executing `location_server.py` and port 11002, e.g. `192.168.0.200:11002`.


# See also
- [SIO01-Device-Simulator](README_DEVSIM.md): A SIO01 simulator

- [io4edge-cli](https://github.com/ci4rail/io4edge-client-go): Contains a command line tool to manage io4edge devices, such as the SIO02. 
