package main

import (

   "os"
   "bufio"
   "strings"
   "io/ioutil"
   postgresdb "github.com/fuzzylemma/scowldb/postgresdb"
)

func populatePostgresDB() {
   pdb := NewPostgresDB("")
   pdb.init()

   scowlFinalPath := "scowl-2020.12.07/final"
   abbreviations := "abbreviations"
   files, err := ioutil.ReadDir(scowlFinalPath)
   check(err)
   addToDB := []string{}
   for _, file := range files {
      if !strings.Contains(file.Name(), abbreviations) && !file.IsDir() {
         addToDB = append(addToDB, scowlFinalPath+"/"+file.Name())
      }
   }

   for _, f := range addToDB {
      fmt.Println(f)
      file, err := os.Open(f)
      check(err)

      scanner := bufio.NewScanner(file)
      scanner.Split(bufio.ScanLines)
      for scanner.Scan() {
         word := scanner.Text()
         if (len(word) == 0) { continue }
         pdb.insertWord(word)
      }
      
   }
}


func main() {
   populatePostgresDB()
}
