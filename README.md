# ql-sensors

Small program who provides **sensors monitoring** from `/sys/class/hwmon` in the shell. Written in **golang** and uses the C library **ncurses** for graphics. Only compatible with Linux distributions.

#### Screenshot
![screenshot](https://raw.githubusercontent.com/qlem/ql-sensors/master/screenshot.png)

### Installation

[Go](https://golang.org/) is require for compilation only. The default installation location is `/usr/local/bin`.

To build and install run:
```
$> make
$> make install
```

To uninstall run:
```
$> make uninstall
```

### Usage

Simply run:
```
$> qlsensors
```
