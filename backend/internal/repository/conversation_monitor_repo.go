package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

type conversationMonitorRepository struct{ db *sql.DB }

func NewConversationMonitorRepository(db *sql.DB) service.ConversationMonitorRepository {
	return &conversationMonitorRepository{db: db}
}

func (r *conversationMonitorRepository) Create(ctx context.Context, in service.ConversationRecordInput) error {
	_, err := r.db.ExecContext(ctx, `INSERT INTO conversation_records
        (request_id,user_id,api_key_id,group_id,provider,endpoint,requested_model,upstream_model,prompt_content,response_content,stream,status,complete,truncated,input_tokens,output_tokens,duration_ms,created_at,completed_at)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,COALESCE($18,NOW()),$19)
        ON CONFLICT (request_id) DO NOTHING`, in.RequestID, in.UserID, in.APIKeyID, in.GroupID, in.Provider, in.Endpoint, in.RequestedModel, in.UpstreamModel, in.PromptContent, in.ResponseContent, in.Stream, in.Status, in.Complete, in.Truncated, in.InputTokens, in.OutputTokens, in.DurationMS, in.CreatedAt, in.CompletedAt)
	return err
}

func conversationWhere(f service.ConversationFilter) (string, []any) {
	parts := []string{"1=1"}
	args := []any{}
	addOne := func(format string, value any) {
		args = append(args, value)
		parts = append(parts, fmt.Sprintf(format, len(args)))
	}
	addTwo := func(format string, first, second any) {
		args = append(args, first)
		firstPos := len(args)
		args = append(args, second)
		parts = append(parts, fmt.Sprintf(format, firstPos, len(args)))
	}
	if f.UserSearch != "" {
		addTwo("(u.username ILIKE '%%' || $%d || '%%' OR u.email ILIKE '%%' || $%d || '%%')", f.UserSearch, f.UserSearch)
	}
	if f.APIKey != "" {
		addTwo("(ak.name ILIKE '%%' || $%d || '%%' OR CAST(cr.api_key_id AS TEXT) = $%d)", f.APIKey, f.APIKey)
	}
	if f.Model != "" {
		addTwo("(cr.requested_model ILIKE '%%' || $%d || '%%' OR cr.upstream_model ILIKE '%%' || $%d || '%%')", f.Model, f.Model)
	}
	if f.Provider != "" {
		addOne("cr.provider = $%d", f.Provider)
	}
	if f.Status != "" {
		addOne("cr.status = $%d", f.Status)
	}
	if f.Keyword != "" {
		addTwo("(cr.prompt_content ILIKE '%%' || $%d || '%%' OR cr.response_content ILIKE '%%' || $%d || '%%')", f.Keyword, f.Keyword)
	}
	if f.Stream != nil {
		addOne("cr.stream = $%d", *f.Stream)
	}
	if f.StartTime != nil {
		addOne("cr.created_at >= $%d", *f.StartTime)
	}
	if f.EndTime != nil {
		addOne("cr.created_at < $%d", *f.EndTime)
	}
	return strings.Join(parts, " AND "), args
}

const conversationSelect = `SELECT cr.id,cr.request_id,cr.user_id,COALESCE(u.username,''),COALESCE(u.email,''),cr.api_key_id,COALESCE(ak.name,''),cr.group_id,COALESCE(g.name,''),cr.account_id,cr.provider,cr.endpoint,cr.requested_model,cr.upstream_model,cr.prompt_content,cr.response_content,cr.stream,cr.status,cr.complete,cr.truncated,cr.input_tokens,cr.output_tokens,cr.duration_ms,cr.created_at,cr.completed_at FROM conversation_records cr LEFT JOIN users u ON u.id=cr.user_id LEFT JOIN api_keys ak ON ak.id=cr.api_key_id LEFT JOIN groups g ON g.id=cr.group_id`

func scanConversation(row interface{ Scan(...any) error }) (*service.ConversationRecord, error) {
	var v service.ConversationRecord
	err := row.Scan(&v.ID, &v.RequestID, &v.UserID, &v.Username, &v.UserEmail, &v.APIKeyID, &v.APIKeyName, &v.GroupID, &v.GroupName, &v.AccountID, &v.Provider, &v.Endpoint, &v.RequestedModel, &v.UpstreamModel, &v.PromptContent, &v.ResponseContent, &v.Stream, &v.Status, &v.Complete, &v.Truncated, &v.InputTokens, &v.OutputTokens, &v.DurationMS, &v.CreatedAt, &v.CompletedAt)
	return &v, err
}

func (r *conversationMonitorRepository) List(ctx context.Context, f service.ConversationFilter) (*service.ConversationPage, error) {
	where, args := conversationWhere(f)
	var total int64
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM conversation_records cr LEFT JOIN users u ON u.id=cr.user_id LEFT JOIN api_keys ak ON ak.id=cr.api_key_id WHERE `+where, args...).Scan(&total); err != nil {
		return nil, err
	}
	page, size := f.Page, f.PageSize
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 20
	}
	if size > 200 {
		size = 200
	}
	offset := (page - 1) * size
	order := "cr.created_at DESC"
	if f.SortOrder == "asc" {
		order = "cr.created_at ASC"
	}
	args = append(args, size, offset)
	rows, err := r.db.QueryContext(ctx, conversationSelect+` WHERE `+where+` ORDER BY `+order+` LIMIT $`+fmt.Sprint(len(args)-1)+` OFFSET $`+fmt.Sprint(len(args)), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]*service.ConversationRecord, 0, size)
	for rows.Next() {
		v, err := scanConversation(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, v)
	}
	return &service.ConversationPage{Items: items, Total: total, Page: page, PageSize: size}, rows.Err()
}

func (r *conversationMonitorRepository) GetByID(ctx context.Context, id int64) (*service.ConversationRecord, error) {
	return scanConversation(r.db.QueryRowContext(ctx, conversationSelect+` WHERE cr.id=$1`, id))
}
func (r *conversationMonitorRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM conversation_records WHERE id=$1`, id)
	return err
}
func (r *conversationMonitorRepository) DeleteByFilter(ctx context.Context, f service.ConversationFilter) (int64, error) {
	where, args := conversationWhere(f)
	res, err := r.db.ExecContext(ctx, `DELETE FROM conversation_records WHERE id IN (SELECT cr.id FROM conversation_records cr LEFT JOIN users u ON u.id=cr.user_id LEFT JOIN api_keys ak ON ak.id=cr.api_key_id WHERE `+where+`)`, args...)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}
