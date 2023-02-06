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
	"github.com/pelletier/go-toml/v2"
	htmlTemplate "html/template"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

//go:embed resources/templates/index.html
var index string

//go:embed resources/templates/stats.html
var statsHTML string

//go:embed resources/static
var staticContent embed.FS

var config install.Configuration

var extensions = map[string]string{
	"css": "text/css",
	"js":  "text/javascript",
	"svg": "image/svg+xml",
	"ogg": "audio/ogg",
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "install" {
		err := install.App()
		if err != nil {
			panic(err)
		}
	}
	nixUser, err := user.Current()
	if err != nil {
		panic(fmt.Errorf("impossible de récupérer l'utilisateur courant %w", err))
	}
	home := nixUser.HomeDir
	switch runtime.GOOS {
	case "darwin":
		// TODO: get the config file
	case "linux":
		configPath := filepath.Join(home, ".config/local-searchengine/config.toml")
		b, err := os.ReadFile(configPath)
		if err != nil {
			panic(err)
		}
		err = toml.Unmarshal(b, &config)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("système d'exploitation non géré %s", runtime.GOOS))
	}

	customization.UpdateBlockList(config.BlockList)

	result, _ := parseHomepage(config.WallpaperPath)

	homePage := &result

	// API
	http.HandleFunc("/api/google", serveWithParser(searchengines.ParseGoogleResponse))
	http.HandleFunc("/api/ddg", serveWithParser(searchengines.ParseDDGResponse))
	http.HandleFunc("/api/log", utils.LogResult)

	// Static files
	http.HandleFunc("/static/", serveStatic)
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
	stats, err := features.LoadStats(config.AppName)
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
	var logo string
	if config.LogoPath == "" {
		logo = "/static/logo.svg#logo"
	} else {
		logo, err = utils.GeneratePathForHTML(&config.LogoPath)
		if err != nil {
			panic(err)
		}
	}
	wp, err := utils.GeneratePathForCss(&wallpaper)
	if err != nil {
		panic(err)
	}
	err = t.Execute(tempWriter, map[string]interface{}{
		"background": wp,
		"bangs":      htmlTemplate.JS(bangs),
		"appName":    config.AppName,
		"logo":       logo,
	})
	if err != nil {
		return "", err
	}
	s := tempWriter.String()
	s = strings.ReplaceAll(s, "<style>", "<link rel=\"stylesheet\" href=\"/static/style.css\"></link>\n  <style>")
	return s, nil
}

func setupCORS(r *http.ResponseWriter) {
	(*r).Header().Set("Access-Control-Allow-Origin", "*")
	(*r).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*r).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	_, extension, _ := strings.Cut(path, ".")
	for e, ext := range extensions {
		if extension == e {
			w.Header().Set("Content-Type", ext)
		}
	}
	content, err := staticContent.ReadFile("resources" + path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	w.Write(content)
}
