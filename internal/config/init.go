package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/laerciocrestani/gitia/internal/ui"
)

func InitInteractive() error {
	sess := ui.New("config", false)
	sess.Header()

	cfg := Default()
	reader := bufio.NewReader(os.Stdin)

	if err := sess.Step("Starting configuration wizard", func() error {
		return nil
	}); err != nil {
		return err
	}

	provider, err := promptChoice(sess, reader, "Provider", []string{"openrouter", "openai", "gemini"}, string(cfg.Provider))
	if err != nil {
		return err
	}
	cfg.Provider = Provider(provider)

	sess.Prompt("API Key: ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	cfg.APIKey = strings.TrimSpace(apiKey)

	defaultModel := defaultModelFor(cfg.Provider)
	sess.Prompt(fmt.Sprintf("Model [%s]: ", defaultModel))
	model, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	model = strings.TrimSpace(model)
	if model == "" {
		model = defaultModel
	}
	cfg.Model = model

	sess.Prompt(fmt.Sprintf("Idioma das mensagens [%s]: ", cfg.Language))
	lang, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	lang = strings.TrimSpace(lang)
	if lang != "" {
		cfg.Language = lang
	}

	sess.Prompt(fmt.Sprintf("Branch base [%s]: ", cfg.BaseBranch))
	base, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	base = strings.TrimSpace(base)
	if base != "" {
		cfg.BaseBranch = base
	}

	sess.Prompt("Co-author trailer (opcional): ")
	coAuthor, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	cfg.CoAuthor = strings.TrimSpace(coAuthor)

	path, err := ConfigPath()
	if err != nil {
		return err
	}

	if err := sess.Step("Saving configuration", func() error {
		return Save(path, cfg)
	}); err != nil {
		return err
	}

	sess.Detail(path)
	sess.Success("Configuration saved ✨")
	return nil
}

func defaultModelFor(p Provider) string {
	switch p {
	case ProviderOpenAI:
		return "gpt-4o-mini"
	case ProviderGemini:
		return "gemini-2.5-flash-lite"
	default:
		return "deepseek/deepseek-chat"
	}
}

func promptChoice(sess *ui.Session, reader *bufio.Reader, label string, options []string, defaultVal string) (string, error) {
	sess.Prompt(fmt.Sprintf("%s (%s) [%s]: ", label, strings.Join(options, ", "), defaultVal))
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return defaultVal, nil
	}
	for _, opt := range options {
		if input == opt {
			return opt, nil
		}
	}
	return "", fmt.Errorf("opção inválida: %q", input)
}
