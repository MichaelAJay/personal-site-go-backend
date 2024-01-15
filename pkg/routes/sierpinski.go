package routes

import (
	"fmt"
	"net/http"
	"strconv"

	triangle "github.com/MichaelAJay/personal-site-go-backend/pkg/sierpinski"
)

func SierpinskiHandler(w http.ResponseWriter, r *http.Request) {
	iterationsParam := r.URL.Query().Get("iterations")
	if iterationsParam == "" {
		http.Error(w, "Missing iterations parameter", http.StatusBadRequest)
		return
	}

	iterations, err := strconv.Atoi(iterationsParam)
	if err != nil {
		http.Error(w, "Invalid iterations parameter", http.StatusBadRequest)
		return
	}

	svg, err := triangle.GenerateSierpinskiSVG(iterations)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error generating SVG: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	fmt.Fprint(w, svg)
}
