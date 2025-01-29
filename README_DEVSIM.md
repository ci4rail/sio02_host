# Devsim - Tracelet Simulator

Devsim is an executable for Windows and Linux that simulates roughly the behaviour of a Tracelet.
Can be used to test host code without a physical tracelet.

Executables can be found [here](https://github.com/ci4rail/tracellet_host/releases).

## Functionality

* Sends tracelet location messages with random data

## Usage

Run `devsim` on one computer and your host code (e.g. `examples/location_server.py`) on another computer.


### Testing with Linux Machines

* `devsim` on computer A at IP: `192.168.0.100`
* `location_server` on computer B at IP: `192.168.0.200`

On computer B, clone this repo

```bash
$ git clone https://github.com/ci4rail/tracelet_host.git --recursive
$ cd tracelet_host
$ export PYTHONPATH=`pwd`
$ cd examples
$ pip3 install -r requirements.txt
$ ./location_server.py
```

On computer A, download the `devsim` binary for the binary for your platform [from the releases](https://github.com/ci4rail/tracelet_host/releases).

```bash
$ tar xvf devsim-<version>-linux-<arch>.tar.gz
$ ./devsim -l 192.168.0.200:11002
```

### Testing Location Messages with Windows Machines

* `devsim` on computer A at IP: `192.168.0.100`
* `location_server` on computer B at IP: `192.168.0.200`

On computer B, clone this repo

```
> git clone https://github.com/ci4rail/tracelet_host.git --recursive
> cd tracelet_host
> set PYTHONPATH=%cd%
> cd examples
> pip3 install -r requirements.txt
> python location_server.py
```

On computer A, download the `devsim` binary for the binary for your platform [from the releases](https://github.com/ci4rail/tracelet_host/releases).

```
> unzip devsim-<version>-windows-<arch>.tar.gz
> devsim.exe -l 192.168.0.200:11002
```


### Running devsim and Host Code on the Same Machine

It is possible to  run `devsim` and the host code on the same computer.

Example for Linux:

In a first terminal:
```bash
$ ./location_server.py
```
Run devsim in a second terminal:
```bash
./devsim -l 127.0.0.1:11002
```
