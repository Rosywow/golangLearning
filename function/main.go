package main

import "fmt"

type Data struct {
	Username string
	Password string
}

func (d Data) TestSame(){
	d.Username = "Rosy"
	d.Password = "None"
}

func (d *Data) TestChange() {
	d.Username = "Rosy"
	d.Password = "None"
}



func main() {
	data := Data{"Tom","123456"}
	fmt.Println(data)

	data.TestSame()
	fmt.Println(data)

	data.TestChange()
	fmt.Println(data)
}

