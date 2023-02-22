# Devsim - SIO02 Device Simulator

Devsim is an executable for Windows and Linux that simulates roughly the behaviour of a SIO02 Tracelet.
Can be used to test host code without a physical SIO02.

Executables can be found [here](https://github.com/ci4rail/SIO02_host/releases).

## Functionality

* Sends tracelet location messages to the localization server every 0.5s
    * Every 1.5s the content of the location message is changed, see `locationGenerator` function in `/devsim/internal/tracelet/location.go`
* Responds to status requests from the localization server
    * Answers with fixed values, see `commandHandler` function in `/devsim/internal/tracelet/location.go`

## Usage

Run `devsim` on one computer and your host code (e.g. `examples/location_server.py`) on another computer.


### Testing with Linux Machines

* `devsim` on computer A at IP: `192.168.0.100`
* `location_server` on computer B at IP: `192.168.0.200`

On computer B, clone this repo

```bash
$ git clone https://github.com/ci4rail/sio02_host.git --recursive
$ cd sio02_host
$ export PYTHONPATH=`pwd`
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
2023/02/22 08:24:27 devsim version: dev
2023/02/22 08:24:27 try to connect to 127.0.0.1:11002
locationClient WriteMessage: receive_ts:{seconds:1677050686 nanos:910927234} tracelet_id:"devsim" location:{gnss:{valid:true latitude:49.425111 longitude:11.077378 altitude:350 eph:0.4 epv:2.5} uwb:{y:1100 z:888 site_id:4660 location_signature:20015998348237 cov_xx:11.1 cov_xy:12.2} speed:9 mileage:50899 temperature:34.5}
locationClient WriteMessage: receive_ts:{seconds:1677050687 nanos:411838137} tracelet_id:"devsim" location:{gnss:{valid:true latitude:49.425111 longitude:11.077378 altitude:350 eph:0.4 epv:2.5} uwb:{y:1100 z:888 site_id:4660 location_signature:20015998348237 cov_xx:11.1 cov_xy:12.2} speed:9 mileage:50899 temperature:34.5}
locationClient WriteMessage: receive_ts:{seconds:1677050687 nanos:912741641} tracelet_id:"devsim" location:{gnss:{latitude:49.425111 longitude:11.077378 altitude:350 eph:0.4 epv:2.5} uwb:{valid:true x:5 y:6.21 z:7.5 site_id:4660 location_signature:20015998348237 cov_xx:11.1 cov_xy:12.2} speed:9 mileage:50899 temperature:34.5}
...
```

On Computer B:
```
new handler Thread-1

message from devsim, ts=2023-02-22 07:24:46.910927
  devsim 2023-02-22 07:24:46.910927
     UWB: valid False 0.00 1100.00 ite:4660
    GNSS: valid True 49.425111 11.077378 0.40
message from devsim, ts=2023-02-22 07:24:47.411838
  devsim 2023-02-22 07:24:47.411838
     UWB: valid False 0.00 1100.00 ite:4660
    GNSS: valid True 49.425111 11.077378 0.40
message from devsim, ts=2023-02-22 07:24:47.912741
  devsim 2023-02-22 07:24:47.912741
     UWB: valid True 5.00 6.21 ite:4660
    GNSS: valid False 49.425111 11.077378 0.40
```

### Testing Location Messages with Windows Machines

* `devsim` on computer A at IP: `192.168.0.100`
* `location_server` on computer B at IP: `192.168.0.200`

On computer B, clone this repo

```
> git clone https://github.com/ci4rail/sio02_host.git --recursive
> cd sio02_host
> set PYTHONPATH=%cd%
> cd examples
> pip3 install -r requirements.txt
> python location_server.py
```

On computer A, download the `devsim` binary for the binary for your platform [from the releases](https://github.com/ci4rail/SIO02_host/releases).

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
