package main

import (
	"fmt"
	"math/rand"
	"time"

	// "github.com/mactsouk/post05"
	"main.go/post05"
)

var MIN = 0
var MAX = 26

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func getString(length int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand))
		temp = temp + newChar
		if i == length {
			break
		}
		i++
	}
	return temp
}

func main() {
	post05.Hostname = "localhost"
	post05.Port = 5433
	post05.Username = "postgres"
	post05.Password = "root"
	post05.Database = "msds"


	// just for returning a random class name
	classname_list := []string{"Lies, Damn Lies, and Statistics", "GoLang, GoingLang, GoneLang", "Python: Going Full Monty", "Artificial Intelligence",
 							"Natural Intelligence", "MSDS RSQS VSNS PSLS", "Intro to probability: Don't tell me the odds!", "Why are all our course syllabi 3 years old?",
							"Data Arts and Data Sciences", "P-Hacking for Dummies", "Nobody likes R", "AwayFromDataScience.com", 
							"Programming without Stack Exchange", "Gradient Ascent: Always Be Climbing"}

	data, err := post05.ListClass()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range data {
		fmt.Println(v)
	}

	SEED := time.Now().Unix()
	rand.Seed(SEED)
	random_classid := fmt.Sprintf("MSDS%v", rand.Int()%1000) // random class id
	random_classname := classname_list[rand.Int()%len(classname_list)] // get a random classname from the list
	random_otherclass := fmt.Sprintf("MSDS%v", rand.Int()%1000)

	t := post05.MSDSCourse{
		CID:      random_classid,
		CNAME:    random_classname,
		CPREREQ:  random_otherclass}

	id := post05.AddClass(t)
	if id == "-1" {
		fmt.Println("There was an error adding class", t.CNAME)
	}

	// err = post05.DeleteClass()
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// // Trying to delete it again!
	// err = post05.DeleteClass(CID)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// id = post05.AddClass(t)
	// if id == -1 {
	// 	fmt.Println("There was an error adding class", t.CNAME)
	// }

}
