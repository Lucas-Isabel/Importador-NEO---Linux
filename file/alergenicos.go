package file

import (
	"bufio"
	"fmt"
	"os"
)

func alergiaAnalyze(arq string) map[string]string {
	campTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer campTxt.Close()

	dictAler := make(map[string]string)
	scanner := bufio.NewScanner(campTxt)
	for scanner.Scan() {
		line := scanner.Text()
		dictAler[line[0:4]] = caracterRemove(line[104:])
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictAler

}
