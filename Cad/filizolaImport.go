package Cad

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/lucasbyte/go-clipse/file"
	"github.com/lucasbyte/go-clipse/models"
)

func CadToItens(arquivo_cad string) string {
	var itens string
	arquivo := strings.Replace(arquivo_cad, "\\", "//", -1)
	if _, err := os.Stat(arquivo); err == nil {
		dict_plu := itensCadAnalize(arquivo)
		err, itens = itensWriter(arquivo, dict_plu)
		if err != nil {
			return ""
		}
	}
	return itens
}

func CadImport(arquivo, recextra, rec_ass string, balancas []models.Balanca) {
	arquivo = strings.Replace(arquivo, "\\", "//", -1)
	recextra = strings.Replace(recextra, "\\", "//", -1)
	rec_ass = strings.Replace(rec_ass, "\\", "//", -1)

	if _, err := os.Stat(arquivo); err == nil {
		time.Sleep(1 * time.Second)

		dict_plu := itensCadAnalize(arquivo)
		dict_nutri := Nutri(arquivo)
		time.Sleep(1 * time.Second)

		itensWriter(arquivo, dict_plu)
		_, err_recextra := os.Stat(recextra)
		_, err_rec_ass := os.Stat(rec_ass)

		if err_rec_ass == err_recextra && err_recextra == nil {
			receitaWriter(recextra, rec_ass)
		}
		if _, err := os.Stat(arquivo); err == nil {
			nutriWriter(arquivo, dict_nutri)
		}
		time.Sleep(1 * time.Second)
		itensCadNew := "SYSTEL-ARQUIVOS/itensSystel-F.txt"
		ReceitaCadNew := "SYSTEL-ARQUIVOS/receitasSystel-F.txt"
		NutriCadNew := "SYSTEL-ARQUIVOS/NutriSystel-F.txt"

		nutri_dict := file.NutriAnalyse(NutriCadNew)

		info := make(map[string]string)

		if _, err := os.Stat(ReceitaCadNew); err == nil {
			info = file.InfoAnalyze(ReceitaCadNew)
		} else {
			fmt.Println("ERRO AO LER ARQUIVO")
		}
		for _, balanca := range balancas {
			file.EnviarPluSeparado(itensCadNew, balanca)
			EnviarInfoSeparadaCad(itensCadNew, nutri_dict, info, balanca)
		}
		fmt.Println("Pronto")
	}
}

func itensCadAnalize(arquivo string) map[string][]string {
	itensDict := make(map[string][]string)
	content, err := ioutil.ReadFile(arquivo)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return itensDict
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) >= 37 && line[6:7] != "0" {
			cod_plu := line[0:6]
			venda := line[6:7]
			descricao := line[7:29]
			valor := line[29:35]
			validade := line[35:38]
			plu := []string{cod_plu, venda, descricao, valor, validade}
			itensDict[cod_plu] = plu
			fmt.Printf("%s<-LINHA \n", itensDict[cod_plu][2])

		}
	}
	return itensDict
}

func Nutri(arquivo string) map[string][]string {
	nutriDict := make(map[string][]string)
	content, err := ioutil.ReadFile(arquivo)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return nutriDict
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) >= 79 && line[39:40] == "@" {
			cod_plu := line[0:6]
			peso := strings.TrimSpace(line[40:44])
			if !strings.ContainsAny(peso, "gG") || len(peso) < 3 {
				peso = "100"
			} else {
				peso = strings.ReplaceAll(peso, "g", "")
				peso = strings.ReplaceAll(peso, "G", "")
				peso = fmt.Sprintf("%03s", peso)
			}
			tipo := line[44:62]
			valores := line[75:169]
			valEn := valores[1:6]
			carb := valores[10:13]
			proten := valores[19:22]
			gordTo := valores[29:32]
			gordSa := valores[37:40]
			gordTr := valores[46:49]
			fibra := valores[56:59]
			sodio := valores[81:86]
			plu := []string{cod_plu, peso, tipo, valEn, carb, proten, gordSa, gordTr, gordTo, fibra, sodio}
			fmt.Printf("%s peso:%s tipo:%s valores:%s\nvalEn = %s\ngorSat = %s\ngordTr = %s\ncarb = %s\nproten = %s\ngordT = %s\nfibra = %s\nsodio = %s\n<-LINHA\n",
				cod_plu, peso, tipo, valores, valEn, gordSa, gordTr, carb, proten, gordTo, fibra, sodio)
			nutriDict[cod_plu] = plu
		}
	}
	return nutriDict
}

func receitaWriter(arquivo string, arquivo_2 string) {
	fmt.Println("HMM")
	receitas := make(map[string]string)
	extras := make(map[string]string)
	plus := make(map[string]bool)

	content, err := ioutil.ReadFile(arquivo)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "@") {
			break
		}
		if len(line) >= 36 {
			plu := line[12:18]
			plus[plu] = true
			fmt.Println(plu)
			rec := strings.TrimSpace(line[24:])
			receitas[plu] = rec
		}
	}

	content2, err := ioutil.ReadFile(arquivo_2)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}
	lines2 := strings.Split(string(content2), "\n")
	for _, line := range lines2 {
		if strings.Contains(line, "@") {
			break
		}
		if len(line) >= 36 {
			plu := line[12:18]
			plus[plu] = true
			fmt.Println(plu)
			rec := strings.TrimSpace(line[24:])
			extras[plu] = rec
		}
	}

	f, err := os.Create("SYSTEL-ARQUIVOS/receitasSystel-F.txt")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer f.Close()

	for plu := range plus {
		receita, ok := receitas[plu]
		if !ok {
			receita = ""
		}
		extra, ok := extras[plu]
		if !ok {
			extra = ""
		}
		if plu != "000000" && (receita != "" || extra != "") {
			f.WriteString(plu + receita + " " + extra + "\n")
		}
	}
}

func itensWriter(arquivo string, dict_plu map[string][]string) (error, string) {
	plus := make(map[string]bool)
	caminho_itens := "SYSTEL-ARQUIVOS/itensSystel-F.txt"
	// Ler o arquivo para obter os PLUs únicos
	content, err := ioutil.ReadFile(arquivo)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return err, ""
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) >= 7 && line[6:7] != "0" {
			cod_plu := line[0:6]
			plus[cod_plu] = true
		}
	}

	// Criar e escrever no arquivo itensSystel.txt
	f, err := os.Create(caminho_itens)
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return err, ""
	}
	defer f.Close()

	for plu := range plus {
		venda := "0"
		if dict_plu[plu][1] == "U" || dict_plu[plu][1] == "u" {
			venda = "1"
		}
		cod_plu := dict_plu[plu][0]
		desc := strings.ReplaceAll(dict_plu[plu][2], "Á", "A")
		desc = strings.ReplaceAll(desc, "Ç", "C")
		desc = strings.ReplaceAll(desc, "Ã", "A")
		valor := dict_plu[plu][3]
		valid := dict_plu[plu][4]

		line := fmt.Sprintf("01%s%s%s%s%s", venda, cod_plu, valor, valid, desc)
		if len(desc) < 50 {
			line += strings.Repeat(" ", 50-len(desc))
		}
		line += fmt.Sprintf("%s0000%s11000000000000000000000000000000000000000000000000000000000000000000000|01|                                                                      0000000000000000000000000||0||0000000000000000000000\n", cod_plu, cod_plu)
		f.WriteString(line)
	}
	return nil, caminho_itens
}

func nutriWriter(arquivo string, dict_nutri map[string][]string) {
	plus := make(map[string]bool)

	content, err := ioutil.ReadFile(arquivo)
	if err != nil {
		fmt.Println("Erro ao ler o arquivo:", err)
		return
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if len(line) >= 6 {
			cod_plu := line[0:6]
			plus[cod_plu] = true
		}
	}

	f, err := os.Create("SYSTEL-ARQUIVOS/nutriSystel-F.txt")
	if err != nil {
		fmt.Println("Erro ao criar o arquivo:", err)
		return
	}
	defer f.Close()

	for plu := range plus {
		valEn := dict_nutri[plu][3]
		carb := dict_nutri[plu][4]
		proten := dict_nutri[plu][5]
		gordSa := dict_nutri[plu][6]
		gordTr := dict_nutri[plu][7]
		gordTo := dict_nutri[plu][8]
		fibra := dict_nutri[plu][9]
		sodio := dict_nutri[plu][10]

		tipo := dict_nutri[plu][2]
		dec := "0"
		if strings.Contains(tipo, "1/2") || strings.Contains(tipo, "1 / 2") {
			dec = "3"
		} else if strings.Contains(tipo, "2/3") || strings.Contains(tipo, "2 / 3") {
			dec = "4"
		}

		qntd := "01"
		peso := dict_nutri[plu][1]

		line := "N" + plu + "000000000000000000000000000000000000000000|0" + "001" + peso + "0" + qntd + dec + "16" +
			valEn + carb + "000000" + proten + gordTo + gordSa + gordTr + fibra + sodio + "00000000000\n"

		if plu != "000000" {
			f.WriteString(line)
		}
	}
}
