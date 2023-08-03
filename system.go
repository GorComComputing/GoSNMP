package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
	//"strings"
	"flag"
	"path/filepath"
	"net"
	
	"net/http"
    	"io/ioutil"
	"bytes"
    	"encoding/json"

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


func cmd_get_param(words []string) string {
	var output string

	// get Target and Port from environment
	envTarget := words[1]
	envPort := words[2]
	if len(envTarget) <= 0 {
		log.Fatalf("not set: IP address")
	}
	if len(envPort) <= 0 {
		log.Fatalf("not set: PORT")
	}
	port, _ := strconv.ParseUint(envPort, 10, 16)

	// Build our own GoSNMP struct, rather than using g.Default.
	// Do verbose logging of packets.
	params := &g.GoSNMP{
		Target:    envTarget,
		Port:      uint16(port),
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Logger:    g.NewLogger(log.New(os.Stdout, "", 0)),
	}
	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	// Function handles for collecting metrics on query latencies.
	var sent time.Time
	params.OnSent = func(x *g.GoSNMP) {
		sent = time.Now()
	}
	params.OnRecv = func(x *g.GoSNMP) {
		log.Println("Query latency in seconds:", time.Since(sent).Seconds())
	}

	oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
	result, err2 := params.Get(oids) // Get() accepts up to g.MAX_OIDS
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



func cmd_get_v3(words []string) string {
	var output string

	// build our own GoSNMP struct, rather than using g.Default
	params := &g.GoSNMP{
		Target:        "127.0.0.1",
		Port:          161,
		Version:       g.Version3,
		SecurityModel: g.UserSecurityModel,
		MsgFlags:      g.AuthPriv,
		Timeout:       time.Duration(30) * time.Second,
		SecurityParameters: &g.UsmSecurityParameters{UserName: "user",
			AuthenticationProtocol:   g.SHA,
			AuthenticationPassphrase: "password",
			PrivacyProtocol:          g.DES,
			PrivacyPassphrase:        "password",
		},
	}
	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
	result, err2 := params.Get(oids) // Get() accepts up to g.MAX_OIDS
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


const (
	On  int = 1
	Off     = 2
)

func cmd_set(words []string) string {
	var output string
	
	var Client = &g.GoSNMP{
		Target:    "127.0.0.1",
		Port:      161,
		Community: "public",
		Version:   g.Version2c,
		Timeout:   time.Duration(2) * time.Second,
		Logger:    g.NewLogger(log.New(os.Stdout, "", 0)),
	}
	err := Client.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer Client.Conn.Close()
	var mySnmpPDU = []g.SnmpPDU{{
		Name:  "1.3.6.1.2.1.25.1.7.0", //"1.3.6.1.4.1.318.1.1.4.4.2.1.3.15",
		Type:  g.Integer,
		Value: Off,
	}}
	setResult, setErr := Client.Set(mySnmpPDU)
	if setErr != nil {
		log.Fatalf("SNMP set() fialed due to err: %v", setErr)
	}
	
	for i, variable := range setResult.Variables {
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


func cmd_get_hex(words []string) string {
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
		switch variable.Type {
		case g.OctetString:
			value := variable.Value.([]byte)
			//if strings.Contains(strconv.Quote(string(value)), "\\x") {
				tmp := ""
                		for i := 0; i < len(value); i++ {
					tmp += fmt.Sprintf("%v", value[i])
					if i != (len(value) - 1) {
						tmp += " "
					}
				}
				//fmt.Printf("Hex-String: %s\n", tmp)
				output += fmt.Sprintf("Hex-String: %s\n", tmp)
			//} else {
				output += fmt.Sprintf("string: %s\n", string(variable.Value.([]byte)))
				//fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
			//}
		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			output += fmt.Sprintf("number: %d\n", g.ToBigInt(variable.Value))
			//fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}
	return output
}



func cmd_walk(words []string) string {
	var output string
	
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("   %s [-community=<community>] host [oid]\n", filepath.Base(os.Args[0]))
		fmt.Printf("     host      - the host to walk/scan\n")
		fmt.Printf("     oid       - the MIB/Oid defining a subtree of values\n\n")
		flag.PrintDefaults()
	}

	var community string
	flag.StringVar(&community, "community", "public", "the community string for device")

	flag.Parse()

	/*if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(1)
	}*/
	target := "127.0.0.1"
	var oid string
	if len(words) > 1 {
		oid = words[1]
	}

	g.Default.Target = target
	g.Default.Community = community
	g.Default.Timeout = time.Duration(10 * time.Second) // Timeout better suited to walking
	err := g.Default.Connect()
	if err != nil {
		fmt.Printf("Connect err: %v\n", err)
		os.Exit(1)
	}
	defer g.Default.Conn.Close()

	err = g.Default.BulkWalk(oid, printValue)
	if err != nil {
		fmt.Printf("Walk Error: %v\n", err)
		os.Exit(1)
	}
	return output
}

func printValue(pdu g.SnmpPDU) error {
	fmt.Printf("%s = ", pdu.Name)

	switch pdu.Type {
	case g.OctetString:
		b := pdu.Value.([]byte)
		fmt.Printf("STRING: %s\n", string(b))
	default:
		fmt.Printf("TYPE %d: %d\n", pdu.Type, g.ToBigInt(pdu.Value))
	}
	return nil
}


func cmd_trap_v1(words []string) string {
	var output string

	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = "127.0.0.1"
	g.Default.Port = 9162
	g.Default.Version = g.Version1

	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	pdu := g.SnmpPDU{
		Name:  "1.3.6.1.2.1.1.6",
		Type:  g.OctetString,
		Value: "Oval Office",
	}

	trap := g.SnmpTrap{
		Variables:    []g.SnmpPDU{pdu},
		Enterprise:   ".1.3.6.1.6.3.1.1.5.1",
		AgentAddress: "127.0.0.1",
		GenericTrap:  0,
		SpecificTrap: 0,
		Timestamp:    300,
	}

	_, err = g.Default.SendTrap(trap)
	if err != nil {
		log.Fatalf("SendTrap() err: %v", err)
	}
	return output
}


func cmd_trap_v2(words []string) string {
	var output string

	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc
	g.Default.Target = "127.0.0.1"
	g.Default.Port = 9162
	g.Default.Version = g.Version2c
	g.Default.Community = "public"
	g.Default.Logger = g.NewLogger(log.New(os.Stdout, "", 0))
	err := g.Default.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer g.Default.Conn.Close()

	pdu := g.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0",
		Type:  g.OctetString,//g.ObjectIdentifier,
		Value: words[1],
	}

	trap := g.SnmpTrap{
		Variables: []g.SnmpPDU{pdu},
	}

	_, err = g.Default.SendTrap(trap)
	if err != nil {
		log.Fatalf("SendTrap() err: %v", err)
	}
	return output
}



func cmd_trap_v3(words []string) string {
	var output string

	// Default is a pointer to a GoSNMP struct that contains sensible defaults
	// eg port 161, community public, etc

	params := &g.GoSNMP{
		Target:        "127.0.0.1",
		Port:          9162,
		Version:       g.Version3,
		Timeout:       time.Duration(30) * time.Second,
		SecurityModel: g.UserSecurityModel,
		MsgFlags:      g.AuthPriv,
		Logger:        g.NewLogger(log.New(os.Stdout, "", 0)),
		SecurityParameters: &g.UsmSecurityParameters{UserName: "user",
			AuthoritativeEngineID:    "1234",
			AuthenticationProtocol:   g.SHA,
			AuthenticationPassphrase: "password",
			PrivacyProtocol:          g.DES,
			PrivacyPassphrase:        "password",
		},
	}
	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	pdu := g.SnmpPDU{
		Name:  ".1.3.6.1.6.3.1.1.4.1.0",
		Type:  g.ObjectIdentifier,
		Value: ".1.3.6.1.6.3.1.1.5.1",
	}

	trap := g.SnmpTrap{
		Variables: []g.SnmpPDU{pdu},
	}

	_, err = params.SendTrap(trap)
	if err != nil {
		log.Fatalf("SendTrap() err: %v", err)
	}
	return output
}


func cmd_trap_srv(words []string) string {
	var output string
	
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		fmt.Printf("   %s\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}

	tl := g.NewTrapListener()
	tl.OnNewTrap = myTrapHandler
	tl.Params = g.Default
	tl.Params.Logger = g.NewLogger(log.New(os.Stdout, "", 0))

	err := tl.Listen("0.0.0.0:9162")
	if err != nil {
		log.Panicf("error in listen: %s", err)
	}
	return output
}

func myTrapHandler(packet *g.SnmpPacket, addr *net.UDPAddr) {
	log.Printf("got trapdata from %s\n", addr.IP)
	for _, v := range packet.Variables {
		switch v.Type {
		case g.OctetString:
			b := v.Value.([]byte)
			fmt.Printf("OID: %s, string: %x\n", v.Name, b)
			
			var words = make([]string, 0)
			words = append(words, string(b))
			curl(words)

		default:
			log.Printf("trap: %+v\n", v)
		}
	}
}


type Payload struct {
	Cmd string `json:"cmd"`
	Level string `json:"level"`
	Ident string `json:"ident"`
	Is_check string `json:"is_check"`
	Object string `json:"object"`
	Source string `json:"source"`
	Body string `json:"body"`
}



// полная функция curl (возвращает string и map)
func curl(words []string) (string, map[string]any) {
	var output string
	var result map[string]any
	
	//data := `{"level":"info_trap","ident":"1234","is_check":"true","object":"88","source":"systemS","body":"body1", "cmd":"ins_evnt"}`
	data := Payload{"ins_evnt", "Trap", words[0], "true", "88", "systemS", "body1"} //"getconfig"
	
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		// handle err
	}
	
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "http://localhost:8085/json", body) //"http://192.168.1.206/cgi-bin/configs.cgi?"
	if err != nil {
		output = "Request FAIL\n"
		return output, nil
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		output = "Request FAIL\n"
		return output, nil
	}
	defer resp.Body.Close()

	body_resp, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		output = "Request FAIL\n"
		return output, nil
	}
	
	json.Unmarshal(body_resp, &result)
	
	/*for key, val := range result {
		str := fmt.Sprintf("%v", val)
		output += string(key) + ": " + str + "\n"
	}*/
	
	output = string(body_resp)

	return output, result 
}
