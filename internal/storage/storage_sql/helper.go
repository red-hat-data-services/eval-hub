package storage_sql

// createDatabaseStatement the order or arguments is:
// database_name
func createDatabaseStatement() string {
	return `CREATE DATABASE IF NOT EXISTS ?;`
}

// createTableStatement	the order or arguments is:
// table_name json_field_type
func createTableStatement() string {
	return `CREATE TABLE IF NOT EXISTS ? (
    id      INTEGER PRIMARY KEY AUTOINCREMENT,
    entity 	? NOT NULL
);`
}

// createAddEntityStatement the order or arguments is:
// table_name entity
func createAddEntityStatement() string {
	return `INSERT INTO ? (entity)
	VALUES (?);`
}
