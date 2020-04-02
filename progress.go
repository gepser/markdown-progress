package progress

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"
)

// Data ... is the collectiion of inputs we need to fill our template
type Data struct {
	BackgroundColor string
	Percentage      int
	Progress        int
	PickedColor     string
}

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

// Progress ... Entrypoint of our Cloud Function
func Progress(w http.ResponseWriter, r *http.Request) {
	var id = fmt.Sprintf(path.Base(r.URL.Path))

	if x, err := strconv.Atoi(id); err == nil {
		data := Data{
			BackgroundColor: grey,
			Percentage:      x,
			Progress:        x - (x / 10),
			PickedColor:     pickColor(x),
		}

		tpl, err := template.ParseFiles("progress.html")
		if err != nil {
			log.Fatalln(err)
		}

		buf := new(bytes.Buffer)

		err = tpl.Execute(buf, data)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("The percentage is: %d\n", x)
		w.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprintf(w, buf.String())
	}
}
