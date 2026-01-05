package main

import (
	"database/sql"
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

	log.Println("Servidor rodando em http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

//
// ================== LOGIN ==================
//

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

	w.Header().Set("Content-Type", "application/json")

	if err != nil || senha != senhaNoBanco {
		w.Write([]byte(`{"success": false, "message": "Usuário ou senha inválidos"}`))
		return
	}

	w.Write([]byte(`{"success": true}`))
}

//
// ================== REGISTRO ==================
//

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

func processarRegisto(w http.ResponseWriter, r *http.Request) {
	novoUsuario := r.FormValue("username")
	novaSenha := r.FormValue("password")

	if novoUsuario == "" || novaSenha == "" {
		http.Error(w, "Preencha todos os campos", http.StatusBadRequest)
		return
	}

	// Verifica se usuário já existe
	var id int
	err := db.QueryRow(
		"SELECT id FROM usuarios WHERE username = ?",
		novoUsuario,
	).Scan(&id)

	if err == nil {
		http.Error(w, "Usuário já existe", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		"INSERT INTO usuarios (username, password) VALUES (?, ?)",
		novoUsuario,
		novaSenha,
	)

	if err != nil {
		http.Error(w, "Erro ao criar conta", http.StatusInternalServerError)
		return
	}

	// Redireciona para login
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

//
// ================== PÁGINA INICIAL ==================
//

func paginaInicio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "templates/pagina_inicio.html")
}
