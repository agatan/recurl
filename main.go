package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

type options struct {
	address  string
	upstream string
}

func main() {
	var op options

	flag.StringVar(&op.address, "addr", ":8080", "listening address")
	flag.StringVar(&op.upstream, "upstream", "", "upstream server address [required]")

	flag.Parse()

	if op.upstream == "" {
		fmt.Fprintln(os.Stderr, "upstream is not given\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	upstreamURL, err := url.Parse(op.upstream)
	if err != nil {
		panic(err)
	}
	rec := NewRecorder(upstreamURL)

	go func() {
		if err := http.ListenAndServe(op.address, rec); err != nil {
			panic(err)
		}
	}()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt)
	<-ch

	if err := json.NewEncoder(os.Stdout).Encode(rec.Exchanges); err != nil {
		panic(err)
	}
}
