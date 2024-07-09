package controllers

import "net/http"

func ErroLeitura(w http.ResponseWriter, r *http.Request) {
	temp.ExecuteTemplate(w, "Erro", nil)
}
