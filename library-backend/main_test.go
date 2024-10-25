package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)



func TestGetAllAlbums(t *testing.T) {
    router := setupRouter()

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/albums", nil)
    router.ServeHTTP(w, req)

    expectedResponse := `[{"id":"1","title":"The Dark Side Of the Moon","artist":"Pink Floyd","price":99.99},{"id":"2","title":"Black Sabbath","artist":"Black Sabbath","price":99.99}]`
    assert.Equal(t, 200, w.Code)
    assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestGetAlbumByIDSuccessfully(t *testing.T) {
    router := setupRouter()

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/albums/1", nil)
    router.ServeHTTP(w, req)

    expectedResponse := `{"id":"1","title":"The Dark Side Of the Moon","artist":"Pink Floyd","price":99.99}`
    assert.Equal(t, 200, w.Code)
    assert.JSONEq(t, expectedResponse, w.Body.String())
}



func TestGetAlbumByIDNotFound(t *testing.T) {
    router := setupRouter()

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("GET", "/albums/99", nil)
    router.ServeHTTP(w, req)

    expectedResponse := `{"message":"album not found"}`
    assert.Equal(t, 404, w.Code)
    assert.JSONEq(t, expectedResponse, w.Body.String())
}

func TestPostAlbumSuccessfully(t *testing.T) {
    router := setupRouter()
    newAlbum := Album {
        ID: "3",
        Title: "Holy Diver",
        Artist: "Ronnie James Dio",
        Price: 9999.99,
    }
    newAlbumJSON, _ := json.Marshal(newAlbum)

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/albums", strings.NewReader(string(newAlbumJSON)))
    router.ServeHTTP(w, req)

    assert.Equal(t, 201, w.Code)
}

func TestPostAlbumBadRequest(t *testing.T) {
    router := setupRouter()
    invalidBodyJSON := `{"aa":"1"}`

    w := httptest.NewRecorder()
    req, _ := http.NewRequest("POST", "/albums", strings.NewReader(invalidBodyJSON))
    router.ServeHTTP(w, req)

    assert.Equal(t, 400, w.Code)
}


func TestPostAlbumConflict(t *testing.T) {
    router := setupRouter()
    newAlbum := Album {
        ID: "3",
        Title: "Holy Diver",
        Artist: "Ronnie James Dio",
        Price: 9999.99,
    }
    newAlbumJSON, _ := json.Marshal(newAlbum)

    w := httptest.NewRecorder()
    http.NewRequest("POST", "/albums", strings.NewReader(string(newAlbumJSON)))
    req, _ := http.NewRequest("POST", "/albums", strings.NewReader(string(newAlbumJSON)))
    router.ServeHTTP(w, req)

    assert.Equal(t, 409, w.Code)
}

