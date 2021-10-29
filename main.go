package main

import (

   "context"
   "os"
   "os/signal"
   "time"
   "fmt"
   "log"
   "bufio"
   "strings"
   "io/ioutil"

   postgresdb "github.com/fuzzylemma/scowldb/pdb"
   server "github.com/fuzzylemma/scowldb/server"
   "github.com/fuzzylemma/scowldb/utils"

)

func populatePostgresDB() {

   pdb := postgresdb.NewPostgresDB("")
   pdb.Init()

   scowlFinalPath := "scowl-2020.12.07/final"
   abbreviations := "abbreviations"
   files, err := ioutil.ReadDir(scowlFinalPath)
   utils.Check(err)
   addToDB := []string{}
   for _, file := range files {
      if !strings.Contains(file.Name(), abbreviations) && !file.IsDir() {
         addToDB = append(addToDB, scowlFinalPath+"/"+file.Name())
      }
   }

   for _, f := range addToDB {
      fmt.Println(f)
      file, err := os.Open(f)
      utils.Check(err)

      scanner := bufio.NewScanner(file)
      scanner.Split(bufio.ScanLines)
      for scanner.Scan() {
         word := scanner.Text()
         if (len(word) == 0) { continue }
         pdb.InsertWord(word)
      }
      
   }
}

func main() {
   var wait time.Duration

   if (false) {
      populatePostgresDB()
   }

   srv := server.NewServer()
   httpSrv := srv.HttpServer()
	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

  	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
   httpSrv.Shutdown(ctx)

	log.Println("shutting down")
}
