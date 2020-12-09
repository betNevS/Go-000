package controller

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/betNevS/Go-000/Week02/dao"

	"github.com/betNevS/Go-000/Week02/model"

	"github.com/betNevS/Go-000/Week02/biz"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Pares request error: ", err)
	}
	resp := &model.Response{}
	name := r.Form.Get("name")
	user, err := biz.FindUser(name)
	if err != nil {
		if errors.Is(err, dao.ErrNotFound) {
			log.Printf("ErrNoRows: %+v", err)
			w.WriteHeader(http.StatusNotFound)
			resp.Code = 404
			resp.Msg = "user not found"
			SendResp(w, resp)
			return
		} else {
			log.Printf("OtherError: %+v", err)
			w.WriteHeader(http.StatusInternalServerError)
			resp.Code = 500
			resp.Msg = "server not available"
			SendResp(w, resp)
			return
		}
	}
	resp.Code = 200
	resp.Msg = "success"
	resp.Data = user
	SendResp(w, resp)
}

func SendResp(w http.ResponseWriter, response *model.Response) {
	w.Header().Add("Content-type", "application/json")
	resp, err := json.Marshal(response)
	if err != nil {
		log.Println("resp json marshal error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}
