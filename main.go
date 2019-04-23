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

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [<dir>]\nOptions are:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *usageInfo == true {
		flag.Usage()
	}

	f, err := os.OpenFile("config.yml", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	handleErr(err)
	b, err := ioutil.ReadAll(f)
	handleErr(err)

	content := string(b)

	if *configure == true {
		if *addURL == "" && *url == "" {
			flag.Usage()
		}

		if len(content) <= 0 {
			data := *addURL + ":" + *url + " "
			_, err = f.WriteString(data)
			handleErr(err)
			f.Close()
			os.Exit(1)
		}

		configLines := strings.Split(content, "\n")
		for i, line := range configLines {
			a := strings.Split(line, ":")
			if *addURL == a[0] {
				configLines[i] = *addURL + ":" + *url
				output := strings.Join(configLines, "\n")
				err = ioutil.WriteFile("config.yml", []byte(output), 0644)
				handleErr(err)
				f.Close()
				os.Exit(1)
			}
		}
		data := *addURL + ":" + *url
		configLines = append(configLines, data)
		output := strings.Join(configLines, "\n")
		err = ioutil.WriteFile("config.yml", []byte(output), 0644)
		handleErr(err)
		f.Close()
		os.Exit(1)
	}

	if *rmURL != "" {
		if len(content) <= 0 {
			fmt.Println("Config file is empty.")
			f.Close()
			os.Exit(1)
		}

		configLines := strings.Split(content, "\n")
		for i, line := range configLines {
			a := strings.Split(line, ":")
			if *rmURL == a[0] {
				configLines = append(configLines[:i], configLines[i+1:]...)
				output := strings.Join(configLines, "\n")
				err = ioutil.WriteFile("config.yml", []byte(output), 0644)
				handleErr(err)
				f.Close()
				os.Exit(1)
			}
		}

		fmt.Println("Can't not find .", *rmURL)
		f.Close()
		os.Exit(1)

	}

	if *list == true {
		configLines := strings.Split(content, "\n")
		fmt.Println("List redirections:")
		for _, line := range configLines {
			a := strings.Split(line, ":")
			fmt.Println("/", a[0], "-->", a[1])
		}
		os.Exit(1)
	}

	if *run == true {
		if *port == "" {
			fmt.Println("Use command: \n ./rails run -p [port_number] \nto start server.")
			os.Exit(1)
		}
		configLines := strings.Split(content, "\n")
		configMap := make(map[string]string)
		for _, line := range configLines {
			a := strings.Split(line, ":")
			configMap[a[0]] = a[1]
		}
		myConfig := &Tconfig{Config: configMap}
		fmt.Fprintf(os.Stdout, "Server is listening on http://localhost:%s \n", *port)
		http.HandleFunc("/", myConfig.handler)
		httpPort := ":" + *port
		log.Fatal(http.ListenAndServe(httpPort, nil))

	}

}

func (cf *Tconfig) handler(w http.ResponseWriter, r *http.Request) {
	redirectUrl := cf.Config[r.URL.Path[1:]]
	fmt.Println(redirectUrl)
	http.Redirect(w, r, "http://"+redirectUrl, http.StatusMovedPermanently)
	return
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
