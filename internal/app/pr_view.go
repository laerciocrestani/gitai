package app

import (
	prpkg "github.com/laerciocrestani/gitai/internal/pr"
	"github.com/laerciocrestani/gitai/internal/ui"
)

func RunPRView() error {
	sess := ui.New("pr view", false)
	sess.Header()

	client, err := prpkg.New()
	if err != nil {
		return err
	}

	var view *prpkg.PRView
	if err := sess.Step("Opening Pull Request", func() error {
		var err error
		view, err = client.OpenInBrowser()
		return err
	}); err != nil {
		return err
	}

	sess.Detail(view.URL)
	sess.Success("PR aberto no browser")
	return nil
}
