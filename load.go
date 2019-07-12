package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
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
	load  float64
	color string
}

var midColorsByLoad = []colorByLoad {
	{0.0, "white"},
	{0.5, "blue"},
	{1.0, "cyan"},
	{2.0, "green"},
	{3.0, "yellow"},
	{4.0, "pink"},
}

func loadToColor(load float64) string {
	for _, lc := range midColorsByLoad {
		if load < lc.load {
			return lc.color
		}
	}
	return "red"
}

func setLedColor(color string) {
	cmd := fmt.Sprintf("ring,80,none,%s", color)
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
	return lv, nil
}

func loadAverageMonitor() {
	for {
		load, err := getLoadAverage()
		if err != nil {
			fmt.Println("monitor error:", err)
		}
		setLedColor(loadToColor(load))
		time.Sleep(loadFeedInterval)
	}
}
