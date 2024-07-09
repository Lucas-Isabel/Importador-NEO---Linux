package controllers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/lucasbyte/go-clipse/Cad"
	"github.com/lucasbyte/go-clipse/file"
	"github.com/lucasbyte/go-clipse/models"
	"github.com/lucasbyte/go-clipse/txitens"
)

var temp *template.Template

// SetTemplates configura o template compilado a partir do embed.FS
func SetTemplates(t *template.Template) {
	temp = t
}

func Index(w http.ResponseWriter, r *http.Request) {
	balancas, err := models.BuscaBalancas()
	if err != nil {
		fmt.Println(err)
	}
	err = temp.ExecuteTemplate(w, "Index", balancas)
	if err != nil {
		temp.ExecuteTemplate(w, "Index", nil)
	}
}

// func ToImport(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Teste")
// 	var balancas_import_checked []models.Balanca

// 	formato, _ := file.LerJSON()

// 	if r.Method == "POST" {
// 		var balancas_import []models.Balanca
// 		fmt.Println("Teste")
// 		balancas, err := models.BuscaBalancas()
// 		for _, balanca := range balancas {
// 			ip := balanca.Ip
// 			checkboxValue := r.FormValue(ip)
// 			fmt.Println(checkboxValue)
// 			fmt.Println(ip)
// 			if checkboxValue == "on" {
// 				balancas_import = append(balancas_import, balanca)
// 			} else {
// 				println("NADA: ", ip)
// 			}
// 		}

// 		if formato.Tipo == "TXITENS" {
// 			arquivos := txitens.ReadTxitensJson()
// 			itensFile := arquivos.Caminhos.Itens
// 			models.Txitens(itensFile, balancas_import)
// 			http.Redirect(w, r, "/", 301)
// 		} else if formato.Tipo == "Cad" {
// 			arquivos := Cad.ReadCadJson()
// 			itensFile := arquivos.Caminhos_Cad.Itens_Cad
// 			receitaFile := arquivos.Caminhos_Cad.Receita_Cad
// 			extra := arquivos.Caminhos_Cad.CampoExtra_Cad
// 			Cad.CadImport(itensFile, receitaFile, extra, balancas_import)
// 			http.Redirect(w, r, "/", 301)
// 		} else {

// 			if err != nil {
// 				fmt.Println(err)
// 			}
// 			arquivos := file.ReadMGVJson()

// 			itensFile := arquivos.Caminhos.Itens
// 			receitaFile := arquivos.Caminhos.Receita
// 			nutriFile := arquivos.Caminhos.Nutricional
// 			campoextra := arquivos.Caminhos.CampoExtra
// 			fornFile := arquivos.Caminhos.Fornecedor
// 			fracionaFile := arquivos.Caminhos.Fracionador
// 			taraFile := arquivos.Caminhos.Tara
// 			conservaFile := arquivos.Caminhos.Conservantes

// 			somente_preco_form := r.FormValue("somente_preco")
// 			somente_preco := false
// 			fmt.Println(somente_preco_form)
// 			fmt.Println(somente_preco)
// 			if somente_preco_form == "on" {
// 				somente_preco = true
// 			} else {
// 				somente_preco = false
// 			}
// 			balancas_import_checked = file.EnviaParaBalancas(itensFile, receitaFile, nutriFile, fracionaFile, fornFile, taraFile, conservaFile, campoextra, balancas_import, somente_preco)
// 			temp.ExecuteTemplate(w, "Log", balancas_import_checked)
// 		}
// 	}

// }

func ToImport(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		var progress int
		var balancas_import []models.Balanca
		var balancas_import_error []models.Balanca
		balancas, err := models.BuscaBalancas()
		if err != nil {
			fmt.Fprintf(w, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		}

		totalSteps := len(balancas) + 1
		currentStep := 0

		checkImport := r.FormValue("Importar-checkbox")
		if checkImport == "on" {

		}

		{
			for _, balanca := range balancas {
				ip := balanca.Ip
				checkboxValue := r.FormValue(ip)
				if checkboxValue == "on" {
					balancas_import = append(balancas_import, balanca)
				}
				currentStep++
				progress = (currentStep * 100) / totalSteps
				fmt.Fprintf(w, "event: progress\ndata: %d\n\n", progress)
				flusher.Flush()
				time.Sleep(500 * time.Millisecond) // Simulate processing time
			}

			time.Sleep(500 * time.Millisecond)
			balancas_conect, err := json.Marshal(balancas_import)
			if err != nil {
				fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
				flusher.Flush()
				return
			}
			fmt.Fprintf(w, "data: {\"conect\": %s}\n\n", balancas_conect)
			flusher.Flush()
			time.Sleep(time.Millisecond * 500)

			progress = 25
			fmt.Fprintf(w, "data: {\"progress\": %d}\n\n", progress)
			flusher.Flush()

			var balancas_import_checked []models.Balanca
			formatos := file.LerTipoJson()
			if formatos == "TXITENS" {
				arquivos := txitens.ReadTxitensJson()
				itensFile := arquivos.Caminhos.Itens
				models.Txitens(itensFile, balancas_import)
				if err != nil {
					progress += 75
					fmt.Fprintf(w, "data: {\"progress\": %d}\n\n", progress)
					flusher.Flush()
				}
			} else if formatos == "Cad" {
				arquivos := Cad.ReadCadJson()
				itensFile := arquivos.Caminhos_Cad.Itens_Cad
				_ = arquivos.Caminhos_Cad.Receita_Cad
				_ = arquivos.Caminhos_Cad.CampoExtra_Cad

				itens := Cad.CadToItens(itensFile)

				if nomeDaPasta := "SYSTEL-ARQUIVOS/"; len(itens) > 3 {
					balancas_import_checked, balancas_import_error = file.Passo2(itens, nomeDaPasta, balancas_import)
					progress += (75 * (len(balancas_import) - len(balancas_import_error)) / len(balancas_import))
					fmt.Fprintf(w, "data: {\"progress\": %d}\n\n", progress)
					flusher.Flush()
				}

			} else {
				arquivos := file.ReadMGVJson()
				itensFile := arquivos.Caminhos.Itens
				receitaFile := arquivos.Caminhos.Receita
				nutriFile := arquivos.Caminhos.Nutricional
				campoextra := arquivos.Caminhos.CampoExtra
				fornFile := arquivos.Caminhos.Fornecedor
				fracionaFile := arquivos.Caminhos.Fracionador
				taraFile := arquivos.Caminhos.Tara
				conservaFile := arquivos.Caminhos.Conservantes
				somente_preco_form := r.FormValue("somente_preco")
				somente_preco := somente_preco_form == "on"
				erroLeitura, arquivo, nomeDaPasta, dict_nutri, info, dict_forn, dict_aler, dict_fraciona, dict_conserva, dict_tara := file.Passo1(itensFile, receitaFile, nutriFile, fracionaFile, fornFile, taraFile, conservaFile, campoextra)
				if erroLeitura != nil {
					progress = -1
					fmt.Fprintf(w, "data: {\"progress\": %d}\n\n", progress)
					flusher.Flush()
					time.Sleep(time.Millisecond * 500)
					//http.Redirect(w, r, "/Leitura500", 404)
					return

				} else {
					progress += 25
					fmt.Fprintf(w, "data: {\"progress\": %d}\n\n", progress)
					flusher.Flush()
					balancas_import_checked, balancas_import_error = file.Passo2(arquivo, nomeDaPasta, balancas_import)
					progress += (25 * (len(balancas_import) - len(balancas_import_error)) / len(balancas_import))
					fmt.Fprintf(w, "data: {\"progress\": %d}\n\n", progress)
					flusher.Flush()

					balancas_import_status, err := json.Marshal(balancas_import_checked)
					if err != nil {
						fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
						flusher.Flush()
						return
					}
					fmt.Fprintf(w, "data: {\"step2\": %s}\n\n", balancas_import_status)
					flusher.Flush()

					time.Sleep(time.Millisecond * 500)
				}
				if err == nil && !somente_preco && len(balancas_import_checked) > 0 {
					balancas_import_status_extras, balancas_extras_error := file.Passo3(arquivo, nomeDaPasta, balancas_import_checked, dict_nutri, info, dict_forn, dict_aler, dict_fraciona, dict_conserva, dict_tara)
					progress += (25 * (len(balancas_import) - len(balancas_extras_error)) / len(balancas_import))
					fmt.Fprintf(w, "data: {\"progress\": %d}\n\n", progress)
					flusher.Flush()
					time.Sleep(time.Millisecond * 500)
					balancas_import_status, err := json.Marshal(balancas_import_status_extras)
					if err != nil {
						fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
						flusher.Flush()
						return
					}
					fmt.Fprintf(w, "data: {\"step3\": %s}\n\n", balancas_import_status)
					flusher.Flush()
				}
				if len(balancas_import_checked) > 0 && somente_preco {
					progress += 25
					fmt.Fprintf(w, "data: {\"progress\": %d}\n\n", progress)
					flusher.Flush()
					time.Sleep(time.Millisecond * 500)
				}
				// balancas_import_checked, balancas_import_error = file.EnviaParaBalancas(itensFile, receitaFile, nutriFile, fracionaFile, fornFile, taraFile, conservaFile, campoextra, balancas_import, somente_preco)
			}

			// // Envio de evento de progresso
			// progress = progress + (60 * (len(balancas_import) - len(balancas_import_error)) / len(balancas_import))

			fmt.Fprintf(w, "data: {\"progress\": %d}\n\n", progress)
			flusher.Flush()

			// Envio de evento de conclusão com balanças importadas
			balancasImportCheckedJSON, err := json.Marshal(balancas_import_checked)
			if err != nil {
				fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
				flusher.Flush()
				return
			}

			// Verifica se houve balanças não importadas e envia os erros
			var errorJSON []byte
			if len(balancas_import_checked) < len(balancas_import) {
				errorJSON, err = json.Marshal(balancas_import_error)
				if err != nil {
					fmt.Fprintf(w, "data: {\"error\": \"%s\"}\n\n", err.Error())
					flusher.Flush()
					return
				}
			}

			time.Sleep(time.Millisecond * 500)

			fmt.Fprintf(w, "data: {\"incomplete\": %s}\n\n", errorJSON)
			flusher.Flush()

			time.Sleep(time.Millisecond * 500)

			fmt.Fprintf(w, "data: {\"complete\": %s}\n\n", balancasImportCheckedJSON)
			flusher.Flush()

			time.Sleep(time.Millisecond * 500)
		}
	}
}

func LogHtml(w http.ResponseWriter, r *http.Request) {

}
