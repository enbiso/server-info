package main

import (
	"strings"

	"github.com/labstack/echo"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// ServerInfo structure
type ServerInfo struct {
	Mem   MemInfo    `json:"mem"`
	Swap  MemInfo    `json:"swap"`
	Temps []TempInfo `json:"temperatures"`
}

func executeAPI(addr string) {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.JSON(200, ServerInfo{
			Mem:   getMemInfo(),
			Swap:  getSwapInfo(),
			Temps: getTempInfo(),
		})
	})
	e.Logger.Fatal(e.Start(addr))
}

// MemInfo structure
type MemInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercentage"`
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
	var temp *TempInfo
	for _, tStat := range temStats {
		if strings.HasSuffix(tStat.SensorKey, "input") {
			if temp != nil {
				temps = append(temps, *temp)
			}
			temp = &TempInfo{
				Temperature: tStat.Temperature,
				Sensor:      strings.Replace(tStat.SensorKey, "_input", "", -1),
			}
		}
		if temp != nil {
			if strings.HasSuffix(tStat.SensorKey, "max") {
				temp.High = tStat.Temperature
			}
			if strings.HasSuffix(tStat.SensorKey, "crit") {
				temp.Critical = tStat.Temperature
			}
		}
	}
	return temps
}

// CPUInfo structure
type CPUInfo struct {
}
