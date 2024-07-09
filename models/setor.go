package models

import (
	"database/sql"
	"fmt"

	"github.com/lucasbyte/go-clipse/db"
)

type Setor struct {
	Id     int
	Codigo int
	Nome   string
}

func NewSetor(cod int, nome string) Setor {
	setor := Setor{
		Codigo: cod,
		Nome:   nome,
	}
	return setor
}

func ExisteSetor(codigo int) (bool, error) {
	db := db.ConectDb()
	defer db.Close()

	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM setor WHERE the_departmet_id = ?", codigo).Scan(&count)
	if err != nil {
		db.Close()
		return false, err
	}

	return count > 0, nil
}

func EnviaSetor(setor Setor) {
	existe, err := ExisteSetor(setor.Codigo)
	if err != nil {
		fmt.Println("Teste")
	}
	if !(existe) {
		db := db.ConectDb()
		insereDadosNoBanco, err := db.Prepare("insert into setor(the_department_id, the_name) values($1, $2)")
		if err != nil {
			fmt.Println(err.Error())
		}

		result, err := insereDadosNoBanco.Exec(setor.Codigo, setor.Nome)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(result)
		}
		defer db.Close()
	}
}

func BuscaSetores() []Setor {
	db := db.ConectDb()

	selectDeTodosOsProdutos, err := db.Query("select * from setor ORDER BY codigo")
	if err != nil {
		fmt.Println(err.Error())
	}

	setor := Setor{}
	setores := []Setor{}

	for selectDeTodosOsProdutos.Next() {
		var nome string
		var codigo, id int
		err = selectDeTodosOsProdutos.Scan(&id, &codigo, &nome)
		if err != nil {
			fmt.Println(err.Error())
		}

		setor.Id = id
		setor.Codigo = codigo
		setor.Nome = nome

		setores = append(setores, setor)
	}
	defer db.Close()
	return setores
}

func Update_setor(nome string, id int) {
	db := db.ConectDb()
	query := "UPDATE setor SET nome = ? WHERE codigo = ?"

	insereDadosNoBanco, err := db.Prepare(query)
	if err != nil {
		fmt.Println(err.Error())
	}

	insereDadosNoBanco.Exec(nome, id)
	defer db.Close()
}

func ExisteSetorPostgres(codigo int, db *sql.DB) (bool, error) {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM department WHERE department_id = $1", codigo).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
