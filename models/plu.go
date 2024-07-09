package models

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
	"unicode/utf8"
)

type Plu struct {
	Id                          int
	Ativo                       rune
	Criaco                      time.Time
	Criado_por                  string
	Atualizado_em               time.Time
	Erp_codigo                  int
	Nome                        string
	Descrição                   string
	Departamento_id             int
	Grupo_id                    int
	Label_format_w              string
	Label_format_u              string
	Barcode_item_code_flag      string
	Barcode_item_code           int
	Print_pimary_barcode        string
	Primary_barcode_type        int
	Primary_bacode_flag_data    string
	Print_secondary_barcode     string
	Secondary_barcode_type      int
	Secondary_bacode_flag_data  string
	Sell_by_date_time           int
	Sell_by_date_source         int
	Sell_by_date                int
	Sell_time_source            int
	Sell_time                   int
	Print_packed_time           rune
	Package_date_source         int
	Package_date                int
	Packge_date_format          string //"dd-mm-YYYY"
	Print_used_by_date          rune
	Used_by_date_source         int
	Used_by_date                int
	Unit_price_selection        int
	Tare                        float32
	Perc                        float32
	Quant                       int
	Upc                         string
	Uom_id                      string //if venda == peso = "1" else peso == "2" = uni
	Discount_schema_id          string
	Isstocked                   rune    //N
	Isirradiated                rune    //N
	Isbom                       rune    //N
	Stock_min                   float32 //0.00000
	Stock_max                   float32 //0.00000
	Versionno                   string
	Discontinued                rune //N
	DiscontinuedBy              time.Time
	Tax_id                      string //@
	Atrribute                   int    //if venda == peso = 1 else peso == 2 = uni
	Advertising_id              string
	Sto_temp_limit1             float32
	Sto_temp_limit2             float32
	Ingredients                 string //limite 2000 chars
	Ingredients_id              int    //integer autoincrement
	Preservation_info           string //limite 2000 chars
	Preservation_info_id        string
	Coupled_product             int
	Pack_indiCador              int
	Extra_field                 string //limite 2000 char
	Extra_field2                string //limite 2000 char
	Image_id                    string //limite 32 char
	Icon_id                     string //limite 32 char
	Lot                         string //20 char
	Born_country                string // 32 char
	Fatten_country              string // 32 char
	Origin_countrt              string // 32 char
	Manufacturer_id             string // 32 char
	Packer_id                   string //32 char
	Distribuidor_id             string // 32 char
	Importer_id                 string
	Exporter_id                 string
	Vendor_id                   string
	Cutting_hall_id             string
	Slaughter_house_id          string
	Supplier_id                 string
	Nut_info_set_id             string // 2
	Print_fix_primary_barcode   rune   // N
	Print_fix_secondary_barcode string // N
}

func NovoPLU(cod, departamento, cod_erp int, descrição string) Plu {
	return Plu{
		Id:                          cod,
		Ativo:                       'Y',
		Criaco:                      time.Now(),
		Criado_por:                  "",
		Atualizado_em:               time.Now(),
		Erp_codigo:                  cod_erp,
		Nome:                        descrição,
		Descrição:                   "",
		Departamento_id:             departamento,
		Grupo_id:                    departamento,
		Label_format_w:              "",
		Label_format_u:              "",
		Barcode_item_code_flag:      "",
		Barcode_item_code:           0,
		Print_pimary_barcode:        "",
		Primary_barcode_type:        0,
		Primary_bacode_flag_data:    "",
		Print_secondary_barcode:     "",
		Secondary_barcode_type:      0,
		Secondary_bacode_flag_data:  "",
		Sell_by_date_time:           0,
		Sell_by_date_source:         0,
		Sell_by_date:                0,
		Sell_time_source:            0,
		Sell_time:                   0,
		Print_packed_time:           'N',
		Package_date_source:         0,
		Package_date:                0,
		Packge_date_format:          "dd-mm-YYYY",
		Print_used_by_date:          'N',
		Used_by_date_source:         0,
		Used_by_date:                0,
		Unit_price_selection:        0,
		Tare:                        0.0,
		Perc:                        0.0,
		Quant:                       0,
		Upc:                         "",
		Uom_id:                      "",
		Discount_schema_id:          "",
		Isstocked:                   'N',
		Isirradiated:                'N',
		Isbom:                       'N',
		Stock_min:                   0.0,
		Stock_max:                   0.0,
		Versionno:                   "",
		Discontinued:                'N',
		DiscontinuedBy:              time.Now(),
		Tax_id:                      "@",
		Atrribute:                   0,
		Advertising_id:              "",
		Sto_temp_limit1:             0.0,
		Sto_temp_limit2:             0.0,
		Ingredients:                 "",
		Ingredients_id:              0,
		Preservation_info:           "",
		Preservation_info_id:        "",
		Coupled_product:             0,
		Pack_indiCador:              0,
		Extra_field:                 "",
		Extra_field2:                "",
		Image_id:                    "",
		Icon_id:                     "",
		Lot:                         "",
		Born_country:                "",
		Fatten_country:              "",
		Origin_countrt:              "",
		Manufacturer_id:             "",
		Packer_id:                   "",
		Distribuidor_id:             "",
		Importer_id:                 "",
		Exporter_id:                 "",
		Vendor_id:                   "",
		Cutting_hall_id:             "",
		Slaughter_house_id:          "",
		Supplier_id:                 "",
		Nut_info_set_id:             "2",
		Print_fix_primary_barcode:   'N',
		Print_fix_secondary_barcode: "N",
	}
}

func EnviaPluSimples(cod, atributo, cod_setor, venda, validade int,
	nome, descricao, imprimeValidade string,
	tare float32) {

}

func InserirProdutoNOTUTF8(c, c_setor, v, val int,
	name, desc, imprimeValidade string, tare float32, db *sql.DB, mapIds map[int]bool) error {
	atributo := 0
	info_venda := 1
	if v > 0 {
		fmt.Println(c, "unitario")
		atributo = 1
		info_venda = 2
	}
	cod := c
	nome := name
	cod_setor := c_setor
	descricao := desc
	imprimeValidade = "Y"
	validade := val

	if name == "" {
		fmt.Println("teste")
	}

	existeSetor, err := ExisteSetorPostgres(cod_setor, db)
	if err != nil {
		fmt.Println("erro: setor", err)
		return err
	}

	if !existeSetor {
		descricao_setor := fmt.Sprint("Setor: ", cod_setor)

		_, err = db.Exec(`
			SELECT public.merge_department(
				$1, 
				$2 
			)
		`, cod_setor, descricao_setor)
		if err != nil {
			fmt.Println("Erro ao chamar a função merge_product: - SETOR: ", c_setor, " ", err)
			return err
		}

	}

	ExistePlu := false
	ExistePlu = !(mapIds[cod])

	if ExistePlu {
		_, err = db.Exec(`
			SELECT public.merge_product(
				$1, 
				$2, 
				$3, 
				$4, 
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11,
				$12
			)
		`, cod, cod, nome, atributo, cod_setor, descricao, imprimeValidade, validade, info_venda, 0, "", "")
		if err != nil {
			fmt.Println("Erro ao chamar a função merge_product: - Plu: ", cod, " ", err)
			return err
		}
		fmt.Println("Função merge_product chamada com sucesso.", cod)
		return nil
	}
	err = UpdatesPluDadosSimples(cod, info_venda, validade, atributo, descricao, db)
	if err != nil {
		println("erro: ", cod, mapIds[cod], " erro: ", err)
		return err
	}
	fmt.Println("Função merge_product chamada com sucesso.")
	return nil
}

func AtualizaPreco(valor string, plu int, db *sql.DB) error {
	float_valor, err := strconv.ParseFloat(valor, 32)
	if err != nil {
		fmt.Println("PREÇO INVALIDO PARA: ", plu)
	}
	float_valor = float_valor / 100

	_, err = db.Exec(`
	SELECT public.merge_product_price(
		$1, 
		$2, 
		$3 
	)
`, "lst1", plu, float_valor)
	if err != nil {
		fmt.Println("Erro ao chamar a função merge_product: - plu: ", plu, "Valor: ", float_valor, ": ", err)
		return err
	}

	fmt.Println("Função merge_product chamada com sucesso.")
	return nil

}

func ExistePluPostgres(codigo int, db *sql.DB) (bool, error) {

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM product WHERE product_id = ?", codigo).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func IdsPlus(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query("SELECT product_id FROM product")
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
		fmt.Println(err)
	}

	return IdsExistentes(ids), nil
}

func IdsExistentes(lista []int) map[int]bool {
	idMap := make(map[int]bool, len(lista))
	for _, id := range lista {
		idMap[id] = true
	}
	return idMap
}

func InserirProduto(c, c_setor, v, val int,
	name, desc, imprimeValidade string, tare float32, db *sql.DB, mapIds map[int]bool) error {
	// Função auxiliar para garantir que a string esteja em UTF-8
	encodeUTF8 := func(s string) (string, error) {
		if !utf8.ValidString(s) {
			return "", fmt.Errorf("string contains invalid UTF-8: %q", s)
		}
		return s, nil
	}

	if c == 3 {
		fmt.Println("log plu")
	}

	// Garantir que as strings estão em UTF-8
	name, err := encodeUTF8(name)
	if err != nil {
		fmt.Println("Erro de codificação UTF-8 no nome:", err)
		return err
	}
	desc, err = encodeUTF8(desc)
	if err != nil {
		fmt.Println("Erro de codificação UTF-8 na descrição:", err)
		return err
	}

	atributo := 0
	info_venda := 1
	if v > 0 {
		fmt.Println(c, "unitario")
		atributo = 1
		info_venda = 2
	}
	cod := c
	nome := name
	cod_setor := c_setor
	descricao := desc
	imprimeValidade = "Y"
	validade := val

	if name == "" {
		fmt.Println("teste")
	}

	existeSetor, err := ExisteSetorPostgres(cod_setor, db)
	if err != nil {
		fmt.Println("erro: setor", err)
		return err
	}

	if !existeSetor {
		descricao_setor := fmt.Sprintf("Setor: %d", cod_setor)

		_, err = db.Exec(`
			SELECT public.merge_department(
				$1, 
				$2 
			)
		`, cod_setor, descricao_setor)
		if err != nil {
			fmt.Println("Erro ao chamar a função merge_product: - SETOR: ", c_setor, " ", err)
			return err
		}
	}

	ExistePlu := false
	ExistePlu = !(mapIds[cod])

	if ExistePlu {
		_, err = db.Exec(`
			SELECT public.merge_product(
				$1, 
				$2, 
				$3, 
				$4, 
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11,
				$12
			)
		`, cod, cod, nome, atributo, cod_setor, descricao, imprimeValidade, validade, info_venda, 0, "", "")
		if err != nil {
			fmt.Println("Erro ao chamar a função merge_product: - Plu: ", cod, " ", err)
			//log.Fatal("ERRO UTF-8")
			return err
		}
		fmt.Println("Função merge_product chamada com sucesso.", cod)
		return nil
	}
	err = UpdatesPluDadosSimples(cod, info_venda, validade, atributo, descricao, db)
	if err != nil {
		println("erro: ", cod, mapIds[cod], " erro: ", err)
		return err
	}
	fmt.Println("Função merge_product chamada com sucesso.")
	return nil
}

func UpdatesPluDadosSimples(plu, venda, validade, attribute int, desc string, db *sql.DB) error {
	_, err := db.Exec(`
	UPDATE product
	SET
		name = $2,
		used_by_date = $3,
		uom_id = $4,
		attribute = $5
	WHERE
		product_id = $1
	`, plu, desc, validade, venda, attribute)
	if err != nil {
		fmt.Println("Erro ao chamar a função merge_product: - PLU: ", plu, " ", err)
		return err
	}
	return nil
}
