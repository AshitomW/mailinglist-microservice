package jsonapi

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"mailinglist-ms/maildb"
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



func CreateEmail(db *sql.DB) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter,req *http.Request){
		if req.Method != "POST" {
			return 
		}

		entry := maildb.EmailEntry{}
		FromJson(req.Body,&entry)


		if err := maildb.CreateEmail(db,entry.Email); err != nil {
			returnErr(writer,err,400)
			return 
		}



		returnJson(writer ,func()(any, error){
			log.Printf("Json CreateEmail: %v\n", entry.Email)
			return maildb.GetEmail(db,entry.Email)
		})

	})
}


func GetEmail(db *sql.DB) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter,req *http.Request){
		if req.Method != "GET" {
			return 
		}

		entry := maildb.EmailEntry{}
		FromJson(req.Body,&entry)
		returnJson(writer ,func()(any, error){
			log.Printf("Json GetEmail: %v\n", entry.Email)
			return maildb.GetEmail(db,entry.Email)
		})

	})
}

func UpdateEmail(db *sql.DB) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter,req *http.Request){
		if req.Method != "PUT" {
			return 
		}

		entry := maildb.EmailEntry{}
		FromJson(req.Body,&entry)


		if err := maildb.UpdateEmail(db,entry); err != nil {
			returnErr(writer,err,400)
			return 
		}



		returnJson(writer ,func()(any, error){
			log.Printf("Json UpdateEmail: %v\n", entry.Email)
			return maildb.GetEmail(db,entry.Email)
		})

	})
}

func DeleteEmail(db *sql.DB) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter,req *http.Request){
		if req.Method != "POST" {
			return 
		}

		entry := maildb.EmailEntry{}
		FromJson(req.Body,&entry)


		if err := maildb.DeleteEmail(db,entry.Email); err != nil {
			returnErr(writer,err,400)
			return 
		}



		returnJson(writer ,func()(any, error){
			log.Printf("Json DeleteEmail: %v\n", entry.Email)
			return maildb.GetEmail(db,entry.Email)
		})

	})
}

func GetEmailBatch(db *sql.DB) http.Handler{
	return http.HandlerFunc(func(writer http.ResponseWriter,req *http.Request){
		if req.Method != "GET" {
			return 
		}
		
		queryOptions := maildb.GetEmailBatchQueryParameters{}
		FromJson(req.Body, &queryOptions)

		if queryOptions.Count <=0 || queryOptions.Page <=0 {
			returnErr(writer,errors.New("Page and count fields are required and must be greater than 0"),400)
			return 
		}




		returnJson(writer, func()(any,error){
			log.Printf("JSOn GetEmailBatch: %v\n ",queryOptions)
			return maildb.GetEmailBatch(db,queryOptions)
		})


	})
}


func Server(db *sql.DB, bind string){
	http.Handle("/email/create",CreateEmail(db))
	http.Handle("/email/get",GetEmail(db))
	http.Handle("/email/get_batch",GetEmailBatch(db))
	http.Handle("/email/update",UpdateEmail(db))
	http.Handle("/email/delete",DeleteEmail(db))


	err := http.ListenAndServe(bind,nil)
	if err != nil {
		log.Fatalf("JSON server error :%v",err)
	}
}
