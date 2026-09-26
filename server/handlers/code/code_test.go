package code_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/bbconfhq/mycodes/database"
	"github.com/bbconfhq/mycodes/handlers"
	"github.com/bbconfhq/mycodes/models"
	"github.com/bbconfhq/mycodes/repository"
	"github.com/gofiber/fiber/v2"
)

func setup(t *testing.T) *fiber.App {
	t.Helper()
	db := database.Connect(filepath.Join(t.TempDir(), "test.db"))
	repository.Initialize(db)
	t.Cleanup(func() {
		if sqlDB, err := db.DB(); err == nil {
			_ = sqlDB.Close()
		}
	})
	app := fiber.New()
	handlers.Initialize(app)
	return app
}

func do(t *testing.T, app *fiber.App, req *http.Request) (int, map[string]interface{}) {
	t.Helper()
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	raw, _ := io.ReadAll(resp.Body)
	body := map[string]interface{}{}
	if err := json.Unmarshal(raw, &body); err != nil {
		t.Fatalf("invalid json %q: %v", raw, err)
	}
	return resp.StatusCode, body
}

func post(t *testing.T, app *fiber.App, payload string, xff string) (int, map[string]interface{}) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/code/", strings.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	if xff != "" {
		req.Header.Set("X-Forwarded-For", xff)
	}
	return do(t, app, req)
}

func TestPostRejectsMissingFields(t *testing.T) {
	app := setup(t)
	status, _ := post(t, app, `{"title":"","name":"n","content":"x"}`, "")
	if status != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", status)
	}
	status, _ = post(t, app, `{"title":"t","name":"n","content":"x","language":""}`, "")
	if status != http.StatusBadRequest {
		t.Errorf("empty language: status = %d, want 400", status)
	}
	status, _ = post(t, app, `{"title":"t","name":"n","content":""}`, "")
	if status != http.StatusBadRequest {
		t.Errorf("status = %d, want 400", status)
	}
}

func TestGetOneMissingIs404(t *testing.T) {
	app := setup(t)
	status, _ := do(t, app, httptest.NewRequest(http.MethodGet, "/api/v1/code/nope", nil))
	if status != http.StatusNotFound {
		t.Errorf("status = %d, want 404", status)
	}
}

func TestExpiredCodeIsHidden(t *testing.T) {
	app := setup(t)
	err := repository.Code.Create(&models.Code{
		ID: "old", Name: "n", Title: "t", Content: "c", Language: "go",
		ExpiredAt: time.Now().Add(-time.Minute),
	})
	if err != nil {
		t.Fatal(err)
	}
	status, _ := do(t, app, httptest.NewRequest(http.MethodGet, "/api/v1/code/old", nil))
	if status != http.StatusNotFound {
		t.Errorf("status = %d, want 404", status)
	}
	_, body := do(t, app, httptest.NewRequest(http.MethodGet, "/api/v1/code/", nil))
	if n := len(body["Data"].([]interface{})); n != 0 {
		t.Errorf("recent list has %d items, want 0", n)
	}
	deleted, err := repository.Code.DeleteExpired()
	if err != nil || deleted != 1 {
		t.Errorf("DeleteExpired = %d, %v; want 1, nil", deleted, err)
	}
}

func TestRecentListOmitsContent(t *testing.T) {
	app := setup(t)
	post(t, app, `{"title":"t","name":"n","content":"secret body","language":"go"}`, "")
	_, body := do(t, app, httptest.NewRequest(http.MethodGet, "/api/v1/code/", nil))
	items := body["Data"].([]interface{})
	if len(items) != 1 {
		t.Fatalf("got %d items", len(items))
	}
	item := items[0].(map[string]interface{})
	if item["content"] != "" {
		t.Errorf("recent item content = %q, want empty", item["content"])
	}
	if item["title"] != "t" || item["ip"] == "" || item["created_at"] == nil {
		t.Errorf("unexpected item: %v", item)
	}
}

func TestPostStoresMaskedIPOfNearestProxyHop(t *testing.T) {
	app := setup(t)
	cases := map[string]string{
		"6.6.6.6, 203.0.113.7": "203.0.0.0",
		"2001:db8:1234:5::1":   "2001:db8:1234::",
		"::ffff:10.1.2.3":      "10.1.0.0",
		"garbage":              "0.0.0.0",
	}
	for xff, want := range cases {
		_, body := post(t, app, `{"title":"t","name":"n","content":"c","language":"go"}`, xff)
		id := body["Data"].(map[string]interface{})["id"].(string)
		stored, err := repository.Code.Get(id)
		if err != nil {
			t.Fatal(err)
		}
		if stored.Ip != want {
			t.Errorf("X-Forwarded-For %q stored as %q, want %q", xff, stored.Ip, want)
		}
	}
}
