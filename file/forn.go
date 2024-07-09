package file

import (
	"bufio"
	"fmt"
	"os"
)

func fornAnalyze(arq string) map[string]string {
	fornTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer fornTxt.Close()

	dictFornecedor := make(map[string]string)
	scanner := bufio.NewScanner(fornTxt)
	for scanner.Scan() {
		line := scanner.Text()
		dictFornecedor[line[0:4]] = caracterRemove(line[104:217])
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictFornecedor
}
