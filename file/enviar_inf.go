package file

import (
	"database/sql"
	"fmt"
	"strings"
)

func EnviarInfs(t []string, p string, db *sql.DB, campos []string) error {
	if len(t) != len(campos) {
		return fmt.Errorf("listas de t e campos devem ter o mesmo tamanho")
	}

	var updates []string
	for i := 0; i < len(t); i++ {
		tara := t[i]
		campo := campos[i]

		tara, err := EncodeToUTF8(tara)
		if err != nil {
			fmt.Println("ERRO: ", err, "item: ", tara)
			return err
		}

		tara = strings.TrimSpace(strings.ReplaceAll(tara, "\n", ""))
		if tara != "" {
			update := fmt.Sprintf("%s = '%s'", campo, tara)
			updates = append(updates, update)
		}
	}

	if len(updates) == 0 {
		return nil // Não há nada para atualizar
	}

	updatesStr := strings.Join(updates, ", ")
	comando := fmt.Sprintf("UPDATE product SET %s WHERE product_id = %s", updatesStr, p)
	_, err := db.Exec(comando)
	if err != nil {
		fmt.Println("Erro-update:", err)
		return err
	}

	return nil
}

func EnviarInf(t, p string, db *sql.DB, campo string) error {
	plu := p
	tara := t
	tara, err := EncodeToUTF8(tara)
	if strings.TrimSpace(strings.ReplaceAll(tara, "\n", "")) != "" {
		if err != nil {
			fmt.Println("ERRO: ", err, "item: ", tara)
			return err
		} else {
			fmt.Println("")
		}
		comando := fmt.Sprintf("UPDATE product set %s = '%s' WHERE product_id = %s", campo, tara, plu)
		_, err = db.Exec(comando)
		if err != nil {
			fmt.Println("Erro-update:", err)
			return err
		}
	}
	return nil
}
