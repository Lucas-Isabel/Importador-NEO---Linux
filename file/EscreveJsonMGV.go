package file

import (
	"fmt"
)

type Arquivos struct {
	Formatacao struct {
		Tipo string `json:"tipo"`
	} `json:"formatacao"`
	Caminhos struct {
		Itens        string `json:"itens"`
		Receita      string `json:"receita"`
		Nutricional  string `json:"nutricional"`
		Fornecedor   string `json:"fornecedor"`
		Fracionador  string `json:"fracionador"`
		Tara         string `json:"tara"`
		CampoExtra   string `json:"campoextra"`
		Conservantes string `json:"conservantes"`
	} `json:"caminhos"`
}

func EditFilePath(tipo, nome_json_extensao, new_caminho string) Arquivos {
	var arquivos Arquivos
	caminho_json := fmt.Sprint(template, nome_json_extensao)
	fmt.Println(caminho_json)
	arquivos = ReadMGVJson()

	fmt.Println(arquivos)
	switch tipo {
	case "itens":
		arquivos.Caminhos.Itens = new_caminho
	case "receita":
		arquivos.Caminhos.Receita = new_caminho
	case "nutricional":
		arquivos.Caminhos.Nutricional = new_caminho
	case "conserva":
		arquivos.Caminhos.Conservantes = new_caminho
	case "fraciona":
		arquivos.Caminhos.Fracionador = new_caminho
	case "tara":
		arquivos.Caminhos.Tara = new_caminho
	case "fornecedor":
		arquivos.Caminhos.Fornecedor = new_caminho
	case "campo-extra-1":
		arquivos.Caminhos.CampoExtra = new_caminho
	default:
		fmt.Println(new_caminho)
	}

	err := writeJson(arquivos)
	if err != nil {
		fmt.Println(err)
	}
	arquivos_new, err := LoadArquivos(caminho_json)
	if err != nil {
		return arquivos
	}
	return arquivos_new
}
