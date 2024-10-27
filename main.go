package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	// start the file reader to pass into the post handler
	fr := FileReader{}
	mux.HandleFunc("GET /posts/{slug}",PostHandler(fr))

	err := http.ListenAndServe(":3030", mux)
	if err != nil {
		log.Fatal(err)
	}
}


// Read the slug from the url
type SlugReader interface {
	Read(slug string) (string, error)
}

type FileReader struct{}

// open and read the file on the disc/source
func (fr FileReader) Read(slug string) (string, error) {
	f, err := os.Open(slug + ".md")
	if err != nil {
		return "", err
	}

	defer f.Close()
	b, err := io.ReadAll(f)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// handle returning of the file contents to the route
func PostHandler(sl SlugReader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		slug := r.PathValue("slug")

		postMarkdown, err := sl.Read(slug)

		if err != nil{
			// TODO: Handle different errors in the future
			http.Error(w, "Post not found", http.StatusNotFound)
		}

		fmt.Fprint(w, postMarkdown)
	}
}