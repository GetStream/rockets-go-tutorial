package seam

import (
	"bytes"
	"fmt"
	"github.com/esimov/caire"
	"net/http"
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
