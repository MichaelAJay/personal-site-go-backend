package routes

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/MichaelAJay/personal-site-go-backend/pkg/custom_utility"
	triangle "github.com/MichaelAJay/personal-site-go-backend/pkg/sierpinski"
	"github.com/gin-gonic/gin"
)

func SierpinskiHandler(c *gin.Context) {
	iterationsParam := c.Query("iterations")
	if iterationsParam == "" {
		c.String(http.StatusBadRequest, "Missing iterations parameter")
		return
	}

	iterations, err := strconv.Atoi(iterationsParam)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid iterations parameter")
		return
	}

	colorString := c.Query("color")
	if colorString == "" {
		colorString = "#000000"
	}
	color, err := custom_utility.ParseSixDigitHexCodeToColor(colorString)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	svg, err := triangle.GenerateSierpinskiSVG(iterations, color)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error generating SVG: %v", err))
		return
	}

	c.Header("Content-Type", "image/svg+xml")
	c.String(http.StatusOK, svg)
}
