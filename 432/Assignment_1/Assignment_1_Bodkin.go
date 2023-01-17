// Assignment_1_Bodkin.go
// 
// Assignment 1 code for MSDS 432
// This code just imports a CSV file that I have downloaded, 
// uses the geocoder package to find the zip codes and puts
// both the zip codes and the day of the week into csv.

package main

import(
	"encoding/csv"
	"fmt"
	"os"
	"time"
	"strconv"
	"strings"

	"github.com/kelvins/geocoder"
)


type entry_line struct{
	crash_date time.Time
	lat float64
	long float64
	weekday string
	zip string
	weather string
	lighting string
	h_a_r string
}

type zip_counter struct{
	crash_c int
	run_c int
	weather_c map[string]int
	lighting_c map[string]int
}

func err_kicker(in_err error){
	if in_err!=nil{
		fmt.Println(in_err)
		os.Exit(1)
	}
	return
}


func main() {
	// get the filename from the input args, if we gave one
	if len(os.Args)<2 {
		fmt.Println("Please provide a filename")
		return
	}

	// csv_ext := "csv"
	if os.Args[1][len(os.Args[1])-3:] != "csv" {
		fmt.Println("Input must be a csv file!")
	}

	file,err := os.Open(os.Args[1]) // open a file pointer
	err_kicker(err) // kick us out of the program if there's an error

	// import the data from the csv
	reader := csv.NewReader(file) // create a csv reader
	records,err := reader.ReadAll()
	err_kicker(err) // kick us out of the program if there's an error

	
	// import map credentials
	geocoder.ApiKey = os.Getenv("Map_Credential")


	// get the indices of the latitude and longitude columns
	// iterate per column
	var lat_idx,long_idx,date_idx,run_idx,weather_idx,lighting_idx int	
	for idx, value := range records[0] {
		switch value{
		case "LATITUDE":
			lat_idx = idx
		case "LONGITUDE":
			long_idx = idx
		case "CRASH_DATE":
			date_idx = idx
		case "HIT_AND_RUN_I":
			run_idx = idx
		case "WEATHER_CONDITION":
			weather_idx = idx
		case "LIGHTING_CONDITION":
			lighting_idx = idx
		
		}

	}

	// fmt.Printf("%d, %d, %d",weather_idx,lighting_idx, run_idx)

	// create a map for each unique lat/long, with a slice as the value 
	var unique_loc = make(map[string]string) // unique lats and longs
	var unique_zip = make(map[string]zip_counter) // counter per zip
	var line_map = make(map[int]entry_line)
	var weekday_c = make(map[string]int)    // counting crashes per weekday in 2021
	// n_records := len(records[1:][1])
	for row_i,row := range records{
		if row_i == 0 {continue}
		// if row_i > 10000 {continue}

		// pull out all of the important data
		latitude := row[lat_idx] 	// latitude
		longitude := row[long_idx] 	// longitude
		weather := row[weather_idx] // weather condition
		lighting := row[lighting_idx] // lighting condition

		date,err := time.Parse("01/02/2006 03:04:05 PM", row[date_idx])
		if err != nil {
			fmt.Println(row_i,err)
			return
		}
		
		// create temporary entry_line instance
		var temp_line entry_line
		temp_line.lat,_ = strconv.ParseFloat(latitude,16)
		temp_line.long,_ = strconv.ParseFloat(longitude,16)
		temp_line.crash_date = date
		temp_line.weekday = date.Weekday().String()
		temp_line.weather = weather
		temp_line.lighting = lighting

		line_map[row_i] = temp_line

		if (date.Year() == 2020) && (date.Month().String() == "December") {

			// create a new string that appends the latitude and longitude
			lat_long := latitude+longitude
			// append on the row, and replace it in the map
			// _, lat_exist := unique_loc[lat_long]

			// update the unique lat/long map with counts etc
			if zip,lat_exist := unique_loc[lat_long]; !lat_exist{
				location := geocoder.Location{
					Latitude: temp_line.lat,
					Longitude: temp_line.long,
				}
				addresses, err := geocoder.GeocodingReverse(location)
				if err != nil {
					// fmt.Println("Rethink yo geocoding")
					continue
				} else {
					address := addresses[1].FormattedAddress
					
					// find the zip in the address -- all Chicago zips start with 606
					zip_start := strings.Index(address, "606")
					if zip_start<0{continue} // skip if there's an error finding the zip code
					
					// find the zip
					zip = address[zip_start:zip_start+5]
			
					// fmt.Println("Zip: ", zip)
					unique_loc[lat_long] = zip
				}

				
			} 
			
			// update the unique zip codes map
			zip := unique_loc[lat_long]
			if temp_zip,zip_exist := unique_zip[zip]; !zip_exist {
				
				// start filling in the fields
				temp_zip.crash_c = 1

				// is this a hit and run?
				if row[run_idx] == "Y" {
					temp_zip.run_c = 1
				} else {
					temp_zip.run_c = 0
				}

				// weather counter
				temp_weather := make(map[string]int)
				temp_weather[row[weather_idx]] = 1
				temp_light := make(map[string]int)
				temp_light[row[lighting_idx]] = 1
				temp_zip.lighting_c = temp_light
				temp_zip.weather_c = temp_weather

				// push the temp into the map
				unique_zip[zip] = temp_zip

			} else {
				// update counters
				temp_zip.crash_c++

				// hit and runs
				if row[run_idx] == "Y" {
					temp_zip.run_c++
				}
				
				// weather and lighting
				temp_zip.weather_c[weather]++
				temp_zip.lighting_c[lighting]++

				// save back to map
				unique_zip[zip] = temp_zip

			}


		}


		// Weekdays for 2021
		if date.Year() == 2021 {
			weekday_c[date.Weekday().String()]++
		}


		// Status Updates
		if row_i%1000 == 0{
			fmt.Printf("Processing record %d\n", row_i)
		}			
	}

	// ------------------------------------------------------
	// Start printing everything out
	// ------------------------------------------------------

	fmt.Println("\n\n")
	fmt.Println("Weekday counts for 2021")
	fmt.Println("-----------------------------")
	for day,day_c := range weekday_c {
		fmt.Printf("\t%s: %d\n", day, day_c)
	}

	fmt.Println("\n\nInformation per zip code in Dec 2020:\n")
	fmt.Printf("Zip\t\tTotal\t\tHit and Run\n")
	fmt.Println("-----------------------------")
	for zip_u,zip_struct := range unique_zip {
		fmt.Printf("%s\t\t%d\t\t%d\n", zip_u, zip_struct.crash_c, zip_struct.run_c)
	}
	

	fmt.Println("\n\nLat/Long converted into Zip Codes\n")
	fmt.Printf("Lat\t\tLong\t\t|\tZip\n")
	fmt.Println("-----------------------------")
	for lat_lon, zip := range unique_loc {
		fmt.Printf("%s\t%s\t|\t%s\n", lat_lon[0:11], lat_lon[12:], zip)
	}

	// fmt.Println(unique_zip)

	file.Close()

}