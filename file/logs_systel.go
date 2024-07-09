package file

import (
	"fmt"
	"os"
)

func adcionarMensagem(mensagem string, add string) string {
	new := mensagem + "\n" + add + "\n"
	return new
}

func logToFile(file, text string) {
	log, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer log.Close()

	if _, err := log.WriteString(text + "\n"); err != nil {
		fmt.Println(err)
	}
}
