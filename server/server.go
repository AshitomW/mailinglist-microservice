package main

import (
	"database/sql"
	"log"
	"mailinglist-ms/jsonapi"
	"mailinglist-ms/maildb"
	"sync"

	"github.com/alexflint/go-arg"
)


var args struct {
	DBPath string `arg:"env:MAILINGLIST_DB"` // specify or from the environment file
	BindJson string `arg:"env:MAILINGLIST_BIND_JSON"`
}



func main(){
	arg.MustParse(&args)

	if args.DBPath == "" {
		args.DBPath = "list.db"
	}
	if args.BindJson == ""{
		args.BindJson = ":8080" // localhost -> no ip 
	}


	log.Printf("using database '%v'\n",args.DBPath)


	db,err := sql.Open("sqlite3",args.DBPath)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	maildb.TryCreate(db)

	var wg sync.WaitGroup


	wg.Add(1)
	go func(){
		log.Printf("Starting JSON API server ....\n")
		jsonapi.Serve(db,args.BindJson)
		wg.Done()
	}()

	wg.Wait()


}
