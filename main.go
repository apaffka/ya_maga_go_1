// чисто для проверки ещё раз
package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func main() {
	errorCount := 0

	for {
		resp, err := http.Get("http://srv.msk01.gigacorp.local/_stats")
		if err != nil || resp.StatusCode != 200 {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				break
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				break
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		fields := strings.Split(strings.TrimSpace(string(body)), ",")
		if len(fields) != 7 {
			errorCount++
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				break
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		errorCount = 0

		values := make([]int64, 7)
		for i := range fields {
			v, err := strconv.ParseInt(fields[i], 10, 64)
			if err != nil {
				errorCount++
				break
			}
			values[i] = v
		}
		if errorCount > 0 {
			if errorCount >= 3 {
				fmt.Println("Unable to fetch server statistic")
				break
			}
			time.Sleep(100 * time.Millisecond)
			continue
		}

		load := values[0]
		totalMem, usedMem := values[1], values[2]
		totalDisk, usedDisk := values[3], values[4]
		totalNet, usedNet := values[5], values[6]

		if load > 30 {
			fmt.Printf("Load Average is too high: %d\n", load)
		}

		if totalMem > 0 {
			memUsage := usedMem * 100 / totalMem
			if memUsage > 80 {
				fmt.Printf("Memory usage too high: %d%%\n", memUsage)
			}
		}

		if totalDisk > 0 {
			diskUsage := usedDisk * 100 / totalDisk
			if diskUsage > 90 {
				freeMb := (totalDisk - usedDisk) / (1024 * 1024)
				fmt.Printf("Free disk space is too low: %d Mb left\n", freeMb)
			}
		}

		if totalNet > 0 {
			netUsage := usedNet * 100 / totalNet
			if netUsage > 90 {
				freeMbit := (totalNet - usedNet) / 1000000
				fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeMbit)
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}