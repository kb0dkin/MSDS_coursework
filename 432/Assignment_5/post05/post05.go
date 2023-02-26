package post05

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
)

// Connection details
var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)

// Userdata is for holding full user data
// Userdata table + Username
type Userdata struct {
	ID          int
	Username    string
	Name        string
	Surname     string
	Description string
}

type MSDSCourse struct {
	CID string `json:"course_ID`
	CNAME string `json:"course_name"`
	CPREREQ string `json:"prerequisite"`
}

func openConnection() (*sql.DB, error) {
	// connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)

	// open database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}



// The function returns the User ID of the username
// -1 if the user does not exist
func exists(username string) int {
	username = strings.ToLower(username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "users" where username = '%s'`, username)
	rows, err := db.Query(statement)

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		userID = id
	}
	defer rows.Close()
	return userID
}

// AddUser adds a new user to the database
// Returns new User ID
// -1 if there was an error
func AddUser(d Userdata) int {
	d.Username = strings.ToLower(d.Username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("User already exists:", Username)
		return -1
	}

	insertStatement := `insert into "users" ("username") values ($1)`
	_, err = db.Exec(insertStatement, d.Username)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	userID = exists(d.Username)
	if userID == -1 {
		return userID
	}

	insertStatement = `insert into "userdata" ("userid", "name", "surname", "description")
	values ($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}

	return userID
}

// DeleteUser deletes an existing user
func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Does the ID exist?
	statement := fmt.Sprintf(`SELECT "username" FROM "users" where id = %d`, id)
	rows, err := db.Query(statement)

	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	defer rows.Close()

	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}

	// Delete from Userdata
	deleteStatement := `delete from "userdata" where userid=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	// Delete from Users
	deleteStatement = `delete from "users" where id=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	return nil
}

// ListUsers lists all users in the database
func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "id","username","name","surname","description"
		FROM "users","userdata"
		WHERE users.id = userdata.userid`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		temp := Userdata{ID: id, Username: username, Name: name, Surname: surname, Description: description}
		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}

// UpdateUser is for updating an existing user
func UpdateUser(d Userdata) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID == -1 {
		return errors.New("User does not exist")
	}
	d.ID = userID
	updateStatement := `update "userdata" set "name"=$1, "surname"=$2, "description"=$3 where "userid"=$4`
	_, err = db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
	if err != nil {
		return err
	}

	return nil
}


// ###################################################################################

// The function returns the User ID of the username
// -1 if the user does not exist
func class_exists(classname string) string {
	classname = strings.ToLower(classname)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return "-1"
	}
	defer db.Close()

	CID := "-1"
	statement := fmt.Sprintf(`SELECT "cid" FROM "msdscoursecatalog" where cname = '%s'`, classname)
	rows, err := db.Query(statement)

	for rows.Next() {
		var id string
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return "-1"
		}
		CID = id
	}
	defer rows.Close()
	return CID
}

// AddUser adds a new user to the database
// Returns new User ID
// -1 if there was an error
func AddClass(d MSDSCourse) string {
	d.CNAME = strings.ToLower(d.CNAME)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return "-1"
	}
	defer db.Close()

	CID := class_exists(d.CNAME)
	if CID != "-1" {
		fmt.Println("Class already exists:", d.CNAME)
		return "-1"
	}

	insertStatement := `insert into "msdscoursecatalog" ("cid", "cname", "cprereq")
	values ($1, $2, $3)`
	_, err = db.Exec(insertStatement, d.CID, d.CNAME, d.CPREREQ)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return "-1"
	}

	fmt.Println("Inserted",d.CID)

	return d.CID
}

// DeleteUser deletes an existing user
func DeleteClass(cid string) error {
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Does the ID exist?
	statement := fmt.Sprintf(`SELECT "cname" FROM "msdscoursecatalog" where cid = %d`, cid)
	rows, err := db.Query(statement)

	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}
	defer rows.Close()

	if class_exists(username) != cid {
		return fmt.Errorf("Class with CID %d does not exist", cid)
	}

	// Delete from msdscoursecatalog
	deleteStatement := `delete from "msdscoursecatalog" where cid=$1`
	_, err = db.Exec(deleteStatement, cid)
	if err != nil {
		return err
	}


	return nil
}

// ListClass lists all users in the database
func ListClass() ([]MSDSCourse, error) {
	Data := []MSDSCourse{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "cid","cname","cprereq"
		FROM "msdscoursecatalog"`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var CID string
		var CNAME string
		var CPREREQ string
		err = rows.Scan(&CID, &CNAME, &CPREREQ)
		temp := MSDSCourse{CID: CID, CNAME: CNAME, CPREREQ:CPREREQ}
		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}

// // UpdateUser is for updating an existing user
// func UpdateClass(d Userdata) error {
// 	db, err := openConnection()
// 	if err != nil {
// 		return err
// 	}
// 	defer db.Close()

// 	userID := exists(d.Username)
// 	if userID == -1 {
// 		return errors.New("User does not exist")
// 	}
// 	d.ID = userID
// 	updateStatement := `update "userdata" set "name"=$1, "surname"=$2, "description"=$3 where "userid"=$4`
// 	_, err = db.Exec(updateStatement, d.Name, d.Surname, d.Description, d.ID)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }



