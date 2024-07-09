package file

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lucasbyte/go-clipse/db"
	"github.com/lucasbyte/go-clipse/models"
)

func EnviarPluSeparado(arq string, balanca models.Balanca) error {
	ip := balanca.Ip
	setores_das_balancas := strings.Trim(balanca.Setores, "[]")
	todos_os_setores := strings.Split(setores_das_balancas, " ")

	item, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer item.Close()

	db := db.ConectDbPQ(ip)
	var enviados []string
	teste := 0
	scanner := bufio.NewScanner(item)
	mapaDeIds, err_conect := models.IdsPlus(db)
	if err_conect != nil {
		return err_conect
	}

	for scanner.Scan() {
		line := scanner.Text()
		setor := line[0:2]
		setor_conv, _ := strconv.Atoi(setor)
		setor = fmt.Sprint(setor_conv)
		if ContainsToStr(setor, todos_os_setores) {
			if !ContainsToStr(setor, enviados) {
				enviados = append(enviados, setor)
			}
			teste++
			codPlu := line[3:9]
			int_codPlu, _ := strconv.Atoi(codPlu)
			venda := line[2:3]
			int_venda, _ := strconv.Atoi(venda)
			validade := line[15:18]
			int_validade, _ := strconv.Atoi(validade)
			int_validade--
			if int_validade < 0 {
				int_validade = 0
			}
			nome := line[18:44]
			//nome = strings.ReplaceAll(nome, "�", " ")
			valor_str := line[9:15]
			plu := int_codPlu
			erro_imp := models.InserirProduto(int_codPlu, setor_conv, int_venda, int_validade, nome, nome, "Y", 0, db, mapaDeIds)
			if erro_imp != nil {
				println(int_codPlu, erro_imp)
			}
			erro_imp = models.AtualizaPreco(valor_str, plu, db)
			if erro_imp != nil {
				err = erro_imp
			}
			// info := d_info[codInfo]
			// if info != "" {
			// 	err_envio_info = enviarInf(plu , codPlu, db, "extra_field1")
			// 	erro_str = fmt.Sprint(err_envio_info)
			// }
			// if strings.Contains(erro_str, "pq: password authentication failed for user") || strings.Contains(erro_str, "No connection could be made because the target machine actively refused it.") {
			// 	fmt.Println("Chegou aqui: ", err_envio_info)
			// 	return err_envio_info
			// }
		}
	}
	defer db.Close()
	fmt.Println(enviados)

	fmt.Println(EncodeToUTF8(ip))
	return err
}
