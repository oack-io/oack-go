package oack

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Comment represents a comment on a monitor.
type Comment struct {
	ID           string  `json:"id"`
	MonitorID    string  `json:"monitor_id"`
	AuthorID     string  `json:"author_id"`
	AuthorName   string  `json:"author_name"`
	AuthorAvatar string  `json:"author_avatar"`
	Body         string  `json:"body"`
	ParentID     *string `json:"parent_id"`
	ReplyCount   int     `json:"reply_count"`
	Resolved     bool    `json:"resolved"`
	ResolvedBy   *string `json:"resolved_by"`
	ResolvedAt   *string `json:"resolved_at"`
	EditedAt     *string `json:"edited_at"`
	CreatedAt    string  `json:"created_at"`
}

// CommentEdit represents a historical edit of a comment.
type CommentEdit struct {
	ID        string `json:"id"`
	CommentID string `json:"comment_id"`
	Body      string `json:"body"`
	EditedBy  string `json:"edited_by"`
	CreatedAt string `json:"created_at"`
}

// CreateCommentParams holds parameters for creating a comment.
type CreateCommentParams struct {
	Body string `json:"body"`
}

// commentBasePath returns the base URL path for comment operations on a monitor.
func commentBasePath(teamID, monitorID string) string {
	return monitorPath(teamID, monitorID) + "/comments"
}

// commentPath returns the URL path for a specific comment.
func commentPath(teamID, monitorID, commentID string) string {
	return commentBasePath(teamID, monitorID) + "/" + commentID
}

// CreateComment creates a new comment on a monitor.
func (c *Client) CreateComment(
	ctx context.Context, teamID, monitorID string, params CreateCommentParams,
) (*Comment, error) {
	body, err := c.do(ctx, http.MethodPost, commentBasePath(teamID, monitorID), params)
	if err != nil {
		return nil, err
	}
	var comment Comment
	if err := json.Unmarshal(body, &comment); err != nil {
		return nil, fmt.Errorf("unmarshal comment: %w", err)
	}
	return &comment, nil
}

// ListComments returns all comments for a monitor.
func (c *Client) ListComments(
	ctx context.Context, teamID, monitorID string,
) ([]Comment, error) {
	body, err := c.do(ctx, http.MethodGet, commentBasePath(teamID, monitorID), nil)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, fmt.Errorf("unmarshal comments: %w", err)
	}
	return comments, nil
}

// EditComment updates the body of a comment.
func (c *Client) EditComment(
	ctx context.Context, teamID, monitorID, commentID, body string,
) (*Comment, error) {
	resp, err := c.do(
		ctx, http.MethodPut,
		commentPath(teamID, monitorID, commentID),
		map[string]string{"body": body},
	)
	if err != nil {
		return nil, err
	}
	var comment Comment
	if err := json.Unmarshal(resp, &comment); err != nil {
		return nil, fmt.Errorf("unmarshal comment: %w", err)
	}
	return &comment, nil
}

// DeleteComment deletes a comment by ID.
func (c *Client) DeleteComment(
	ctx context.Context, teamID, monitorID, commentID string,
) error {
	_, err := c.do(
		ctx, http.MethodDelete,
		commentPath(teamID, monitorID, commentID), nil,
	)
	return err
}

// ReplyToComment creates a reply to an existing comment.
func (c *Client) ReplyToComment(
	ctx context.Context, teamID, monitorID, commentID, body string,
) (*Comment, error) {
	resp, err := c.do(
		ctx, http.MethodPost,
		commentPath(teamID, monitorID, commentID)+"/replies",
		map[string]string{"body": body},
	)
	if err != nil {
		return nil, err
	}
	var comment Comment
	if err := json.Unmarshal(resp, &comment); err != nil {
		return nil, fmt.Errorf("unmarshal comment: %w", err)
	}
	return &comment, nil
}

// ListReplies returns all replies to a comment.
func (c *Client) ListReplies(
	ctx context.Context, teamID, monitorID, commentID string,
) ([]Comment, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		commentPath(teamID, monitorID, commentID)+"/replies", nil,
	)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, fmt.Errorf("unmarshal comments: %w", err)
	}
	return comments, nil
}

// ResolveComment marks a comment as resolved.
func (c *Client) ResolveComment(
	ctx context.Context, teamID, monitorID, commentID string,
) (*Comment, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		commentPath(teamID, monitorID, commentID)+"/resolve", nil,
	)
	if err != nil {
		return nil, err
	}
	var comment Comment
	if err := json.Unmarshal(body, &comment); err != nil {
		return nil, fmt.Errorf("unmarshal comment: %w", err)
	}
	return &comment, nil
}

// ReopenComment reopens a previously resolved comment.
func (c *Client) ReopenComment(
	ctx context.Context, teamID, monitorID, commentID string,
) (*Comment, error) {
	body, err := c.do(
		ctx, http.MethodPost,
		commentPath(teamID, monitorID, commentID)+"/reopen", nil,
	)
	if err != nil {
		return nil, err
	}
	var comment Comment
	if err := json.Unmarshal(body, &comment); err != nil {
		return nil, fmt.Errorf("unmarshal comment: %w", err)
	}
	return &comment, nil
}

// ListCommentEdits returns the edit history of a comment.
func (c *Client) ListCommentEdits(
	ctx context.Context, teamID, monitorID, commentID string,
) ([]CommentEdit, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		commentPath(teamID, monitorID, commentID)+"/edits", nil,
	)
	if err != nil {
		return nil, err
	}
	var edits []CommentEdit
	if err := json.Unmarshal(body, &edits); err != nil {
		return nil, fmt.Errorf("unmarshal comment edits: %w", err)
	}
	return edits, nil
}

// ListCommentsByTeam returns all comments across monitors in a team.
func (c *Client) ListCommentsByTeam(ctx context.Context, teamID string) ([]Comment, error) {
	body, err := c.do(ctx, http.MethodGet, "/api/v1/teams/"+teamID+"/comments", nil)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, fmt.Errorf("unmarshal comments: %w", err)
	}
	return comments, nil
}

// ListCommentsByAccount returns all comments across monitors in an account.
func (c *Client) ListCommentsByAccount(
	ctx context.Context, accountID string,
) ([]Comment, error) {
	body, err := c.do(
		ctx, http.MethodGet,
		"/api/v1/accounts/"+accountID+"/comments", nil,
	)
	if err != nil {
		return nil, err
	}
	var comments []Comment
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, fmt.Errorf("unmarshal comments: %w", err)
	}
	return comments, nil
}
