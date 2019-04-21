package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Tconfig struct {
	Config map[string]string
}

func main() {
	usageInfo := flag.Bool("h", false, "Prints usage Information")
	addURL := flag.String("a", "", "Add alias name")
	rmURL := flag.String("d", "", "Delete alias name")
	url := flag.String("u", "", "Redirect Url")
	port := flag.String("p", "", "Port numer")
	run := flag.Bool("r", false, "Use to run server")
	configure := flag.Bool("c", false, "Use to change config")
	list := flag.Bool("l", false, "List redirections")
	_, _, _, _, _ = *usageInfo, *addURL, *rmURL, *url, *port

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [<dir>]\nOptions are:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *usageInfo == true {
		flag.Usage()
	}

	f, err := os.OpenFile("config.txt", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	handleErr(err)
	b, err := ioutil.ReadAll(f)
	handleErr(err)

	content := string(b)
	current_config := strings.Split(content, " ")
	config_map := make(map[string]string)

	if *configure == true {
		if *addURL == "" && *url == "" {
			flag.Usage()
		}

		if len(current_config) >= 1 {
			for i, pair := range current_config {
				a := strings.Split(pair, ":")
				if *addURL == a[0] {
					current_config[i] = *addURL + ":" + *url
					output := strings.Join(current_config, " ")
					err = ioutil.WriteFile("config.txt", []byte(output), 0644)
					if err != nil {
						log.Fatalln(err)
					}
					f.Close()
					os.Exit(1)
				}
			}
			data := *addURL + ":" + *url + " "
			if _, err = f.WriteString(data); err != nil {
				panic(err)
			}

		} else {
			data := *addURL + ":" + *url + " "
			if _, err = f.WriteString(data); err != nil {
				panic(err)
			}
		}
	}

	if *rmURL != "" {
		if len(current_config) >= 1 {
			for _, pair := range current_config {
				a := strings.Split(pair, ":")
				if *rmURL == a[0] {
					current_config = remove(current_config, *rmURL)

					output := strings.Join(current_config, " ")
					err = ioutil.WriteFile("config.txt", []byte(output), 0644)
					if err != nil {
						log.Fatalln(err)
					}
					f.Close()
					os.Exit(1)
				}
			}
		}
	}

	if len(current_config) <= 1 {
		fmt.Println("No config.")
		os.Exit(1)
	}
	for _, pair := range current_config {
		a := strings.Split(pair, ":")
		config_map[a[0]] = a[1]
	}

	if *list == true {
		for _, pair := range current_config {
			fmt.Println(pair)
		}

	}

	if *run == true {
		if *port == "" {
			fmt.Println("Use command: \n ./rails run -p [port_number] \nto start server.")
			os.Exit(1)
		}
		flag.Parse()
		fmt.Println(*port)
		myConfig := &Tconfig{Config: config_map}
		fmt.Fprintf(os.Stdout, "Server is listening on http://localhost:%s \n", *port)
		http.HandleFunc("/", myConfig.handler)
		http_port := ":" + *port
		log.Fatal(http.ListenAndServe(http_port, nil))

	}

}

func (cf *Tconfig) handler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "This is the %s page!", cf.Config[r.URL.Path[1:]])
}

// func handleConfig() bool {
// 	return true
// }

func handleErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	fmt.Println()
	return s
}
