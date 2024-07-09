package file

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lucasbyte/go-clipse/db"
	"github.com/lucasbyte/go-clipse/models"
)

func ExisteSetor(codigo int) (bool, error) {
	db := db.ConectDb()
	defer db.Close()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM setor WHERE codigo = ?", codigo).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func DeleteAllSetores() error {
	db := db.ConectDb()
	defer db.Close()
	_, err := db.Exec("DELETE FROM setor")
	if err != nil {
		return err
	}
	fmt.Println("Todos os setores foram deletados com sucesso.")
	time.Sleep(time.Second)
	return nil
}

func AdicionaSetores(itens_caminho string) {
	fmt.Println("Adicionando Setores")
	itens_file, err := os.Open(itens_caminho)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(itens_file)

	db := db.ConectDb()

	for scanner.Scan() {
		line := scanner.Text()

		setor := line[0:2]
		int_setor, _ := strconv.Atoi(setor)
		existe, err := ExisteSetor(int_setor)
		if !existe && err == nil {
			codigo_setor := int_setor
			insereDadosNoBanco, _ := db.Prepare("insert into setor(codigo, nome) values($1, $2)")
			nome_setor := fmt.Sprintf("Setor: " + setor)
			result, err := insereDadosNoBanco.Exec(codigo_setor, nome_setor)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(result)
			}
		}
	}
	defer db.Close()
}

func AdicionaSetoresTxitens(itens_caminho string) {
	fmt.Println("Adicionando Setores")
	itens_file, err := os.Open(itens_caminho)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(itens_file)

	db := db.ConectDb()

	for scanner.Scan() {
		line := scanner.Text()

		setor := line[0:2]
		int_setor, _ := strconv.Atoi(setor)
		existe, err := ExisteSetor(int_setor)
		if !existe && err == nil {
			codigo_setor := int_setor
			insereDadosNoBanco, _ := db.Prepare("insert into setor(codigo, nome) values($1, $2)")
			nome_setor := fmt.Sprintf("Setor: " + setor)
			result, err := insereDadosNoBanco.Exec(codigo_setor, nome_setor)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println(result)
			}
		}
	}
	defer db.Close()
}

func VerificarSetores() ([]int, []models.Setor) {
	var lista_setores []int
	setores := []models.Setor{}
	fmt.Println("SETORES: ", setores)
	db := db.ConectDb()

	selectDeTodosOsSetores, err := db.Query("select * from setor ORDER BY codigo")
	if err != nil {
		fmt.Println(err.Error())
	}

	setor := models.Setor{}

	for selectDeTodosOsSetores.Next() {
		var id, codigo int
		var nome string

		err = selectDeTodosOsSetores.Scan(&id, &codigo, &nome)
		if err != nil {
			fmt.Println(err.Error())
		}

		setor.Id = id
		setor.Nome = nome
		setor.Codigo = codigo

		setores = append(setores, setor)
	}
	defer db.Close()
	fmt.Println("SETORES: ", setores)
	return lista_setores, setores
}
