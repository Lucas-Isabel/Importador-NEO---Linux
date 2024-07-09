package file

import (
	"fmt"
	"os"
	"strings"
)

func VerificaPasta() string {
	nomePasta := "SYSTEL-ARQUIVOS"

	// Verifica se a pasta existe
	if _, err := os.Stat(nomePasta); os.IsNotExist(err) {
		// Cria a pasta se não existir
		err := os.Mkdir(nomePasta, 0755)
		if err != nil {
			fmt.Println("Erro ao criar a pasta:", err)
			return ""
		}
		fmt.Printf("A pasta \"%s\" foi criada com sucesso.\n", nomePasta)
	} else {
		fmt.Printf("A pasta \"%s\" já existe.\n", nomePasta)
	}
	outraPasta()
	return nomePasta
}

func ArrumaExtensão(Filename string, format string, Newformat string) string {
	arquivo := strings.Replace(Filename, format, Newformat, -1)
	return arquivo
}

func outraPasta() string {
	nomePasta := "templates"

	// Verifica se a pasta existe
	if _, err := os.Stat(nomePasta); os.IsNotExist(err) {
		// Cria a pasta se não existir
		err := os.Mkdir(nomePasta, 0755)
		if err != nil {
			fmt.Println("Erro ao criar a pasta:", err)
			return ""
		}
		fmt.Printf("A pasta \"%s\" foi criada com sucesso.\n", nomePasta)
	} else {
		fmt.Printf("A pasta \"%s\" já existe.\n", nomePasta)
	}
	return nomePasta
}

// func fileExists(filename string) bool {
// 	_, err := os.Stat(filename)
// 	return err == nil
// }

// func obterDataHoraCriacao(arquivo string) (string, error) {
// 	// Obtém informações do arquivo
// 	info, err := os.Stat(arquivo)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Obtém a data e hora de criação
// 	dataHoraCriacao := info.ModTime().Format(time.RFC3339)

// 	return dataHoraCriacao, nil
// }
