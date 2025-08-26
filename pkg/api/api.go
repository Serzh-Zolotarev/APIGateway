package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

const (
	newsServiceURL     = "http://localhost:8080/news"
	commentsServiceURL = "http://localhost:8082"
	censorServiceURL   = "http://localhost:8083"
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
	Comments    []Comment
}

type Pagination struct {
	Page       int
	MaxPages   int
	NewsOnPage int
}

type NewsList struct {
	Pagination Pagination
	News       []NewsShortDetailed
}
type NewsShortDetailed struct {
	ID         int    `json:"id" xml:"id"`
	Title      string `json:"title" xml:"title"`
	AuthorName string `json:"author_name" xml:"author_name"`
}

type Comment struct {
	ID        int    `json:"id"`
	NewsID    int    `json:"news_id"`
	CommentID int    `json:"comment_id"`
	Content   string `json:"content"`
	Author    string `json:"author"`
	PubTime   int64  `json:"pub_time"`
}

func New() *API {
	api := API{}
	api.router = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (a *API) endpoints() {
	a.router.Use(requestIdMiddleware)
	// вывод списка новостей
	a.router.HandleFunc("/", a.newsHandler).Methods(http.MethodGet, http.MethodOptions)
	// получение детальной новости
	a.router.HandleFunc("/{n}", a.newsByIdHandler).Methods(http.MethodGet, http.MethodOptions)
	// добавление комментария
	a.router.HandleFunc("/{n}/addComment", a.newsAddCommentHandler).Methods(http.MethodPost, http.MethodOptions)
	a.router.Use(loggingMiddleware)
}

func (a *API) Router() *mux.Router {
	return a.router
}

func (a *API) newsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	s := query.Get("s")
	page := query.Get("page")

	if len(page) == 0 {
		page = "1"
	}

	url := fmt.Sprintf("%s?page=%s", newsServiceURL, page)

	if len(s) != 0 {
		url += fmt.Sprintf("&s=%s", s)
	}

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	news := NewsList{}

	err = json.Unmarshal(body, &news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func (a *API) newsByIdHandler(w http.ResponseWriter, r *http.Request) {
	news := NewsFullDetailed{}
	comments := []Comment{}
	errChan := make(chan error)

	var wg sync.WaitGroup // Declare a WaitGroup

	wg.Add(2)

	urlPath := r.URL.String()

	go func(wg *sync.WaitGroup, urlPath string, news *NewsFullDetailed) {
		defer wg.Done()

		resp, err := http.Get(newsServiceURL + urlPath)
		if err != nil {
			errChan <- err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errChan <- err
		}

		err = json.Unmarshal(body, &news)
		if err != nil {
			errChan <- err
		}
	}(&wg, urlPath, &news)

	go func(wg *sync.WaitGroup, urlPath string, comments *[]Comment) {
		defer wg.Done()

		resp, err := http.Get(commentsServiceURL + urlPath)
		if err != nil {
			errChan <- err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			errChan <- err
		}

		err = json.Unmarshal(body, &comments)
		if err != nil {
			errChan <- err
		}
	}(&wg, urlPath, &comments)

	go func() {
		for err := range errChan {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}()

	wg.Wait()

	news.Comments = comments

	bytes, err := json.Marshal(news)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
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
	marshal, err := json.Marshal(comm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	censorResp, err := http.Post(censorServiceURL+"/", "application/json", bytes.NewReader(marshal))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if censorResp.StatusCode != http.StatusOK {
		http.Error(w, "failed censoring", censorResp.StatusCode)
		return
	}

	resp, err := http.Post(commentsServiceURL+"/addComment", "application/json", bytes.NewReader(marshal))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
