package service

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

const ConversationMonitorSettingKey = "conversation_monitor_config"

type ConversationMonitorConfig struct {
	Enabled             bool `json:"enabled"`
	CaptureAllGroups    bool `json:"capture_all_groups"`
	MaxPromptBytes      int  `json:"max_prompt_bytes"`
	MaxResponseBytes    int  `json:"max_response_bytes"`
	ManualDeleteEnabled bool `json:"manual_delete_enabled"`
}

func DefaultConversationMonitorConfig() ConversationMonitorConfig {
	return ConversationMonitorConfig{MaxPromptBytes: 64 * 1024, MaxResponseBytes: 256 * 1024, ManualDeleteEnabled: true}
}

type ConversationRecord struct {
	ID              int64      `json:"id"`
	RequestID       string     `json:"request_id"`
	UserID          int64      `json:"user_id"`
	Username        string     `json:"username,omitempty"`
	UserEmail       string     `json:"user_email,omitempty"`
	APIKeyID        *int64     `json:"api_key_id,omitempty"`
	APIKeyName      string     `json:"api_key_name,omitempty"`
	GroupID         *int64     `json:"group_id,omitempty"`
	GroupName       string     `json:"group_name,omitempty"`
	AccountID       *int64     `json:"account_id,omitempty"`
	Provider        string     `json:"provider"`
	Endpoint        string     `json:"endpoint"`
	RequestedModel  string     `json:"requested_model"`
	UpstreamModel   string     `json:"upstream_model"`
	PromptContent   string     `json:"prompt_content"`
	ResponseContent string     `json:"response_content"`
	Stream          bool       `json:"stream"`
	Status          string     `json:"status"`
	Complete        bool       `json:"complete"`
	Truncated       bool       `json:"truncated"`
	InputTokens     int64      `json:"input_tokens"`
	OutputTokens    int64      `json:"output_tokens"`
	DurationMS      int64      `json:"duration_ms"`
	CreatedAt       time.Time  `json:"created_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

type ConversationRecordInput struct {
	RequestID                                         string
	UserID                                            int64
	APIKeyID                                          *int64
	GroupID                                           *int64
	Provider, Endpoint, RequestedModel, UpstreamModel string
	PromptContent, ResponseContent                    string
	Stream, Complete, Truncated                       bool
	Status                                            string
	InputTokens, OutputTokens, DurationMS             int64
	CreatedAt, CompletedAt                            *time.Time
}

type ConversationFilter struct {
	Page, PageSize                                       int
	UserSearch, APIKey, Model, Provider, Status, Keyword string
	Stream                                               *bool
	StartTime, EndTime                                   *time.Time
	SortBy, SortOrder                                    string
}

type ConversationPage struct {
	Items    []*ConversationRecord `json:"items"`
	Total    int64                 `json:"total"`
	Page     int                   `json:"page"`
	PageSize int                   `json:"page_size"`
}

type ConversationMonitorRepository interface {
	Create(context.Context, ConversationRecordInput) error
	List(context.Context, ConversationFilter) (*ConversationPage, error)
	GetByID(context.Context, int64) (*ConversationRecord, error)
	Delete(context.Context, int64) error
	DeleteByFilter(context.Context, ConversationFilter) (int64, error)
}

type ConversationMonitorService struct {
	repo        ConversationMonitorRepository
	settingRepo SettingRepository
}

func NewConversationMonitorService(repo ConversationMonitorRepository, settingRepo SettingRepository) *ConversationMonitorService {
	return &ConversationMonitorService{repo: repo, settingRepo: settingRepo}
}

func normalizeConversationConfig(cfg ConversationMonitorConfig) ConversationMonitorConfig {
	d := DefaultConversationMonitorConfig()
	if cfg.MaxPromptBytes <= 0 {
		cfg.MaxPromptBytes = d.MaxPromptBytes
	}
	if cfg.MaxResponseBytes <= 0 {
		cfg.MaxResponseBytes = d.MaxResponseBytes
	}
	if cfg.MaxPromptBytes > 1024*1024 {
		cfg.MaxPromptBytes = 1024 * 1024
	}
	if cfg.MaxResponseBytes > 4*1024*1024 {
		cfg.MaxResponseBytes = 4 * 1024 * 1024
	}
	return cfg
}

func (s *ConversationMonitorService) GetConfig(ctx context.Context) (ConversationMonitorConfig, error) {
	cfg := DefaultConversationMonitorConfig()
	if s == nil || s.settingRepo == nil {
		return cfg, nil
	}
	raw, err := s.settingRepo.GetValue(ctx, ConversationMonitorSettingKey)
	if errors.Is(err, ErrSettingNotFound) {
		return cfg, nil
	}
	if err != nil {
		return cfg, err
	}
	if strings.TrimSpace(raw) == "" {
		return cfg, nil
	}
	if err := json.Unmarshal([]byte(raw), &cfg); err != nil {
		return cfg, err
	}
	return normalizeConversationConfig(cfg), nil
}

func (s *ConversationMonitorService) UpdateConfig(ctx context.Context, cfg ConversationMonitorConfig) (ConversationMonitorConfig, error) {
	cfg = normalizeConversationConfig(cfg)
	b, err := json.Marshal(cfg)
	if err != nil {
		return cfg, err
	}
	if err := s.settingRepo.Set(ctx, ConversationMonitorSettingKey, string(b)); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func (s *ConversationMonitorService) Record(ctx context.Context, input ConversationRecordInput) error {
	if s == nil || s.repo == nil {
		return nil
	}
	cfg, err := s.GetConfig(ctx)
	if err != nil || !cfg.Enabled {
		return err
	}
	return s.repo.Create(ctx, input)
}

func (s *ConversationMonitorService) List(ctx context.Context, f ConversationFilter) (*ConversationPage, error) {
	return s.repo.List(ctx, f)
}
func (s *ConversationMonitorService) GetByID(ctx context.Context, id int64) (*ConversationRecord, error) {
	return s.repo.GetByID(ctx, id)
}
func (s *ConversationMonitorService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}
func (s *ConversationMonitorService) DeleteByFilter(ctx context.Context, f ConversationFilter) (int64, error) {
	return s.repo.DeleteByFilter(ctx, f)
}
