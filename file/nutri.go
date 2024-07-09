package file

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Newnutri(line string, codNutriMGVArray, codNutriArray []string) string {
	if len(line) > 10 {
		codNutri := line[1:7]
		boo := true
		if len(line) > 106 {
			boo = line[7:110] != "000000000000000000000000000000000000000000|000000000000000000000000000000000000000000000000000000000000"
		}
		if len(line) < 50 {
			line = line[0:49] + "\n"
			line += "|" + strings.Repeat("0", 3)
			porcao := line[7:11]
			if parseInt(porcao) <= 0 {
				porcao = "0100"
			}
			line += porcao + "0" + line[12:26]
			line += strings.Repeat("0", 6) + line[26:50]
			line += strings.Repeat("0", 9)
		} else {
			fmt.Println(line)
			line = strings.ReplaceAll(line, "|", "0000|")
		}

		if parseInt(line[61:63]) > 28 {
			line = line[:61] + "16" + line[63:]
		}

		if containsTo(codNutriMGVArray, codNutri) && !containsTo(codNutriArray, codNutri) && boo {
			//codNutriArray = append(codNutriArray, codNutri)
			fmt.Println("oi", line)
			return fmt.Sprint(line)
			//fmt.Fprint(nutriFile, line)
		} else {
			fmt.Println(codNutri, boo)
		}
	}
	return ""
}

func NutriAnalyse(arq string) map[string][13]string {
	//var valores [13]string
	campNutri, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer campNutri.Close()

	dictNutri := make(map[string][13]string)
	scanner := bufio.NewScanner(campNutri)
	for scanner.Scan() {

		var valores [13]string
		line := scanner.Text()
		fmt.Println(line)

		if len(line) > 50 {
			valores = parseLine(line)
			dictNutri[line[1:7]] = valores
		}

	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		return nil
	}
	//valoresDoPlu = valoresDoPlu.append(valoresDoPlu, valores)
	return dictNutri

}

func parseLine(line string) [13]string {
	var item [13]string

	fmt.Println("Dados do arquivo: ", line)

	// Extraindo os campos da linha
	plu := line[1:7]
	fmt.Println(plu)
	line = line[53:]
	//|0000100005000034106500000001100410160000390148000000000000000000000000000000000000000000
	if plu == "001005" {
		fmt.Println(plu)
	}
	quantidadePorEmbalagem_v, _ := strconv.Atoi(line[2:5])
	quantidadePorPorcao_v, _ := strconv.Atoi(line[5:8])
	carboidratos_v, _ := stringTofloat(line[18:22], 10)
	acucares_to_v, _ := stringTofloat(line[22:25], 10)
	acucares_ad_v, _ := stringTofloat(line[25:28], 10)
	proteinas_v, _ := stringTofloat(line[28:31], 10)
	gorduras_totais_v, _ := stringTofloat(line[31:34], 10)
	gordura_saturada_v, _ := stringTofloat(line[34:37], 10)
	gordura_trans_v, _ := stringTofloat(line[37:40], 10)
	fibra_v, _ := stringTofloat(line[40:43], 10)
	sodio_v, _ := stringTofloat(line[43:48], 10)

	quantidadePorEmbalagem := fmt.Sprint(quantidadePorEmbalagem_v)
	quantidadePorPorcao := fmt.Sprint(quantidadePorPorcao_v)
	carboidratos := fmt.Sprint(carboidratos_v)
	acucares_to := fmt.Sprint(acucares_to_v)
	acucares_ad := fmt.Sprint(acucares_ad_v)
	proteinas := fmt.Sprint(proteinas_v)
	gorduras_totais := fmt.Sprint(gorduras_totais_v)
	gordura_saturada := fmt.Sprint(gordura_saturada_v)
	gordura_trans := fmt.Sprint(gordura_trans_v)
	fibra := fmt.Sprint(fibra_v)
	sodio := fmt.Sprint(sodio_v)
	parte_inteira := line[9:11]
	parte_decimal := line[11:12]

	parcionamento := "0"

	valor_energetico_v, _ := strconv.Atoi(line[14:18])

	valor_energetico := fmt.Sprint(valor_energetico_v)

	parte_decimal_int, err := strconv.Atoi(parte_decimal)
	if err != nil {
		fmt.Println(err)
	}
	parte_inteira_int, err := strconv.Atoi(parte_inteira)
	if err != nil {
		fmt.Println(err)
	}
	valores := []string{
		"", "1/4", "1/3", "1/2", "2/3", "3/4",
	}
	if (parte_decimal_int) < len(valores) {
		parcionamento = fmt.Sprint(parte_inteira_int, " ", valores[parte_decimal_int])
	}
	if parte_inteira_int == 0 {
		parcionamento = fmt.Sprint(valores[parte_decimal_int])
	}
	medida_caseira := "Porção(ões)"

	medidas := []string{
		"Colher(es) de Sopa",
		"Colher(es) de Café",
		"Colher(es) de Chá",
		"Xícara(s)",
		"De Xícara(s)",
		"Unidade(s)",
		"Pacote(s)",
		"Fatia(s)",
		"Fatia(s) Fina(s)",
		"Pedaço(s)",
		"Folha(s)",
		"Pão(es)",
		"Biscoito(s)",
		"Bisnaguinha(s)",
		"Disco(s)",
		"Copo(s)",
		"Porção(ões)",
		"Tablete(s)",
		"Sache(s)",
		"Almôndega(s)",
		"Bife(s)",
		"Filé(s)",
		"Concha(s)",
		"Bala(s)",
		"Prato(s) Fundo(s)",
		"Pitada(s)",
		"Lata(s)",
		"Xícara de Chá",
		"Prato raso",
	}
	if value, _ := strconv.Atoi(line[12:14]); value <= len(valores) {
		medida_caseira = fmt.Sprint(medidas[value])
	}

	unidade_file := fmt.Sprint(strconv.Atoi(line[9:10]))
	unidades := []string{
		" g", " ml",
	}
	if value, _ := strconv.Atoi(unidade_file); value <= len(unidades) {
		quantidadePorPorcao = quantidadePorPorcao + fmt.Sprint(unidades[value])
	}

	// 	1  "PORÇÃO"
	// 10 "Gorduras saturadas (g)"
	// 11 "Gordura trans (g)"
	// 12 "Fibra alimentar"
	// 13 "SÓDIO"
	// 2  "Porção:"
	// 3  "000 g"
	// 4  "Valor energético (kcal)"
	// 5  "Caboidratos totais (g)"
	// 6  "Açúcares totais (g)"
	// 7  "Açucares adicionados (g)"
	// 8  "Proteinas (g)"
	// 9  "Gorduras totais (g)"

	item[12] = sodio
	item[0] = quantidadePorEmbalagem
	item[2] = quantidadePorPorcao
	item[1] = fmt.Sprintf("(%s %s)", parcionamento, medida_caseira)
	item[3] = valor_energetico
	item[4] = carboidratos
	item[5] = acucares_to
	item[6] = acucares_ad
	item[7] = proteinas
	item[8] = gorduras_totais
	item[9] = gordura_saturada
	item[10] = gordura_trans
	item[11] = fibra

	return item
}
