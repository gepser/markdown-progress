package progress

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
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

const gcloudFuncSourceDir = "serverless_function_source_code"

func fixDir() {
	fileInfo, err := os.Stat(gcloudFuncSourceDir)
	if err == nil && fileInfo.IsDir() {
		_ = os.Chdir(gcloudFuncSourceDir)
	}
}

func pickColor(percentage int, successColor string, warningColor string, dangerColor string) string {
	pickedColor := green
	if successColor != "" {
		pickedColor = "#" + successColor
	}

	if percentage >= 0 && percentage < 33 {
		if dangerColor != "" {
			pickedColor = "#" + dangerColor
		} else {
			pickedColor = red
		}
	} else if percentage >= 33 && percentage < 70 {
		if warningColor != "" {
			pickedColor = "#" + warningColor
		} else {
			pickedColor = yellow
		}
	}

	return pickedColor
}

func init() {
	fixDir()
	functions.HTTP("Progress", Progress)
}

// Progress ... Entrypoint of our Cloud Function
func Progress(w http.ResponseWriter, r *http.Request) {
	var id = fmt.Sprintf(path.Base(r.URL.Path))

	if percentage, err := strconv.Atoi(id); err == nil {

		// Read (with the intention to overwrite) success, warning, and danger colors if provided
		successColor := r.URL.Query().Get("successColor")
		warningColor := r.URL.Query().Get("warningColor")
		dangerColor := r.URL.Query().Get("dangerColor")

		data := Data{
			BackgroundColor: grey,
			Percentage:      percentage,
			Progress:        percentage - (percentage / 10),
			PickedColor:     pickColor(percentage, successColor, warningColor, dangerColor),
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

		fmt.Printf("The percentage is: %d\n", percentage)
		w.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprintf(w, buf.String())
	}
}
