# ql-sensors

Small program who provides sensors input values monitoring from **/sys/class/hwmon** in shell. Written in **[golang](https://golang.org/)** and uses the C library **ncurses** for graphics.

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
