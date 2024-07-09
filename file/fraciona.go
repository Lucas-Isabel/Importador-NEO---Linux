package file

import (
	"bufio"
	"fmt"
	"os"
)

func fracionaAnalyze(arq string) map[string]string {
	fracionaTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer fracionaTxt.Close()

	dictFraciona := make(map[string]string)
	scanner := bufio.NewScanner(fracionaTxt)
	for scanner.Scan() {
		line := scanner.Text()
		dictFraciona[line[0:4]] = caracterRemove(line[104:])
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictFraciona
}
