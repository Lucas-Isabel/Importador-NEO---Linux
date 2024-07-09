package file

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/lib/pq"
	"github.com/lucasbyte/go-clipse/models"
)

func EnviaParaBalancas(itens, receita, nutri, frac, forn, tara, cons, campext1 string,
	Balancas []models.Balanca, somente_preco bool) ([]models.Balanca, []models.Balanca) {
	var balancas_enviadas []models.Balanca
	balancas_enviadas = append(balancas_enviadas, Balancas...)
	var balancas_nao_enviadas []models.Balanca
	itens = strings.Replace(itens, "\\", "//", -1)
	receita = strings.Replace(receita, "\\", "//", -1)
	nutri = strings.Replace(nutri, "\\", "//", -1)
	frac = strings.Replace(frac, "\\", "//", -1)
	forn = strings.Replace(forn, "\\", "//", -1)
	tara = strings.Replace(tara, "\\", "//", -1)
	cons = strings.Replace(cons, "\\", "//", -1)
	campext1 = strings.Replace(campext1, "\\", "//", -1)

	// var Ehigual bool = false
	// var dataEhora string
	// var criacaoDoarquivo string = ""
	var mensagem string
	{
		arquivo := itens

		// if fileExists(arquivo) {
		// 	dataEhora, _ = obterDataHoraCriacao(arquivo)
		// 	Ehigual = criacaoDoarquivo == dataEhora
		// } else if fileExists(ArrumaExtensão(arquivo, "txt", "bak")) {
		// 	arquivo = ArrumaExtensão(arquivo, "txt", "bak")
		// 	dataEhora, _ = obterDataHoraCriacao(arquivo)
		// 	Ehigual = criacaoDoarquivo == dataEhora
		// }

		// if !Ehigual {
		// 	criacaoDoarquivo = dataEhora
		if _, err := os.Stat(arquivo); err == nil {
			time.Sleep(1 * time.Second)

			var dict_conserva, dict_fraciona, dict_aler, dict_forn map[string]string
			var info map[string]string
			var dict_tara map[string]float64

			if _, err := os.Stat(tara); err == nil {
				dict_tara = taraAnalyze(tara)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}
			if _, err := os.Stat(cons); err == nil {
				dict_conserva = conservaAnalyze(cons)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}

			if _, err := os.Stat(frac); err == nil {
				dict_fraciona = fracionaAnalyze(frac)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}

			if _, err := os.Stat(campext1); err == nil {
				dict_aler = alergiaAnalyze(campext1)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}

			if _, err := os.Stat(forn); err == nil {
				dict_forn = fornAnalyze(forn)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}

			if _, err := os.Stat(receita); err == nil {
				info = InfoAnalyze(receita)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}
			arquivoInfonutri := nutri
			if _, err := os.Stat(nutri); err == nil {
				arquivoInfonutri = nutri
			} else {
				nutri = ""
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}
			fmt.Print(mensagem)

			nomeDaPasta := VerificaPasta()

			new_itens := fmt.Sprintf("%s/itens.TXT", nomeDaPasta)
			itensAnalyze(arquivo, new_itens)
			mgv7File, err := os.Open(new_itens)
			if err != nil {
				fmt.Println(err)
			}

			defer mgv7File.Close()

			infonutriFile, err := os.Open(arquivoInfonutri)
			if err != nil {
				fmt.Println(err)
			}
			defer infonutriFile.Close()

			// Create output files
			nutriFile, err := os.Create(fmt.Sprintf("%s/nutriSystel.TXT", nomeDaPasta))
			if err != nil {
				fmt.Println(err)
			}
			defer nutriFile.Close()

			systelFile, err := os.Create(fmt.Sprintf("%s/itensSystel.TXT", nomeDaPasta))
			if err != nil {
				fmt.Println(err)
			}
			defer systelFile.Close()

			// Initialize arrays
			//var codPluArray []string
			var codNutriArray []string
			var codNutriMGVArray []string
			//var receitaArray []string

			// Read and process mgv7File
			scanner := bufio.NewScanner(mgv7File)
			for scanner.Scan() {
				line := scanner.Text()

				codPlu := line[3:9]
				//codPluArray = append(codPluArray, codPlu)

				codNutriMGV := line[78:84]
				codNutriMGVArray = append(codNutriMGVArray, codNutriMGV)

				//codReceita := line[68:74]
				//receitaArray = append(receitaArray, codReceita)

				textModified := line[0:43] + strings.Repeat(" ", 25) + codPlu + line[74:150] +
					"000000|01|                                                                      0000000000000000000000000||0||0000000000000000000000"
				fmt.Fprintln(systelFile, textModified)
			}

			// Read and process infonutriFile
			scanner = bufio.NewScanner(infonutriFile)
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) > 10 {
					text_nutri := Newnutri(line, codNutriMGVArray, codNutriArray)
					fmt.Println(len(text_nutri))
					if text_nutri != "" {
						fmt.Fprintln(nutriFile, text_nutri)
					}
				}
			}

			dict_nutri := NutriAnalyse(fmt.Sprintf("%s/nutriSystel.TXT", nomeDaPasta))

			fmt.Println("Data analysis completed successfully.")

			//ip_db := strings.TrimSpace(ip)
			//err_import = infoSystelWriter(arquivo, info, dict_forn, dict_aler, dict_fraciona, dict_conserva, ip_db)
			// Balancas, err := models.BuscaBalancas()
			// if err != nil {
			// 	fmt.Println(err_import)
			// }

			var ips_nao_enviados []string
			var ips_enviados []string

			var text string
			var err_import error
			arquivo = fmt.Sprintf("%s/itensSystel.TXT", nomeDaPasta)
			for _, ip := range Balancas {
				err_import = EnviarPluSeparado(arquivo, ip)
				if err_import != nil {
					fmt.Println(err_import)
					text = fmt.Sprintf("informação de produtos enviados corretamente para: %s", ip.Ip)
					ips_nao_enviados = append(ips_nao_enviados, ip.Ip)
					balancas_enviadas = models.RemoveElement(balancas_enviadas, ip)
					balancas_nao_enviadas = append(balancas_nao_enviadas, ip)
				} else {
					text = fmt.Sprintf("informação de produtos enviados corretamente para: %s", ip.Ip)
					ips_enviados = append(ips_enviados, ip.Ip)
				}
				if !somente_preco {
					err_import = EnviarInfoSeparada(arquivo, dict_nutri, info, dict_forn, dict_aler, dict_fraciona, dict_conserva, dict_tara, ip)
					if err_import != nil {
						fmt.Println(err_import)
						text = text + "\n" + fmt.Sprintf("informações extras não enviadas corretamente para: %s", ip.Ip)
						ips_enviados = append(ips_nao_enviados, ip.Ip)
						fmt.Println(balancas_enviadas)
						balancas_enviadas = models.RemoveElement(balancas_enviadas, ip)
						balancas_nao_enviadas = append(balancas_nao_enviadas, ip)
						fmt.Println(balancas_enviadas)
					} else {
						text = text + "\n" + fmt.Sprintf("informações extras enviadas corretamente para: %s", ip.Ip)
						ips_enviados = append(ips_enviados, ip.Ip)
					}
				}
			}

			fmt.Println("Chegou aqui tbm: ", err_import)
			if err_import != nil {
				logToFile("log-erro-conexao-go.txt", fmt.Sprintf("erro ao importar para: %s erro: %s\n ", fmt.Sprint(ips_nao_enviados), fmt.Sprint(err_import)))
			}
			time.Sleep(time.Second)
			logToFile("log.txt", fmt.Sprintf("importou para: %s \n ", fmt.Sprint(ips_enviados)))
			fmt.Println("pronto")
			logToFile("log-completo.txt", text)
		}
	}
	return balancas_enviadas, balancas_nao_enviadas
}

func Passo1(itens, receita, nutri, frac, forn, tara, cons, campext1 string) (err error,
	arquivo, nomeDaPasta string,
	dict_nutri map[string][13]string,
	info, dict_forn, dict_aler, dict_fraciona, dict_conserva map[string]string,
	dict_tara map[string]float64) {

	itens = strings.Replace(itens, "\\", "//", -1)
	receita = strings.Replace(receita, "\\", "//", -1)
	nutri = strings.Replace(nutri, "\\", "//", -1)
	frac = strings.Replace(frac, "\\", "//", -1)
	forn = strings.Replace(forn, "\\", "//", -1)
	tara = strings.Replace(tara, "\\", "//", -1)
	cons = strings.Replace(cons, "\\", "//", -1)
	campext1 = strings.Replace(campext1, "\\", "//", -1)

	// var Ehigual bool = false
	// var dataEhora string
	// var criacaoDoarquivo string = ""
	var mensagem string
	{
		var dict_conserva, dict_fraciona, dict_aler, dict_forn map[string]string
		var info map[string]string
		var dict_tara map[string]float64
		var dict_nutri map[string][13]string
		arquivo := itens

		// if fileExists(arquivo) {
		// 	dataEhora, _ = obterDataHoraCriacao(arquivo)
		// 	Ehigual = criacaoDoarquivo == dataEhora
		// } else if fileExists(ArrumaExtensão(arquivo, "txt", "bak")) {
		// 	arquivo = ArrumaExtensão(arquivo, "txt", "bak")
		// 	dataEhora, _ = obterDataHoraCriacao(arquivo)
		// 	Ehigual = criacaoDoarquivo == dataEhora
		// }

		// if !Ehigual {
		// 	criacaoDoarquivo = dataEhora
		_, erro := os.Stat(arquivo)
		endereco := &erro
		if erro == nil {
			time.Sleep(1 * time.Second)

			if _, err := os.Stat(tara); err == nil {
				dict_tara = taraAnalyze(tara)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}
			if _, err := os.Stat(cons); err == nil {
				dict_conserva = conservaAnalyze(cons)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}

			if _, err := os.Stat(frac); err == nil {
				dict_fraciona = fracionaAnalyze(frac)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}

			if _, err := os.Stat(campext1); err == nil {
				dict_aler = alergiaAnalyze(campext1)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}

			if _, err := os.Stat(forn); err == nil {
				dict_forn = fornAnalyze(forn)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}

			if _, err := os.Stat(receita); err == nil {
				info = InfoAnalyze(receita)
			} else {
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}
			arquivoInfonutri := nutri
			if _, err := os.Stat(nutri); err == nil {
				arquivoInfonutri = nutri
			} else {
				nutri = ""
				mensagem = adcionarMensagem(mensagem, "Sem arquivo ou caminho não encontrado")
			}
			fmt.Print(mensagem)

			nomeDaPasta := VerificaPasta()

			new_itens := fmt.Sprintf("%s/itens.TXT", nomeDaPasta)
			*endereco = itensAnalyze(arquivo, new_itens)
			mgv7File, err := os.Open(new_itens)
			if err != nil {
				fmt.Println(err)
			}

			defer mgv7File.Close()

			infonutriFile, err := os.Open(arquivoInfonutri)
			if err != nil {
				fmt.Println(err)
			}
			defer infonutriFile.Close()

			// Create output files
			nutriFile, err := os.Create(fmt.Sprintf("%s/nutriSystel.TXT", nomeDaPasta))
			if err != nil {
				fmt.Println(err)
			}
			defer nutriFile.Close()

			systelFile, err := os.Create(fmt.Sprintf("%s/itensSystel.TXT", nomeDaPasta))
			if err != nil {
				fmt.Println(err)
			}
			defer systelFile.Close()

			// Initialize arrays
			//var codPluArray []string
			var codNutriArray []string
			var codNutriMGVArray []string
			//var receitaArray []string

			// Read and process mgv7File
			scanner := bufio.NewScanner(mgv7File)
			for scanner.Scan() {
				line := scanner.Text()

				codPlu := line[3:9]
				//codPluArray = append(codPluArray, codPlu)

				codNutriMGV := line[78:84]
				codNutriMGVArray = append(codNutriMGVArray, codNutriMGV)

				//codReceita := line[68:74]
				//receitaArray = append(receitaArray, codReceita)

				textModified := line[0:43] + strings.Repeat(" ", 25) + codPlu + line[74:150] +
					"000000|01|                                                                      0000000000000000000000000||0||0000000000000000000000"

				//textModified = strings.ReplaceAll(textModified, "�", "   ")
				fmt.Fprintln(systelFile, textModified)
			}

			// Read and process infonutriFile
			scanner = bufio.NewScanner(infonutriFile)
			for scanner.Scan() {
				line := scanner.Text()
				if len(line) > 10 {
					text_nutri := Newnutri(line, codNutriMGVArray, codNutriArray)
					fmt.Println(len(text_nutri))
					if text_nutri != "" {
						fmt.Fprintln(nutriFile, text_nutri)
					}
				}
			}

			dict_nutri := NutriAnalyse(fmt.Sprintf("%s/nutriSystel.TXT", nomeDaPasta))

			fmt.Println("Data analysis completed successfully.")
			return *endereco, arquivo, nomeDaPasta, dict_nutri, info, dict_forn, dict_aler, dict_fraciona, dict_conserva, dict_tara
		} else {

			fmt.Println(err)
		}
		return err, arquivo, nomeDaPasta, dict_nutri, info, dict_forn, dict_aler, dict_fraciona, dict_conserva, dict_tara
	}

}

func Passo2(arquivo, nomeDaPasta string, Balancas []models.Balanca) ([]models.Balanca, []models.Balanca) {

	var balancas_enviadas []models.Balanca
	balancas_enviadas = append(balancas_enviadas, Balancas...)
	var balancas_nao_enviadas []models.Balanca
	var ips_nao_enviados []string
	var ips_enviados []string

	var text string
	var err_import error
	arquivo = fmt.Sprintf("%s/itensSystel.TXT", nomeDaPasta)
	for _, ip := range Balancas {
		err_import = EnviarPluSeparado(arquivo, ip)
		if err_import != nil {
			fmt.Println(err_import)
			text = fmt.Sprintf("informação de produtos enviados corretamente para: %s", ip.Ip)
			ips_nao_enviados = append(ips_nao_enviados, ip.Ip)
			balancas_enviadas = models.RemoveElement(balancas_enviadas, ip)
			balancas_nao_enviadas = append(balancas_nao_enviadas, ip)
		} else {
			text = fmt.Sprintf("informação de produtos enviados corretamente para: %s", ip.Ip)
			ips_enviados = append(ips_enviados, ip.Ip)
		}
	}
	fmt.Println("Chegou aqui tbm: ", err_import)
	if err_import != nil {
		logToFile("log-erro-conexao-go.txt", fmt.Sprintf("erro ao importar para: %s erro: %s\n ", fmt.Sprint(ips_nao_enviados), fmt.Sprint(err_import)))
	}
	time.Sleep(time.Second)
	logToFile("log.txt", fmt.Sprintf("importou para: %s \n ", fmt.Sprint(ips_enviados)))
	fmt.Println("pronto")
	logToFile("log-completo.txt", text)
	return balancas_enviadas, balancas_nao_enviadas

}

func Passo3(
	arquivo, nomeDaPasta string,
	Balancas []models.Balanca,
	dict_nutri map[string][13]string,
	info, dict_forn, dict_aler, dict_fraciona, dict_conserva map[string]string,
	dict_tara map[string]float64,
) ([]models.Balanca, []models.Balanca) {

	var wg sync.WaitGroup
	balancasEnviadasCh := make(chan models.Balanca, len(Balancas))
	balancasNaoEnviadasCh := make(chan models.Balanca, len(Balancas))
	var text string
	var ipsNaoEnviados []string
	var ipsEnviados []string
	arquivo = fmt.Sprintf("%s/itensSystel.TXT", nomeDaPasta)

	for _, ip := range Balancas {
		wg.Add(1)
		go func(ip models.Balanca) {
			defer wg.Done()
			errImport := EnviarInfoSeparada(arquivo, dict_nutri, info, dict_forn, dict_aler, dict_fraciona, dict_conserva, dict_tara, ip)
			if errImport != nil {
				fmt.Println(errImport)
				text += fmt.Sprintf("\ninformações extras não enviadas corretamente para: %s", ip.Ip)
				ipsNaoEnviados = append(ipsNaoEnviados, ip.Ip)
				balancasNaoEnviadasCh <- ip
			} else {
				text += fmt.Sprintf("\ninformações extras enviadas corretamente para: %s", ip.Ip)
				ipsEnviados = append(ipsEnviados, ip.Ip)
				balancasEnviadasCh <- ip
			}
		}(ip)
	}

	wg.Wait()
	close(balancasEnviadasCh)
	close(balancasNaoEnviadasCh)

	var balancasEnviadas, balancasNaoEnviadas []models.Balanca
	for balanca := range balancasEnviadasCh {
		balancasEnviadas = append(balancasEnviadas, balanca)
	}
	for balanca := range balancasNaoEnviadasCh {
		balancasNaoEnviadas = append(balancasNaoEnviadas, balanca)
	}

	if len(ipsNaoEnviados) > 0 {
		logToFile("log-erro-conexao-go.txt", fmt.Sprintf("erro ao importar para: %s\n", fmt.Sprint(ipsNaoEnviados)))
	}
	time.Sleep(time.Second)
	logToFile("log.txt", fmt.Sprintf("importou para: %s\n", fmt.Sprint(ipsEnviados)))
	fmt.Println("pronto")
	logToFile("log-completo.txt", text)

	return balancasEnviadas, balancasNaoEnviadas
}

// Helper function to convert string to integer

// Helper function to check if an integer is in an array
