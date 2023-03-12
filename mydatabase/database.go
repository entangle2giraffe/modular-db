package mydatabase

import (
  "database/sql"
  "fmt"
  "os"

  _ "github.com/lib/pq"
  "github.com/joho/godotenv"
  "github.com/google/uuid"
)

// Define a Database struct to manage the database connection
type Database struct {
  db *sql.DB
}

func NewDatabaseFromEnv() (*Database, error) {
  err := godotenv.Load(".env")
  if err != nil {
    return nil, fmt.Errorf("failed to load environment variables: %w", err)
  }

  connStr := fmt.Sprintf(
    "user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
    os.Getenv("DB_USER"),
    os.Getenv("DB_PASSWORD"),
    os.Getenv("DB_NAME"),
    os.Getenv("DB_HOST"),
    os.Getenv("DB_PORT"),
  )
  
  db, err := sql.Open("postgres", connStr)
  if err != nil {
    return nil, fmt.Errorf("failed to connect to database: %w", err)
  }

  err = db.Ping()
  if err != nil {
    return nil, fmt.Errorf("failed to ping database: %w", err)
  }

  return &Database{db}, nil
}

func (d *Database) Close() {
  d.db.Close()
}

func (d *Database) Create(name string, age int) error  {
  id := uuid.New().String()
  _, err := d.db.Exec(`
    INSERT INTO mytable (id, name, age) VALUES ($1, $2, $3)
  `, id, name, age)
  if err != nil {
    return fmt.Errorf("failed to create record: %w", err)
  }

  return nil
}

func (d *Database) Update(id uuid.UUID, name string, age int) error {
  _, err := d.db.Exec(`
    UPDATE mytable SET name=$1, age=$2 WHERE id=$3
  `, name, age, id)
  if err != nil {
    return fmt.Errorf("failed to update record: %w", err)
  }

  return nil
}

func (d *Database) Delete(id uuid.UUID) error {
  _, err := d.db.Exec(`
    DELETE FROM mytable WHERE id=$1
  `,id)
  if err != nil {
    return fmt.Errorf("failed to delete record: %w", err)
  }

  return nil
}

func (d *Database) GetById(id uuid.UUID) (*Record, error) {
  row := d.db.QueryRow(`
    SELECT id, name, age FROM mytable WHERE id=$1
  `, id)

  var record Record
  err := row.Scan(&record.Id, &record.Name, &record.Age)
  if err != nil {
    if err == sql.ErrNoRows {
      return nil, nil
    }
    return nil, fmt.Errorf("failed to get record: %w", err)
  }

  return &record, nil
}

func (d *Database) GetAll() ([]*Record, error) {
  rows, err := d.db.Query(`
    SELECT id, name, age FROM mytable
  `)
  if err != nil {
    return nil, fmt.Errorf("failed to get all records: %w", err)
  }
  defer rows.Close()

  var records []*Record
  for rows.Next() {
    var record Record
    err := rows.Scan(&record.Id, &record.Name, &record.Age)
    if err != nil {
      return nil, fmt.Errorf("failed to scan record: %w", err)
    }
    records = append(records, &record)
  }

  return records, nil
}

type Record struct {
  Id uuid.UUID
  Name string
  Age int
}
