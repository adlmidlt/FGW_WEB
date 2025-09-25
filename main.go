package main

import (
	"context"
	"fgw_web/internal/config"
	"fgw_web/internal/database"
	"fgw_web/internal/server"
	"html/template"
	"log"
	"net/http"
)

const templateHtmlHome = "web/html/index.html"
const pathToYamlFile = "internal/config/database.yml"

func main() {
	var cfg config.Config
	if err := cfg.LoadConfigDatabase(pathToYamlFile); err != nil {
		log.Fatalf("Ошибка загрузки конфигурационных данных: %v", err)
	}

	db, err := database.NewPgxPool(context.Background(), &cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

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
