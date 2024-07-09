package controllers

import (
	"fmt"
	"net/http"

	"github.com/lucasbyte/go-clipse/Cad"
	"github.com/lucasbyte/go-clipse/file"
	"github.com/lucasbyte/go-clipse/models"
)

func EnviarDadosPg(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		itens := r.FormValue("itens")
		file.DeleteAllSetores()
		receita := r.FormValue("mgv-receita")
		nutri := r.FormValue("mgv-nutri")
		frac := r.FormValue("mgv-fracionador")
		forn := r.FormValue("mgv-fornecedor")
		tara := r.FormValue("mgv-tara")
		cons := r.FormValue("mgv-conservantes")
		campext1 := r.FormValue("campo-extra-1")

		caminhos_dos_arquivos := make(map[string]string)

		caminhos_dos_arquivos["itens"] = itens
		caminhos_dos_arquivos["receita"] = receita
		caminhos_dos_arquivos["nutricional"] = nutri
		caminhos_dos_arquivos["conserva"] = cons
		caminhos_dos_arquivos["fraciona"] = frac
		caminhos_dos_arquivos["tara"] = tara
		caminhos_dos_arquivos["fornecedor"] = forn
		caminhos_dos_arquivos["campo-extra-1"] = campext1

		for k, caminho := range caminhos_dos_arquivos {
			file.EditFilePath(k, "itensmgv.json", caminho)
		}

		file.EscreverJSON("ITENSMGV")
		// file.EnviaParaBalancas(itens, receita, nutri, frac, forn, tara, cons, campext1)
		file.AdicionaSetores(itens)
	}
	setores := models.BuscaSetores()
	temp.ExecuteTemplate(w, "Setor", setores)
}

func EnviarTxitensPg(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		itens := r.FormValue("txitens_itens")
		file.DeleteAllSetores()
		// itens = strings.Replace(itens, "\\", "//", -1)
		// fmt.Println("Txitens:", itens)
		// balancas, err := models.BuscaBalancas()
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// models.Txitens(itens, balancas)
		file.EscreverJSON("TXITENS")
		file.AdicionaSetoresTxitens(itens)
	}
	setores := models.BuscaSetores()
	temp.ExecuteTemplate(w, "Setor", setores)
}

func EnviarCadPg(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		itens := r.FormValue("Cad-itens")
		file.DeleteAllSetores()
		// itens = strings.Replace(itens, "\\", "//", -1)
		// fmt.Println("Txitens:", itens)
		// balancas, err := models.BuscaBalancas()
		// if err != nil {
		// 	fmt.Println(err)
		// }
		// models.Txitens(itens, balancas)
		new_itens := Cad.CadToItens(itens)
		file.AdicionaSetores(new_itens)
		file.EscreverJSON("Cad")
		fmt.Println(itens)
	}
	setores := models.BuscaSetores()
	temp.ExecuteTemplate(w, "Setor", setores)
}
