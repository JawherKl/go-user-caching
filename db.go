package main

import (
    "github.com/jackc/pgx/v4"
    "context"
    "log"
    "os"
)

var dbConn *pgx.Conn

func InitPostgres() {
    var err error
    dbConn, err = pgx.Connect(context.Background(), "postgres://test2:test@172.20.0.18:5432/gouserdb")
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
        os.Exit(1)
    }

    log.Println("Connected to PostgreSQL")
}

func GetUserById(userId string) (*User, error) {
    var user User
    log.Printf("Fetching user with ID: %s", userId)  // Add logging
    err := dbConn.QueryRow(context.Background(), "SELECT id, name, email FROM users WHERE id=$1", userId).Scan(&user.ID, &user.Name, &user.Email)

    if err != nil {
        log.Printf("Error fetching user: %v", err)  // Log error if any
        return nil, err
    }

    log.Printf("Fetched user: %+v", user)  // Log the fetched user
    return &user, nil
}

func FetchPaginatedUsers(limit, offset int) ([]User, error) {
    rows, err := dbConn.Query(context.Background(), "SELECT id, name, email FROM users LIMIT $1 OFFSET $2", limit, offset)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err = rows.Scan(&user.ID, &user.Name, &user.Email)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}

func UpdateUser(id, name, email string) error {
    _, err := dbConn.Exec(context.Background(), "UPDATE users SET name=$1, email=$2 WHERE id=$3", name, email, id)
    return err
}

func CreateUser(name, email string) error {
    _, err := dbConn.Exec(context.Background(), "INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
    return err
}

func SearchUsers(query string) ([]User, error) {
    rows, err := dbConn.Query(context.Background(), "SELECT id, name, email FROM users WHERE name ILIKE $1 OR email ILIKE $1", "%"+query+"%")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err = rows.Scan(&user.ID, &user.Name, &user.Email)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }
    return users, nil
}