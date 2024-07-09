package models

import (
	"database/sql"
	"fmt"
	"math"
	"strconv"
)

func enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment int, hasjoinedcolumns, info_set_id, value string, existe bool, db *sql.DB) error {
	fmt.Println(plu, value, pos_column, aligment)
	_, err := db.Exec(`
	UPDATE product
	SET
		nut_info_set_id = '2'
	WHERE
		product_id = $1
`, plu)
	if err != nil {
		fmt.Println("Erro ao chamar a função merge_nutri: - Nutri do PLU: ", plu, "\nERRO: ", err)
		return err
	}

	fmt.Println("Função merge_product chamada com sucesso. valor: ", value)

	if !existe {
		err := ApenasUpdateNutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, info_set_id, value, db)
		if err != nil {
			fmt.Println(err)
			return err
		}
		return nil
	}

	_, err = db.Exec(`
   INSERT INTO nut_info_el_instance(
           nut_info_el_instance_id, 
           isactive, created, createdby, updated, updatedby, 
           nut_info_set_id, product_id, value, hasjoinedrows, 
           pos_row, pos_rowto, 
           hasjoinedcolumns, 
           pos_column, pos_columnto, 
           alignment)
   VALUES (
           get_uuid(), 
           'Y', LOCALTIMESTAMP, 'importador', LOCALTIMESTAMP, 'importador', 
           '2', $1, $2, 'N',
           $3,
           $4, 
           $5, 
           $6, $7, 
           $8
   );
`, plu, value, pos_row, 0, hasjoinedcolumns, pos_column, pos_columnto, aligment)
	if err != nil {
		fmt.Println("Erro ao chamar a função merge_nutri: - Nutri do PLU: ", plu, "\nERRO: ", err)
		return err
	}

	return nil
}

func ApenasUpdateNutri(plu, pos_column, pos_row, pos_columnto, aligment int, hasjoinedcolumns, info_set_id, value string, db *sql.DB) error {
	// Atualizar registro existente
	queryUpdate := `
UPDATE nut_info_el_instance
SET value = $4
WHERE product_id = $1 AND pos_row = $2 AND pos_column = $3
`
	_, err := db.Exec(queryUpdate, plu, pos_row, pos_column, value)
	if err != nil {
		fmt.Println("erro: ", err)
		return err
	}
	return nil
}

func MapNutriPlusPG(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query("SELECT product_id FROM nut_info_el_instance")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []int

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			fmt.Println(err)
		}
		ids = append(ids, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return NutriExistentes(ids), nil
}

func NutriExistentes(lista []int) map[int]bool {
	idMap := make(map[int]bool, len(lista))
	for _, id := range lista {
		idMap[id] = true
	}
	return idMap
}

func retornaInfoLinhaNutri(num_linha_nutri int) string {
	switch num_linha_nutri {
	case 1:
		return "PORÇÃO:"
	case 10:
		return "Gorduras saturadas (g)"
	case 11:
		return "Gordura trans (g)"
	case 12:
		return "Fibra alimentar"
	case 13:
		return "SÓDIO"
	case 2:
		return "Porção:"
	case 3:
		return "000 g"
	case 4:
		return "Valor energético (kcal)"
	case 5:
		return "Caboidratos totais (g)"
	case 6:
		return "Açúcares totais (g)"
	case 7:
		return "Açucares adicionados (g)"
	case 8:
		return "Proteinas (g)"
	case 9:
		return "Gorduras totais (g)"
	default:
		return ""
	}
}

func EnviaNutriPG2(codPlu string, valores_e_medida [13]string, existe bool, db *sql.DB) {
	fmt.Println("Valores e medidas: ", valores_e_medida)
	plu, _ := strconv.Atoi(codPlu)
	count := 0
	for _, v := range valores_e_medida {
		count++
		linha := retornaInfoLinhaNutri(count)
		fmt.Println(linha, count)
		hasjoinedcolumns := ""
		pos_column, pos_columnto, aligment := 0, 0, 0
		pos_row := count
		switch count {
		case 1:
			if value, _ := strconv.Atoi(v); value == 0 {
				v = "1"
			}
			hasjoinedcolumns = "Y"
			pos_column, pos_columnto, aligment = 1, 3, 0
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 2:
			// 2380	"2 Colher(es) de Sopa"	"N"	2	0	"Y"	2	3	0
			hasjoinedcolumns = "Y"
			pos_column, pos_columnto, aligment = 2, 3, 0
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 3:
			// 2380	"030g"					"N"	2	0	"Y"	1	1	0
			// 2380	"030g"					"N"	3	0	"N"	1	0	2
			// 2380	"%VD*"					"N"	3	0	"N"	2	0	2
			hasjoinedcolumns = "Y"
			pos_column, pos_columnto, aligment = 1, 1, 0
			enviar_unico_nutri(plu, pos_column, 2, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 2, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", "%VD", existe, db)
		case 4:
			// 2380	"6%"					"N"	4	0	"N"	2	0	2
			// 2380	"129"					"N"	4	0	"N"	1	0	2
			perc := TointPerc(v, 0.1)
			str_perc := fmt.Sprint(perc)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 2, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", str_perc, existe, db)

			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 5:
			// 2380	"0%"					"N"	5	0	"N"	2	0	2
			perc := TointPerc(v, 0.3)
			str_perc := fmt.Sprint(perc)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 2, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", str_perc, existe, db)
			//2380	"0"						"N"	5	0	"N"	1	0	2
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 6:
			// 2380	"0"						"N"	5	0	"N"	1	0	2
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 7:
			perc := TointPerc(v, 2)
			str_perc := fmt.Sprint(perc)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 2, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", str_perc, existe, db)
			// 2380	"0"						"N"	5	0	"N"	1	0	2
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 8:
			perc := TointPerc(v, 2)
			str_perc := fmt.Sprint(perc)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 2, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", str_perc, existe, db)
			// 2380	"0"						"N"	5	0	"N"	1	0	2
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 9:
			perc := TointPerc(v, 2)
			str_perc := fmt.Sprint(perc)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 2, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", str_perc, existe, db)
			// 2380	"0"						"N"	5	0	"N"	1	0	2
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 10:
			perc := TointPerc(v, 5)
			str_perc := fmt.Sprint(perc)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 2, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", str_perc, existe, db)
			// 2380	"0"						"N"	5	0	"N"	1	0	2
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 11:
			perc := TointPerc(v, 50)
			str_perc := fmt.Sprint(perc)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 2, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", str_perc, existe, db)
			// 2380	"0"						"N"	5	0	"N"	1	0	2
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 12:
			perc := TointPerc(v, 4)
			str_perc := fmt.Sprint(perc)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 2, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", str_perc, existe, db)
			// 2380	"0"						"N"	5	0	"N"	1	0	2
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		case 13:
			perc := TointPerc(v, 0.1)
			str_perc := fmt.Sprint(perc)
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 2, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", str_perc, existe, db)
			// 2380	"0"						"N"	5	0	"N"	1	0	2
			hasjoinedcolumns = "N"
			pos_column, pos_columnto, aligment = 1, 0, 2
			enviar_unico_nutri(plu, pos_column, pos_row, pos_columnto, aligment, hasjoinedcolumns, "2", v, existe, db)
		default:
			fmt.Print()
		}
	}
}

func TointPerc(str string, multp float64) string {
	valor, err := strconv.ParseFloat(str, 32)
	if err != nil {
		return "0"
	}

	valor = valor * multp
	valor = math.Round(valor)
	intValor := int(valor)
	str_perc := fmt.Sprint(intValor)
	str_perc = str_perc + " %"
	return str_perc
}

var processamentoMap = map[int]ProcessamentoFunc{
	1:  processaCaso1,
	2:  processaCaso2,
	3:  processaCaso3,
	4:  processaCaso4,
	5:  processaCaso5,
	6:  processaCaso6,
	7:  processaCaso7,
	8:  processaCaso8,
	9:  processaCaso9,
	10: processaCaso10,
	11: processaCaso11,
	12: processaCaso12,
	13: processaCaso13,
}

func EnviaNutriPG(codPlu string, valores_e_medida [13]string, existe bool, db *sql.DB) {
	fmt.Println("Valores e medidas: ", valores_e_medida)
	plu, _ := strconv.Atoi(codPlu)
	count := 0
	if plu == 7907 {
		fmt.Println("Teste")
	}
	for _, v := range valores_e_medida {
		count++
		if processo, ok := processamentoMap[count]; ok {

			processo(plu, count, v, existe, db)
		} else {
			fmt.Println("Processamento não encontrado para o caso:", count)
		}
	}
}

// Define um tipo para as funções de processamento
type ProcessamentoFunc func(plu, posRow int, v string, existe bool, db *sql.DB)

// Função para cada caso de processamento
func processaCaso1(plu, posRow int, v string, existe bool, db *sql.DB) {
	if value, _ := strconv.Atoi(v); value == 0 {
		v = "1"
	}
	enviar_unico_nutri(plu, 1, posRow, 3, 0, "Y", "2", v, existe, db)
}

func processaCaso2(plu, posRow int, v string, existe bool, db *sql.DB) {
	enviar_unico_nutri(plu, 2, posRow, 3, 0, "Y", "2", v, existe, db)
}

func processaCaso3(plu, posRow int, v string, existe bool, db *sql.DB) {
	enviar_unico_nutri(plu, 1, 2, 1, 0, "Y", "2", v, existe, db)
	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
	enviar_unico_nutri(plu, 2, posRow, 0, 2, "N", "2", "%VD", existe, db)
}

func processaCaso4(plu, posRow int, v string, existe bool, db *sql.DB) {
	perc := TointPerc(v, 0.05)
	strPerc := fmt.Sprint(perc)
	enviar_unico_nutri(plu, 2, posRow, 0, 2, "N", "2", strPerc, existe, db)

	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
}

func processaCaso5(plu, posRow int, v string, existe bool, db *sql.DB) {
	perc := TointPerc(v, 0.333333333)
	strPerc := fmt.Sprint(perc)
	enviar_unico_nutri(plu, 2, posRow, 0, 2, "N", "2", strPerc, existe, db)

	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
}

func processaCaso6(plu, posRow int, v string, existe bool, db *sql.DB) {
	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
}

func processaCaso7(plu, posRow int, v string, existe bool, db *sql.DB) {
	perc := TointPerc(v, 2)
	strPerc := fmt.Sprint(perc)
	enviar_unico_nutri(plu, 2, posRow, 0, 2, "N", "2", strPerc, existe, db)

	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
}

func processaCaso8(plu, posRow int, v string, existe bool, db *sql.DB) {
	perc := TointPerc(v, 2)
	strPerc := fmt.Sprint(perc)
	enviar_unico_nutri(plu, 2, posRow, 0, 2, "N", "2", strPerc, existe, db)

	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
}

func processaCaso9(plu, posRow int, v string, existe bool, db *sql.DB) {
	perc := TointPerc(v, (100.0 / 65.0))
	strPerc := fmt.Sprint(perc)
	enviar_unico_nutri(plu, 2, posRow, 0, 2, "N", "2", strPerc, existe, db)

	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
}

func processaCaso10(plu, posRow int, v string, existe bool, db *sql.DB) {
	perc := TointPerc(v, (100.0 / 20.0))
	strPerc := fmt.Sprint(perc)
	enviar_unico_nutri(plu, 2, posRow, 0, 2, "N", "2", strPerc, existe, db)

	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
}

func processaCaso11(plu, posRow int, v string, existe bool, db *sql.DB) {
	perc := TointPerc(v, 50)
	strPerc := fmt.Sprint(perc)
	enviar_unico_nutri(plu, 2, posRow, 0, 2, "N", "2", strPerc, existe, db)

	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
}

func processaCaso12(plu, posRow int, v string, existe bool, db *sql.DB) {
	perc := TointPerc(v, 4)
	strPerc := fmt.Sprint(perc)
	enviar_unico_nutri(plu, 2, posRow, 0, 2, "N", "2", strPerc, existe, db)

	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
}

func processaCaso13(plu, posRow int, v string, existe bool, db *sql.DB) {
	perc := TointPerc(v, 0.05)
	strPerc := fmt.Sprint(perc)
	enviar_unico_nutri(plu, 2, posRow, 0, 2, "N", "2", strPerc, existe, db)

	enviar_unico_nutri(plu, 1, posRow, 0, 2, "N", "2", v, existe, db)
}
