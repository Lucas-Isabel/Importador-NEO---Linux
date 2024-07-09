package file

import (
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func caracterRemove(txt string) string {
	texto := strings.ReplaceAll(txt, "\n", "")
	texto = strings.ReplaceAll(texto, "    ", "")
	texto = strings.ReplaceAll(texto, ", ", ",")
	texto = strings.ReplaceAll(texto, " (", "(")
	return texto
}

func EncodeToUTF8(str string) (string, error) {
	// Como o encoding/charmap fazia a conversão do Windows-1252 para UTF-8,
	// você pode usar a função ReplaceAll diretamente para substituir os caracteres.
	utf8Str := strings.ReplaceAll(str, "\x92", "'")
	utf8Str = strings.ReplaceAll(utf8Str, "\x96", "-")
	return replaceSpecialChars(utf8Str), nil
}

func saoNumeros(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

func stringTofloat(val string, divisor float64) (float64, error) {
	valor, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	valor = valor / divisor
	return valor, nil
}

func replaceSpecialChars(txt string) string {
	var sb strings.Builder
	for _, r := range txt {
		if unicode.Is(unicode.Mn, r) {
			// Se o caractere for um marcador de não-combinação, ignore-o
			continue
		}
		if utf8.ValidRune(r) {
			// Se o caractere for válido UTF-8, adicione-o ao buffer
			sb.WriteRune(r)
		}
	}
	return sb.String()
}
