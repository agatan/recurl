package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/signal"
)

type options struct {
	address  string
	upstream string
	output   string
}

func main() {
	var op options

	flag.StringVar(&op.address, "addr", ":8080", "listening address")
	flag.StringVar(&op.upstream, "upstream", "", "upstream server address [required]")
	flag.StringVar(&op.output, "o", "", "output file (default: stdout)")

	flag.Parse()

	if op.upstream == "" {
		fmt.Fprintln(os.Stderr, "upstream is not given\nUsage:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	var w io.Writer
	if op.output == "" || op.output == "stdout" {
		w = os.Stdout
	} else {
		f, err := os.OpenFile(op.output, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			panic(err)
		}
		w = f
		defer f.Close()
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

	if err := json.NewEncoder(w).Encode(rec.Exchanges); err != nil {
		panic(err)
	}
}
