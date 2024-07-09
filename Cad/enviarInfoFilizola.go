package Cad

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lucasbyte/go-clipse/file"
	"github.com/lucasbyte/go-clipse/models"
)

func EnviarInfoSeparadaCad(arq string, d_nutri map[string][13]string, d_info map[string]string, balanca models.Balanca) error {
	ip := balanca.Ip
	setores_das_balancas := strings.Trim(balanca.Setores, "[]")
	todos_os_setores := strings.Split(setores_das_balancas, " ")

	item, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer item.Close()

	passw, user := "1234", "user"
	if ip == "localhost" {
		passw, user = "Systel#4316", "systel"
	}

	db_acess := fmt.Sprintf("user=%s password=%s host=%s dbname=cuora sslmode=disable", user, passw, ip)
	db, err := sql.Open("postgres", db_acess)
	if err != nil {
		fmt.Println("ERRO: ", err)
		return err
	}
	defer db.Close()
	var enviados []string
	scanner := bufio.NewScanner(item)
	mapNutriPlu, err := models.MapNutriPlusPG(db)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		setor := line[0:2]
		setor_conv, _ := strconv.Atoi(setor)
		setor = fmt.Sprint(setor_conv)
		//fmt.Println("esse é o setor: ", setor, "esses sao os setores: ", todos_os_setores)
		if file.ContainsToStr(setor, todos_os_setores) {
			if !file.ContainsToStr(setor, enviados) {
				enviados = append(enviados, setor)
			}
			codPlu := line[3:9]
			codInfo := line[68:74]

			codNutri := line[78:84]
			if codPlu == "002112" {
				fmt.Println(codNutri)
			}

			nutri := d_nutri[codNutri]
			info := d_info[codInfo]
			existeNutri := false

			plu_int, _ := strconv.Atoi(codPlu)
			existeNutri = !mapNutriPlu[plu_int]

			if codNutri != "000000" && !allElementsEmpty(nutri) {
				models.EnviaNutriPG(codPlu, nutri, existeNutri, db)
			}

			var err_envio_info error
			var erro_str string
			if info != "" {
				err_envio_info = file.EnviarInf(info, codPlu, db, "extra_field1")
				erro_str = fmt.Sprint(err_envio_info)
			}
			if strings.Contains(erro_str, "pq: password authentication failed for user") || strings.Contains(erro_str, "No connection could be made because the target machine actively refused it.") {
				fmt.Println("Chegou aqui: ", err_envio_info)
				return err_envio_info
			}
		}
	}
	fmt.Println(enviados)

	fmt.Println(file.EncodeToUTF8(ip))
	return nil
}

func allElementsEmpty(list [13]string) bool {
	sum := 0
	for _, element := range list {
		num, err := strconv.Atoi(element)
		if err != nil {
			num = 0
		}
		sum = sum + num
	}
	return sum == 0
}

func anyElementsEmpty(list [13]string) bool {
	for _, element := range list {
		if element == "" {
			return true
		}
	}
	return false
}
