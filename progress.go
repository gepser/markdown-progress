package progress

import (
	"bytes"
	"log"
	"math"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// Data ... is the collection of inputs we need to fill our template
type Data struct {
	BackgroundColor string
	Label           string
	Progress        int
	PickedColor     string
}

const (
	gcloudFuncSourceDir = "serverless_function_source_code"
	minPercentage       = 0.0
	maxPercentage       = 100.0
	totalBarWidth       = 90.0
	cacheControlValue   = "public, max-age=300"
	maxLabelRunes       = 64
)

var (
	grey   = "#555"
	red    = "#d9534f"
	yellow = "#f0ad4e"
	green  = "#5cb85c"

	hexColorPattern     = regexp.MustCompile(`^[0-9a-fA-F]{6}$`)
	progressTemplate    *template.Template
	progressTemplateErr error
)

func fixDir() {
	fileInfo, err := os.Stat(gcloudFuncSourceDir)
	if err == nil && fileInfo.IsDir() {
		_ = os.Chdir(gcloudFuncSourceDir)
	}
}

func clampPercentage(percentage float64) float64 {
	if percentage < minPercentage {
		return minPercentage
	}

	if percentage > maxPercentage {
		return maxPercentage
	}

	return percentage
}

func percentageToWidth(percentage float64) int {
	return int((totalBarWidth * percentage) / maxPercentage)
}

func parseOptionalColor(raw string) (string, bool) {
	if raw == "" {
		return "", true
	}

	if !hexColorPattern.MatchString(raw) {
		return "", false
	}

	return "#" + strings.ToLower(raw), true
}

func pickColor(percentage float64, successColor string, warningColor string, dangerColor string) string {
	pickedColor := green
	if successColor != "" {
		pickedColor = successColor
	}

	if percentage >= 0 && percentage < 33 {
		if dangerColor != "" {
			pickedColor = dangerColor
		} else {
			pickedColor = red
		}
	} else if percentage >= 33 && percentage < 70 {
		if warningColor != "" {
			pickedColor = warningColor
		} else {
			pickedColor = yellow
		}
	}

	return pickedColor
}

func formatNumber(value float64) string {
	formatted := strconv.FormatFloat(value, 'f', 2, 64)
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimRight(formatted, ".")
	if formatted == "-0" || formatted == "" {
		return "0"
	}

	return formatted
}

func formatPercentLabel(value float64) string {
	return formatNumber(value) + "%"
}

func init() {
	fixDir()
	progressTemplate, progressTemplateErr = template.ParseFiles("progress.html")
	functions.HTTP("Progress", Progress)
}

// Progress ... Entrypoint of our Cloud Function
func Progress(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	statusCode := http.StatusOK
	defer func() {
		log.Printf(
			"method=%s path=%s status=%d duration_ms=%d",
			r.Method,
			r.URL.Path,
			statusCode,
			time.Since(start).Milliseconds(),
		)
	}()

	if r.Method != http.MethodGet && r.Method != http.MethodHead {
		statusCode = http.StatusMethodNotAllowed
		w.Header().Set("Allow", "GET, HEAD")
		http.Error(w, "method not allowed", statusCode)
		return
	}

	if progressTemplateErr != nil {
		statusCode = http.StatusInternalServerError
		log.Printf("unable to parse progress template: %v", progressTemplateErr)
		http.Error(w, "template unavailable", statusCode)
		return
	}

	id := path.Base(strings.TrimSuffix(r.URL.Path, "/"))
	value, err := strconv.ParseFloat(id, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		statusCode = http.StatusBadRequest
		http.Error(w, "percentage must be a number", statusCode)
		return
	}

	// Read and validate colors if provided.
	successColor, ok := parseOptionalColor(r.URL.Query().Get("successColor"))
	if !ok {
		statusCode = http.StatusBadRequest
		http.Error(w, "successColor must be a 6-character hex value", statusCode)
		return
	}

	warningColor, ok := parseOptionalColor(r.URL.Query().Get("warningColor"))
	if !ok {
		statusCode = http.StatusBadRequest
		http.Error(w, "warningColor must be a 6-character hex value", statusCode)
		return
	}

	dangerColor, ok := parseOptionalColor(r.URL.Query().Get("dangerColor"))
	if !ok {
		statusCode = http.StatusBadRequest
		http.Error(w, "dangerColor must be a 6-character hex value", statusCode)
		return
	}

	barColor, ok := parseOptionalColor(r.URL.Query().Get("barColor"))
	if !ok {
		statusCode = http.StatusBadRequest
		http.Error(w, "barColor must be a 6-character hex value", statusCode)
		return
	}

	customLabel := r.URL.Query().Get("label")
	if utf8.RuneCountInString(customLabel) > maxLabelRunes {
		statusCode = http.StatusBadRequest
		http.Error(w, "label is too long (max 64 characters)", statusCode)
		return
	}

	minRaw := r.URL.Query().Get("min")
	maxRaw := r.URL.Query().Get("max")
	hasMin := minRaw != ""
	hasMax := maxRaw != ""

	if hasMin != hasMax {
		statusCode = http.StatusBadRequest
		http.Error(w, "min and max must be provided together", statusCode)
		return
	}

	percentage := clampPercentage(value)
	if hasMin {
		minValue, minErr := strconv.ParseFloat(minRaw, 64)
		maxValue, maxErr := strconv.ParseFloat(maxRaw, 64)
		if minErr != nil || maxErr != nil || math.IsNaN(minValue) || math.IsNaN(maxValue) || math.IsInf(minValue, 0) || math.IsInf(maxValue, 0) {
			statusCode = http.StatusBadRequest
			http.Error(w, "min and max must be numeric values", statusCode)
			return
		}

		if maxValue <= minValue {
			statusCode = http.StatusBadRequest
			http.Error(w, "max must be greater than min", statusCode)
			return
		}

		normalized := ((value - minValue) / (maxValue - minValue)) * maxPercentage
		percentage = clampPercentage(normalized)
	}

	pickedColor := pickColor(percentage, successColor, warningColor, dangerColor)
	if barColor != "" {
		pickedColor = barColor
	}

	label := formatPercentLabel(percentage)
	if hasMin {
		label = formatNumber(value)
	}
	if customLabel != "" {
		label = customLabel
	}

	data := Data{
		BackgroundColor: grey,
		Label:           label,
		Progress:        percentageToWidth(percentage),
		PickedColor:     pickedColor,
	}

	buf := new(bytes.Buffer)
	err = progressTemplate.Execute(buf, data)
	if err != nil {
		statusCode = http.StatusInternalServerError
		log.Printf("unable to render progress template: %v", err)
		http.Error(w, "failed to render SVG", statusCode)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml")
	w.Header().Set("Cache-Control", cacheControlValue)

	if r.Method == http.MethodHead {
		return
	}

	_, _ = w.Write(buf.Bytes())
}
