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

	// 🔹 ROTAS DINÂMICAS PRIMEIRO
	http.HandleFunc("/login", processarLogin)
	http.HandleFunc("/register", processarRegisto)
	http.HandleFunc("/pagina_inicio", paginaInicio)

	// 🔹 ARQUIVOS ESTÁTICOS DEPOIS
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	fmt.Println("Servidor a correr em http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

func paginaInicio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, "pagina_inicio.html")
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
		fmt.Fprint(w, "Utilizador não encontrado.")
		return
	}

	if senha == senhaNoBanco {
		http.Redirect(w, r, "/pagina_inicio", http.StatusSeeOther)
		return
	}

	fmt.Fprint(w, "Senha incorreta.")
}

// --- NOVA FUNÇÃO DE REGISTO ---
func processarRegisto(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	novoUsuario := r.FormValue("username")
	novaSenha := r.FormValue("password")

	// Verificar se os campos não estão vazios
	if novoUsuario == "" || novaSenha == "" {
		fmt.Fprintf(w, "Por favor, preencha todos os campos.")
		return
	}

	// Inserir no Banco de Dados
	// Usamos 'Exec' para comandos que não retornam linhas (como INSERT, UPDATE)
	_, err := db.Exec("INSERT INTO usuarios (username, password) VALUES (?, ?)", novoUsuario, novaSenha)

	if err != nil {
		fmt.Println("Erro SQL:", err)
		fmt.Fprintf(w, "Erro ao criar utilizador. Talvez o nome já exista?")
		return
	}

	fmt.Fprintf(w, "Conta criada com sucesso! <a href='/login.html'>Clique aqui para fazer login</a>")
}
