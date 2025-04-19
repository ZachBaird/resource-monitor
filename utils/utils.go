package utils

import (
	"html/template"
	"net/http"
	"resource-monitor/types"
)

func ServeTemplateFile(fileName string, w http.ResponseWriter, status *types.AppStatus, templates *template.Template) {
	err := templates.ExecuteTemplate(w, fileName, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
