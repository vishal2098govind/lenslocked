package main

import (
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib"
	db "github.com/vishal2098govind/lenslocked/models/db"
	userM "github.com/vishal2098govind/lenslocked/models/user"
)

func main() {
	cfg := db.DefaultPostgresConfig()

	db, err := db.Open(cfg)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	us := userM.UserService{
		DB: db,
	}

	res, err := us.Authenticate(userM.AuthenticateRequest{
		Email:    "vishal.govind@gmil.com",
		Password: "vishal2098",
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	// row := db.QueryRow("")
	// row.Scan()

	// err = db.Ping()
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("Connected!")

	// 	_, err = db.Exec(`
	// 	CREATE TABLE IF NOT EXISTS users (
	// 		id SERIAL PRIMARY KEY,
	// 		name TEXT,
	// 		email TEXT UNIQUE NOT NULL
	// 	);

	// 	CREATE TABLE IF NOT EXISTS orders (
	// 		id SERIAL PRIMARY KEY,
	// 		user_id INT NOT NULL,
	// 		amount INT,
	// 		description TEXT
	// 	);
	// `)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("Tables created")

	// name := "vishal govind"
	// email := "vishal@govind3.io"

	// _, err = db.Exec(fmt.Sprintf(`
	// 	INSERT INTO users (name, email)
	// 	VALUES (%s, %s);
	// `, name, email))

	// _, err = db.Exec(`
	// 	INSERT INTO users (name, email)
	// 	VALUES ($1, $2);
	// `, name, email)

	// if err != nil {
	// 	panic(err)
	// }

	// row = db.QueryRow(`
	// 	INSERT INTO users (name, email)
	// 	VALUES ($1, $2)
	// 	RETURNING id, name;
	// `, name, email)

	// var id int
	// var rname string
	// err = row.Scan(&id, &rname)
	// if err != nil {
	// 	panic(err)
	// }

	// _, err = db.Exec(`
	// 	SELECT * FROM users;
	// `)

	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Printf("User created. id = %v, name = %v\n", id, rname)

	// userID := 6
	// for i := 1; i <= 5; i++ {
	// 	amount := i * 100
	// 	desc := fmt.Sprintf("Fake order #%d", i)
	// 	_, err := db.Exec(`
	// INSERT INTO orders(user_id, amount, description)
	// VALUES($1, $2, $3)`, userID, amount, desc)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
	// fmt.Println("Created fake orders.")
}
