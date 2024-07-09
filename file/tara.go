package file

import (
	"bufio"
	"fmt"
	"os"
)

func taraAnalyze(arq string) map[string]float64 {
	taraTxt, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer taraTxt.Close()

	dictTara := make(map[string]float64)
	scanner := bufio.NewScanner(taraTxt)
	for scanner.Scan() {
		line := scanner.Text()
		mgv5 := false
		mgv6 := false
		mgv7 := false
		var tara float64
		if len(line) > 10 {
			mgv5 = saoNumeros(line[:10]) && !saoNumeros(line[10:])
			mgv6 = saoNumeros(line[:11])
			mgv7 = !(saoNumeros(line[:1])) && saoNumeros(line[1:6])
		}
		if mgv7 {
			line = line[1:]

			taraStr := line[5:11]
			tara, err = stringTofloat(taraStr, 1000)
			if err != nil {
				print(err)
			}
		} else if mgv6 {
			taraStr := line[4:10]
			tara, err = stringTofloat(taraStr, 1000)
			if err != nil {
				print(err)
			}
		} else if mgv5 {
			taraStr := line[4:10]
			tara, err = stringTofloat(taraStr, 1000)
			if err != nil {
				print(err)
			}
		}

		dictTara[line[0:4]] = tara
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}

	return dictTara
}
