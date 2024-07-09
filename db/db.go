package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func createDatabaseIfNotExists(db *sql.DB) error {
	// Verifica se o arquivo do banco de dados já existe
	_, err := os.Stat(".db")
	if os.IsNotExist(err) {
		// Se o arquivo não existir, cria o banco de dados e a tabela
		_, err := db.Exec(`
            CREATE TABLE IF NOT EXISTS produtos (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
				plu INTEGER UNIQUE CHECK(plu <= 200),
            	descricao TEXT CHECK(length(descricao) <= 15),
				venda INTEGER CHECK(venda <= 10),
                validade INTEGER CHECK(validade <= 200),
				preco DOUBLE CHECK(preco < 1000)
            );
        `)
		if err != nil {
			return err
		}
		fmt.Println("Banco de dados e tabela criados com sucesso.")
	} else if err != nil {
		// Em caso de erro ao verificar a existência do arquivo
		return err
	}

	_, err = os.Stat(".db")
	if os.IsNotExist(err) {
		// Se o arquivo não existir, cria o banco de dados e a tabela
		_, err := db.Exec(`
            CREATE TABLE IF NOT EXISTS balancas (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
				ip TEXT UNIQUE,
            	nome TEXT CHECK(length(nome) <= 35),
				setores TEXT,
				nomeDosSetores TEXT
            );
        `)
		if err != nil {
			return err
		}
		fmt.Println("Banco de dados e tabela criados com sucesso.")
	} else if err != nil {
		// Em caso de erro ao verificar a existência do arquivo
		return err
	}

	_, err = os.Stat(".db")
	if os.IsNotExist(err) {
		// Se o arquivo não existir, cria o banco de dados e a tabela
		_, err := db.Exec(`
            CREATE TABLE IF NOT EXISTS setor (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
				codigo INTEGER UNIQUE,
            	nome TEXT CHECK(length(nome) <= 35)
            );
        `)
		_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS setor_ip (
                id INTEGER PRIMARY KEY AUTOINCREMENT,
				codigo_setor INTEGER,
            	ip TEXT
            );
        `)
		if err != nil {
			return err
		}
		fmt.Println("Banco de dados e tabela criados com sucesso.")
	} else if err != nil {
		// Em caso de erro ao verificar a existência do arquivo
		return err
	}

	return nil
}

func init() {
	// Abre a conexão com o banco de dados
	db, err := sql.Open("sqlite3", "import-br.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()

	err = createDatabaseIfNotExists(db)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Resto do código...
}

func ConectDb() *sql.DB {
	db, err := sql.Open("sqlite3", "import-br.db")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

func Existe(tabela, campo string, codigo int) (bool, error) {
	db := ConectDb()
	defer db.Close()

	query := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s = ?", tabela, campo)
	var count int
	err := db.QueryRow(query, codigo).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
