package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/GetStream/rockets/seam"
	"github.com/GetStream/rockets/unsplash"
	"github.com/flosch/pongo2"
)
import b64 "encoding/base64"

const (
	url = "https://www.outerplaces.com/media/k2/items/cache/aa810e73f519481c15cdf02790d21ac8_L.jpg"
)

type Task struct {
	Position int
	URL      string
}

type TaskResult struct {
	Position int
	Resized  []byte
}

var spacexTemplate = pongo2.Must(pongo2.FromFile("spacex.html"))

func main() {
	fmt.Println("Cooling the engines, checkout http://localhost:3000/occupymars")

	http.HandleFunc("/occupymars", func(w http.ResponseWriter, r *http.Request) {
		resized, err := seam.ContentAwareResize(url)
		if err != nil {
			fmt.Errorf("things broke, %s", err)
		}

		w.Header().Set("Content-Type", "image/jpeg")
		io.Copy(w, bytes.NewReader(resized))

	})

	http.HandleFunc("/spacex", func(w http.ResponseWriter, r *http.Request) {
		response, err := unsplash.LoadRockets()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		err = spacexTemplate.ExecuteWriter(pongo2.Context{"response": response}, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	http.HandleFunc("/spacex_seams", func(w http.ResponseWriter, r *http.Request) {
		response, err := unsplash.LoadRockets()

		results := make(chan *TaskResult)
		urlsChannel := make(chan *Task, 9)
		for i, r := range response.Results[:9] {
			urlsChannel <- &Task{i, r.URLs["small"]}
		}

		for w := 1; w <= 3; w++ {
			go worker(w, urlsChannel, results)
		}

		close(urlsChannel)

		for a := 1; a <= 9; a++ {
			taskResult := <-results
			sEnc := b64.StdEncoding.EncodeToString(taskResult.Resized)
			response.Results[taskResult.Position].Resized = sEnc
		}

		err = spacexTemplate.ExecuteWriter(pongo2.Context{"response": response}, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	log.Fatal(http.ListenAndServe(":3000", nil))

}

func worker(id int, jobs <-chan *Task, results chan<- *TaskResult) {
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		resized, _ := seam.ContentAwareResize(j.URL)

		results <- &TaskResult{j.Position, resized}
	}
}
