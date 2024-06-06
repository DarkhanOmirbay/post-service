package post

import (
	"context"
	"fmt"
	pv1 "github.com/DarkhanOmirbay/proto/proto/gen/go/post"
	"log/slog"
	"strconv"
	"time"
)

type PostProvider interface {
	CreatePost(ctx context.Context, userId int64, title, content string) (*pv1.Post, error)
	UpdatePost(ctx context.Context, userId, post_id int64, title, content string) (*pv1.Post, error)
	DeletePost(ctx context.Context, user_id, post_id int64) (string, error)
	CreateComment(ctx context.Context, user_id, post_id int64, content string) (*pv1.Comment, error)
	CreateLike(ctx context.Context, user_id, post_id int64) (string, error)
}
type Post struct {
	log          *slog.Logger
	postProvider PostProvider
	tokenTTL     time.Duration
}

func New(log *slog.Logger, tokenTTL time.Duration, postProvider PostProvider) *Post {
	return &Post{log: log, postProvider: postProvider, tokenTTL: tokenTTL}
}
func (p *Post) WritePost(ctx context.Context, userId int64, title, content string) (*pv1.Post, error) {
	const op = "post.post.WritePost"
	log := p.log.With(
		slog.String("op", op),
		slog.String("user id ", strconv.FormatInt(userId, 10)),
	)
	log.Info("adding post to db")
	post, err := p.postProvider.CreatePost(ctx, userId, title, content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return post, nil
}
func (p *Post) EditPost(ctx context.Context, userId, post_id int64, title, content string) (*pv1.Post, error) {
	const op = "post.post.EditPost"
	log := p.log.With(
		slog.String("op", op),
		slog.String("user id ", strconv.FormatInt(userId, 10)),
		slog.String("post id", strconv.FormatInt(post_id, 10)),
	)
	log.Info("editing post")
	post, err := p.postProvider.UpdatePost(ctx, userId, post_id, title, content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return post, nil
}
func (p *Post) DeletePost(ctx context.Context, user_id, post_id int64) (string, error) {
	const op = "post.post.DeletePost"
	log := p.log.With(
		slog.String("op", op),
		slog.String("user id ", strconv.FormatInt(user_id, 10)),
		slog.String("post id", strconv.FormatInt(post_id, 10)),
	)
	log.Info("deleting post")
	msg, err := p.postProvider.DeletePost(ctx, user_id, post_id)
	if err != nil {
		return "error", fmt.Errorf("%s: %w", op, err)
	}
	return msg, nil
}
func (p *Post) CommentPost(ctx context.Context, user_id, post_id int64, content string) (*pv1.Comment, error) {
	const op = "post.post.CommentPost"
	log := p.log.With(
		slog.String("op", op),
		slog.String("user id ", strconv.FormatInt(user_id, 10)),
		slog.String("post id", strconv.FormatInt(post_id, 10)),
	)
	log.Info("commenting post")
	comment, err := p.postProvider.CreateComment(ctx, user_id, post_id, content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return comment, nil
}
func (p *Post) LikePost(ctx context.Context, user_id, post_id int64) (string, error) {
	const op = "post.post.LikePost"
	log := p.log.With(
		slog.String("op", op),
		slog.String("user id ", strconv.FormatInt(user_id, 10)),
		slog.String("post id", strconv.FormatInt(post_id, 10)),
	)
	log.Info("like post")
	msg, err := p.postProvider.CreateLike(ctx, user_id, post_id)
	if err != nil {
		return "err", fmt.Errorf("%s: %w", op, err)
	}
	return msg, nil
}
