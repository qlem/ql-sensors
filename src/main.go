package main

import (
	"bufio"
	"container/list"
	"os"
	"os/signal"
	"time"
)

func recoverTerm(term *Terminal) {
	deleteWindow(term)
	endWindow()
	resetTerminal(term)
}

func loop(term *Terminal, sensors *list.List) {
	reader := bufio.NewReader(os.Stdin)
	run := true
	ref := time.Unix(0, 0)
	for run {
		buffer := make([]byte, 10)
		_, _ = reader.Read(buffer)
		if buffer[0] == 'q' {
			run = false
		}

		now := time.Now()
		delta := now.Sub(ref)
		if delta.Seconds() >= 1 {
			windowClear(term)
			refreshSensorList(sensors)
			print_(term, sensors)
			windowRefresh(term)
			ref = time.Now()
		}
	}
}

func main() {
	term := new(Terminal)
	sensors := list.New()

	// handle ^C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		recoverTerm(term)
		os.Exit(1)
	}()

	saveInitialConfig(term)
	initTerminal()
	initWindow(term, 0, 0)
	nonCanonicalMode()
	loop(term, sensors)
	recoverTerm(term)
}
