package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type StationData struct {
	min   float64
	max   float64
	total float64
	count int
}

func getMin(a, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func getMax(a, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func getFile(fileName string) (file *os.File) {
	// Open File
	file, fileError := os.Open(fileName)
	if fileError != nil {
		panic(fileError)
	}

	return file
}

func main() {
	if len(os.Args) < 2 {
		panic("No arguments")
	}
	fileName := os.Args[1]

	started := time.Now()
	file := getFile(fileName)

	iterateThroughFile(file)
	fmt.Printf("%0.6f", time.Since(started).Seconds())

	file.Close()
}

func consumer(channel chan []byte) {
	for {
		<-channel
	}
}

func iterateThroughFile(file *os.File) {

	channel := make(chan []byte) // SYNC
	go consumer(channel)

	// Start processing the records
	//stations := make(map[string]*StationData)
	bufferSize := 256
	buffer := make([]byte, bufferSize*bufferSize)
	for {
		_, err := file.Read(buffer)

		data := make([]byte, bufferSize)
		copy(data, buffer)
		channel <- data

		if err == io.EOF {
			break
		}

		if err != nil {
			panic(err)
		}
	}

}

func printRecords(stations map[string]*StationData) {
	for key, value := range stations {
		average := value.total / float64(value.count)
		fmt.Printf("%s=%f/%f/%f\n", key, value.min, value.max, average)
	}
}

func processRecord(line string, stations map[string]*StationData) {
	parts := strings.Split(line, ";")
	stationName := parts[0]
	// To do, change float to int
	temp, parseError := strconv.ParseFloat(parts[1], 64)

	if parseError != nil {
		fmt.Printf("Error parsing float on line, " + line)
		panic(parseError)
	}

	station := stations[stationName]
	if station == nil {
		stations[stationName] = &StationData{temp, temp, temp, 1}
	} else {
		station.count += 1
		station.total += temp
		if station.min > temp {
			station.min = temp
		}
		if station.max < temp {
			station.max = temp
		}

	}
}

type StnData struct {
	Name  string
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func referenceCode() {
	data := make(map[string]*StnData)

	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ";")
		name := parts[0]
		tempStr := strings.Trim(parts[1], "\n")

		temperature, err := strconv.ParseFloat(tempStr, 64)
		if err != nil {
			panic(err)
		}

		station, ok := data[name]
		if !ok {
			data[name] = &StnData{name, temperature, temperature, temperature, 1}
		} else {
			if temperature < station.Min {
				station.Min = temperature
			}
			if temperature > station.Max {
				station.Max = temperature
			}
			station.Sum += temperature
			station.Count++
		}
	}

	printResult(data)
}

func printResult(data map[string]*StnData) {
	result := make(map[string]*StnData, len(data))
	keys := make([]string, 0, len(data))
	for _, v := range data {
		keys = append(keys, v.Name)
		result[v.Name] = v
	}
	sort.Strings(keys)

	print("{")
	for _, k := range keys {
		v := result[k]
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", k, v.Min, v.Sum/float64(v.Count), v.Max)
	}
	print("}\n")
}
