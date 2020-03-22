package main

import (
	"fmt"
	"net/http"
)

var progress = "90"
var percentage = "100"
var grey = "#555"
var yellow = "#f0ad4e"
var green = "#5cb85c"

var begining = `
<html>
<svg width="90.0" height="20" xmlns="http://www.w3.org/2000/svg">
  <linearGradient id="a" x2="0" y2="100%%">
    <stop offset="0" stop-color="#bbb" stop-opacity=".2"/>
    <stop offset="1" stop-opacity=".1"/>
  </linearGradient>
`
var backgroundBar = fmt.Sprintf(`<rect rx="4" x="0" width="90.0" height="20" fill="%s"/>`, grey)

var percentageBar = fmt.Sprintf(`<rect rx="4" x="0" width="%s" height="20" fill="%s"/>`, progress, green)
var next = `<rect rx="4" width="90.0" height="20" fill="url(#a)"/>
  <g fill="#fff" text-anchor="middle" font-family="DejaVu Sans,Verdana,Geneva,sans-serif" font-size="11">
  <text x="45.0" y="14">
`
var final = `%%
    </text>
  </g>
</svg>
</html>
`

var stringfinal = fmt.Sprintf("%s%s%s%s%s%s", begining, backgroundBar, percentageBar, next, percentage, final)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
		fmt.Fprintf(w, stringfinal)
	})

	http.ListenAndServe(":8080", nil)
}
