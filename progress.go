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

type Data struct {
	BackgroundColor string
	Percentage      float64
	Progress        float64
	PickedColor     string
	Label           string
	ShowLabel       bool
	BarWidth        float64
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

// Template helper to check if a float is an integer
func isInt(f float64) bool {
	return f == float64(int(f))
}

// Template helper for division
func div(a, b float64) float64 {
	if b == 0 {
		return 0
	}
	return a / b
}

func Progress(w http.ResponseWriter, r *http.Request) {
	id := path.Base(r.URL.Path)

	value, err := strconv.ParseFloat(id, 64)
	if err != nil {
		http.Error(w, "Invalid value", http.StatusBadRequest)
		return
	}

	successColor := r.URL.Query().Get("successColor")
	warningColor := r.URL.Query().Get("warningColor")
	dangerColor := r.URL.Query().Get("dangerColor")
	barColor := r.URL.Query().Get("barColor")
	customLabel := r.URL.Query().Get("label")
	minStr := r.URL.Query().Get("min")
	maxStr := r.URL.Query().Get("max")

	// Determine if we're using custom label mode (data bar mode)
	showLabel := customLabel != ""
	
	// Calculate percentage and bar width
	var percentage float64
	var barWidth float64
	const defaultWidth = 90.0
	
	if showLabel && minStr != "" && maxStr != "" {
		// Data bar mode: scale proportionally based on min/max
		min, errMin := strconv.ParseFloat(minStr, 64)
		max, errMax := strconv.ParseFloat(maxStr, 64)
		
		if errMin != nil || errMax != nil || min >= max {
			http.Error(w, "Invalid min/max values", http.StatusBadRequest)
			return
		}
		
		// Calculate percentage for color picking (0-100)
		if max > min {
			percentage = ((value - min) / (max - min)) * 100
		} else {
			percentage = 0
		}
		
		// Calculate bar width proportionally (leave some padding)
		if max > min {
			barWidth = ((value - min) / (max - min)) * (defaultWidth * 0.95)
		} else {
			barWidth = 0
		}
	} else {
		// Original percentage mode
		percentage = value
		if percentage > 100 {
			percentage = 100
		}
		barWidth = percentage - (percentage / 10)
	}

	// Determine bar color
	var pickedColor string
	if barColor != "" {
		// Use custom bar color if provided (for data bars)
		pickedColor = "#" + barColor
	} else {
		// Use percentage-based color logic
		pickedColor = pickColor(percentage, successColor, warningColor, dangerColor)
	}

	data := Data{
		BackgroundColor: grey,
		Percentage:      percentage,
		Progress:        barWidth,
		BarWidth:        defaultWidth,
		PickedColor:     pickedColor,
		Label:           customLabel,
		ShowLabel:       showLabel,
	}

	// Parse template with custom functions
	tpl := template.New("progress.html").Funcs(template.FuncMap{
		"isInt": isInt,
		"div":   div,
	})
	tpl, err = tpl.ParseFiles("progress.html")
	if err != nil {
		log.Fatalln("Error parsing template:", err)
	}

	buf := new(bytes.Buffer)
	if err := tpl.Execute(buf, data); err != nil {
		log.Fatalln("Error executing template:", err)
	}

	fmt.Printf("The percentage is: %v\n", percentage)
	w.Header().Set("Content-Type", "image/svg+xml")
	_, _ = w.Write(buf.Bytes())
}