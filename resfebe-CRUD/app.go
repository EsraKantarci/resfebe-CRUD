package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
)

// App testte çağrılan .App
type App struct {
	Router *mux.Router
	Logger http.Handler
	DB     *sql.DB
}

// curl 127.0.0.1:8000/api/products/list
func (a *App) getResfebeler(w http.ResponseWriter, r *http.Request) {
	rows, err := a.DB.Query("SELECT * FROM resfebeler")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	defer rows.Close()

	var resfebeler []Resfebe
	for rows.Next() {
		var p Resfebe
		err := rows.Scan(&p.ImageID, &p.Word, &p.Difficulty, &p.Category, &p.Date, &p.Language)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}

		resfebeler = append(resfebeler, p)
	}

	_ = json.NewEncoder(w).Encode(resfebeler)
}

// curl --header "Content-Type: application/json" --request POST --data '{"name": "ABC", "manufacturer": "ACME"}' \
// 		127.0.0.1:8000/api/products/new
func (a *App) createResfebe(w http.ResponseWriter, r *http.Request) {
	var p Resfebe
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	defer r.Body.Close()

	_, err := a.DB.Query("INSERT INTO resfebeler (word, imagePath, difficulty, category, date, language) VALUES (?, ?, ?, ?, ?, ?)", p.Word, p.ImagePath, p.Difficulty, p.Category, p.Date, p.Language)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithMessage(w, http.StatusCreated, "New row added.")
}

// curl 127.0.0.1:8000/api/products/10
func (a *App) getResfebe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resfebe ID")
		return
	}

	p := Resfebe{ImageID: id}
	row := a.DB.QueryRow("SELECT word, imagePath, difficulty, category, date, language FROM resfebeler WHERE id=?", p.ImageID)
	if err := row.Scan(&p.Word, &p.ImagePath, &p.Difficulty, &p.Category, &p.Date, &p.Language); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// curl --request PUT --data '{"name": "ABC", "manufacturer": "ACME"}' 127.0.0.1:8000/api/products/11
func (a *App) updateResfebe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//queryler aslında http handlerlar, putu request r olarak alıyor, w de response bizim döndüklerimiz
	//requestteki parametreler mux.Vars ile alıyoruz. key value dönüyor, idyi kontrol ediyoruz
	//id'i al ama string olarak al integera dönüştür a to i ile.
	//go'da null = nill
	//go'da try catch yok, errorları strconv.Atoi id ve err dönüyor, hata varsa err yoksa id döner.
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resfebe ID") //code=statusbadrequest, error: invalid prod...
		return
	}
	//update'tesin, data alıp data yazcan. geçici resfebe p oluştur, sonra bodyde json geldi.
	//key value'ya çevirsin okusun diye decode ediyorsun
	//body'deki veriyi &p'ye göm. error'u initailaize et, error nil değilse error döndür

	var p Resfebe
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	//defer son çalıştırıyor, linkin açık kalmasını engelliyor.
	defer r.Body.Close()
	p.ImageID = id
	//en sona id gönder
	_, err = a.DB.Query("UPDATE resfebeler SET word=?, imagePath=?, difficulty=?, category=?, date=?, language=? WHERE id=?", p.Word, p.ImagePath, p.Difficulty, p.Category, p.Date, p.Language, p.ImageID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

// curl --request DELETE 127.0.0.1:8000/api/products/10
func (a *App) deleteResfebe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//buradaki id ile mux handle'daki id aynı, yani _id ya da ImageID kullanma
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid resfebe ID")
		return
	}

	_, err = a.DB.Query("DELETE FROM resfebeler WHERE id=?", id)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithMessage(w, http.StatusOK, "Deleted.")
}

/* path ekle, route'a pathleri yaz, veritabanının tipi değişecek
https://astaxie.gitbooks.io/build-web-application-with-golang/en/05.3.html
*/

func (a *App) InitializeRoutes() {
	a.Router.HandleFunc("/api/resfebeler/list", a.getResfebeler).Methods("GET")
	a.Router.HandleFunc("/api/resfebeler/new", a.createResfebe).Methods("POST")
	a.Router.HandleFunc("/api/resfebeler/{id:[0-9]+}", a.getResfebe).Methods("GET")
	a.Router.HandleFunc("/api/resfebeler/{id:[0-9]+}", a.updateResfebe).Methods("PUT")
	a.Router.HandleFunc("/api/resfebeler/{id:[0-9]+}", a.deleteResfebe).Methods("DELETE")

}

func (a *App) Initialize() {
	// dataSource := username + ":" + password + "@tcp(" + server + ":" + port + ")/" + dbName
	a.DB, err = sql.Open("sqlite3", "C:\\db\\resfebe.db")
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.Logger = handlers.CombinedLoggingHandler(os.Stdout, a.Router)
	a.InitializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(":"+viper.GetString("Server.port"), a.Logger))
}
