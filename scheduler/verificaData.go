package scheduler

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/djherbis/times"
	"github.com/lucasbyte/go-clipse/Cad"
	"github.com/lucasbyte/go-clipse/file"
	"github.com/lucasbyte/go-clipse/txitens"
)

func Atualizou2() (string, bool, error) {
	var Ehigual bool = false
	var dataEhora string
	var criacaoDoarquivo string = ""
	itensFile := ""
	caminho := "SYSTEL-ARQUIVOS/output.txt"
	tipo := file.LerTipoJson()
	upper := strings.ToUpper(tipo)
	switch upper {
	case "TXITENS":
		arquivos := txitens.ReadTxitensJson()
		itensFile = arquivos.Caminhos.Itens
	case "CAD":
		arquivos := Cad.ReadCadJson()
		itensFile = arquivos.Caminhos_Cad.Itens_Cad
	case "ITENSMGV":
		arquivos := file.ReadMGVJson()
		itensFile = arquivos.Caminhos.Itens
	}

	if _, err := os.Stat(itensFile); err == nil {
		dataEhora, _ = obterDataHoraCriacao(itensFile)

	}

	if _, err := os.Stat(itensFile); err != nil {
		return itensFile, false, err
	}

	if _, err := os.Stat(caminho); err == nil {
		file, err := os.Open(caminho)
		if err != nil {
			fmt.Println("Erro ao abrir o arquivo:", err)
		}
		// Criar um novo leitor de buffer
		reader := bufio.NewReader(file)

		// Ler uma linha do arquivo
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// Caso de fim de arquivo, possivelmente sem quebra de linha
				fmt.Println(line)
			} else {
				fmt.Println("Erro ao ler a linha:", err)
			}
		} else {
			// Imprimir a linha lida
			fmt.Println(line)
		}
		criacaoDoarquivo = strings.TrimSpace(line)
	}

	Ehigual = criacaoDoarquivo == dataEhora

	fmt.Println(Ehigual)

	if Ehigual {
		return itensFile, false, nil
	}

	file, err := os.Create("SYSTEL-ARQUIVOS/output.txt")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return itensFile, false, err
	}
	defer file.Close()

	// Escrever uma linha no arquivo
	fmt.Fprintln(file, dataEhora)

	return itensFile, true, nil
}

func obterDataHoraCriacao(arquivo string) (string, error) {

	// Obtém informações de tempo avançadas
	t, err := times.Stat(arquivo)
	if err != nil {
		return "", err
	}

	// Obtém a data e hora de criação
	dataHoraCriacao := t.ChangeTime().Format("2006-01-02 15:04:05")
	return dataHoraCriacao, nil
}

func Atualizou() (string, bool, error) {
	var Ehigual bool = false
	itensFile := ""
	caminho := "SYSTEL-ARQUIVOS/output.txt"
	new := "SYSTEL-ARQUIVOS/comparable.txt"
	tipo := file.LerTipoJson()
	upper := strings.ToUpper(tipo)
	switch upper {
	case "TXITENS":
		arquivos := txitens.ReadTxitensJson()
		itensFile = arquivos.Caminhos.Itens
	case "CAD":
		arquivos := Cad.ReadCadJson()
		itensFile = arquivos.Caminhos_Cad.Itens_Cad
	case "ITENSMGV":
		arquivos := file.ReadMGVJson()
		itensFile = arquivos.Caminhos.Itens
	}

	if _, err := os.Stat(itensFile); err != nil {
		return itensFile, false, err
	}

	if _, err := os.Stat(caminho); err == nil {
		file, err := os.Open(caminho)
		if err != nil {
			fmt.Println("Erro ao abrir o arquivo:", err)
		}
		// Criar um novo leitor de buffer
		reader := bufio.NewReader(file)

		// Ler uma linha do arquivo
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// Caso de fim de arquivo, possivelmente sem quebra de linha
				fmt.Println(line)
			} else {
				fmt.Println("Erro ao ler a linha:", err)
			}
		} else {
			// Imprimir a linha lida
			fmt.Println(line)
		}
	}

	copyFile(itensFile, new)

	Ehigual, err := compareFiles(new, caminho)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(Ehigual)

	if Ehigual {
		return itensFile, false, nil
	}

	err = copyFile(new, caminho)
	if err != nil {
		fmt.Println(err)
	}
	// Escrever uma linha no arquivo

	return itensFile, true, nil
}

func SetAuto(bole bool) {
	if bole {
		file, err := os.Create("SYSTEL-ARQUIVOS/scheduler.txt")
		if err != nil {
			fmt.Println("Erro ao criar o arquivo:", err)
			return
		}
		defer file.Close()
		copyFile("templates/formato.json", "SYSTEL-ARQUIVOS/scheduler.txt")
	} else {

		file, err := os.Create("SYSTEL-ARQUIVOS/scheduler.txt")
		if err != nil {
			fmt.Println("Erro ao criar o arquivo:", err)
			return
		}
		defer file.Close()
		file.WriteString("none")
	}
}

func ReadAuto() (bolean bool) {
	bolean = false
	caminho := "SYSTEL-ARQUIVOS/scheduler.txt"
	line := ""
	if _, err := os.Stat(caminho); err != nil {
		SetAuto(false)
		return
	}

	if _, err := os.Open(caminho); err == nil {
		file, err := os.Open(caminho)
		if err != nil {
			fmt.Println("Erro ao abrir o arquivo:", err)
		}
		// Criar um novo leitor de buffer
		reader := bufio.NewReader(file)

		// Ler uma linha do arquivo
		line, err = reader.ReadString('\n')
		if err != nil {
			if err.Error() == "EOF" {
				// Caso de fim de arquivo, possivelmente sem quebra de linha
				fmt.Println(line)
			} else {
				fmt.Println("Erro ao ler a linha:", err)
			}
		} else {
			// Imprimir a linha lida
			fmt.Println(line)
		}
	}
	bolean = !(line == "none")
	return
}

// Função para ler o conteúdo de um arquivo e retorná-lo como uma slice de bytes
func readFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ioutil.ReadAll(file)
}

// Função para comparar o conteúdo de dois arquivos
func compareFiles(file1Path, file2Path string) (bool, error) {
	content1, err := readFile(file1Path)
	if err != nil {
		return false, err
	}

	content2, err := readFile(file2Path)
	if err != nil {
		return false, err
	}

	return bytes.Equal(content1, content2), nil
}

// Função para copiar um arquivo de uma pasta para outra
func copyFile(src, dst string) (err error) {
	// Abrir o arquivo de origem
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Criar o arquivo de destino
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	// Copiar o conteúdo do arquivo de origem para o arquivo de destino
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}
	return
}
