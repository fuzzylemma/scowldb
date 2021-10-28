package postgresdb

import (

   "log"
   "fmt"

   "database/sql"
   _ "github.com/lib/pq"
   "github.com/spf13/viper"
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


func check (err error) {
   if err != nil {
      log.Fatal(err)
   }
}

func NewPostgresDB(pathToConfig string) *PostgresDB {
   if pathToConfig == "" {
      pathToConfig = "./config"
   }
   viper.SetConfigName("config")
   viper.AddConfigPath(pathToConfig)
   err := viper.ReadInConfig()
   if err != nil {
      log.Fatal(err)
   }
   host := viper.GetString("sql.host")
   port := viper.GetString("sql.port")
   dbname := viper.GetString("sql.dbname")
   user := viper.GetString("sql.user")
   password := viper.GetString("sql.password")
   return &PostgresDB{host, port, dbname, user, password}
}

func (pdb *PostgresDB) init() {
   pdb.CreateWordTable()
}

func (pdb *PostgresDB) createTable(tableStr string) {
   db, err := pdb.OpenConnection()
   defer db.Close()
   check(err)
   _, e := db.Exec(tableStr)
   check(e)
}

func (pdb *PostgresDB) CreateWordTable() {
   pdb.createTable(WORD_TABLE)
}

func (pdb *PostgresDB) psqlConnString() string {
   return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", pdb.Host, pdb.Port, pdb.User, pdb.Password, pdb.Name)
}

func (pdb *PostgresDB) OpenConnection() (*sql.DB, error) {
   psqlConn := pdb.psqlConnString()
   return sql.Open("postgres", psqlConn)
}

func (pdb *PostgresDB) insertWord(word string) {
   db, err := pdb.OpenConnection()
   defer db.Close()
   check(err)

   id := 0 
   if err = db.QueryRow("SELECT MAX(id) FROM words;").Scan(&id); err != nil {
      if err == sql.ErrNoRows {
         check(err)
      }
   }
   id += 1
   insertQuery := fmt.Sprintf(`INSERT INTO words (id, word) VALUES ($1, $2);`)
   _, e := db.Exec(insertQuery, id, word)
   if e != nil {
      fmt.Println(id, insertQuery)
      fmt.Println(e)
   }
}
