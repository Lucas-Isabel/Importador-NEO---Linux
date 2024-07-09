package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lucasbyte/go-clipse/models"
)

func Setores(w http.ResponseWriter, r *http.Request) {
	setores := models.BuscaSetores()
	temp.ExecuteTemplate(w, "Setor", setores)
}

func Update_setor(w http.ResponseWriter, r *http.Request) {
	idDoSetor := ""
	nome := ""

	idDoSetor = r.URL.Query().Get("identifiCador1")
	nome = r.URL.Query().Get("identifiCador2")

	// Verificar se os parâmetros estão sendo passados corretamente
	fmt.Println("ID do Setor:", idDoSetor)
	fmt.Println("Nome do Setor:", nome)

	int_idDoSetor, err := strconv.Atoi(idDoSetor)
	if err != nil {
		fmt.Println("Erro ao converter a string para inteiro:", err)
		return
	}

	// Verificar se a conversão para inteiro está correta
	fmt.Println("ID do Setor (int):", int_idDoSetor)

	// Verificar se o modelo está sendo atualizado corretamente
	models.Update_setor(nome, int_idDoSetor)

	// Redirecionamento após a atualização bem-sucedida
	setores := models.BuscaSetores()
	temp.ExecuteTemplate(w, "Setor", setores)
}
