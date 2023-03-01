package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const StatsFile = Datas + "stats.csv"

type logError struct {
	Message string `json:"message"`
}

type logRequest struct {
	Query  string `json:"q"`
	Found  bool   `json:"found"`
	Google bool   `json:"googleFound"`
	URL    string `json:"url,omitempty"`
}

func (data logRequest) AppendToFile(w io.Writer) error {
	csvw := csv.NewWriter(w)
	defer csvw.Flush()
	return csvw.Write([]string{
		data.Query,
		LogBoolToString(data.Found),
		LogBoolToString(data.Google),
		data.URL,
		time.Now().Format("2006-01-02"),
	})
}

func LogBoolToString(b bool) string {
	v := "1"
	if !b {
		v = "0"
	}
	return v
}

func LogResult(w http.ResponseWriter, r *http.Request) {
	// Parse the request
	var data logRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)
	if err != nil {
		LogErrorResponse(w, err)
		return
	}

	// Open the stats.csv file
	u, err := user.Current()
	if err != nil {
		LogErrorResponse(w, err)
		return
	}
	f, err := os.OpenFile(filepath.Join(u.HomeDir, StatsFile), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		LogErrorResponse(w, err)
		return
	}
	defer f.Close()

	// Write the CSV line
	err = data.AppendToFile(f)
	if err != nil {
		LogErrorResponse(w, err)
		return
	}
}

func LogErrorResponse(w http.ResponseWriter, err error) {
	fmt.Println(err)
	data, _ := json.Marshal(logError{Message: err.Error()})
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("content-type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		panic("Cannot reply to error on LogErrorResponse")
	}
}
