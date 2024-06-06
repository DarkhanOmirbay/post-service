package storage

import (
	"context"
	"database/sql"
	"fmt"
	pv1 "github.com/DarkhanOmirbay/proto/proto/gen/go/post"
	_ "github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
	"time"
)

type PostStorage struct {
	db *sql.DB
}

func (ps *PostStorage) CreatePost(ctx context.Context, userId int64, title, content string) (*pv1.Post, error) {
	const op = "storage.post.CreatePost"
	fail := func(e error) error {
		return fmt.Errorf("%s: %w", op, e)
	}

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fail(err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	stmt, err := tx.Prepare("INSERT INTO posts(user_id, title, content) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at")
	if err != nil {
		return nil, fail(err)
	}
	defer stmt.Close()

	var postId int64
	var createdAt, updatedAt time.Time
	if err := stmt.QueryRow(userId, title, content).Scan(&postId, &createdAt, &updatedAt); err != nil {
		return nil, fail(err)
	}

	post := &pv1.Post{
		Id:        postId,
		UserId:    userId,
		Title:     title,
		Content:   content,
		CreatedAt: createdAt.String(),
		UpdatedAt: updatedAt.String(),
	}

	if err := tx.Commit(); err != nil {
		return nil, fail(err)
	}

	return post, nil
}

func (ps *PostStorage) UpdatePost(ctx context.Context, userId, postID int64, title, content string) (*pv1.Post, error) {
	const op = "storage.post.UpdatePost"
	fail := func(e error) error {
		return fmt.Errorf("%s: %w", op, e)
	}

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fail(err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	// Получение текущей информации о посте
	currentPost, err := ps.GetPostByID(ctx, postID)
	if err != nil {
		return nil, fail(err)
	}

	// Проверка прав доступа для редактирования поста
	if currentPost.UserId != userId {
		return nil, status.Error(codes.PermissionDenied, "you don't have permission to edit this post")
	}

	// Подготовка запроса на обновление поста
	stmt, err := tx.Prepare("UPDATE posts SET title = $1, content = $2, updated_at = $3 WHERE id = $4 RETURNING id, created_at")
	if err != nil {
		return nil, fail(err)
	}
	defer stmt.Close()

	// Выполнение запроса на обновление
	var updatedPostID int64
	var createdAt time.Time
	if err := stmt.QueryRow(title, content, time.Now(), postID).Scan(&updatedPostID, &createdAt); err != nil {
		return nil, fail(err)
	}

	// Получение обновленной информации о посте
	updatedPost, err := ps.GetPostByID(ctx, updatedPostID)
	if err != nil {
		return nil, fail(err)
	}

	// Фиксация изменений
	if err := tx.Commit(); err != nil {
		return nil, fail(err)
	}

	return updatedPost, nil
}
func (ps *PostStorage) DeletePost(ctx context.Context, userID, postID int64) (string, error) {
	const op = "storage.post.DeletePost"
	fail := func(e error) error {
		return fmt.Errorf("%s: %w", op, e)
	}

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fail(err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	// Проверка прав доступа для удаления поста
	currentPost, err := ps.GetPostByID(ctx, postID)
	if err != nil {
		return "", fail(err)
	}

	if currentPost.UserId != userID {
		return "", status.Error(codes.PermissionDenied, "you don't have permission to delete this post")
	}

	// Подготовка и выполнение запроса на удаление поста
	stmt, err := tx.Prepare("DELETE FROM posts WHERE id = $1")
	if err != nil {
		return "", fail(err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(postID); err != nil {
		return "", fail(err)
	}

	// Фиксация изменений
	if err := tx.Commit(); err != nil {
		return "", fail(err)
	}

	return fmt.Sprintf("Post with ID %d has been deleted", postID), nil
}
func (ps *PostStorage) CreateComment(ctx context.Context, userID, postID int64, content string) (*pv1.Comment, error) {
	const op = "storage.post.CreateComment"
	fail := func(e error) error {
		return fmt.Errorf("%s: %w", op, e)
	}

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fail(err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	// Проверка существования поста
	_, err = ps.GetPostByID(ctx, postID)
	if err != nil {
		return nil, fail(err)
	}

	// Вставка нового комментария в базу данных
	var commentID int64
	err = tx.QueryRowContext(ctx, "INSERT INTO comments(user_id, post_id, content) VALUES ($1, $2, $3) RETURNING id", userID, postID, content).Scan(&commentID)
	if err != nil {
		return nil, fail(err)
	}

	// Фиксация изменений
	if err := tx.Commit(); err != nil {
		return nil, fail(err)
	}

	// Создание и возвращение объекта комментария
	comment := &pv1.Comment{
		Id:        commentID,
		UserId:    userID,
		PostId:    postID,
		Content:   content,
		CreatedAt: strconv.FormatInt(time.Now().Unix(), 10),
	}

	return comment, nil
}
func (ps *PostStorage) CreateLike(ctx context.Context, userID, postID int64) (string, error) {
	const op = "storage.post.CreateLike"
	fail := func(e error) error {
		return fmt.Errorf("%s: %w", op, e)
	}

	tx, err := ps.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fail(err)
	}
	defer func() {
		if err := tx.Rollback(); err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	// Проверка существования поста
	_, err = ps.GetPostByID(ctx, postID)
	if err != nil {
		return "", fail(err)
	}

	// Проверка, существует ли уже лайк от данного пользователя для этого поста
	var existingLikeID int64
	err = tx.QueryRowContext(ctx, "SELECT id FROM likes WHERE user_id = $1 AND post_id = $2", userID, postID).Scan(&existingLikeID)
	if err == nil {
		return "", status.Error(codes.AlreadyExists, "like already exists")
	}
	if err != sql.ErrNoRows {
		return "", fail(err)
	}

	// Вставка нового лайка в базу данных
	_, err = tx.ExecContext(ctx, "INSERT INTO likes(user_id, post_id) VALUES ($1, $2)", userID, postID)
	if err != nil {
		return "", fail(err)
	}

	// Фиксация изменений
	if err := tx.Commit(); err != nil {
		return "", fail(err)
	}

	return fmt.Sprintf("User %d liked post %d", userID, postID), nil
}

func (ps *PostStorage) GetPostByID(ctx context.Context, postID int64) (*pv1.Post, error) {
	const op = "storage.post.GetPostByID"
	fail := func(e error) (*pv1.Post, error) {
		return nil, fmt.Errorf("%s: %w", op, e)
	}

	var post pv1.Post
	err := ps.db.QueryRowContext(ctx, "SELECT id, user_id, title, content, created_at, updated_at FROM posts WHERE id = $1", postID).Scan(
		&post.Id, &post.UserId, &post.Title, &post.Content, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		return fail(err)
	}

	return &post, nil
}
