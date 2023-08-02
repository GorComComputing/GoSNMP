package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	//"io"
	"net/http"
	
	"strings"
)



// /api handler
func http_pars(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)	// enable CORS
	
	// parameters from POST or GET
        r.ParseForm()
	words := []string{}

	for param, values := range r.Form {   	  // range over map
  		for _, value := range values {    // range over []string
     			if param == "cmd" {
				words = strings.Fields(value)
			} else {
				words = append(words, string(param) + "=" + string(value))
			}
  		}
	}
	
	
		isUnion := false
    		union := ""
    		var count []int
    		var length []int
    		var Unions []string
    		for i, val := range words {
    			if (val[0] == 39 && val[len(val)-1] != 39) || (val[0] == 39 && len(val) == 1) && isUnion != true {
    				isUnion = true
    				union += val[1:] + " "
    				count = append(count, i)
    				continue
    			}
    			if val[len(val)-1] != 39 && isUnion == true {
    				union += val + " "
    				continue
    			}
    			if val[len(val)-1] == 39 {
    				isUnion = false
    				union += val[:len(val)-1]
    				length = append(length, i - count[len(count)-1])
    				Unions = append(Unions, union)
    				union = ""
    				continue
    			}
    		}
    		
    		
    		x := 0
    		for i, val := range count {
    			words[val+x] = Unions[i] 
    			copy(words[val+x+1:], words[val+length[i]+1:])
    			x -= length[i]
		}
		words = words[:len(words)+x]
		
	
	out := interpretator(words)
	if len(out) > 0 {
		fmt.Fprintf(w, out)
	}
}


// Event handler
func json_pars(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)	// enable CORS
	
	// Parsing JSON
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(body))
	var req map[string]string
	err = json.Unmarshal(body, &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(req)
	words := []string{}
	if req["cmd"] != "" {
		fmt.Println("CMD")
		words = strings.Fields(req["cmd"])
		for param, value := range req {    // range over []string
     			if param != "cmd" {
				words = append(words, string(param) + "=" + string(value))
			}
  		}
	} else {
		fmt.Println("not CMD")
		for param, value := range req {    // range over []string
				words = append(words, string(param) + "=" + string(value))
  		}
	}
	
  	//fmt.Println(words)
  	
  	isUnion := false
    		union := ""
    		var count []int
    		var length []int
    		var Unions []string
    		for i, val := range words {
    			if (val[0] == 39 && val[len(val)-1] != 39) || (val[0] == 39 && len(val) == 1) && isUnion != true {
    				isUnion = true
    				union += val[1:] + " "
    				count = append(count, i)
    				continue
    			}
    			if val[len(val)-1] != 39 && isUnion == true {
    				union += val + " "
    				continue
    			}
    			if val[len(val)-1] == 39 {
    				isUnion = false
    				union += val[:len(val)-1]
    				length = append(length, i - count[len(count)-1])
    				Unions = append(Unions, union)
    				union = ""
    				continue
    			}
    		}
    		
    		
    		x := 0
    		for i, val := range count {
    			words[val+x] = Unions[i] 
    			copy(words[val+x+1:], words[val+length[i]+1:])
    			x -= length[i]
		}
		words = words[:len(words)+x]
		
	out := interpretator(words)
	if len(out) > 0 {
		fmt.Fprintf(w, out)
	}
}


// Enable CORS
func enableCors(w *http.ResponseWriter) {
        (*w).Header().Set("Access-Control-Allow-Origin", "*")
}




