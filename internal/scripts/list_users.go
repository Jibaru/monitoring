package scripts

import (
	"context"
	"monitoring/internal/domain"
	"strings"
)

type ListUsersReq struct {
	RootUserID string `json:"rootUserId"`
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	SortOrder  string `json:"sortOrder"`
	SearchTerm string `json:"searchTerm"`
}

type ListUsersResp struct {
	Users []domain.User `json:"users"`
}

type ListUsersScript struct {
	userRepo domain.UserRepo
}

func NewListUsersScript(userRepo domain.UserRepo) *ListUsersScript {
	return &ListUsersScript{userRepo: userRepo}
}

func (s *ListUsersScript) Exec(ctx context.Context, req ListUsersReq) (*ListUsersResp, error) {
	rootUserID, err := domain.NewID(req.RootUserID)
	if err != nil {
		return nil, err
	}

	filters := []domain.Filter{
		domain.NewFilter("rootUserId", domain.Equals, rootUserID),
	}

	if strings.TrimSpace(req.SearchTerm) != "" {
		filters = append(filters, domain.NewFilter("email", domain.Like, req.SearchTerm))
	}

	criteria := domain.NewCriteria(
		filters,
		domain.NewPagination(req.Limit, (req.Page-1)*req.Limit),
		domain.NewSort("email", domain.SortOrder(req.SortOrder)),
	)

	users, err := s.userRepo.ListUsers(ctx, criteria)
	if err != nil {
		return nil, err
	}

	return &ListUsersResp{Users: users}, nil
}
