package main

import (
	"encoding/json"
	"net/http"
)

func writeJson(w http.ResponseWriter,status int,data any) error{
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter,r *http.Request,data any) error{
	maxByte:=1_048_578 //restrict it to 1 MB
	r.Body=http.MaxBytesReader(w,r.Body,int64(maxByte))
	decoder:=json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(data)
}

func writeJSONError(w http.ResponseWriter,r *http.Request,status int,message string) error{
	type envelope struct{
		Error string `json:"error"`
	}
	return writeJson(w,status,&envelope{Error: message})
}