package routes

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"

	"github.com/lucasbyte/go-clipse/controllers"
)

func CarregaRotas(fs embed.FS) {
	// Carregar e compilar templates embutidos
	tmpl, err := template.ParseFS(fs, "templates/*.html")
	if err != nil {
		fmt.Printf("failed to parse templates: %v\n", err)
		return
	}

	// Passar os templates compilados para os controladores
	controllers.SetTemplates(tmpl)

	http.HandleFunc("/", controllers.Index)
	http.HandleFunc("/file", controllers.File)
	http.HandleFunc("/balancas", controllers.Balanca)
	http.HandleFunc("/add-balanca", controllers.New_balanca)
	http.HandleFunc("/new-balanca", controllers.Add_balanca)
	http.HandleFunc("/edit-balanca", controllers.Edit_balanca)
	http.HandleFunc("/delete", controllers.Delete_bal)
	http.HandleFunc("/updatesetor", controllers.Update_setor)
	http.HandleFunc("/setor", controllers.Setores)
	http.HandleFunc("/push", controllers.EnviarDadosPg)
	http.HandleFunc("/importar", controllers.ToImport)
	http.HandleFunc("/txitens", controllers.EnviarTxitensPg)
	http.HandleFunc("/auto_add_files", controllers.PuxarArquivos)
	http.HandleFunc("/edit_filepath", controllers.EditaArquivo)
	http.HandleFunc("/push-Cad", controllers.EnviarCadPg)
	http.HandleFunc("/erro", controllers.ErroLeitura)
	http.HandleFunc("/can-import", controllers.CanImport)
	http.HandleFunc("/set-scheduler", controllers.CheckImportAuto)
	http.HandleFunc("/can-scheduler", controllers.StartPageImport)

	content := http.FileServer(http.FS(fs))
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Serving static file: %s\n", r.URL.Path)
		http.StripPrefix("/static", content).ServeHTTP(w, r)
	})
}
