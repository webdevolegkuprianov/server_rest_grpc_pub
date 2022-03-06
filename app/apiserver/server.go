package apiserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	logger "github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/logger"
	"github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/model"
	"github.com/webdevolegkuprianov/server_rest_grpc/app/apiserver/store"
)

//errors
var (
	errIncorrectEmailOrPassword = errors.New("incorrect auth")
	errReg                      = errors.New("service registration error")
	errJwt                      = errors.New("token error")
	errFindUser                 = errors.New("user not found")
)

//server configure
type server struct {
	router *mux.Router
	store  store.Store
	config *model.Service
}

func newServer(store store.Store, config *model.Service) *server {
	s := &server{
		router: mux.NewRouter(),
		store:  store,
		config: config,
	}
	s.configureRouter()
	return s
}

//write new token struct
func newToken(token string, exp time.Time) *model.Token_exp {
	return &model.Token_exp{
		Token: token,
		Exp:   exp,
	}
}

//write http error
func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})

}

//write http response
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	//open
	s.router.HandleFunc("/authentication", s.handleAuth()).Methods("POST")
	//private
	auth := s.router.PathPrefix("/auth").Subrouter()
	auth.Use(s.middleWare)
	//proto
	auth.HandleFunc("/stream1", s.handleStream1()).Methods("POST")
	auth.HandleFunc("/stream2", s.handleStream2()).Methods("POST")
	auth.HandleFunc("/stream3", s.handleStream3()).Methods("POST")
}

//handle Auth
func (s *server) handleAuth() http.HandlerFunc {

	var req model.User1

	return func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			s.error(w, r, http.StatusBadRequest, errReg)
			logger.ErrorLogger.Println(err)
			return
		}

		u, err := s.store.User().FindUser(req.Email, req.Password)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			logger.ErrorLogger.Println(err)
			return
		}

		token, datetime_exp, err := s.store.User().CreateToken(uint64(u.ID), s.config)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, errJwt)
			logger.ErrorLogger.Println(err)
			return
		}
		token_data := newToken(token, datetime_exp)
		s.respond(w, r, http.StatusOK, token_data)
		logger.InfoLogger.Println("token issued success")

	}

}

//Middleware
func (s *server) middleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		//extract user_id
		user_id, err := s.store.User().ExtractTokenMetadata(r, s.config)
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errJwt)
			logger.ErrorLogger.Println(err)
			return
		}

		if err := s.store.User().FindUserid(user_id.UserId); err != nil {
			s.error(w, r, http.StatusUnauthorized, errFindUser)
			logger.ErrorLogger.Println(err)
			return
		}

		next.ServeHTTP(w, r)

	})

}

//handle stream1
func (s *server) handleStream1() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var bytesBody bytes.Buffer
		_, err := bytesBody.ReadFrom(r.Body)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}

		if err := s.store.Data().SendStreamMessages1(bytesBody.Bytes()); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}
		logger.InfoLogger.Println("data_sent_stream1")

	}

}

//handle stream2
func (s *server) handleStream2() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var bytesBody bytes.Buffer
		_, err := bytesBody.ReadFrom(r.Body)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}

		if err := s.store.Data().SendStreamMessages2(bytesBody.Bytes()); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}
		logger.InfoLogger.Println("data_sent_stream2")

	}

}

//handle stream3
func (s *server) handleStream3() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var bytesBody bytes.Buffer
		_, err := bytesBody.ReadFrom(r.Body)
		if err != nil {
			logger.ErrorLogger.Println(err)
		}

		if err := s.store.Data().SendStreamMessages3(bytesBody.Bytes()); err != nil {
			logger.ErrorLogger.Println(err)
			return
		}
		logger.InfoLogger.Println("data_sent_stream3")

	}

}
