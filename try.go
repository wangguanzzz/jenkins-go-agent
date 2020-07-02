package main

import (
	"flag"
	"fmt"
)

func main() {

	name := flag.String("name", "Name", "The name of the tag to attach to the instance")
	value := flag.String("value", "JenkinsAgent", "The value of the tag to attach to the instance")
	t := flag.Int("tag", 123, "The value of the tag to attach to the instance")
	flag.Parse()
	fmt.Println(*name, *value, *t)

}
