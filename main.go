/*
  padbeam - Display etherpads on a beamer
  Copyright (C) 2016 padbeam authors and contributers

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	var padurl string
	var addr string
	flag.StringVar(&padurl, "u", "", "url of the pad")
	flag.StringVar(&addr, "b", ":8080", "the address to bind to")
	flag.Parse()

	// check if the pad url was specified
	if padurl == "" {
		fmt.Println("You have to specify a pad url. See --help for more information.")
		return
	}

	// forceReload enforces a location.reload()
	var forceReload bool = false

	// handle requests to /reload
	http.HandleFunc("/reload", func(w http.ResponseWriter, r *http.Request) {
		forceReload = true
		fmt.Fprint(w, "OK")
	})

	// handle requests to /pad
	http.HandleFunc("/pad", func(w http.ResponseWriter, r *http.Request) {
		if forceReload {
			forceReload = false
			fmt.Fprint(w, "RELOAD")
			return
		}

		resp, err := http.Get(padurl + "/export/txt")
		if err != nil {
			fmt.Fprint(w, "ERROR")
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		fmt.Fprint(w, string(body))
	})

	// handle all other requests
	http.Handle("/", http.FileServer(http.Dir("./html")))

	// start the server
	log.Println("Server is starting on \"" + addr + "\"")
	log.Println("Pad URL: \"" + padurl + "\"")
	log.Fatal(http.ListenAndServe(addr, nil))
}
