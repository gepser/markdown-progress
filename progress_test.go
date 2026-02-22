package progress

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func executeProgressRequest(path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(http.MethodGet, path, nil)
	rr := httptest.NewRecorder()
	Progress(rr, req)
	return rr
}

func executeProgressRequestWithMethod(method string, path string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, nil)
	rr := httptest.NewRecorder()
	Progress(rr, req)
	return rr
}

func TestPickColorThresholds(t *testing.T) {
	if got := pickColor(10, "", "", ""); got != red {
		t.Fatalf("expected red for low percentage, got %q", got)
	}

	if got := pickColor(50, "", "", ""); got != yellow {
		t.Fatalf("expected yellow for medium percentage, got %q", got)
	}

	if got := pickColor(80, "", "", ""); got != green {
		t.Fatalf("expected green for high percentage, got %q", got)
	}
}

func TestPickColorWithCustomColors(t *testing.T) {
	if got := pickColor(10, "#00aa00", "#ffcc00", "#aa0000"); got != "#aa0000" {
		t.Fatalf("expected custom danger color, got %q", got)
	}

	if got := pickColor(50, "#00aa00", "#ffcc00", "#aa0000"); got != "#ffcc00" {
		t.Fatalf("expected custom warning color, got %q", got)
	}

	if got := pickColor(80, "#00aa00", "#ffcc00", "#aa0000"); got != "#00aa00" {
		t.Fatalf("expected custom success color, got %q", got)
	}
}

func TestProgressReturnsSVGForValidInput(t *testing.T) {
	rr := executeProgressRequest("/progress/76")

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if got := rr.Header().Get("Content-Type"); got != "image/svg+xml" {
		t.Fatalf("expected image/svg+xml content type, got %q", got)
	}

	if got := rr.Header().Get("Cache-Control"); got != cacheControlValue {
		t.Fatalf("expected cache control %q, got %q", cacheControlValue, got)
	}

	body := rr.Body.String()
	if !strings.Contains(body, "76%") {
		t.Fatalf("expected percentage text in SVG body, got %q", body)
	}

	if !strings.Contains(body, `width="68" height="20" fill="#5cb85c"`) {
		t.Fatalf("expected progress width/color in SVG body, got %q", body)
	}
}

func TestProgressSupportsHEAD(t *testing.T) {
	rr := executeProgressRequestWithMethod(http.MethodHead, "/progress/76")

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	if got := rr.Header().Get("Content-Type"); got != "image/svg+xml" {
		t.Fatalf("expected image/svg+xml content type, got %q", got)
	}

	if got := rr.Header().Get("Cache-Control"); got != cacheControlValue {
		t.Fatalf("expected cache control %q, got %q", cacheControlValue, got)
	}

	if rr.Body.Len() != 0 {
		t.Fatalf("expected empty body for HEAD, got %q", rr.Body.String())
	}
}

func TestProgressRejectsUnsupportedMethod(t *testing.T) {
	rr := executeProgressRequestWithMethod(http.MethodPost, "/progress/50")

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", rr.Code)
	}

	if got := rr.Header().Get("Allow"); got != "GET, HEAD" {
		t.Fatalf("expected Allow header for GET and HEAD, got %q", got)
	}
}

func TestProgressReturnsBadRequestForInvalidPercentage(t *testing.T) {
	rr := executeProgressRequest("/progress/not-a-number")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "percentage must be an integer") {
		t.Fatalf("expected invalid percentage message, got %q", rr.Body.String())
	}
}

func TestProgressReturnsBadRequestForInvalidColor(t *testing.T) {
	rr := executeProgressRequest("/progress/50?successColor=nothex")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rr.Code)
	}

	if !strings.Contains(rr.Body.String(), "successColor must be a 6-character hex value") {
		t.Fatalf("expected invalid color message, got %q", rr.Body.String())
	}
}

func TestProgressClampsPercentageToRange(t *testing.T) {
	high := executeProgressRequest("/progress/150")
	if high.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", high.Code)
	}

	highBody := high.Body.String()
	if !strings.Contains(highBody, "100%") {
		t.Fatalf("expected clamped high percentage in body, got %q", highBody)
	}

	if !strings.Contains(highBody, `width="90" height="20" fill="#5cb85c"`) {
		t.Fatalf("expected full width for clamped high percentage, got %q", highBody)
	}

	low := executeProgressRequest("/progress/-10")
	if low.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", low.Code)
	}

	lowBody := low.Body.String()
	if !strings.Contains(lowBody, "0%") {
		t.Fatalf("expected clamped low percentage in body, got %q", lowBody)
	}

	if !strings.Contains(lowBody, `width="0" height="20" fill="#d9534f"`) {
		t.Fatalf("expected empty width for clamped low percentage, got %q", lowBody)
	}
}
