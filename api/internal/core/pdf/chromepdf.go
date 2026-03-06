package pdf

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// NewChromedpContext creates a chromedp context. If CHROME_PATH is set (Docker Chromium), uses it.
// Identical pattern to mims-api-service/helpers/helper.go NewChromedpContext.
func NewChromedpContext(parent context.Context) (context.Context, context.CancelFunc) {
	if path := os.Getenv("CHROME_PATH"); path != "" {
		opts := append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.ExecPath(path),
			chromedp.Flag("no-sandbox", true),
			chromedp.Flag("disable-gpu", true),
			chromedp.Flag("disable-dev-shm-usage", true),
		)
		allocCtx, allocCancel := chromedp.NewExecAllocator(parent, opts...)
		ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
		return ctx, func() { cancel(); allocCancel() }
	}
	return chromedp.NewContext(parent, chromedp.WithLogf(log.Printf))
}

// RunWithTimeout wraps tasks with a timeout context. Mirrors mims RunWithTimeOut.
func RunWithTimeout(ctx *context.Context, timeout time.Duration, tasks chromedp.Tasks) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		timeoutCtx, cancel := context.WithTimeout(ctx, timeout*time.Second)
		defer cancel()
		return tasks.Do(timeoutCtx)
	}
}

// PrintToPDF renders raw HTML string and prints to PDF bytes.
// Mirrors mims-api-service/helpers/helper.go PrintToPDF exactly:
//   - Navigate to about:blank
//   - Inject HTML via page.SetDocumentContent (waits for EventLoadEventFired)
//   - Optionally waits for div#success-pagejs (for pagedjs rendered pages)
//   - Calls page.PrintToPDF with PrintBackground + PreferCSSPageSize + DisplayHeaderFooter
func PrintToPDF(html string, res *[]byte, waitForPagedJS bool) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate("about:blank"),
		chromedp.ActionFunc(func(ctx context.Context) error {
			lctx, cancel := context.WithCancel(ctx)
			defer cancel()

			var wg sync.WaitGroup
			wg.Add(1)

			chromedp.ListenTarget(lctx, func(ev interface{}) {
				if _, ok := ev.(*page.EventLoadEventFired); ok {
					cancel()
					wg.Done()
				}
			})

			frameTree, err := page.GetFrameTree().Do(ctx)
			if err != nil {
				return err
			}
			if err := page.SetDocumentContent(frameTree.Frame.ID, html).Do(ctx); err != nil {
				return err
			}

			delay := 5
			if waitForPagedJS {
				delay = 10
			}
			// Wait for ready signal (e.g. fonts loaded then div#success-pagejs shown) or timeout
			defer chromedp.Run(
				ctx,
				RunWithTimeout(&ctx, time.Duration(delay), chromedp.Tasks{
					chromedp.WaitVisible("success-pagejs", chromedp.ByID),
				}),
			)

			wg.Wait()
			return nil
		}),

		chromedp.ActionFunc(func(ctx context.Context) error {
			buf, _, err := page.PrintToPDF().
				WithPrintBackground(true).
				WithPreferCSSPageSize(true).
				WithDisplayHeaderFooter(false).
				Do(ctx)
			if err != nil {
				return err
			}
			*res = buf
			return nil
		}),
	}
}
