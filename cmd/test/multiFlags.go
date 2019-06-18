package main

import (
	"flag"
	"fmt"
)

// MultiFlagVar check for multiple settings for a flag
type MultiFlagVar struct {
	Values []string
}

func (f *MultiFlagVar) String() string {
	return fmt.Sprint(f.Values)
}

// Set a flag value
func (f *MultiFlagVar) Set(value string) error {
	f.Values = append(f.Values, value)
	return nil
}

func main() {
	var v MultiFlagVar
	flag.Var(&v, "testV", "multiple values build []string")
	flag.Parse()
	fmt.Printf("v = %s\n", v)
}
