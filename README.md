# ql-sensors

Small program who provides **sensors monitoring** from `/sys/class/hwmon` in the shell. Written in **[golang](https://golang.org/)** and uses the C library **ncurses** for graphics. Only compatible with Linux distributions.

#### Screenshot
![screenshot](https://raw.githubusercontent.com/qlem/ql-sensors/master/screenshot.png)

### Installation

The default installation location is `/usr/local/bin`.
```
$> make
$> make install
```

Run this for uninstall:
```
$> make uninstall
```

### Usage

Simply run:
```
$> qlsensors
```
