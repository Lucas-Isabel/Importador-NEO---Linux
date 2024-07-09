package file

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func itensAnalyze(arq string, arq_new string) error {
	arquivoItensMgv, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer arquivoItensMgv.Close()

	arquivoMGV7, err := os.Create(arq_new)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer arquivoMGV7.Close()

	scanner := bufio.NewScanner(arquivoItensMgv)
	for scanner.Scan() {
		linha := scanner.Text()
		textoModifiCado := linha
		if len(linha) < 70 {
			textoModifiCado = strings.ReplaceAll(linha, "\n", " ")
			arquivoMGV7.WriteString(textoModifiCado)
		} else {
			arquivoMGV7.WriteString(textoModifiCado + "\n")
		}
		//textoModifiCado = strings.ReplaceAll(textoModifiCado, "�", " ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
