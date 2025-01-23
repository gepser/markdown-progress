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

// Data ... is the collection of inputs to fill our template
type Data struct {
	BackgroundColor string
	Percentage      float64
	Progress        float64
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

func pickColor(percentage float64, successColor, warningColor, dangerColor string) string {
	pickedColor := green
	if successColor != "" {
		pickedColor = "#" + successColor
	}

	switch {
	case percentage >= 0 && percentage < 33:
		if dangerColor != "" {
			pickedColor = "#" + dangerColor
		} else {
			pickedColor = red
		}
	case percentage >= 33 && percentage < 70:
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
	id := path.Base(r.URL.Path)

	percentage, err := strconv.ParseFloat(id, 64)
	if err != nil {
		http.Error(w, "Invalid percentage value", http.StatusBadRequest)
		return
	}

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
		log.Fatalln("Error parsing template:", err)
	}

	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, data); err != nil {
		log.Fatalln("Error executing template:", err)
	}

	fmt.Printf("The percentage is: %.2f\n", percentage)
	w.Header().Set("Content-Type", "image/svg+xml")
	_, _ = w.Write(buf.Bytes())
}