package templviews

import (
	"net/http"

	"github.com/a-h/templ"
)

func Handler() http.Handler {
	// Create a new file server handler
	comp := hello("vvv")
	return templ.Handler(comp)
	// Return the handler
}