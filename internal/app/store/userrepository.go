package store

import (
	"github.com/VIPowERuS/nsu_postman/internal/app/model"
)

// UserRepository ...
type UserRepository struct {
	store *Store
}

// Post ...
type Post struct {
	ID      int
	Header  string
	Author  int
	Content string
	Date    string
}

// FindByEmail takes user's data from database and return "user" struct by email
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow("SELECT id, email, encrypted_password, access FROM users WHERE email = $1",
		email).Scan(&u.ID, &u.Email, &u.EncryptedPassword, &u.Access); err != nil {
		return nil, err
	}
	return u, nil
}

// FindPost use post id and department to find necessary post in database and return it
func (r *UserRepository) FindPost(ID string, department string) (*Post, error) { // GOOD
	post := &Post{}
	if err := r.store.db.QueryRow("SELECT header, author, content, date FROM "+department+" WHERE id = $1",
		ID).Scan(&post.Header, &post.Author, &post.Content, &post.Date); err != nil {
		return nil, err
	}
	return post, nil
}

// FindAllDepartmentPosts returns all posts from department's database
func (r *UserRepository) FindAllDepartmentPosts(department string) ([]Post, error) { // GOOD
	rows, err := r.store.db.Query("SELECT * FROM " + department)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]Post, 0)
	for rows.Next() {
		var post Post
		err := rows.Scan(&post.ID, &post.Header, &post.Author, &post.Content, &post.Date)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

// AddPost to department's database and returns new post id
func (r *UserRepository) AddPost(post Post, department string) (int64, error) { // GOOD
	result, err := r.store.db.Exec("INSERT INTO "+department+" (header, author, content, date) VALUES ($1, $2, $3, '02/02/2021');",
		post.Header, post.Author, post.Content)
	if err != nil {
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return postID, nil
}

// ChangePost update already created post in database
func (r *UserRepository) ChangePost(post Post, department string) error { // GOOD
	_, err := r.store.db.Exec("UPDATE "+department+" SET header = $1, author = $2, content = $3 WHERE id = $4;", // UPDATE!!!!!!!!!!!!!!!!!!!!!!!!!!!
		post.Header, post.Author, post.Content, post.ID)
	return err
}

// DeletePost from database by id and department
func (r *UserRepository) DeletePost(ID string, department string) error { // GOOD
	if _, err := r.store.db.Exec("DELETE FROM "+department+" WHERE id = $1", ID); err != nil {
		return err
	}
	return nil
}
