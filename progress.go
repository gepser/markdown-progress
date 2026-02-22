package progress

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// Data ... is the collection of inputs we need to fill our template
type Data struct {
	BackgroundColor string
	Percentage      int
	Progress        int
	PickedColor     string
}

const (
	gcloudFuncSourceDir = "serverless_function_source_code"
	minPercentage       = 0
	maxPercentage       = 100
	totalBarWidth       = 90
	cacheControlValue   = "public, max-age=300"
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

func clampPercentage(percentage int) int {
	if percentage < minPercentage {
		return minPercentage
	}

	if percentage > maxPercentage {
		return maxPercentage
	}

	return percentage
}

func percentageToWidth(percentage int) int {
	return (totalBarWidth * percentage) / maxPercentage
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

func pickColor(percentage int, successColor string, warningColor string, dangerColor string) string {
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
	percentage, err := strconv.Atoi(id)
	if err != nil {
		statusCode = http.StatusBadRequest
		http.Error(w, "percentage must be an integer", statusCode)
		return
	}

	percentage = clampPercentage(percentage)

	// Read and validate success, warning, and danger colors if provided.
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

	data := Data{
		BackgroundColor: grey,
		Percentage:      percentage,
		Progress:        percentageToWidth(percentage),
		PickedColor:     pickColor(percentage, successColor, warningColor, dangerColor),
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
