package controllers

import (
	"net/http"

	"github.com/lucasbyte/go-clipse/Cad"
	"github.com/lucasbyte/go-clipse/file"
	"github.com/lucasbyte/go-clipse/txitens"
)

type Tipos struct {
	Tipo    file.FormatacaoSelect
	Toledo  file.Arquivos
	Txitens txitens.Arquivos_Txitens
	Cad     Cad.Arquivos_Cad
}

func File(w http.ResponseWriter, r *http.Request) {

}

func EditaArquivo(w http.ResponseWriter, r *http.Request) {

}

func EditaCad(w http.ResponseWriter, r *http.Request) {

}

func PuxarArquivos(w http.ResponseWriter, r *http.Request) {

}
