package main

import (
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
}

func getTempInfo() []TempInfo {
	temStats, _ := host.SensorsTemperatures()
	temps := []TempInfo{}
	for _, tStat := range temStats {
		temps = append(temps, TempInfo{
			Sensor:      tStat.SensorKey,
			Temperature: tStat.Temperature,
		})
	}
	return temps
}

// CPUInfo structure
type CPUInfo struct {
}
