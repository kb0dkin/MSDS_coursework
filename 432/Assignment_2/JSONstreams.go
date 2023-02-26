package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"runtime"
	"time"
	"os"
	"strconv"
	"runtime/pprof"
)

type Data struct {
	Key string `json:"key"`
	Val int    `json:"value"`
}

var DataRecords []Data

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

var MIN = 0
var MAX = 26

func getString(l int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == l {
			break
		}
		i++
	}
	return temp
}

// DeSerialize decodes a serialized slice with JSON records
func DeSerialize(e *json.Decoder, slice interface{}) error {
	return e.Decode(slice)
}

// Serialize serializes a slice with JSON records
func Serialize(e *json.Encoder, slice interface{}) error {
	return e.Encode(slice)
}

func main() {

	// number of records
	n_records := 10000
	if len(os.Args) > 1 {
		n_records,_ = strconv.Atoi(os.Args[1])
	} 
	
	// Create sample data
	var i int
	var t Data
	for i = 0; i < n_records; i++ {
		t = Data{
			Key: getString(5),
			Val: random(1, 100),
		}
		DataRecords = append(DataRecords, t)
	}

	// bytes.Buffer is both an io.Reader and io.Writer

	// hw assignment info:

	// setup the cpu profiler
	f, err:= os.Create("JSONstreams.prof")
	if err != nil {
		fmt.Println(err)
		return
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()


	t_s0 := time.Now() // current time in ms
	var m_s0 runtime.MemStats
	runtime.ReadMemStats(&m_s0)
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	err = Serialize(encoder, DataRecords)
	if err != nil {
		fmt.Println(err)
		return
	}
	t_se := time.Since(t_s0) // elapsed time
	var m_s1 runtime.MemStats
	runtime.ReadMemStats(&m_s1)
	

	decoder := json.NewDecoder(buf)
	var temp []Data
	t_d0 := time.Now() // current time
	var m_d0 runtime.MemStats // creating a new variable
	runtime.ReadMemStats(&m_d0) // and filling it with info
	err = DeSerialize(decoder, &temp)
	if err != nil {
		fmt.Println(err)
		return
	}
	t_de := time.Since(t_d0) // Time elapsed
	var m_d1 runtime.MemStats
	runtime.ReadMemStats(&m_d1)

	// ------------------------------------------
	// Printing out the summaries in a pretty way
	
	fmt.Printf("\nSummaries for %v records:\n", n_records)
	
	fmt.Println("\nTime Elapsed:")
	fmt.Printf("\tSerialization: %v\n", t_se)
	fmt.Printf("\tDeserialization: %v\n", t_de)

	fmt.Println("\nMemory Usage:")
	fmt.Printf("\tSerialization: %v kilobytes\n", (m_s1.Alloc - m_s0.Alloc))
	fmt.Printf("\tDeserialization: %v kilobytes", (m_d1.Alloc - m_d0.Alloc))

}
