package server

import (

   "log"
   "net/http"
   "math/big"
   "strconv"
   "encoding/json"

   "github.com/holiman/uint256"
   "github.com/fuzzylemma/scowldb/pdb"
   "github.com/gorilla/mux"

)


func GetRandomWord(w http.ResponseWriter, r *http.Request) {
   scowldb := pdb.NewPostgresDB("")
   vars := mux.Vars(r)
   vrf := vars["vrf"]
   var bigUint big.Int
   bigUint.SetString(vrf, 10) 
   vrfNum, overflow := uint256.FromBig(&bigUint)
   if overflow {
      log.Print("Overflow on random number")
      http.Error(w, "Overflow error", http.StatusInternalServerError)
      return
   }
   maxId, err := scowldb.MaxId()
   if err != nil {
      http.Error(w, "Error during query", http.StatusInternalServerError)
      return
   }
   maxIdNum, overflow := uint256.FromBig(big.NewInt(maxId))
   if overflow {
      log.Print("Overflow on random number")
      http.Error(w, "Overflow error", http.StatusInternalServerError)
      return
   }

   var temp uint256.Int
   temp.Mod(vrfNum, maxIdNum)
   id := int64(temp.Uint64() + uint64(1))
   randomWord, err := scowldb.QueryById(id)

   var resp = make(map[string]string)
   resp["word"] = randomWord
   response, err := json.Marshal(&resp)
   if err != nil {
      log.Print(err)
      http.Error(w, "Response encoding error", http.StatusInternalServerError)
      return
   }
   w.Header().Add("Content-Type", "application/json")
   w.WriteHeader(http.StatusCreated)
   w.Write(response)
}

func GetWordById(w http.ResponseWriter, r *http.Request) {
   scowldb := pdb.NewPostgresDB("")

   vars := mux.Vars(r)
   id, err := strconv.ParseInt(vars["id"], 10, 64)
   if err != nil {
      log.Print(err)
      http.Error(w, "Error parsing `id` parameter", http.StatusInternalServerError)
      return
   }
   word, err := scowldb.QueryById(id)
   if err != nil {
      log.Print(err)
      http.Error(w, "Error query database", http.StatusInternalServerError)
      return
   }

   var resp = make(map[int64]string)
   resp[id] = word 
   response, err := json.Marshal(&resp)
   if err != nil {
      log.Print(err)
      http.Error(w, "Response encoding error", http.StatusInternalServerError)
      return
   }
   w.Header().Add("Content-Type", "application/json")
   w.WriteHeader(http.StatusCreated)
   w.Write(response)
}

func GetIdByWord(w http.ResponseWriter, r *http.Request) {
   scowldb := pdb.NewPostgresDB("")

   vars := mux.Vars(r)
   word := vars["word"]
   id, err := scowldb.QueryByWord(word)
   if err != nil {
      log.Print(err)
      http.Error(w, "Error query database", http.StatusInternalServerError)
      return
   }
 
   var resp = make(map[string]int64)
   resp[word] = id
   response, err := json.Marshal(&resp)
   if err != nil {
      log.Print(err)
      http.Error(w, "Response encoding error", http.StatusInternalServerError)
      return
   }
   w.Header().Add("Content-Type", "application/json")
   w.WriteHeader(http.StatusCreated)
   w.Write(response)

}

func WordCount(w http.ResponseWriter, r *http.Request) {
   scowldb := pdb.NewPostgresDB("")

   count, err := scowldb.WordCount()
   if err != nil {
      log.Print(err)
      http.Error(w, "Error query database", http.StatusInternalServerError)
      return
   }
 
   var resp = make(map[string]int64)
   resp["count"] = count 
   response, err := json.Marshal(&resp)
   if err != nil {
      log.Print(err)
      http.Error(w, "Response encoding error", http.StatusInternalServerError)
      return
   }
   w.Header().Add("Content-Type", "application/json")
   w.WriteHeader(http.StatusCreated)
   w.Write(response)
}

func MaxId(w http.ResponseWriter, r *http.Request) {
   scowldb := pdb.NewPostgresDB("")

   maxid, err := scowldb.MaxId()
   if err != nil {
      log.Print(err)
      http.Error(w, "Error query database", http.StatusInternalServerError)
      return
   }
 
   var resp = make(map[string]int64)
   resp["maxid"] = maxid
   response, err := json.Marshal(&resp)
   if err != nil {
      log.Print(err)
      http.Error(w, "Response encoding error", http.StatusInternalServerError)
      return
   }
   w.Header().Add("Content-Type", "application/json")
   w.WriteHeader(http.StatusCreated)
   w.Write(response)
}
