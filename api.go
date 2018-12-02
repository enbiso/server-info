package main

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// ServerInfo structure
type ServerInfo struct {
	Mem   MemInfo    `json:"mem"`
	Swap  MemInfo    `json:"swap"`
	Temps []TempInfo `json:"temperatures"`
	CPUs  []CPUInfo  `json:"cpus"`
}

func executeAPI(addr string) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, ServerInfo{
			Mem:   getMemInfo(),
			Swap:  getSwapInfo(),
			Temps: getTempInfo(),
			CPUs:  getCPUInfo(),
		})
	})
	e.Logger.Fatal(e.Start(addr))
}

// MemInfo structure
type MemInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}

func getMemInfo() MemInfo {
	v, _ := mem.VirtualMemory()
	return MemInfo{
		Total:       v.Total,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}
}

func getSwapInfo() MemInfo {
	s, _ := mem.SwapMemory()
	return MemInfo{
		Total:       s.Total,
		Used:        s.Used,
		UsedPercent: s.UsedPercent,
	}
}

// TempInfo structure
type TempInfo struct {
	Sensor      string  `json:"sensor"`
	Temperature float64 `json:"temperature"`
	High        float64 `json:"high"`
	Critical    float64 `json:"critical"`
}

func getTempInfo() []TempInfo {
	temStats, _ := host.SensorsTemperatures()
	temps := []TempInfo{}
	temp := TempInfo{}
	for _, tStat := range temStats {
		if strings.HasSuffix(tStat.SensorKey, "input") {
			if temp.Sensor != "" {
				temps = append(temps, temp)
			}
			temp = TempInfo{
				Sensor: tStat.SensorKey[:strings.LastIndexByte(tStat.SensorKey, '_')],
			}
			temp.Temperature = tStat.Temperature
		}
		if strings.HasSuffix(tStat.SensorKey, "max") {
			temp.High = tStat.Temperature
		}
		if strings.HasSuffix(tStat.SensorKey, "crit") {
			temp.Critical = tStat.Temperature
		}
	}
	// Add the last temp
	if temp.Sensor != "" {
		temps = append(temps, temp)
	}
	return temps
}

// CPUInfo structure
type CPUInfo struct {
	Name        string  `json:"Name"`
	User        float64 `json:"user"`
	System      float64 `json:"system"`
	Idle        float64 `json:"idle"`
	UsedPercent float64 `json:"usedPercent"`
}

func getCPUInfo() []CPUInfo {
	cs, _ := cpu.Times(true)
	cInfos := []CPUInfo{}
	for _, c := range cs {
		cInfos = append(cInfos, CPUInfo{
			Name:        c.CPU,
			User:        c.User,
			System:      c.System,
			Idle:        c.Idle,
			UsedPercent: ((c.User + c.System) / (c.User + c.System + c.Idle)) * 100,
		})
	}
	return cInfos
}
