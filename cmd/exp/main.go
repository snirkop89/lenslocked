package main

import (
	"html/template"
	"os"
)

type User struct {
	Name string
	Bio  string
	Age  int
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	u := User{
		Name: "Johnny",
		Bio:  "Engineer",
		Age:  123,
	}

	err = t.Execute(os.Stdout, u)
	if err != nil {
		panic(err)
	}
}
