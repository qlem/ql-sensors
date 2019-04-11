# ql-sensors

Small program who displays sensors input values from **/sys/class/hwmon**. Written in **[golang](https://golang.org/)** and uses the C library **ncurses** for graphics.

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

### Screenshot

![screenshot]()