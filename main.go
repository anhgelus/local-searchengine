package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/anhgelus/local-searchengine/src/customization"
	"github.com/anhgelus/local-searchengine/src/features"
	"github.com/anhgelus/local-searchengine/src/install"
	"github.com/anhgelus/local-searchengine/src/searchengines"
	"github.com/anhgelus/local-searchengine/src/utils"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed resources/templates/index.html
var index string

//go:embed resources/templates/stats.html
var statsHTML string

//go:embed static/*
var staticContent embed.FS

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "install" {
		err := install.InstallApp()
		if err != nil {
			panic(err)
		}
		return
	}

	result, _ := parseHomepage("")

	homePage := &result
	c := customization.BingWallpaperFetcher()
	go func() {
		for wallpaper := range c {
			result, err := parseHomepage(wallpaper)
			if err == nil {
				*homePage = result
			}
			return
		}
	}()

	// API
	http.HandleFunc("/api/google", serveWithParser(searchengines.ParseGoogleResponse))
	http.HandleFunc("/api/ddg", serveWithParser(searchengines.ParseDDGResponse))
	http.HandleFunc("/api/log", utils.LogResult)

	// Static files
	http.Handle("/static/", http.FileServer(http.FS(staticContent)))
	http.HandleFunc("/weather", serveWeather)
	http.HandleFunc("/", serveHome(homePage))
	http.HandleFunc("/stats", serveStats)

	// Start the server
	fmt.Println("Listening on http://localhost:8042")
	log.Fatal(http.ListenAndServe(":8042", nil))
}

func serveWithParser(fn func(string) ([]searchengines.SearchResult, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupCORS(&w)
		q := r.URL.Query().Get("q")
		results, err := fn(features.ParseFilterBangs(q))
		if err != nil {
			serveError(w, err)
			return
		}
		data, err := json.Marshal(results)
		if err != nil {
			serveError(w, err)
			return
		}
		w.Write(data)
	}
}

func serveError(r http.ResponseWriter, err error) {
	r.WriteHeader(http.StatusInternalServerError)
	data := map[string]string{
		"message": err.Error(),
	}
	body, _ := json.Marshal(data)
	r.Write(body)
}

func serveRedirect(w http.ResponseWriter, url string) {
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

func serveHome(homePage *string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		q := r.URL.Query().Get("q")
		redirect := features.ParseRedirectBangs(q)
		if redirect != "" {
			serveRedirect(w, redirect)
			return
		}
		w.Write([]byte(*homePage))
	}
}

func serveWeather(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	url, err := features.ExtractUrlFromYrNoDk(q)
	if err != nil {
		serveError(w, err)
		return
	}
	serveRedirect(w, url)
}

func serveStats(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("stats.html").Parse(statsHTML)
	if err != nil {
		serveError(w, err)
		return
	}
	stats, err := features.LoadStats()
	if err != nil {
		serveError(w, err)
		return
	}
	t.Execute(w, stats)
}

func parseHomepage(wallpaper string) (string, error) {
	bangs, err := json.Marshal(features.RedirectBangs)
	if err != nil {
		return "", err
	}
	t, err := template.New("index.html").Parse(index)
	if err != nil {
		return "", err
	}

	tempWriter := new(strings.Builder)
	err = t.Execute(tempWriter, map[string]interface{}{
		"background": wallpaper,
		"bangs":      template.JS(bangs),
	})
	if err != nil {
		return "", err
	}
	s := tempWriter.String()
	s = strings.ReplaceAll(s, "/assets/app.ts", "/static/app.js")
	s = strings.ReplaceAll(s, "<style>", "<link rel=\"stylesheet\" href=\"/static/style.css\"></link>\n  <style>")
	return s, nil
}

func setupCORS(r *http.ResponseWriter) {
	(*r).Header().Set("Access-Control-Allow-Origin", "*")
	(*r).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*r).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}
