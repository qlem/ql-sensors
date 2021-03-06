package main

// #cgo LDFLAGS: -lncurses
// #include <stdlib.h>
// #include <unistd.h>
// #include <curses.h>
// #include <termios.h>
//
// static void printWindow(WINDOW *win, int y, int x, char *str) {
//		mvwprintw(win, y, x, "%s", str);
// }
import "C"
import "unsafe"

type Terminal struct {
	window      *C.WINDOW
	initialConf C.struct_termios
}

func setReverse(term *Terminal) {
	C.wattron(term.window, C.A_REVERSE)
}

func unsetReverse(term *Terminal) {
	C.wattroff(term.window, C.A_REVERSE)
}

func setBold(term *Terminal) {
	C.wattron(term.window, C.A_BOLD)
}

func unsetBold(term *Terminal) {
	C.wattroff(term.window, C.A_BOLD)
}

func printWindow(term *Terminal, y int, x int, str string) {
	cstr := C.CString(str)
	C.printWindow(term.window, C.int(y), C.int(x), cstr)
	C.free(unsafe.Pointer(cstr))
}

func initWindow(term *Terminal, nbLines int, nbColumns int) {
	term.window = C.newwin(C.int(nbLines), C.int(nbColumns), 0, 0)
}

func windowClear(term *Terminal) {
	C.wclear(term.window)
}

func windowRefresh(term *Terminal) {
	C.wrefresh(term.window)
}

func deleteWindow(term *Terminal) {
	C.delwin(term.window)
}

func endWindow() {
	C.endwin()
}

func nonCanonicalMode() {
	var config C.struct_termios

	C.tcgetattr(C.fileno(C.stdin), &config)
	config.c_lflag &^= C.ICANON
	config.c_lflag &^= C.ECHO
	config.c_cc[C.VMIN] = 0
	config.c_cc[C.VTIME] = 1
	C.tcsetattr(C.fileno(C.stdin), C.TCSANOW, &config)
}

func resetTerminal(term *Terminal) {
	C.tcsetattr(C.fileno(C.stdin), C.TCSANOW, &term.initialConf)
}

func saveInitialConfig(term *Terminal) {
	C.tcgetattr(C.fileno(C.stdin), &term.initialConf)
}

func initTerminal() {
	C.initscr()
	C.cbreak()
	C.noecho()
	C.keypad(C.stdscr, true)
	C.curs_set(0)
}
