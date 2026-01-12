package main

import (
	"fmt"
	"time"
	"unsafe"
)

func main() {
	var t time.Time
	var s string
	var i int

	fmt.Printf("time.Time size: %d bytes\n", unsafe.Sizeof(t))
	fmt.Printf("string size: %d bytes\n", unsafe.Sizeof(s))
	fmt.Printf("int size: %d bytes\n", unsafe.Sizeof(i))

	// Let's also check the struct before and after optimization
	type Before struct {
		CompanyName       string
		FullName          string
		Email             string
		Phone             string
		Status            string
		LicenceIssueDate  time.Time
		LicenceExpiryDate time.Time
		CreatedAt         time.Time
		ID                int
	}

	type After struct {
		LicenceIssueDate  time.Time
		LicenceExpiryDate time.Time
		CreatedAt         time.Time
		CompanyName       string
		FullName          string
		Email             string
		Phone             string
		Status            string
		ID                int
	}

	var before Before
	var after After

	fmt.Printf("\nStruct sizes:\n")
	fmt.Printf("Before optimization: %d bytes\n", unsafe.Sizeof(before))
	fmt.Printf("After optimization: %d bytes\n", unsafe.Sizeof(after))
}
