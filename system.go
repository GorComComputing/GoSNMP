package main

import (
	"fmt"
	"log"
	//"strconv"

	g "github.com/gosnmp/gosnmp"
)

func cmd_get(words []string) string {
	var output string

	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = "127.0.0.1"
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
	result, err2 := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
	}

	for i, variable := range result.Variables {
		output += fmt.Sprintf("%d: oid: %s ", i, variable.Name)

		//output += string(i) + ": oid: " + string(variable.Name) + " "

		// the Value of each variable returned by Get() implements
		// interface{}. You could do a type switch...
		switch variable.Type {
		case g.OctetString:
			output += fmt.Sprintf("string: %s\n", string(variable.Value.([]byte)))
			//output += "string: " + string(variable.Value.([]byte)) + "\n" 
			
		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			output += fmt.Sprintf("number: %d\n", g.ToBigInt(variable.Value))
			//output += "number: " + g.ToBigInt(variable.Value) + "\n" 
		}
	}
	return output
}
