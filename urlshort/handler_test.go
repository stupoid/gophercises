package urlshort

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/boltdb/bolt"
)

func TestMapHandler(t *testing.T) {
	fallbackResponse := "fallback"
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fallbackResponse)
	})
	pathToUrls := map[string]string{"/foo": "/bar"}

	t.Run("/foo redirects to /bar", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/foo", nil)
		recorder := httptest.NewRecorder()
		mapHandler := MapHandler(pathToUrls, fallback)
		mapHandler(recorder, request)
		response := recorder.Result()
		assertLocation(t, response, "/bar")
		assertStatusCode(t, response, http.StatusTemporaryRedirect)
	})

	t.Run("fallback response on unknown routes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
		recorder := httptest.NewRecorder()
		mapHandler := MapHandler(pathToUrls, fallback)
		mapHandler(recorder, request)
		response := recorder.Result()
		assertStatusCode(t, response, http.StatusOK)
		assertBody(t, response, fallbackResponse)
	})
}

func TestYAMLHandler(t *testing.T) {
	fallbackResponse := "fallback"
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fallbackResponse)
	})
	yaml := `
- path: /foo
  url: /bar
`

	t.Run("/foo redirects to /bar", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/foo", nil)
		recorder := httptest.NewRecorder()
		yamlHandler, err := YAMLHandler([]byte(yaml), fallback)
		if err != nil {
			t.Fatal("Could not handle request", err)
		}
		yamlHandler(recorder, request)
		response := recorder.Result()
		assertLocation(t, response, "/bar")
		assertStatusCode(t, response, http.StatusTemporaryRedirect)
	})

	t.Run("fallback response on unknown routes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
		recorder := httptest.NewRecorder()
		yamlHandler, err := YAMLHandler([]byte(yaml), fallback)
		if err != nil {
			t.Fatal("Could not handle request", err)
		}
		yamlHandler(recorder, request)
		response := recorder.Result()
		assertStatusCode(t, response, http.StatusOK)
		assertBody(t, response, fallbackResponse)
	})
}
func TestJSONHandler(t *testing.T) {
	fallbackResponse := "fallback"
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fallbackResponse)
	})
	json := `
[
	{
		"path": "/foo",
		"url": "/bar"
	}
]
`

	t.Run("/foo redirects to /bar", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/foo", nil)
		recorder := httptest.NewRecorder()
		jsonHandler, err := JSONHandler([]byte(json), fallback)
		if err != nil {
			t.Fatal("Could not handle request", err)
		}
		jsonHandler(recorder, request)
		response := recorder.Result()
		assertLocation(t, response, "/bar")
		assertStatusCode(t, response, http.StatusTemporaryRedirect)
	})

	t.Run("fallback response on unknown routes", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
		recorder := httptest.NewRecorder()
		jsonHandler, err := JSONHandler([]byte(json), fallback)
		if err != nil {
			t.Fatal("Could not handle request", err)
		}
		jsonHandler(recorder, request)
		response := recorder.Result()
		assertStatusCode(t, response, http.StatusOK)
		assertBody(t, response, fallbackResponse)
	})

}

func TestBoltHandler(t *testing.T) {
	fallbackResponse := "fallback"
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, fallbackResponse)
	})
	db, err := bolt.Open("test.db", 0600, nil)
	if err != nil {
		t.Fatal("Could not create db", err)
	}
	defer db.Close()
	defer os.Remove("test.db")

	db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte("TestBucket"))
		if err != nil {
			t.Fatal("Could not create bucket", err)
		}
		err = bucket.Put([]byte("/foo"), []byte("bar"))
		if err != nil {
			t.Fatal("Could not add entry to bucket", err)
		}

		t.Run("/foo redirects to /bar", func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/foo", nil)
			recorder := httptest.NewRecorder()
			boltHandler := BoltHandler(bucket, fallback)
			boltHandler(recorder, request)
			response := recorder.Result()
			assertLocation(t, response, "/bar")
			assertStatusCode(t, response, http.StatusTemporaryRedirect)
		})

		t.Run("fallback response on unknown routes", func(t *testing.T) {
			request, _ := http.NewRequest(http.MethodGet, "/unknown", nil)
			recorder := httptest.NewRecorder()
			boltHandler := BoltHandler(bucket, fallback)
			boltHandler(recorder, request)
			response := recorder.Result()
			assertStatusCode(t, response, http.StatusOK)
			assertBody(t, response, fallbackResponse)
		})

		return nil
	})

}

func assertBody(t testing.TB, response *http.Response, want string) {
	t.Helper()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		t.Fatal("Could not read Body", err)
	}
	got := string(body)
	if got != want {
		t.Errorf("Expected response Body to be %s, got %s", want, got)
	}
}

func assertStatusCode(t testing.TB, response *http.Response, want int) {
	t.Helper()
	got := response.StatusCode
	if got != want {
		t.Errorf("Expected response StatusCode to be %d, got %d", want, got)
	}
}

func assertLocation(t testing.TB, response *http.Response, want string) {
	t.Helper()
	url, err := response.Location()
	if err != nil {
		t.Fatal("Could not read Location", err)
	}
	got := url.String()
	if got != want {
		t.Errorf("Expected response Location to be %s, got %s", want, got)
	}
}
