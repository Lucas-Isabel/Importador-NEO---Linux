package file

import (
	"bufio"
	"fmt"
	"os"
)

func InfoAnalyze(arq string) map[string]string {
	infoTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer infoTxt.Close()

	dictInfo := make(map[string]string)
	scanner := bufio.NewScanner(infoTxt)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("linha:       ", line)
		if len(line) > 100 {
			dictInfo[line[0:6]] = caracterRemove(line[106:])
		} else if len(line) > 6 {
			dictInfo[line[0:6]] = caracterRemove(line)
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictInfo
}
