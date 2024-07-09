package Cad

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	dialogOpen bool
	template   = "template/"
)

type Arquivos_Cad struct {
	Formatacao struct {
		Tipo string `json:"tipo"`
	} `json:"formatacao_Cad"`
	Caminhos_Cad struct {
		Itens_Cad      string `json:"itens"`
		Receita_Cad    string `json:"receita"`
		CampoExtra_Cad string `json:"campoextra"`
	} `json:"caminhos_Cad"`
}

func writeJson(arquivosJSON Arquivos_Cad) error {
	arquivosJSONBytes, err := json.Marshal(arquivosJSON)
	if err != nil {
		fmt.Println(err)
		dialogOpen = false
		return err
	}
	// Cria e abre o arquivo Cad.json para escrita
	jsonFilePath := "templates/Cad.json"
	jsonFile, err := os.Create(jsonFilePath)
	if err != nil {
		fmt.Println(err)
		dialogOpen = false
	}
	defer jsonFile.Close()

	// Escreve os bytes da estrutura JSON no arquivo
	_, err = jsonFile.Write(arquivosJSONBytes)
	if err != nil {
		fmt.Println(err)
		dialogOpen = false
		return err
	}

	fmt.Printf("Arquivo Cad.json salvo em: %s\n", jsonFilePath)

	dialogOpen = false
	return nil
}

func ReadCadJson() Arquivos_Cad {
	// Ler o conteúdo do arquivo "arquivo.json"
	// Instanciar a struct file.Arquivos
	var arquivos Arquivos_Cad

	jsonFile, err := ioutil.ReadFile("templates/Cad.json")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return arquivos
	}

	// Converter o JSON para a struct file.Arquivos
	err = json.Unmarshal(jsonFile, &arquivos)
	if err != nil {
		fmt.Println("Erro ao fazer unmarshal do JSON:", err)
		return arquivos
	}

	// Exibir os valores da struct file.Arquivos
	fmt.Println("Tipo de formatação:", arquivos.Formatacao.Tipo)
	fmt.Println("Caminho para Itens:", arquivos.Caminhos_Cad.Itens_Cad)
	fmt.Println("Caminho para Receita:", arquivos.Caminhos_Cad.Receita_Cad)
	fmt.Println("Caminho para Campo Extra:", arquivos.Caminhos_Cad.CampoExtra_Cad)
	return arquivos
}

// Função para carregar os dados do arquivo JSON na struct file.Arquivos
func LoadArquivosCad(filePath string) (Arquivos_Cad, error) {
	var arquivos Arquivos_Cad

	// Abrir o arquivo JSON
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return arquivos, err
	}
	defer jsonFile.Close()

	// Ler o conteúdo do arquivo
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return arquivos, err
	}

	// Deserializar o JSON na struct
	err = json.Unmarshal(byteValue, &arquivos)
	if err != nil {
		return arquivos, err
	}

	return arquivos, nil
}

func EditFilePath(tipo, nome_com_extensao, new_caminho string) Arquivos_Cad {
	var arquivos Arquivos_Cad
	caminho_json := fmt.Sprint(template, nome_com_extensao)
	fmt.Println(caminho_json)
	arquivos = ReadCadJson()

	fmt.Println(arquivos)
	switch tipo {
	case "Cad-itens":
		arquivos.Caminhos_Cad.Itens_Cad = new_caminho
	case "receita":
		arquivos.Caminhos_Cad.Receita_Cad = new_caminho
	case "campoextra":
		arquivos.Caminhos_Cad.CampoExtra_Cad = new_caminho
	default:
		fmt.Println(new_caminho)
	}

	err := writeJson(arquivos)
	if err != nil {
		fmt.Println(err)
	}
	arquivos_new, err := LoadArquivosCad(caminho_json)
	if err != nil {
		return arquivos
	}
	return arquivos_new
}
