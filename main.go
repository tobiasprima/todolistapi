package main

import (
	"fmt"
	"log"
)

func main(){
	log.Println("test")
}

func sayHello(name string) string {
	return fmt.Sprintf("Hello %s", name)
}