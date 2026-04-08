package xiaohongshu

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

type NavigateAction struct {
	page *rod.Page
}

func NewNavigate(page *rod.Page) *NavigateAction {
	return &NavigateAction{page: page}
}

func (n *NavigateAction) ToExplorePage(ctx context.Context) error {
	page := n.page.Context(ctx)

	if err := page.Navigate("https://www.xiaohongshu.com/explore"); err != nil {
		return fmt.Errorf("failed to navigate to explore page: %w", err)
	}
	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for page load: %w", err)
	}
	if _, err := page.Element(`div#app`); err != nil {
		return fmt.Errorf("failed to find app element: %w", err)
	}

	return nil
}

func (n *NavigateAction) ToProfilePage(ctx context.Context) error {
	page := n.page.Context(ctx)

	// First navigate to explore page
	if err := n.ToExplorePage(ctx); err != nil {
		return err
	}

	if err := page.WaitStable(time.Second); err != nil {
		return fmt.Errorf("failed to wait for page stable: %w", err)
	}

	// Find and click the "我" channel link in sidebar
	profileLink, err := page.Element(`div.main-container li.user.side-bar-component a.link-wrapper span.channel`)
	if err != nil {
		return fmt.Errorf("failed to find profile link: %w", err)
	}
	if err := profileLink.Click(proto.InputMouseButtonLeft, 1); err != nil {
		return fmt.Errorf("failed to click profile link: %w", err)
	}

	// Wait for navigation to complete
	if err := page.WaitLoad(); err != nil {
		return fmt.Errorf("failed to wait for page load after click: %w", err)
	}

	return nil
}
