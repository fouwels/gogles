package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kaelanfouwels/gogles/ioman"
)

const _bind = ":8080"

//Serve web application
func Serve(cherr chan<- error, io *ioman.IOMan) {

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {

		data := io.GetDataPacket()
		j := json.NewEncoder(w)
		j.Encode(data)
	})

	err := http.ListenAndServe(_bind, nil)
	if err != nil {
		cherr <- fmt.Errorf("Failed to listen and serve on bind %v: %w", _bind, err)
	}
}
