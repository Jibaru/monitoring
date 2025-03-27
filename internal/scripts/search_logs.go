package scripts

import (
	"context"
	"fmt"
	"time"

	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/packages/param"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"monitoring/internal/persistence"
)

type SearchLogsReq struct {
	AppID  string    `json:"appId"`
	From   time.Time `form:"from"`
	To     time.Time `form:"to"`
	Search string    `form:"search"`
	Page   int       `form:"page"`
	Limit  int       `form:"limit"`
}

type SearchLogsResp struct {
	Logs []persistence.Log `json:"logs"`
}

type SearchLogsScript struct {
	db           *mongo.Database
	openaiClient *openai.Client
}

func NewSearchLogsScript(db *mongo.Database) *SearchLogsScript {
	client := openai.NewClient(openai.DefaultClientOptions()...)
	return &SearchLogsScript{db: db, openaiClient: &client}
}

func (s *SearchLogsScript) Exec(ctx context.Context, req SearchLogsReq) (*SearchLogsResp, error) {
	// Prompt that will be used to generate a mongo pipeline to search logs
	prompt := fmt.Sprintf("create a mongo pipeline in mongo from %s to %s containing \"%s\"", req.From, req.To, req.Search)
	openaiResp, err := s.openaiClient.Completions.New(ctx, openai.CompletionNewParams{
		Model: "gpt-4",
		Prompt: openai.CompletionNewParamsPromptUnion{
			OfString: param.NewOpt(prompt),
		},
		MaxTokens: param.NewOpt[int64](100),
	})
	if err != nil {
		return nil, err
	}

	pipelineStr := openaiResp.Choices[0].Text
	pipeline, err := ParseMatch(pipelineStr)
	if err != nil {
		return nil, err
	}

	logs, err := persistence.SearchLogs(ctx, s.db, req.AppID, req.From, req.To, req.Page, req.Limit, pipeline)
	if err != nil {
		return nil, err
	}
	return &SearchLogsResp{Logs: logs}, nil
}

func ParseMatch(matchStage string) (bson.M, error) {
	return nil, nil
}
