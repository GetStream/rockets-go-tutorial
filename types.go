package main

import (
	"fmt"
)

func main() {
	message := "hello world"
	friends := []string{"john", "amy"}
	friends = append(friends, "jack")
	populationByCity := map[string]int{"Boulder": 108090, "Amsterdam": 821752}
	populationByCity["Palo Alto"] = 67024
	fmt.Println(message, friends, populationByCity)

}