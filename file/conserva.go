package file

import (
	"bufio"
	"fmt"
	"os"
)

func conservaAnalyze(arq string) map[string]string {
	conservaTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer conservaTxt.Close()

	dictConserva := make(map[string]string)
	scanner := bufio.NewScanner(conservaTxt)
	for scanner.Scan() {
		line := scanner.Text()
		key := line[0:4]
		value := caracterRemove(line[104:])
		if key != "0003" {
			value = caracterRemove(line[4:40])
		}
		dictConserva[key] = value
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictConserva
}
