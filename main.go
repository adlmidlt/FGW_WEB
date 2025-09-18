package main

import (
	"fgw_web/internal/config"
	"fgw_web/internal/database"
	"fgw_web/internal/server"
	"html/template"
	"net/http"
)

const templateHtmlHome = "web/html/index.html"

func main() {
	cfg := config.LoadConfig()
	db := database.SetupDatabase(cfg)
	defer database.CloseDatabase(db)

	mux := http.NewServeMux()
	indexPage(mux)

	server.RunServerWithGracefulShutdown(mux)
}

// indexPage - начальная страница.
func indexPage(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(templateHtmlHome)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
