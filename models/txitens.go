package models

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/lucasbyte/go-clipse/db"
)

func Txitens(caminho string, balancas []Balanca) error {
	//arquivo, err := os.Open("C:\\Users\\MAXWELL\\Desktop\\ARQUIVOS DE ITENS\\txitens.txt")

	arquivo, err := os.Open(caminho)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return err
	}
	defer arquivo.Close()

	scanner := bufio.NewScanner(arquivo)

	produtos := make([]Produto, 0)

	for scanner.Scan() {
		linha := scanner.Text()

		if len(linha) < 17 {
			fmt.Println("Linha do arquivo TXT inválida:", linha)
			continue
		}

		venda, _ := strconv.Atoi(string(linha[4]))
		plu, _ := strconv.Atoi(linha[5:11])
		setores := linha[0:2]
		precoStr := linha[11:17]
		validadeStr := linha[17:20]
		validade, _ := strconv.Atoi(validadeStr)
		descricao := ""
		if len(linha) < 35 {
			descricao = strings.TrimSpace(linha[20:])
		} else {
			descricao = strings.TrimSpace(linha[20:35])
		}
		fmt.Println(linha[5:11], venda, precoStr, validadeStr, descricao)
		produto := Produto{
			Id:        plu, // Seu código para definir o ID do produto
			Plu:       plu,
			Setores:   setores,
			Descricao: descricao,
			Preco:     precoStr,
			Venda:     venda,
			Validade:  validade,
		}

		produtos = append(produtos, produto)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Erro ao ler o arquivo TXT:", err)
		return err
	}
	for _, balanca := range balancas {
		setores_das_balancas := strings.Trim(balanca.Setores, "[]")
		todos_os_setores := strings.Split(setores_das_balancas, " ")
		db := db.ConectDbPQ(balanca.Ip)

		mapIds, err_conect := IdsPlus(db)
		if err_conect != nil {
			return err_conect
		}
		var enviados []string
		for _, p := range produtos {
			if p.Plu > 0 {
				if ContainsToStr(p.Setores, todos_os_setores) {
					if !ContainsToStr(p.Setores, enviados) {
						existe := false
						int_setor, _ := strconv.Atoi(p.Setores)
						if existe {
							fmt.Println("Envio apenas de preco e descrição")
						} else {
							fmt.Println(InserirProduto(p.Plu, int_setor, p.Venda, p.Validade, p.Descricao, " ", "Y", 0, db, mapIds))
							fmt.Println(AtualizaPreco(p.Preco, p.Plu, db))
						}
					}
				}
			}
		}
		db.Close()
	}
	return nil
}

func ContainsToStr(target string, arr []string) bool {
	for _, val := range arr {
		if val == target {
			return true
		}
	}
	return false
}
