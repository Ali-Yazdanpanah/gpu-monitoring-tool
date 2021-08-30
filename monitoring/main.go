package main

import (
	"fmt"
	"log"

	"os"
	"time"

	"encoding/csv"

	"github.com/NVIDIA/go-nvml/pkg/nvml"
)

func checkError(message string, err error) {
	if err != nil {
		log.Fatalf(message, err)
	}
}

func main() {
	ret := nvml.Init()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to initialize NVML: %v", nvml.ErrorString(ret))
	}
	defer func() {
		ret := nvml.Shutdown()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to shutdown NVML: %v", nvml.ErrorString(ret))
		}
	}()

	_, ret = nvml.DeviceGetCount()
	if ret != nvml.SUCCESS {
		log.Fatalf("Unable to get device count: %v", nvml.ErrorString(ret))
	}

	for i := 0; i < count; i++ {
		go monitor(i)
	}

}

func buildFileName() string {
	return time.Now().Format("2006-01-02 15:04:05") + ".csv"
}

func monitor(index int) {
	// Create results file
	fileName := buildFileName()
	file, err := os.Create(fileName)
	checkError("Cannot create file", err)
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	fmt.Println("Recording started!\n")
	for {
		device, ret := nvml.DeviceGetHandleByIndex(index)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get device at index %d: %v", index, nvml.ErrorString(ret))
		}

		uuid, ret := device.GetUUID()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get uuid of device at index %d: %v", index, nvml.ErrorString(ret))
		}

		// fmt.Printf("Device id: %v\n", uuid)
		powerUsage, ret := device.GetPowerUsage()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get power usage of device at index %d: %v", index, nvml.ErrorString(ret))
		}
		// fmt.Printf("Power usage: %v\n", powerUsage)

		//fanSpeed, ret := device.GetFanSpeed()
		// if ret != nvml.SUCCESS {
		//     log.Fatalf("Unable to get fan speed of device at index %d: %v", index, nvml.ErrorString(ret))
		// }
		// clock, ret := device.GetClock()
		// if ret != nvml.SUCCESS {
		// 	log.Fatalf("Unable to get clock of device at index %d: %v", index, nvml.ErrorString(ret))
		// }
		// fmt.Printf("Clock: %v\n", clock)

		clockInfo, ret := device.GetClockInfo(0)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get clock info of device at index %d: %v", index, nvml.ErrorString(ret))
		}
		// fmt.Printf("Clock info: %v\n", clockInfo)

		utilizationRate, ret := device.GetUtilizationRates()
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get utilization rates of device at index %d: %v", index, nvml.ErrorString(ret))
		}
		// fmt.Printf("Utilization rates:\n")
		// fmt.Printf("    GPU: %v\n", utilizationRate.Gpu)
		// fmt.Printf("    MEM: %v\n", utilizationRate.Memory)
		temp, ret := device.GetTemperature(0)
		if ret != nvml.SUCCESS {
			log.Fatalf("Unable to get temperature of device at index %d: %v", index, nvml.ErrorString(ret))
		}
		// fmt.Printf("Temperature: %v\n", temp)
		entry := []string{uuid, fmt.Sprint(powerUsage), fmt.Sprint(clockInfo), fmt.Sprint(utilizationRate.Gpu), fmt.Sprint(utilizationRate.Memory), fmt.Sprint(temp)}
		err := writer.Write(entry)
		checkError("Cannot write to file", err)

		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println("Recording Stoped!\n")
}

