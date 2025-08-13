package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type API struct {
	router *mux.Router
}

type NewsFullDetailed struct {
	ID          int
	Title       string
	Content     string
	AuthorID    int
	AuthorName  string
	CreatedAt   int64
	PublishedAt int64
}

type NewsShortDetailed struct {
	ID         int
	Title      string
	AuthorName string
}

type Comment struct {
	ID          int    `json:"id"`
	Content     string `json:"content"`
	NewsID      int    `json:"news_id"`
	AuthorID    int    `json:"author_id"`
	AuthorName  string `json:"author_name"`
	CreatedAt   int64  `json:"created_at"`
	PublishedAt int64  `json:"published_at"`
}

func New() *API {
	api := API{}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (a *API) endpoints() {
	// вывод списка новостей
	a.router.HandleFunc("/news", a.newsHandler).Methods(http.MethodGet, http.MethodOptions)
	// фильтр новостей
	a.router.HandleFunc("/news/filter", a.newsFilteredHandler).Methods(http.MethodGet, http.MethodOptions)
	// получение детальной новости
	a.router.HandleFunc("/news/{n}", a.newsByIdHandler).Methods(http.MethodGet, http.MethodOptions)
	// добавление комментария
	a.router.HandleFunc("/news/{n}/addComment", a.newsAddCommentHandler).Methods(http.MethodPost, http.MethodOptions)
}

func (a *API) Router() *mux.Router {
	return a.router
}

func (a *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	news := []NewsShortDetailed{
		{
			ID:         1,
			Title:      "My new role",
			AuthorName: "Leonard Nimoy",
		},
		{
			ID:         2,
			Title:      "Top 10 songs to listen",
			AuthorName: "Howard Shore",
		},
		{
			ID:         3,
			Title:      "Brand New Book Published",
			AuthorName: "Lev Tolstoy",
		},
	}

	bytes, err := json.Marshal(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (a *API) newsFilteredHandler(w http.ResponseWriter, r *http.Request) {
	news := []NewsShortDetailed{
		{
			ID:         2,
			Title:      "Top 10 songs to listen",
			AuthorName: "Howard Shore",
		},
		{
			ID:         3,
			Title:      "Brand New Book Published",
			AuthorName: "Lev Tolstoy",
		},
	}

	bytes, err := json.Marshal(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (a *API) newsByIdHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["n"]
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	news := NewsFullDetailed{
		ID:          n,
		Title:       "Top 10 songs to listen",
		Content:     "The top 10 songs this week, according to various music charts, include \"Die With A Smile\" by Lady Gaga & Bruno Mars, \"Espresso\" by Sabrina Carpenter, \"Please Please Please\" by Sabrina Carpenter, \"A Bar Song (Tipsy)\" by Shaboozey, \"BIRDS OF A FEATHER\" and \"WILDFLOWER\" by Billie Eilish, \"luther\" by Kendrick Lamar & SZA, \"Ordinary\" by Alex Warren, \"Messy\" by Lola Young, and \"Golden\" by HUNTR/X.",
		AuthorID:    2,
		AuthorName:  "Howard Shore",
		CreatedAt:   1755093923,
		PublishedAt: 1755093923,
	}

	bytes, err := json.Marshal(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}

func (a *API) newsAddCommentHandler(w http.ResponseWriter, r *http.Request) {
	s := mux.Vars(r)["n"]
	n, err := strconv.Atoi(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var comm Comment
	err = json.NewDecoder(r.Body).Decode(&comm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	comm.NewsID = n

	bytes, err := json.Marshal(comm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}
