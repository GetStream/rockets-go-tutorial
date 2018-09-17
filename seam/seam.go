package seam

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/esimov/caire"
	"github.com/pkg/errors"
)

func ContentAwareResize(url string) ([]byte, error) {
	fmt.Printf("Download starting for url %s", url)
	response, err := http.Get(url)
	defer response.Body.Close()

	if err != nil {
		return nil, errors.Wrap(err, "Failed to read the image")
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
		return nil, errors.Wrap(err, "Failed to apply seam carving to the image")
	}
	fmt.Printf("Seam carving completed for %s", url)

	return converted.Bytes(), nil
}
