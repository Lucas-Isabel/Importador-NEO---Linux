package models

import (
	"fmt"
	"time"

	"github.com/lucasbyte/go-clipse/db"
)

type Balanca struct {
	Id            int
	Ip            string
	Nome          string
	Setores       string
	Lista_Setores string
}

func CriaNovaBalanca(ip, nome, setores, nomesDosSetores string) error {
	db := db.ConectDb()
	defer db.Close()

	insereDadosNoBanco, err := db.Prepare("insert into balancas(ip, nome, setores, nomeDosSetores) values($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	result, err := insereDadosNoBanco.Exec(ip, nome, setores, nomesDosSetores)
	if err != nil {
		return err
	}
	fmt.Println(result)
	return nil
}

func AtualizaBalanca(ip, nome, setores, nomesDosSetores string) error {
	db := db.ConectDb()
	defer db.Close() // Assegura que a conexão será fechada mesmo em caso de erro

	// atualizaDadosNoBanco, err := db.Prepare("UPDATE balancas SET nome = $2, setores = $3, nomeDosSetores = $4 WHERE ip = $1")
	// if err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }

	// result, err := atualizaDadosNoBanco.Exec(ip, nome, setores, nomesDosSetores)
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(result)
	// return nil

	fmt.Println(DeleteBalanca(ip))
	time.Sleep(time.Second)
	err := CriaNovaBalanca(ip, nome, setores, nomesDosSetores)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func BuscaBalancas() ([]Balanca, error) {
	db := db.ConectDb()
	defer db.Close()

	selectDeTodosOsProdutos, err := db.Query("select * from balancas ORDER BY ip")
	if err != nil {
		return nil, err
	}

	bal := Balanca{}
	balancas := []Balanca{}

	for selectDeTodosOsProdutos.Next() {
		var ip, descricao, setores, lista_setores string
		var id int
		err = selectDeTodosOsProdutos.Scan(&id, &ip, &descricao, &setores, &lista_setores)
		if err != nil {
			return nil, err
		}

		bal.Id = id
		bal.Ip = ip
		bal.Nome = descricao
		bal.Setores = setores
		bal.Lista_Setores = lista_setores

		balancas = append(balancas, bal)
	}
	return balancas, nil
}

func BuscaBalancasPorIps(ips []string) ([]Balanca, error) {
	db := db.ConectDb()
	defer db.Close()
	var balancas []Balanca
	for _, ip := range ips {
		bal, err := BuscaBalanca(ip)
		if err != nil {
			continue
		}
		balancas = append(balancas, bal)
	}
	return balancas, nil
}

func BuscaBalanca(ip string) (Balanca, error) {
	db := db.ConectDb()
	defer db.Close()
	bal := Balanca{}

	query := "SELECT id, ip, nome, setores, nomeDosSetores FROM balancas WHERE ip = $1"
	err := db.QueryRow(query, ip).Scan(&bal.Id, &bal.Ip, &bal.Nome, &bal.Setores, &bal.Lista_Setores)
	if err != nil {
		return bal, err
	}

	return bal, nil
}

func DeleteBalanca(id string) (string, error) {
	db := db.ConectDb()
	defer db.Close()

	fmt.Println("delete:", id)
	query := "DELETE FROM balancas WHERE ip=$1"

	insereDadosNoBanco, err := db.Prepare(query)
	if err != nil {
		db.Close()
		return fmt.Sprintln(err), err
	}
	fmt.Println("delete:", id)

	result, err := insereDadosNoBanco.Exec(id)
	if err != nil {
		db.Close()
		return fmt.Sprintln(err), err
	}
	v, err := result.RowsAffected()
	db.Close()

	return fmt.Sprint(v), err
}

func RemoveElement(items []Balanca, itemToRemove Balanca) []Balanca {
	for i, v := range items {
		if v == itemToRemove {
			return append(items[:i], items[i+1:]...)
		}
	}
	return items
}
