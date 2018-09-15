package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/esimov/caire"
)

const (
	IMAGE_URL string = "https://bit.ly/2N8Ra4q"
)

func ContentAwareResize(url string) ([]byte, error) {
	fmt.Printf("Download starting for url %s", url)
	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return nil, err
	}

	converted := &bytes.Buffer{}
	fmt.Printf("Download complete %s", url)

	shrinkFactor := 30
	fmt.Printf("Resize in progress %s, shrinking width by %d percent...", url, shrinkFactor)
	p := &caire.Processor{
		NewWidth:   shrinkFactor,
		Percentage: true,
	}

	err = p.Process(response.Body, converted)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Seam carving completed for %s", url)

	return converted.Bytes(), nil
}

func main() {
	fmt.Println("Ready for liftoff! Checkout http://localhost:3000/occupymars")

	http.HandleFunc("/occupymars", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("resize") > "" {
			resized, err := ContentAwareResize(IMAGE_URL)
			if err != nil {
				fmt.Errorf("things broke, %s", err)
				return
			}

			w.Header().Set("Content-Type", "image/jpeg")
			io.Copy(w, bytes.NewReader(resized))
		} else {
			fmt.Fprintf(w, "<html><div>Original image:</div> <img src=\"%s\" /><br/><a href=\"?resize=1\">Resize using Seam Carving</a></html>", IMAGE_URL)
		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
