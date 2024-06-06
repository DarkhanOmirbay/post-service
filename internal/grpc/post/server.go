package post

import (
	"context"
	"fmt"
	pv1 "github.com/DarkhanOmirbay/proto/proto/gen/go/post"
	ssov1 "github.com/DarkhanOmirbay/proto/proto/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Post interface {
	WritePost(ctx context.Context, user_id int64, title, content string) (*pv1.Post, error)
	EditPost(ctx context.Context, user_id, post_id int64, title, content string) (*pv1.Post, error)
	DeletePost(ctx context.Context, user_id, post_id int64) (string, error)
	CommentPost(ctx context.Context, user_id, post_id int64, content string) (*pv1.Comment, error)
	LikePost(ctx context.Context, user_id, post_id int64) (string, error)
}
type serverAPI struct {
	pv1.UnimplementedPostServiceServer
	post       Post
	authClient ssov1.AuthClient
}

func Register(gRPCServer *grpc.Server, post Post, client ssov1.AuthClient) {
	pv1.RegisterPostServiceServer(gRPCServer, &serverAPI{post: post, authClient: client})
}
func (s *serverAPI) authenticate(token string) (int64, error) {
	const op = "post.server.authenticate"

	req := &ssov1.IsAuthenticatedRequest{Token: token}
	resp, err := s.authClient.IsAuthenticated(context.Background(), req)
	if err != nil || !resp.IsAuthenticated {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return resp.UserId, nil
}
func (s *serverAPI) WritePost(ctx context.Context, in *pv1.WritePostRequest) (*pv1.PostResponse, error) {
	const op = "post.server.WritePost"
	if in.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if in.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}
	userId, err := s.authenticate(in.Token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	post, err := s.post.WritePost(ctx, userId, in.Title, in.Content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &pv1.PostResponse{Post: post}, nil
}
func (s *serverAPI) EditPost(ctx context.Context, in *pv1.EditPostRequest) (*pv1.PostResponse, error) {
	const op = "post.server.EditPost"
	if in.Title == "" {
		return nil, status.Error(codes.InvalidArgument, "title is required")
	}
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if in.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}
	if in.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	userId, err := s.authenticate(in.Token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	post, err := s.post.EditPost(ctx, userId, in.Id, in.Title, in.Content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &pv1.PostResponse{Post: post}, nil
}
func (s *serverAPI) DeletePost(ctx context.Context, in *pv1.DeletePostRequest) (*pv1.DeletePostResponse, error) {
	const op = "post.server.DeletePost"
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if in.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	userId, err := s.authenticate(in.Token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	msg, err := s.post.DeletePost(ctx, userId, in.Id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &pv1.DeletePostResponse{Msg: msg}, nil
}
func (s *serverAPI) CommentPost(ctx context.Context, in *pv1.CommentPostRequest) (*pv1.CommentResponse, error) {
	const op = "post.server.CommentPost"
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if in.PostId == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if in.Content == "" {
		return nil, status.Error(codes.InvalidArgument, "content is required")
	}
	userId, err := s.authenticate(in.Token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	comment, err := s.post.CommentPost(ctx, userId, in.PostId, in.Content)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &pv1.CommentResponse{Comment: comment}, nil

}
func (s *serverAPI) LikePost(ctx context.Context, in *pv1.LikePostRequest) (*pv1.LikePostResponse, error) {
	const op = "post.server.LikePost"
	if in.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}
	if in.PostId == 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	userId, err := s.authenticate(in.Token)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	msg, err := s.post.LikePost(ctx, userId, int64(in.PostId))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)

	}
	return &pv1.LikePostResponse{Msg: msg}, nil
}

//show post or get all post by sorting pagination etc
