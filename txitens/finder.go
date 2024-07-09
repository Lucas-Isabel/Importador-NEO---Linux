package txitens

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

type Arquivos_Txitens struct {
	Formatacao struct {
		Tipo string `json:"tipo"`
	} `json:"formatacao_txitens"`
	Caminhos struct {
		Itens       string `json:"itens"`
		Receita     string `json:"receita"`
		Nutricional string `json:"nutricional"`
	} `json:"caminhos_txitens"`
}

func writeJson(arquivosJSON Arquivos_Txitens) error {
	arquivosJSONBytes, err := json.Marshal(arquivosJSON)
	if err != nil {
		fmt.Println(err)
		dialogOpen = false
		return err
	}
	// Cria e abre o arquivo Txitens.json para escrita
	jsonFilePath := "templates/Txitens.json"
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

	fmt.Printf("Arquivo Txitens.json salvo em: %s\n", jsonFilePath)

	dialogOpen = false
	return nil
}

func ReadTxitensJson() Arquivos_Txitens {
	// Ler o conteúdo do arquivo "arquivo.json"
	// Instanciar a struct file.Arquivos
	var arquivos Arquivos_Txitens

	jsonFile, err := ioutil.ReadFile("templates/Txitens.json")
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
	fmt.Println("Caminho para Itens:", arquivos.Caminhos.Itens)
	fmt.Println("Caminho para Receita:", arquivos.Caminhos.Receita)
	fmt.Println("Caminho para Nutricional:", arquivos.Caminhos.Nutricional)
	return arquivos
}

// Função para carregar os dados do arquivo JSON na struct file.Arquivos
func LoadArquivosTxitens(filePath string) (Arquivos_Txitens, error) {
	var arquivos Arquivos_Txitens

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

func EditFilePath(tipo, nome_com_extensao, new_caminho string) Arquivos_Txitens {
	var arquivos Arquivos_Txitens
	caminho_json := fmt.Sprint(template, nome_com_extensao)
	fmt.Println(caminho_json)
	arquivos = ReadTxitensJson()

	fmt.Println(arquivos)
	switch tipo {
	case "txitens_itens":
		arquivos.Caminhos.Itens = new_caminho
	case "receita":
		arquivos.Caminhos.Receita = new_caminho
	case "nutricional":
		arquivos.Caminhos.Nutricional = new_caminho
	default:
		fmt.Println(new_caminho)
	}

	err := writeJson(arquivos)
	if err != nil {
		fmt.Println(err)
	}
	arquivos_new, err := LoadArquivosTxitens(caminho_json)
	if err != nil {
		return arquivos
	}
	return arquivos_new
}
