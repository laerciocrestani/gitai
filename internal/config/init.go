package config

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/laerciocrestani/gitai/internal/ui"
)

func InitInteractive() error {
	sess := ui.New("config", false)
	sess.Header()

	existing, savePath, err := LoadExisting()
	if err != nil {
		return err
	}

	cfg := *existing
	reader := ui.StdinReader()
	hadConfig := hasSavedConfig(existing)

	sess.SectionFirst("Configuração")
	if hadConfig {
		sess.Info("Configuração atual detectada — Enter mantém cada valor entre colchetes")
	} else {
		sess.Info("Configure o provedor e o modelo de IA; a chave API vem em seguida")
	}

	prevProvider := cfg.Provider

	provider, err := sess.Select(reader, ui.SelectConfig{
		Label:   "Provedor",
		Options: []string{"openrouter", "openai", "gemini"},
		Default: string(cfg.Provider),
	})
	if err != nil {
		return err
	}
	cfg.Provider = Provider(provider)

	modelDefault := cfg.Model
	if cfg.Provider != prevProvider || modelDefault == "" {
		modelDefault = defaultModelFor(cfg.Provider)
	}
	modelKeep := cfg.Model
	if cfg.Provider != prevProvider {
		modelKeep = modelDefault
	}

	modelOptions := modelSuggestions(cfg.Provider)
	if modelKeep != "" && !slices.Contains(modelOptions, modelKeep) {
		modelOptions = append([]string{modelKeep}, modelOptions...)
	}

	model, err := sess.Select(reader, ui.SelectConfig{
		Label:      "Modelo",
		Options:    modelOptions,
		Default:    modelKeep,
		AllowOther: true,
	})
	if err != nil {
		return err
	}
	if model == "" {
		model = modelDefault
	}
	cfg.Model = model

	apiKey, err := promptAPIKey(sess, reader, cfg.Provider, cfg.APIKey)
	if err != nil {
		return err
	}
	cfg.APIKey = apiKey

	sess.Section("Preferências")

	lang, err := promptKeep(sess, reader, "Idioma das mensagens", cfg.Language, cfg.Language)
	if err != nil {
		return err
	}
	cfg.Language = lang

	base, err := promptKeep(sess, reader, "Branch base", cfg.BaseBranch, cfg.BaseBranch)
	if err != nil {
		return err
	}
	cfg.BaseBranch = base

	coAuthorDefault := cfg.CoAuthor
	if coAuthorDefault == "" {
		coAuthorDefault = "(vazio)"
	}
	coAuthor, err := promptKeep(sess, reader, "Co-author trailer (opcional)", coAuthorDefault, cfg.CoAuthor)
	if err != nil {
		return err
	}
	cfg.CoAuthor = coAuthor

	fmt.Fprintln(os.Stderr)
	sess.Info("Limpar o terminal antes de cada comando deixa só a saída do GitAi visível,")
	sess.Info("sem misturar com histórico anterior no console.")
	clear, err := promptYesNo(sess, reader, "Ativar limpeza do terminal?", cfg.ClearScreen)
	if err != nil {
		return err
	}
	cfg.ClearScreen = clear

	if strings.TrimSpace(cfg.APIKey) == "" && strings.TrimSpace(os.Getenv(EnvAPIKey)) == "" {
		return fmt.Errorf("chave API obrigatória — defina no wizard ou na variável %s", EnvAPIKey)
	}

	if err := sess.Step("Saving configuration", func() error {
		return Save(savePath, cfg)
	}); err != nil {
		return err
	}

	sess.Detail(savePath)
	sess.Success("Configuration saved ✨")
	return nil
}

func hasSavedConfig(cfg *Config) bool {
	if strings.TrimSpace(cfg.APIKey) != "" {
		return true
	}
	path, err := ConfigPath()
	if err != nil {
		return false
	}
	if _, err := os.Stat(path); err == nil {
		return true
	}
	localPath := LocalConfigPath()
	if localPath == "" {
		return false
	}
	_, err = os.Stat(localPath)
	return err == nil
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

func modelSuggestions(p Provider) []string {
	switch p {
	case ProviderOpenAI:
		return []string{"gpt-4o-mini", "gpt-4o", "gpt-4.1-mini"}
	case ProviderGemini:
		return []string{"gemini-2.5-flash-lite", "gemini-2.5-flash", "gemini-2.0-flash"}
	default:
		return []string{"deepseek/deepseek-chat", "anthropic/claude-sonnet-4", "google/gemini-2.5-flash-lite"}
	}
}

func apiKeyHint(p Provider) string {
	switch p {
	case ProviderOpenAI:
		return "https://platform.openai.com/api-keys"
	case ProviderGemini:
		return "https://aistudio.google.com/apikey"
	default:
		return "https://openrouter.ai/keys"
	}
}

func promptKeep(sess *ui.Session, reader *bufio.Reader, label, displayDefault, current string) (string, error) {
	sess.Prompt(fmt.Sprintf("%s [%s]: ", label, displayDefault))
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return current, nil
	}
	return input, nil
}

func promptAPIKey(sess *ui.Session, reader *bufio.Reader, provider Provider, current string) (string, error) {
	sess.Info("Chave em " + apiKeyHint(provider))
	current = strings.TrimSpace(current)
	if current == "" {
		sess.Prompt("Chave API: ")
	} else {
		sess.Prompt(fmt.Sprintf("Chave API [%s, Enter mantém]: ", MaskAPIKey(current)))
	}
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return current, nil
	}
	return input, nil
}

func promptYesNo(sess *ui.Session, reader *bufio.Reader, label string, current bool) (bool, error) {
	defaultLabel := "n"
	if current {
		defaultLabel = "s"
	}
	sess.Prompt(fmt.Sprintf("%s (s/n) [%s]: ", label, defaultLabel))
	input, err := reader.ReadString('\n')
	if err != nil {
		return false, err
	}
	input = strings.TrimSpace(strings.ToLower(input))
	if input == "" {
		return current, nil
	}
	switch input {
	case "s", "sim", "y", "yes":
		return true, nil
	case "n", "nao", "não", "no":
		return false, nil
	default:
		return false, fmt.Errorf("resposta inválida: %q (use s ou n)", input)
	}
}
