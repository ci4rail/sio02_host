# Devsim - SIO02 Device Simulator

Devsim is an executable for Windows and Linux that simulates roughly the behaviour of a SIO02 Tracelet.
Can be used to test host code without a physical SIO02.

Executables can be found [here](https://github.com/ci4rail/SIO02_host/releases).

## Functionality

TODO
* Sends tracelet location messages to a TCP server
    * simulating a moving vehicle with approx 10km/h
* Acts as a status server for a monitoring system
    * Status server is announced via zeroconf/mdns

## Usage

Please run `devsim` on one computer and your host code (e.g. `examples/location_server.py`) on another computer.


### Testing with Linux Machines

* `devsim` on computer A at IP: `192.168.0.100`
* `location_server` on computer B at IP: `192.168.0.200`

On computer B, have the `examples` folder of this repo

```bash
$ cd examples
$ pip3 install -r requirements.txt
$ ./location_server.py
```

On computer A, download the `devsim` binary for the binary for your platform [from the releases](https://github.com/ci4rail/SIO02_host/releases).

```bash
$ tar xvf devsim-<version>-linux-<arch>.tar.gz
$ ./devsim -l 192.168.0.200:11002
```

You should see something like this:

On Computer A:
```
2022/01/07 20:18:06 devsim version: dev
2022/01/07 20:18:06 mdns advertisting service on IPs [192.168.0.100]
2022/01/07 20:18:06 locationGenerator: havePos=true x=-100.00 y=-100.00
2022/01/07 20:18:07 locationGenerator: havePos=true x=-96.80 y=-96.80
2022/01/07 20:18:08 locationGenerator: havePos=true x=-93.60 y=-93.60
...
```

On Computer B:
```
devsim 2022-01-07 20:18:07.516261 -96.80 3.20 site:12345 sign: 5124095577148911
devsim 2022-01-07 20:18:08.520163 -93.60 6.40 site:12345 sign: 5124095577148911
```

On Computer B, open a second shell and execute the status client

```bash
$ ./status_client.py devsim
192.168.0.100 10000
0 power ups
has server connection
has valid time
has valid position
eloc module status is ok
```

### Testing Location Messages with Windows Machines

* `devsim` on computer A at IP: `192.168.0.100`
* `location_server` on computer B at IP: `192.168.0.200`

On computer B, have the `examples` folder of this repo

```
> cd examples
> pip3 install -r requirements.txt
> python location_server.py
```

On computer A, download the `devsim` binary for the binary for your platform [from the releases](https://github.com/ci4rail/SIO02_host/releases).

```
> unzip devsim-<version>-windows-<arch>.tar.gz
> devsim.exe -l 192.168.0.200:11002
```

On Computer B, open a second command prompt and execute the status client

```
> python status_client.py devsim
```

The outputs should be similar as for Linux above.

### Running devsim and Host Code on the Same Machine

It is possible to  run `devsim` and the host code on the same computer if
* the OS is Linux
* the OS is Windows and the status server within `devsim` is not needed

Example for Linux:

In a first terminal:
```bash
$ ./location_server.py
```
Run devsim in a second terminal:
```bash
./devsim --mdns-ip=127.0.0.1 -l 127.0.0.1:11002
```
