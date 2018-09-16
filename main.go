package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/GetStream/rockets-go-tutorial/seam"
	"github.com/GetStream/rockets-go-tutorial/unsplash"
	"github.com/flosch/pongo2"
)
import b64 "encoding/base64"

const (
	IMAGE_URL   string = "https://bit.ly/2N8Ra4q"
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
	fmt.Println("Ready for liftoff! Checkout \n http://localhost:3000/occupymars \n http://localhost:3000/spacex \n http://localhost:3000/spacex_seams")

	http.HandleFunc("/occupymars", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("resize") > "" {
			resized, err := seam.ContentAwareResize(IMAGE_URL)
			if err != nil {
				fmt.Errorf("things broke, %s", err)
			}

			w.Header().Set("Content-Type", "image/jpeg")
			io.Copy(w, bytes.NewReader(resized))
		} else {
			fmt.Fprintf(w, "<html><div>Original image:</div> <img src=\"%s\" /><br/><a href=\"?resize=1\">Resize using Seam Carving</a></html>", IMAGE_URL)
		}
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

		resultChannel := make(chan *TaskResult)
		taskChannel := make(chan *Task, 8)
		for i, r := range response.Results[:8] {
			taskChannel <- &Task{i, r.URLs["small"]}
		}

		for w := 1; w <= 4; w++ {
			go worker(w, taskChannel, resultChannel)
		}

		close(taskChannel)

		for a := 1; a <= 8; a++ {
			taskResult := <-resultChannel
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

func worker(id int, taskChannel <-chan *Task, resultChannel chan<- *TaskResult) {
	for j := range taskChannel {
		fmt.Println("worker", id, "started  job", j)
		resized, _ := seam.ContentAwareResize(j.URL)

		resultChannel <- &TaskResult{j.Position, resized}
	}
}
