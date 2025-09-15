package jsonapi

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

func setJsonHeader(writer http.ResponseWriter){
	writer.Header().Set("Content-Type","application/json; charset=utf-8")
}


func FromJson[T any](body io.Reader, target T){
	buffer := new(bytes.Buffer)
	buffer.ReadFrom(body)
	json.Unmarshal(buffer.Bytes(),&target)
}

func returnJson[T any](writer http.ResponseWriter, withData func()(T,error)){
	setJsonHeader(writer)
	data, serverErr := withData()

	if serverErr != nil {
		writer.WriteHeader(500)
		serverErrJson, err := json.Marshal(&serverErr)
		if err != nil {
			log.Print(err)
			return 
		}
		writer.Write(serverErrJson)
		return 
	}

	dataJson, err := json.Marshal(&data)
	if err != nil {
		log.Print(err)
		writer.WriteHeader(500)
		return 
	}

	writer.Write(dataJson)
}


func returnErr(writer http.ResponseWriter, err error, code int){
	returnJson(writer, func()(any,error){
		errorMessage := struct {
			Err string 
		}{
			Err: err.Error(),
		}
		writer.WriteHeader(code)
		return errorMessage,nil
	})	
}
