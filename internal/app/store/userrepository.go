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
	ID         int
	Header     string
	Department string
	Content    string
	Date       string
}

// FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow("SELECT id, email, encrypted_password, access FROM users WHERE email = $1",
		email).Scan(&u.ID, &u.Email, &u.EncryptedPassword, &u.Access); err != nil {
		return nil, err
	}
	return u, nil
}

// FindPost ...
func (r *UserRepository) FindPost(ID string) (*Post, error) {
	post := &Post{}
	if err := r.store.db.QueryRow("SELECT header, department, content, date FROM posts WHERE id = $1",
		ID).Scan(&post.Header, &post.Department, &post.Content, &post.Date); err != nil {
		return nil, err
	}
	return post, nil
}

// FindAllPosts ...
func (r *UserRepository) FindAllPosts() ([]Post, error) {
	rows, err := r.store.db.Query("SELECT * FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	posts := make([]Post, 0)
	for rows.Next() {
		//post := new(Post)
		var post Post
		err := rows.Scan(&post.ID, &post.Header, &post.Department, &post.Content, &post.Date)
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

// AddPost ...
func (r *UserRepository) AddPost(post Post) (int64, error) {
	result, err := r.store.db.Exec("INSERT INTO posts (header, department, content, date) VALUES ($1, $2, $3, '01/02/2021');",
		post.Header, post.Department, post.Content)
	if err != nil {
		return 0, err
	}
	postID, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return postID, nil
}

// ChangePost ...
func (r *UserRepository) ChangePost(post Post) error {
	_, err := r.store.db.Exec("UPDATE posts SET header = $1, department = $2, content = $3 WHERE id = $4;", // UPDATE!!!!!!!!!!!!!!!!!!!!!!!!!!!
		post.Header, post.Department, post.Content, post.ID)
	return err
}
