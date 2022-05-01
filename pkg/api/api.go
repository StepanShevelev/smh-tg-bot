package api

import (
	"encoding/json"
	"errors"
	cfg "github.com/StepanShevelev/smh-tg-bot/pkg/config"
	ipstack "github.com/StepanShevelev/smh-tg-bot/pkg/ipstackapi"
	"net/http"
	"strconv"

	mydb "github.com/StepanShevelev/smh-tg-bot/pkg/db"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitBackendApi(config *cfg.Config) {
	http.HandleFunc("/API/get_users", returnAllUsers)
	http.HandleFunc("/API/get_user", returnSingleUser)
	http.HandleFunc("/API/get_history_by_tg", returnSingleUserHistory)
	http.HandleFunc("/API/delete_history_by_tg", deleteHistoryField)
}

func returnAllUsers(w http.ResponseWriter, r *http.Request) {
	if !isMethodGET(w, r) {
		return
	}
	users, okUsers := apiGiveUsers(w)
	if !okUsers {
		return
	}
	sendData(users, w)
}

func returnSingleUser(w http.ResponseWriter, r *http.Request) {
	if !isMethodGET(w, r) {
		return
	}
	userId, okId := parseId(w, r)
	if !okId {
		return
	}

	user, okUser := apiGiveUserById(userId, w)
	if !okUser {
		return
	}
	sendData(user, w)
}

func returnSingleUserHistory(w http.ResponseWriter, r *http.Request) {
	if !isMethodGET(w, r) {
		return
	}
	userId, okId := parseId(w, r)
	if !okId {
		return
	}

	infoList, okInfo := apiGiveUserHistoryRet(userId, w)
	if !okInfo {
		return
	}
	sendData(infoList, w)
}

func deleteHistoryField(w http.ResponseWriter, r *http.Request) {
	if !isMethodGET(w, r) {
		return
	}
	id, okId := parseId(w, r)
	if !okId {
		return
	}

	hist, okHist := apiGiveUserHistoryByID(id, w)
	if !okHist {
		return
	}

	mydb.Client.Unscoped().Delete(&hist)
	w.WriteHeader(http.StatusOK)
}

func sendData(data interface{}, w http.ResponseWriter) {
	b, err := json.Marshal(data)
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "can't marshal json"}`))
		return
	}
	w.Write(b)
	w.WriteHeader(http.StatusOK)
}

func isMethodGET(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "method not found"}`))
		return false
	}
	return true
}

func parseId(w http.ResponseWriter, r *http.Request) (int, bool) {
	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "arguments params is missing"}`))
		return 0, false
	}
	userId, err := strconv.Atoi(keys[0])
	if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "can't pars id"}`))
		return 0, false
	}
	return userId, true
}

func apiGiveUsers(w http.ResponseWriter) ([]mydb.User, bool) {

	users, err := mydb.Client.GiveUsers()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[]`))
		return []mydb.User{}, false
	} else if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
		return []mydb.User{}, false
	}
	return users, true
}

func apiGiveUserById(userId int, w http.ResponseWriter) (*mydb.User, bool) {

	user, err := mydb.Client.GiveUserByID(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "unrecognized user"}`))
		return nil, false
	} else if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
		return nil, false
	}
	return user, true
}

func apiGiveUserHistoryRet(userId int, w http.ResponseWriter) ([]ipstack.InfoIP, bool) {

	infoList, err := mydb.Client.GiveUserHistoryRet(userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusOK)
		return infoList, true
	} else if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
		return nil, false
	}
	return infoList, true
}

func apiGiveUserHistoryByID(id int, w http.ResponseWriter) (*mydb.UserHistory, bool) {

	hist, err := mydb.Client.GiveUserHistoryByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad history ID"}`))
		return nil, false
	} else if err != nil {
		log.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": ""}`))
		return nil, false
	}
	return hist, true
}
