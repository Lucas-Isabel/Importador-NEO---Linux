package file

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lucasbyte/go-clipse/models"
)

func EnviarInfoSeparada(arq string, d_nutri map[string][13]string, d_info, d_forn, d_aler, d_fra, d_con map[string]string, d_tara map[string]float64, balanca models.Balanca) error {
	ip := balanca.Ip
	if ip == "localhost" {
		ip = "127.0.0.1" // Use IPv4 address for localhost
	}
	setores_das_balancas := strings.Trim(balanca.Setores, "[]")
	todos_os_setores := strings.Split(setores_das_balancas, " ")

	item, err := os.Open(arq)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer item.Close()

	passw, user := "1234", "user"
	if ip == "127.0.0.1" {
		passw, user = "Systel#4316", "systel"
	}

	db_acess := fmt.Sprintf("user=%s password=%s host=%s dbname=cuora sslmode=disable", user, passw, ip)

	scanner := bufio.NewScanner(item)

	db, err := sql.Open("postgres", db_acess)
	if err != nil {
		fmt.Println("ERRO: ", err)
		return err
	}
	defer db.Close()

	mapNutriPlu, err := models.MapNutriPlusPG(db)
	if err != nil {
		return err
	}

	for scanner.Scan() {
		line := scanner.Text()
		setor := strings.TrimSpace(line[0:2])
		setorConv, _ := strconv.Atoi(setor)
		setor = fmt.Sprint(setorConv)

		if ContainsToStr(setor, todos_os_setores) {
			processLine(line, d_nutri, d_info, d_forn, d_aler, d_fra, d_con, db, mapNutriPlu)
		}
	}

	fmt.Println(EncodeToUTF8(ip))
	return nil
}

func processLine(line string, d_nutri map[string][13]string, d_info, d_forn, d_aler, d_fra, d_con map[string]string, db *sql.DB, mapNutriPlu map[int]bool) {
	lote := line[90:102]
	codPlu := line[3:9]
	codInfo := line[68:74]
	codAler := line[126:130]
	codForn := line[86:90]
	codFrac := line[122:126]
	codCons := line[134:138]
	codNutri := line[78:84]

	nutri := d_nutri[codNutri]
	info := d_info[codInfo]
	aler := d_aler[codAler]
	forn := d_forn[codForn]
	cons := d_con[codCons]
	frac := d_fra[codFrac]

	pluInt, _ := strconv.Atoi(codPlu)
	existeNutri := !mapNutriPlu[pluInt]

	if codNutri != "000000" && !allElementsEmpty(nutri) {
		models.EnviaNutriPG(codPlu, nutri, existeNutri, db)
	}

	if len(forn) < 158 && forn != "" {
		forn += strings.Repeat(" ", 158-len(forn))
	}

	fornFrac := forn + " " + frac

	loteInt, _ := strconv.Atoi(lote)
	sendInfos(loteInt, codPlu, info, aler, cons, fornFrac, db)
}

func sendInfo(loteInt int, codPlu, info, aler, cons, fornFrac string, db *sql.DB) {
	if loteInt > 0 {
		err := EnviarInf(strconv.Itoa(loteInt), codPlu, db, "lot")
		if err != nil {
			handleError(err)
			return
		}
	}
	if info != "" {
		err := EnviarInf(info, codPlu, db, "extra_field1")
		if err != nil {
			handleError(err)
			return
		}
	}
	if aler != "" {
		err := EnviarInf(aler, codPlu, db, "extra_field2")
		if err != nil {
			handleError(err)
			return
		}
	}
	if cons != "" {
		err := EnviarInf(cons, codPlu, db, "preservation_info")
		if err != nil {
			handleError(err)
			return
		}
	}
	if fornFrac != "" {
		err := EnviarInf(fornFrac, codPlu, db, "ingredients")
		if err != nil {
			handleError(err)
			return
		}
	}
}

func sendInfos(loteInt int, codPlu, info, aler, cons, fornFrac string, db *sql.DB) {
	var valores []string
	var campos []string

	if loteInt > 0 {
		valores = append(valores, strconv.Itoa(loteInt))
		campos = append(campos, "lot")
	}
	if info != "" {
		valores = append(valores, info)
		campos = append(campos, "extra_field1")
	}
	if aler != "" {
		valores = append(valores, aler)
		campos = append(campos, "extra_field2")
	}
	if cons != "" {
		valores = append(valores, cons)
		campos = append(campos, "preservation_info")
	}
	if fornFrac != "" {
		valores = append(valores, fornFrac)
		campos = append(campos, "ingredients")
	}

	if len(valores) == 0 {
		return // Não há nada para enviar
	}

	err := EnviarInfs(valores, codPlu, db, campos)
	if err != nil {
		handleError(err)
		return
	}
}

func handleError(err error) {
	if strings.Contains(err.Error(), "pq: password authentication failed for user") || strings.Contains(err.Error(), "No connection could be made because the target machine actively refused it.") {
		fmt.Println("Chegou aqui: ", err)
	}
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
