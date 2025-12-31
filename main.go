package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func main() {
	var err error

	dsn := "root:@tcp(127.0.0.1:3306)/meu_site"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err = db.Ping(); err != nil {
		log.Fatal("Erro ao conectar ao banco: ", err)
	}

	// 🔹 ROTAS
	http.HandleFunc("/", loginPage)
	http.HandleFunc("/login", processarLogin)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/pagina_inicio", paginaInicio)

	// 🔹 ARQUIVOS ESTÁTICOS
	http.Handle("/static/",
		http.StripPrefix("/static/",
			http.FileServer(http.Dir("static")),
		),
	)

	fmt.Println("Servidor a correr em http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func processarLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	usuario := r.FormValue("username")
	senha := r.FormValue("password")

	var senhaNoBanco string
	err := db.QueryRow(
		"SELECT password FROM usuarios WHERE username = ?",
		usuario,
	).Scan(&senhaNoBanco)

	if err != nil {
		http.Error(w, "Utilizador não encontrado", http.StatusUnauthorized)
		return
	}

	if senha == senhaNoBanco {
		http.Redirect(w, r, "/pagina_inicio", http.StatusSeeOther)
		return
	}

	http.Error(w, "Senha incorreta", http.StatusUnauthorized)
}

//REGISTRO
func processarRegisto(w http.ResponseWriter, r *http.Request) {
	novoUsuario := r.FormValue("username")
	novaSenha := r.FormValue("password")

	if novoUsuario == "" || novaSenha == "" {
		http.Error(w, "Preencha todos os campos", http.StatusBadRequest)
		return
	}

		//VERIFICAÇÃO SE JA EXISTE O USUARIO selecionando apenas o id para ser mais rapdio
			var idExistente int 
			err := db.QueryRow("SELECT id FROM usuarios WHERE username=?", novoUsuario).Scan(&idExistente)

			if err == nil {
				fmt.Fprintf(w, "Erro: O nome de utilizador '%s' já está em uso.", novoUsuario)
				return
			}
			


			_, err = db.Exec("INSERT INTO usuarios (username, password) VALUES (?, ?)", novoUsuario, novaSenha)

			if err != nil {
				fmt.Println("Erro SQL:", err)
				fmt.Fprintf(w, "Erro ao criar conta.")
				return
			}
		
			fmt.Fprintf(
				w,
				"Conta criada com sucesso! <a href='/'>Clique aqui para fazer login</a>",
			)
			
		}

	



// LOGIN DE PAGE
func loginPage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "templates/login.html")
}

//PAGINA INICIO
func paginaInicio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "templates/pagina_inicio.html")
}


//REGISTRO DE HANDLER
func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "templates/register.html")
		return
	}

	if r.Method == http.MethodPost {
		processarRegisto(w, r)
		return
	}

	http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
}
