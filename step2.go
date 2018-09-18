package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/esimov/caire"
	"github.com/pkg/errors"
)

const (
	IMAGE_URL string = "https://bit.ly/2QGPDkr"
)

func ContentAwareResize(url string) ([]byte, error) {
	fmt.Printf("Download starting for url %s\n", url)
	response, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to read the image")
	}
	defer response.Body.Close()

	converted := &bytes.Buffer{}
	fmt.Printf("Download complete %s", url)

	shrinkFactor := 30
	fmt.Printf("Resize in progress %s, shrinking width by %d percent...\n", url, shrinkFactor)
	p := &caire.Processor{
		NewWidth:   shrinkFactor,
		Percentage: true,
	}

	err = p.Process(response.Body, converted)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to apply seam carving to the image")
	}
	fmt.Printf("Seam carving completed for %s\n", url)

	return converted.Bytes(), nil
}

func main() {
	fmt.Println("Ready for liftoff! Checkout http://localhost:3000/occupymars")

	http.HandleFunc("/occupymars", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("resize") > "" {
			resized, err := ContentAwareResize(IMAGE_URL)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("Content-Type", "image/jpeg")
			io.Copy(w, bytes.NewReader(resized))
		} else {
			fmt.Fprintf(w, "<html><div>Original image:</div> <img src=\"%s\" /><br/><a href=\"?resize=1\">Resize using Seam Carving</a></html>", IMAGE_URL)
		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
