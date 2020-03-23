package main

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
)

var grey = "#555"
var red = "#d9534f"
var yellow = "#f0ad4e"
var green = "#5cb85c"

func pickColor(percentage int) string {
	pickedColor := red

	if percentage >= 0 && percentage < 33 {
		pickedColor = red
	} else if percentage >= 33 && percentage < 70 {
		pickedColor = yellow
	} else {
		pickedColor = green
	}
	return pickedColor
}

func buildSVG(percentage int) string {
	progress := percentage - (percentage / 10)
	begining := `
    <html>
    <svg width="90.0" height="20" xmlns="http://www.w3.org/2000/svg">
      <linearGradient id="a" x2="0" y2="100%%">
        <stop offset="0" stop-color="#bbb" stop-opacity=".2"/>
        <stop offset="1" stop-opacity=".1"/>
      </linearGradient>
  `
	backgroundBar := fmt.Sprintf(`
    <rect rx="4" x="0" width="90.0" height="20" fill="%s"/>
    `, grey)

	percentageBar := fmt.Sprintf(`
    <rect rx="4" x="0" width="%d" height="20" fill="%s"/>
    `, progress, pickColor(percentage))

	next := `<rect rx="4" width="90.0" height="20" fill="url(#a)"/>
    <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
    <text x="45.0" y="14">
  `
	final := `%%
        </text>
      </g>
    </svg>
  </html>
  `

	return fmt.Sprintf("%s%s%s%s%d%s", begining, backgroundBar, percentageBar, next, percentage, final)
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var id = fmt.Sprintf(path.Base(r.URL.Path))

		if x, err := strconv.Atoi(id); err == nil {
			fmt.Printf("The percentage is: %d\n", x)
			fmt.Fprintf(w, buildSVG(x))
		}
	})

	http.ListenAndServe(":8080", nil)
}
