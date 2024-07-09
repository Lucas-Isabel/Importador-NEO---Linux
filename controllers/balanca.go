package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/lucasbyte/go-clipse/file"
	"github.com/lucasbyte/go-clipse/models"
)

func New_balanca(w http.ResponseWriter, r *http.Request) {
	_, setores := file.VerificarSetores()
	err := temp.ExecuteTemplate(w, "NewBalanca", setores)
	if err != nil {
		temp.ExecuteTemplate(w, "NewBalanca", nil)
	}
}

func Add_balanca(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ip := r.FormValue("codigo")
		nome := strings.ToUpper(r.FormValue("descricao"))
		setores_form := r.Form["setores"]
		cod_setores, str_setores := separeToComma(setores_form)

		models.CriaNovaBalanca(ip, nome, cod_setores, str_setores)
	}
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func Balanca(w http.ResponseWriter, r *http.Request) {
	balancas, err := models.BuscaBalancas()
	if err != nil {
		fmt.Println(err)
	}
	temp.ExecuteTemplate(w, "/", balancas)
}

type setores_e_balanca struct {
	Setores []models.Setor
	Balanca models.Balanca
}

func Edit_balanca(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		ip := r.FormValue("codigo")
		nome := strings.ToUpper(r.FormValue("descricao"))
		setores_form := r.Form["setores"]
		cod_setores, str_setores := separeToComma(setores_form)

		// Tratamento de erro de conexão com o banco de dados
		if err := models.AtualizaBalanca(ip, nome, cod_setores, str_setores); err != nil {
			http.Error(w, "Erro ao atualizar balança: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// balancas, err := models.BuscaBalancas()
		// if err != nil {
		// 	http.Error(w, "Erro ao buscar balanças: "+err.Error(), http.StatusInternalServerError)
		// 	return
		// }

		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	} else {
		idDaBalanca := r.URL.Query().Get("identifiCador")

		_, setores := file.VerificarSetores()
		balanca, err := models.BuscaBalanca(idDaBalanca)
		if err != nil {
			http.Error(w, "Erro ao buscar balança: "+err.Error(), http.StatusInternalServerError)
			return
		}

		Setor_e_balanca := setores_e_balanca{
			Setores: setores,
			Balanca: balanca,
		}

		temp.ExecuteTemplate(w, "EditBalanca", Setor_e_balanca)
	}
}

// func Update_balanca(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == "POST" {
// 		ip := r.FormValue("codigo")
// 		nome := strings.ToUpper(r.FormValue("descricao"))
// 		setores_form := r.Form["setores"]
// 		cod_setores, str_setores := separeToComma(setores_form)
// 		models.AtualizaBalanca(ip, nome, cod_setores, str_setores)
// 		balancas := models.BuscaBalancas(),

// 		temp.ExecuteTemplate(w, "Balanca", balancas)
// 	}

// }

func Delete_bal(w http.ResponseWriter, r *http.Request) {
	idDaBalanca := ""
	idDaBalanca = r.URL.Query().Get("identifiCador")
	time.Sleep(time.Second)
	fmt.Println(idDaBalanca)
	fmt.Println(models.DeleteBalanca(idDaBalanca))
	time.Sleep(time.Second)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func separeToComma(lista []string) (string, string) {
	var valoresAntes []string
	var valoresDepois []string
	for _, value := range lista {
		valorAntesDepois := strings.Split(value, ";")
		// Adicione os valores às listas correspondentes
		valoresAntes = append(valoresAntes, valorAntesDepois[0])
		valoresDepois = append(valoresDepois, valorAntesDepois[1])
	}

	return fmt.Sprint(valoresAntes), fmt.Sprint(valoresDepois)
}
