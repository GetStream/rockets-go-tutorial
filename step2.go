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
	IMAGE_URL   string = "https://www.outerplaces.com/media/k2/items/cache/aa810e73f519481c15cdf02790d21ac8_L.jpg"
)

func ContentAwareResize(url string) ([]byte, error) {
	fmt.Printf("Download starting for url %s", url)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	converted := &bytes.Buffer{}
	fmt.Printf("Download complete %s", url)

	shrinkFactor := 30
	fmt.Printf("Resize in progress %s, shrinking width by %d percent...", url, shrinkFactor)
	p := &caire.Processor{
		NewWidth:       shrinkFactor,
		Percentage:     true,
	}

	err = p.Process(response.Body, converted)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Seam carving completed for %s", url)

	return converted.Bytes(), nil
}

func main() {
	fmt.Println("Ready for liftoff! checkout http://localhost:3000/occupymars")

	http.HandleFunc("/occupymars", func(w http.ResponseWriter, r *http.Request) {
		resized, _ := ContentAwareResize(IMAGE_URL)

		w.Header().Set("Content-Type", "image/jpeg")
		io.Copy(w, bytes.NewReader(resized))

	})

	log.Fatal(http.ListenAndServe(":3000", nil))

}
