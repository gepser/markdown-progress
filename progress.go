package progress

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path"
	"strconv"
	"text/template"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// Data ... is the collection of inputs we need to fill our template
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
	pickedColor := green

	if percentage >= 0 && percentage < 33 {
		pickedColor = red
	} else if percentage >= 33 && percentage < 70 {
		pickedColor = yellow
	}

	return pickedColor
}

func init() {
	functions.HTTP("Progress", Progress)
}

// Progress ... Entrypoint of our Cloud Function
func Progress(w http.ResponseWriter, r *http.Request) {
	var id = fmt.Sprintf(path.Base(r.URL.Path))

	if percentage, err := strconv.Atoi(id); err == nil {
		data := Data{
			BackgroundColor: grey,
			Percentage:      percentage,
			Progress:        percentage - (percentage / 10),
			PickedColor:     pickColor(percentage),
		}

		//To test it locally we should interchange these two paths
		//tpl, err := template.ParseFiles("progress.html")
		tpl, err := template.ParseFiles("src/progress/progress.html")
		if err != nil {
			log.Fatalln(err)
		}

		buf := new(bytes.Buffer)

		err = tpl.Execute(buf, data)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("The percentage is: %d\n", percentage)
		w.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprintf(w, buf.String())
	}
}
