package pdb 

import (
   "os"
   "fmt"
   "log"
   "errors"
   "database/sql"
   _ "github.com/lib/pq"
   "github.com/spf13/viper"
   "github.com/fuzzylemma/scowldb/utils"
)

const WORD_TABLE string = `
   CREATE TABLE IF NOT EXISTS words (
      id INT NOT NULL UNIQUE PRIMARY KEY,
      word TEXT NOT NULL UNIQUE 
   );
`
type PostgresDB struct {
   Host string
   Port string
   Name string
   User string
   Password string
}



func NewPostgresDB(pathToConfig string) *PostgresDB {
   var host, port, dbname, user, password string
   if pathToConfig == "" {
      host = os.Getenv("POSTGRES_HOST")
      port = os.Getenv("POSTGRES_PORT")
      dbname = os.Getenv("POSTGRES_DB")
      user = os.Getenv("POSTGRES_USER")
      password = os.Getenv("POSTGRES_PASSWORD")
   } else {
      viper.SetConfigName("config")
      viper.AddConfigPath(pathToConfig)
      err := viper.ReadInConfig()
      utils.Check(err)

      host = viper.GetString("sql.host")
      port = viper.GetString("sql.port")
      dbname = viper.GetString("sql.dbname")
      user = viper.GetString("sql.user")
      password = viper.GetString("sql.password")
   }
   return &PostgresDB{host, port, dbname, user, password}
}

func (pdb *PostgresDB) Init() {
   pdb.CreateWordTable()
}

func (pdb *PostgresDB) createTable(tableStr string) {
   db, err := pdb.OpenConnection()
   defer db.Close()
   utils.Check(err)
   _, e := db.Exec(tableStr)
   utils.Check(e)
}

func (pdb *PostgresDB) CreateWordTable() {
   pdb.createTable(WORD_TABLE)
}

func (pdb *PostgresDB) psqlConnString() string {
   if pdb.Port == "" {
      return fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable", pdb.Host, pdb.User, pdb.Password, pdb.Name)
   }
   return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", pdb.Host, pdb.Port, pdb.User, pdb.Password, pdb.Name)
}

func (pdb *PostgresDB) OpenConnection() (*sql.DB, error) {
   psqlConn := pdb.psqlConnString()
   return sql.Open("postgres", psqlConn)
}

func (pdb *PostgresDB) InsertWord(word string) {
   db, err := pdb.OpenConnection()
   defer db.Close()
   utils.Check(err)

   id, err := pdb.MaxId() 
   utils.Check(err)
   id += 1

   insertQuery := fmt.Sprintf(`INSERT INTO words (id, word) VALUES ($1, $2);`)
   _, e := db.Exec(insertQuery, id, word)
   if e != nil {
      fmt.Println(id, word, e)
   }
}

/** Returns highest id in database */
func (pdb *PostgresDB) MaxId() (int64, error) {
   db, err := pdb.OpenConnection()
   defer db.Close()
   utils.Check(err)

   id := 0 
   if row := db.QueryRow("SELECT MAX(id) FROM words;"); row.Err() == nil {
      if err := row.Scan(&id); err != nil {
         if err != sql.ErrNoRows {
            return 0, nil
         }
      }
   }
   return int64(id), nil
}
/** Returns number of words in database */
func (pdb *PostgresDB) WordCount() (int64, error) {
   db, err := pdb.OpenConnection()
   defer db.Close()
   utils.Check(err)

   count := 0
   if row := db.QueryRow("SELECT count(*) FROM words;"); row.Err() == nil {
      if err := row.Scan(&count); err != nil {
         if err != sql.ErrNoRows {
            return 0, nil
         }
      }
   }
   return int64(count), nil 
}

/** Returns id given word */
func (pdb *PostgresDB) QueryByWord(word string) (int64, error) {
   db, err := pdb.OpenConnection()
   defer db.Close()
   utils.Check(err)

   var id int64
   rows, err := db.Query(`SELECT id FROM words where word = $1`, word)
   if err != nil {
      return -1, err
   }
 
   if rows.Next() {
      rows.Scan(&id)
   } else {
      log.Print("Word '", word + "' not found.")
      return 0, errors.New("Word not found in database.")
   }
   return id, nil
}

/** Returns word given id */
func (pdb *PostgresDB) QueryById(id int64) (string, error) {
   db, err := pdb.OpenConnection()
   defer db.Close()
   utils.Check(err)

   var word string 
   rows, err := db.Query(`SELECT word FROM words where id = $1`, id)
   if err != nil {
      return "", err
   }
   if rows.Next() {
      rows.Scan(&word)
   } else {
      log.Print("Id '", id, "' not found.")
      return "", errors.New("Id not found in database.")
   }
   return word, nil 
}
