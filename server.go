package horo

import (
	"fmt"
	"net/http"
	"time"

	"github.com/boltdb/bolt"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
)

// NewServer creates app server instance
func NewServer(app *App) *negroni.Negroni {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()

	api.HandleFunc("/daily/{date}/{type}.xml", func(w http.ResponseWriter, r *http.Request) {
		// Init DB
		db, err := bolt.Open(app.DBPath, 0600, &bolt.Options{Timeout: 5 * time.Second})
		if err != nil {
			app.Log.Errorln(err)
			http.Error(w, "Application Error", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		dateStr := mux.Vars(r)["date"]
		_, err = time.Parse("20060102", dateStr)
		if err != nil {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		lang := r.FormValue("lang")
		
		if lang == "" {
			lang = "en"
		}
		
 		if lang != "en" && lang != "ru" && lang != "fr" && lang !="es" && lang != "de" && lang !="it" && lang != "jp" && lang != "cn" && lang != "pt" {
 			lang = "en"
 		}

		horoType := mux.Vars(r)["type"]
		horoKey := []byte(fmt.Sprintf("horo:%s:%s:%s", lang, horoType, dateStr))
		fmt.Println(string(horoKey))
		horoXML := []byte("")
		db.View(func(tx *bolt.Tx) error {
			horoXML = tx.Bucket([]byte("horoscopes")).Get(horoKey)
			return nil
		})

		if len(horoXML) <= 0 {
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "text/xml")
		w.Write(horoXML)
	})

	n := negroni.Classic()
	n.UseHandler(r)

	return n
}
