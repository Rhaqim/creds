package main

import (
	"fmt"
	"reflect"
)

type User struct {
	ID       int
	Username string
	Password string
}

func main() {
	user := User{
		ID:       1,
		Username: "admin",
		Password: "admin",
	}

	// Get the type of the user variable
	userType := reflect.TypeOf(user)
	fmt.Println("Type of user:", userType)

	// Get the value of the user variable
	userValue := reflect.ValueOf(user)
	fmt.Println("Value of user:", userValue)

	// Get the type of the user variable
	userType = userValue.Type()
	fmt.Println("Type of user:", userType)

	// Get the kind of the user variable
	userKind := userValue.Kind()
	fmt.Println("Kind of user:", userKind)

	// Get the number of fields in the user variable
	numFields := userType.NumField()
	fmt.Println("Number of fields in user:", numFields)

	// Get the field at index 0
	field0 := userType.Field(0)
	fmt.Println("Field at index 0:", field0)

	// Get the field by name
	fieldUsername, _ := userType.FieldByName("Username")
	fmt.Println("Field by name:", fieldUsername)

}
