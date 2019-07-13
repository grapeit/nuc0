package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	ledDriverFile = "/proc/acpi/nuc_led"
	loadAvgFile = "/proc/loadavg"
	loadFeedInterval = 5 * time.Second
)

type colorByLoad struct {
	load float64
	color string
	brightness int
}

var blinkLoadThreshold = 4.0
var blinkLedMode = "none"

var colorsByLoad = []colorByLoad {
	{0, "white", 80},
	{0.02, "blue", 8},
	{0.05, "blue", 20},
	{0.10, "blue", 40},
	{0.15, "blue", 80},
	{0.25, "cyan", 50},
	{0.5, "green", 80},
	{0.75, "yellow", 60},
	{1.0, "pink", 60},
	{0, "red", 80},
}

func init() {
	cores := float64(runtime.NumCPU())
	blinkLoadThreshold *= cores
	for i, _ := range colorsByLoad[:len(colorsByLoad) - 1] {
		colorsByLoad[i].load *= cores
	}
}

func loadToColor(load float64) (string, int) {
	for _, lc := range colorsByLoad[:len(colorsByLoad) - 1] {
		if load < lc.load {
			return lc.color, lc.brightness
		}
	}
	last := colorsByLoad[len(colorsByLoad) - 1]
	return last.color, last.brightness
}

func setRingLed(color string, brightness int) {
	cmd := fmt.Sprintf("ring,%d,%s,%s", brightness, blinkLedMode, color)
	err := ioutil.WriteFile(ledDriverFile, []byte(cmd), 0644)
	if err != nil {
		fmt.Println("set error:", err)
	}
}

func getLoadAverage() (float64, error) {
	load, err := ioutil.ReadFile(loadAvgFile)
	if err != nil {
		return -1, err
	}
	ls := strings.Split(string(load), " ")
	if len(ls) < 1 {
		return -1, errors.New("no data")
	}
	lv, err := strconv.ParseFloat(ls[0], 64)
	if err != nil {
		return -1, err
	}
	if lv > blinkLoadThreshold {
		blinkLedMode = "fade_fast"
	} else {
		blinkLedMode = "none"
	}
	return lv, nil
}

func loadAverageMonitor() {
	for {
		load, err := getLoadAverage()
		if err != nil {
			fmt.Println("getLoadAverage error", err)
		}
		setRingLed(loadToColor(load))
		time.Sleep(loadFeedInterval)
	}
}
