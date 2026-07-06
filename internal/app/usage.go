package app

import (
	"time"

	"github.com/laerciocrestani/gitia/internal/ai"
	"github.com/laerciocrestani/gitia/internal/config"
	gitpkg "github.com/laerciocrestani/gitia/internal/git"
	"github.com/laerciocrestani/gitia/internal/usage"
)

func recordAIUsage(command string, cfg *config.Config, summary ai.UsageSummary) {
	if len(summary.Records) == 0 {
		return
	}

	project := "unknown"
	if repo, err := gitpkg.New(); err == nil {
		project = repo.ProjectName()
	}

	for _, r := range summary.Records {
		_ = usage.Log(usage.Entry{
			Timestamp:    time.Now().UTC(),
			Command:      command,
			Project:      project,
			Provider:     string(cfg.Provider),
			Model:        cfg.Model,
			Label:        r.Label,
			InputTokens:  r.PromptTokens,
			OutputTokens: r.CompletionTokens,
			CostUSD:      r.CostUSD,
		})
	}
}
