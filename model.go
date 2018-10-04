package main

import (
	"database/sql"
	"log"
	"time"
)

// The person Type (more like an object)
type Category struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Created_at  time.Time `json:"created_at"`
	Updated_at  time.Time `json:"updated_at"`
}

func getCategories(db *sql.DB, start, count int) ([]Category, error) {
	rows, err := db.Query("SELECT * FROM categories LIMIT $1 OFFSET $2", count, start)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	defer rows.Close()

	categories := []Category{}

	for rows.Next() {
		var p Category
		if err := rows.Scan(&p.ID, &p.Description, &p.Created_at, &p.Updated_at); err != nil {
			return nil, err
		}
		categories = append(categories, p)
	}

	return categories, nil
}

func (p *Category) getCategory(db *sql.DB) error {
	return db.QueryRow("SELECT * FROM categories WHERE id=$1", p.ID).Scan(&p.ID, &p.Description, &p.Created_at, &p.Updated_at)
}

func (p *Category) createCategory(db *sql.DB) error {
	// postgres doesn't return the last inserted ID so this is the workaround
	err := db.QueryRow(
		"INSERT INTO categories(description, created_at, updated_at) VALUES($1, $2, $3) RETURNING id",
		p.Description, p.Created_at, p.Updated_at).Scan(&p.ID)
	return err
}

func (p *Category) deleteCategory(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM categories WHERE id=$1", p.ID)
	return err
}

func (p *Category) updateCategory(db *sql.DB) error {
	_, err := db.Exec("UPDATE categories SET description=$1, updated_at=$2 WHERE id=$3", p.Description, p.Updated_at, p.ID)
	return err
}
