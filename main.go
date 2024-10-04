package main

import (
    "github.com/gorilla/mux"
    "log"
    "net/http"
    "strconv"
    "time"
    "encoding/json"
)

func main() {
    InitRedis()       // Initialize Redis
    InitPostgres()    // Initialize PostgreSQL

    r := mux.NewRouter()

    // Routes
    r.HandleFunc("/user/{id}", GetUserHandler).Methods("GET")
    r.HandleFunc("/users", GetPaginatedUsersHandler).Methods("GET")
    r.HandleFunc("/user/{id}", UpdateUserHandler).Methods("PUT")
    r.HandleFunc("/user", CreateUserHandler).Methods("POST")
    r.HandleFunc("/search-users", SearchUsersHandler).Methods("GET")

    // Apply rate limiter middleware
    r.Use(RateLimiterMiddleware)

    // Start server
    srv := &http.Server{
        Handler:      r,
        Addr:         "0.0.0.0:3000",
        WriteTimeout: 15 * time.Second,
        ReadTimeout:  15 * time.Second,
    }

    log.Println("Server is running on port 3000")
    log.Fatal(srv.ListenAndServe())
}

// GetUserHandler returns a single user by ID
func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userId := vars["id"]

    log.Printf("Requested User ID: %s", userId)  // Log the user ID

    cacheKey := "user:" + userId

    // Try to get user from Redis cache
    cachedUser, _ := GetFromCache(cacheKey)
    if cachedUser != "" {
        log.Println("Cache hit")
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(cachedUser))
        return
    }

    // Fetch from PostgreSQL
    user, err := GetUserById(userId)
    if err != nil {
        http.Error(w, "User not found", http.StatusNotFound)
        return
    }

    // Set result in cache
    userJson, _ := json.Marshal(user)
    SetToCache(cacheKey, string(userJson))

    log.Println("Cache miss, fetching from PostgreSQL")
    w.Header().Set("Content-Type", "application/json")
    w.Write(userJson)
}

// GetPaginatedUsersHandler fetches paginated list of users
func GetPaginatedUsersHandler(w http.ResponseWriter, r *http.Request) {
    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")

    page, _ := strconv.Atoi(pageStr)
    limit, _ := strconv.Atoi(limitStr)

    offset := (page - 1) * limit
    cacheKey := "users:page:" + pageStr + ":limit:" + limitStr

    // Try to get users from cache
    cachedUsers, _ := GetFromCache(cacheKey)
    if cachedUsers != "" {
        log.Println("Cache hit")
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(cachedUsers))
        return
    }

    // Fetch users from PostgreSQL
    users, err := FetchPaginatedUsers(limit, offset)
    if err != nil {
        http.Error(w, "Error fetching users", http.StatusInternalServerError)
        return
    }

    usersJson, _ := json.Marshal(users)
    SetToCache(cacheKey, string(usersJson))

    log.Println("Cache miss, fetching from PostgreSQL")
    w.Header().Set("Content-Type", "application/json")
    w.Write(usersJson)
}

// UpdateUserHandler updates a user profile and invalidates cache
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    userId := vars["id"]

    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err = UpdateUser(userId, user.Name, user.Email)
    if err != nil {
        http.Error(w, "Error updating user", http.StatusInternalServerError)
        return
    }

    InvalidateCache("users:page:*") // Invalidate paginated users cache
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"message": "User updated"}`))
}

// CreateUserHandler creates a new user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

    err = CreateUser(user.Name, user.Email)
    if err != nil {
        http.Error(w, "Error creating user", http.StatusInternalServerError)
        return
    }

    InvalidateCache("users:page:*") // Invalidate cache after creating a new user
    w.Header().Set("Content-Type", "application/json")
    w.Write([]byte(`{"message": "User created"}`))
}

// SearchUsersHandler searches users by name or email
func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")
    cacheKey := "users:search:" + query

    cachedSearch, _ := GetFromCache(cacheKey)
    if cachedSearch != "" {
        log.Println("Cache hit for search")
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(cachedSearch))
        return
    }

    users, err := SearchUsers(query)
    if err != nil {
        http.Error(w, "Error searching users", http.StatusInternalServerError)
        return
    }

    usersJson, _ := json.Marshal(users)
    SetToCache(cacheKey, string(usersJson))

    log.Println("Cache miss, fetching search results from PostgreSQL")
    w.Header().Set("Content-Type", "application/json")
    w.Write(usersJson)
}
