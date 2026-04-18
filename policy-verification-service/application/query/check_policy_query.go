package query

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/janicaleksander/cloud/common/event"
	"github.com/janicaleksander/cloud/policyverificationservice/application/interfaces"
	"github.com/janicaleksander/cloud/policyverificationservice/domain"
	"github.com/mehdihadeli/go-mediatr"
)

type CheckPolicyQuery struct {
	ClaimID      string
	UserID       string
	VIN          string
	AccidentDate time.Time
	URLs         []string
}

type CheckPolicyQueryResponse struct {
	Result bool `json:"result"`
}
type CheckPolicyQueryHandler struct {
	repo      domain.PolicyRepository
	publisher interfaces.PolicyEventPublisher
}

func NewCheckPolicyQueryHandler(r domain.PolicyRepository, p interfaces.PolicyEventPublisher) *CheckPolicyQueryHandler {
	return &CheckPolicyQueryHandler{repo: r, publisher: p}
}
func (h *CheckPolicyQueryHandler) SelfRegister() error {
	return mediatr.RegisterRequestHandler[*CheckPolicyQuery, *CheckPolicyQueryResponse](h)
}

func (h *CheckPolicyQueryHandler) Handle(ctx context.Context, query *CheckPolicyQuery) (*CheckPolicyQueryResponse, error) {
	uid, err := uuid.Parse(query.UserID)
	if err != nil {
		slog.Error("Invalid user ID format", "userID", query.UserID, "error", err)
		return &CheckPolicyQueryResponse{Result: false}, nil
	}
	hasPolicy, policy := h.repo.IfUserHasPolicy(context.Background(), uid, query.VIN)

	if !hasPolicy {
		err := h.publisher.Publish("events", event.PolicyDeniedEvent{
			ClaimID: query.ClaimID,
			Reason:  string(domain.PolicyNotFound),
		})
		if err != nil {
			slog.Error("Failed to publish PolicyDeniedEvent for claimID", "claimID", query.ClaimID, "error", err)
		}
		return &CheckPolicyQueryResponse{Result: false}, nil
	}

	valid, reason := policy.IsValid(query.AccidentDate)

	if valid {
		err := h.publisher.Publish("events", event.PolicyVerifiedEvent{
			ClaimID:    query.ClaimID,
			StorageURL: query.URLs,
		})
		if err != nil {
			slog.Error("Failed to publish PolicyVerifiedEvent for claimID", "claimID", query.ClaimID, "error", err)
		}
	} else {
		err := h.publisher.Publish("events", event.PolicyDeniedEvent{
			ClaimID: query.ClaimID,
			Reason:  string(reason),
		})
		if err != nil {
			slog.Error("Failed to publish PolicyDeniedEvent for claimID", "claimID", query.ClaimID, "error", err)
		}
	}
	return &CheckPolicyQueryResponse{Result: valid}, nil
}
