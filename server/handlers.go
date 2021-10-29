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

const QueryError    = "Error querying database"
const OverflowError = "Overflow error"
const EncodingError = "Encoding response error"
const ParseIntError = "Error parsing 'id' parameter"

type WordResponse struct {
   Word  string    `json:"word"`
   Id    int64     `json:"id"`
}

type ErrorResponse struct {
   Error    string   `json:"error"`
}

func errorResponse(w *http.ResponseWriter, err error, errorMessage string) {
   log.Print(err)
   errResp, _ := json.Marshal(ErrorResponse{errorMessage})
   http.Error(*w, string(errResp), http.StatusInternalServerError)
}

func GetRandomWord(w http.ResponseWriter, r *http.Request) {
   scowldb := pdb.NewPostgresDB("")
   vars := mux.Vars(r)
   vrf := vars["vrf"]
   var bigUint big.Int
   bigUint.SetString(vrf, 10) 
   vrfNum, overflow := uint256.FromBig(&bigUint)
   if overflow {
      errorResponse(&w, nil, OverflowError)
      return
   }
   maxId, err := scowldb.MaxId()
   if err != nil {
      errorResponse(&w, err, QueryError)
      return
   }
   maxIdNum, overflow := uint256.FromBig(big.NewInt(maxId))
   if overflow {
      errorResponse(&w, nil, OverflowError)
      return
   }

   var temp uint256.Int
   temp.Mod(vrfNum, maxIdNum)
   id := int64(temp.Uint64() + uint64(1))
   randomWord, err := scowldb.QueryById(id)

   var resp = WordResponse{randomWord, id} 
   response, err := json.Marshal(&resp)
   if err != nil {
      errorResponse(&w, err, EncodingError)
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
      errorResponse(&w, err, ParseIntError)
      return
   }
   word, err := scowldb.QueryById(id)
   if err != nil {
      errorResponse(&w, err, QueryError)
      return
   }

   var resp = WordResponse{word, id} 
   response, err := json.Marshal(resp)
   if err != nil {
      errorResponse(&w, err, EncodingError)
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
      errorResponse(&w, err, QueryError)
      return
   }
 
   var resp = WordResponse{word, id} 
   response, err := json.Marshal(&resp)
   if err != nil {
      errorResponse(&w, err, EncodingError)
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
      errorResponse(&w, err, QueryError)
      return
   }
 
   var resp = make(map[string]int64)
   resp["count"] = count 
   response, err := json.Marshal(&resp)
   if err != nil {
      errorResponse(&w, err, EncodingError)
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
      errorResponse(&w, err, QueryError)
      return
   }
 
   var resp = make(map[string]int64)
   resp["maxid"] = maxid
   response, err := json.Marshal(&resp)
   if err != nil {
      errorResponse(&w, err, EncodingError)
      return
   }
   w.Header().Add("Content-Type", "application/json")
   w.WriteHeader(http.StatusCreated)
   w.Write(response)
}
