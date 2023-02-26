package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Ignoring empty fields in JSON
type MSDSCourse struct {
	CID string `json:"course_ID"`
	CNAME string `json:"course_name"`
	CPREREQ int    `json:"prerequisite"`
}

func MSDSCourse_Returner(course_i int) MSDSCourse {
	// initialize the return value
	var ret_course MSDSCourse

	switch course_i {
		case 0:
			ret_course.CID = "MSDS 42"
			ret_course.CNAME = "Life, The Universe, and Everything"
			ret_course.CPREREQ = 0

		case 1:
			ret_course.CID = "MSDS ??"
			ret_course.CNAME = "The Question"
			ret_course.CPREREQ = 42

		case 2:
			ret_course.CID = "MSDS 2"
			ret_course.CNAME = "Zaphod's Heads have something to say"
			ret_course.CPREREQ = 1

		case 3:
			ret_course.CID = "MSDS 60000"
			ret_course.CNAME = "The Art of Flying: forgetting to fall"
			ret_course.CPREREQ = 59999

		case 4:
			ret_course.CID = "MSDS Infinity"
			ret_course.CNAME = "The Heart of Gold"
			ret_course.CPREREQ = 42
		
	}

	return ret_course
}


func main() {

	// as an array
	var course_array [5]MSDSCourse
	// as a slice
	var course_slice []MSDSCourse
	// as a map
	var course_map = make(map[int]MSDSCourse)
	for ii := 0; ii < 5; ii++ {
		temp_course := MSDSCourse_Returner(ii)
		
		course_array[ii] = temp_course
		course_slice = append(course_slice, temp_course)
		course_map[ii] = temp_course
	}
	
	// Ignoring empty fields in JSON
	t_0 := time.Now().UnixMilli()
	MSDSCourseVar, err := json.Marshal(&course_array)
	t_1 := time.Now().UnixMilli()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("\narray elapsed time: %v ms\n", t_1 - t_0 )
		fmt.Printf("Courses Marshalled: %s",MSDSCourseVar)
	}

	t_0 = time.Now().UnixMilli()
	MSDSCourseVar, err = json.Marshal(&course_slice)
	t_1 = time.Now().UnixMilli()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("\nslice elapsed time: %v ms\n", t_1 - t_0 )
		fmt.Printf("Courses Marshalled: %s",MSDSCourseVar)
	}

	t_0 = time.Now().UnixMilli()
	MSDSCourseVar, err = json.Marshal(&course_map)
	t_1 = time.Now().UnixMilli()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("\nmap elapsed time: %v ms\n", t_1 - t_0 )
		fmt.Printf("Courses Marshalled: %s", MSDSCourseVar)
	}
}
