package pdf

import (
	"context"
	"log"
	"os"
	"sync"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/cdproto/runtime"
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
//   - Waits for document.fonts.ready so external fonts (e.g. Google Fonts) are loaded
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

			wg.Wait()
			return nil
		}),

		// Wait for all fonts (including Google Fonts) to finish loading before printing.
		chromedp.ActionFunc(func(ctx context.Context) error {
			fontCtx, fontCancel := context.WithTimeout(ctx, 15*time.Second)
			defer fontCancel()
			// Evaluate document.fonts.ready as a Promise via Runtime.awaitPromise pattern
			var ready bool
			err := chromedp.Run(fontCtx, chromedp.Evaluate(
				`document.fonts.ready.then(() => true)`,
				&ready,
				func(p *runtime.EvaluateParams) *runtime.EvaluateParams {
					return p.WithAwaitPromise(true)
				},
			))
			if err != nil {
				// Non-fatal: if fonts.ready times out, still proceed with print
				_ = err
			}
			// Extra buffer for rendering (layout, images)
			time.Sleep(500 * time.Millisecond)
			return nil
		}),

		// Only wait for Paged.js ready signal when requested (e.g. timeline report).
		// Quotation HTML has no #success-pagejs, so skip to avoid timeout/failure.
		chromedp.ActionFunc(func(ctx context.Context) error {
			if waitForPagedJS {
				delay := 10 * time.Second
				timeoutCtx, timeoutCancel := context.WithTimeout(ctx, delay)
				defer timeoutCancel()
				if err := chromedp.Run(timeoutCtx, chromedp.WaitVisible("success-pagejs", chromedp.ByID)); err != nil {
					return err
				}
			}
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
