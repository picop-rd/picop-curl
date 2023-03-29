package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hiroyaonoe/bcop-go/contrib/net/http/bcophttp"
	"github.com/hiroyaonoe/bcop-go/propagation"
	"github.com/hiroyaonoe/bcop-go/protocol/header"
)

func main() {
	envID := flag.String("env-id", "", "BCoP env-id")
	url := flag.String("url", "", "URL")
	data := flag.String("data", "", "HTTP POST data")
	method := flag.String("method", "GET", "HTTP method")

	flag.Parse()

	h := header.NewV1()
	h.Set(propagation.EnvIDHeader, *envID)
	ctx := context.Background()
	ctx = propagation.EnvID{}.Extract(ctx, propagation.NewBCoPCarrier(h))

	client := &http.Client{
		Transport: bcophttp.NewTransport(nil, propagation.EnvID{}),
	}
	req, err := http.NewRequestWithContext(ctx, *method, *url, strings.NewReader(*data))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)

	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		fmt.Println(sc.Text())
	}
	if err := sc.Err(); err != nil {
		log.Fatal(err)
	}
}
