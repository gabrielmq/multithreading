package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func main() {
	viaCepCh := make(chan string)
	apiCepCh := make(chan string)

	go getZipcodeFrom("http://viacep.com.br/ws/02861030/json/", viaCepCh)
	go getZipcodeFrom("https://cdn.apicep.com/file/apicep/02861-030.json", apiCepCh)

	// Fica esperando o channel que responder mais rápido para esvaziar ele
	select {
	case viacep := <-viaCepCh:
		fmt.Printf("ViaCep\n %v", viacep)

	case apiCep := <-apiCepCh:
		fmt.Printf("ApiCep\n %v", apiCep)

	// define um timeout para caso nenhum channel tenha sido preenchido em 1s
	case <-time.After(time.Second):
		fmt.Println("Timeout")
	}
}

func getZipcodeFrom(url string, ch chan<- string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	ch <- string(body)
}
