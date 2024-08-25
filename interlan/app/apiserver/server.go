package apiserver

import (
	"encoding/json"
	"fmt"
	"medods/storage"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type server struct {
	router  *mux.Router
	storage storage.Storage
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func newServer(storage storage.Storage) *server {
	s := &server{
		router:  mux.NewRouter(),
		storage: storage,
	}
	s.configureRouter()
	return s

}

func (s *server) configureRouter() {
	s.router.HandleFunc("/auth/{userid:[0-9]+}", s.createTokens()).Methods("POST")
	s.router.HandleFunc("/refresh", s.refreshToken()).Methods("POST")
}

func (s *server) createTokens() http.HandlerFunc {
	type request struct {
		Email string
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		vars := mux.Vars(r)
		strUseridId := vars["userid"]
		userid, err := strconv.ParseInt(strUseridId, 10, 64)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
		}
		ip, err := getIP(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		tokens, err := s.storage.Auth().CreateTokens(userid, ip, req.Email)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusCreated, tokens)
	})
}

func (s *server) refreshToken() http.HandlerFunc {
	type request struct {
		UserId       int64  `json:"userid"`
		RefreshToken string `json:"refreshtoken"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		ip, err := getIP(r)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		tokens, err := s.storage.Auth().RefreshAccess(req.UserId, ip, req.RefreshToken)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusCreated, tokens)
	})
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func getIP(r *http.Request) (string, error) {
	//Get IP from the X-REAL-IP header
	ip := r.Header.Get("X-REAL-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}

	//Get IP from X-FORWARDED-FOR header
	ips := r.Header.Get("X-FORWARDED-FOR")
	splitIps := strings.Split(ips, ",")
	for _, ip := range splitIps {
		netIP := net.ParseIP(ip)
		if netIP != nil {
			return ip, nil
		}
	}

	//Get IP from RemoteAddr
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	netIP = net.ParseIP(ip)
	if netIP != nil {
		return ip, nil
	}
	return "", fmt.Errorf("no valid ip found")
}
