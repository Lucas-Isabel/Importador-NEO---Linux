package file

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	dialogOpen bool
	template   = "templates/"
)

func init() {
	nomePasta := "templates"
	itensmgv_json := "templates/itensmgv.json"
	Cad_json := "templates/Cad.json"
	txitens_json := "templates/txitens.json"
	formato_json := "templates/formato.json"

	jsons := []string{
		itensmgv_json,
		Cad_json,
		txitens_json,
		formato_json,
	}

	// Verifica se a pasta existe
	if _, err := os.Stat(nomePasta); os.IsNotExist(err) {
		// Cria a pasta se não existir
		err := os.Mkdir(nomePasta, 0755)
		if err != nil {
			fmt.Println("Erro ao criar a pasta:", err)
		}
		fmt.Printf("A pasta \"%s\" foi criada com sucesso.\n", nomePasta)
	} else {
		fmt.Printf("A pasta \"%s\" já existe.\n", nomePasta)
	}
	for c, json := range jsons {
		if _, err := os.Stat(json); os.IsNotExist(err) {
			jsonString := `{"formatacao":{"tipo":"itensmgv"},"caminhos":{"itens":"","receita":"","nutricional":"","fornecedor":"","fracionador":"","tara":"","campoextra":"","conservantes":""}}`

			if c == 2 {
				jsonString = `{"formatacao_txitens":{"tipo":"txitens"},"caminhos_txitens":{"itens":"","receita":"","nutri":""}}`
			}
			if c == 3 {
				jsonString = `{"formatacao_Cad":{"tipo":"Cad"},"caminhos_Cad":{"itens":"","receita":"","campoextra":""}}`
			}

			// Criar ou abrir o arquivo
			file, err := os.Create(json)
			if err != nil {
				fmt.Println("Erro ao criar o arquivo:", err)
				return
			}
			defer file.Close()

			// Escrever a string JSON no arquivo
			_, err = file.WriteString(jsonString)
			if err != nil {
				fmt.Println("Erro ao escrever no arquivo:", err)
				return
			}

			fmt.Println("Dados JSON escritos com sucesso no arquivo 'data.json'")

		}
	}
}

type FormatacaoSelect struct {
	Tipo string `json:"tipo"`
}

type Dados struct {
	Formatacao FormatacaoSelect `json:"formatacao"`
}

func EscreverJSON(tipo string) error {
	// Convertendo a struct para JSON
	formatacao := FormatacaoSelect{
		Tipo: tipo,
	}
	jsonData, err := json.Marshal(map[string]FormatacaoSelect{"formatacao": formatacao})
	if err != nil {
		return fmt.Errorf("erro ao converter para JSON: %v", err)
	}

	// Criando o arquivo JSON
	file, err := os.Create("templates/formato.json")
	if err != nil {
		return fmt.Errorf("erro ao criar o arquivo: %v", err)
	}
	defer file.Close()

	// Escrevendo no arquivo
	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo: %v", err)
	}

	return nil
}

// LerJSON lê a estrutura FormatacaoSelect de um arquivo JSON
func LerJSON() (FormatacaoSelect, error) {
	// Abrindo o arquivo JSON
	file, err := os.Open("templates/formato.json")
	if err != nil {
		return FormatacaoSelect{}, fmt.Errorf("erro ao abrir o arquivo: %v", err)
	}
	defer file.Close()

	// Decodificando o JSON para a struct
	var data map[string]FormatacaoSelect
	err = json.NewDecoder(file).Decode(&data)
	if err != nil {
		return FormatacaoSelect{}, fmt.Errorf("erro ao decodificar o JSON: %v", err)
	}

	return data["formatacao"], nil
}

func LerTipoJson() string {
	filePath := "templates/formato.json"

	// Abrir o arquivo
	file, err := os.Open(filePath)
	if err != nil {
		return ""
	}
	defer file.Close()

	// Ler o conteúdo do arquivo
	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		return ""
	}

	// Variável para armazenar o JSON decodificado
	var data Dados

	// Decodificar o JSON
	err = json.Unmarshal(fileContent, &data)
	if err != nil {
		return ""
	}

	// Imprimir o valor do campo "tipo"
	return (data.Formatacao.Tipo)
}

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

func writeJson(arquivosJSON Arquivos) error {
	arquivosJSONBytes, err := json.Marshal(arquivosJSON)
	if err != nil {
		fmt.Println(err)
		dialogOpen = false
		return err
	}
	// Cria e abre o arquivo itensmgv.json para escrita
	jsonFilePath := "templates/itensmgv.json"
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

	fmt.Printf("Arquivo itensmgv.json salvo em: %s\n", jsonFilePath)

	dialogOpen = false
	return nil
}

func ReadMGVJson() Arquivos {
	// Ler o conteúdo do arquivo "arquivo.json"
	// Instanciar a struct Arquivos
	var arquivos Arquivos

	jsonFile, err := ioutil.ReadFile("templates/itensmgv.json")
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return arquivos
	}

	// Converter o JSON para a struct Arquivos
	err = json.Unmarshal(jsonFile, &arquivos)
	if err != nil {
		fmt.Println("Erro ao fazer unmarshal do JSON:", err)
		return arquivos
	}

	// Exibir os valores da struct Arquivos
	fmt.Println("Tipo de formatação:", arquivos.Formatacao.Tipo)
	fmt.Println("Caminho para Itens:", arquivos.Caminhos.Itens)
	fmt.Println("Caminho para Receita:", arquivos.Caminhos.Receita)
	fmt.Println("Caminho para Informações Nutricionais:", arquivos.Caminhos.Nutricional)
	fmt.Println("Caminho para Fornecedor:", arquivos.Caminhos.Fornecedor)
	fmt.Println("Caminho para Fracionador:", arquivos.Caminhos.Fracionador)
	fmt.Println("Caminho para Tara:", arquivos.Caminhos.Tara)
	fmt.Println("Caminho para Campo Extra:", arquivos.Caminhos.CampoExtra)
	fmt.Println("Caminho para Conservantes:", arquivos.Caminhos.Conservantes)
	return arquivos
}

// Função para carregar os dados do arquivo JSON na struct Arquivos
func LoadArquivos(filePath string) (Arquivos, error) {
	var arquivos Arquivos

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
