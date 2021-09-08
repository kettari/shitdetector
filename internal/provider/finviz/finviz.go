package finviz

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"github.com/kettari/shitdetector/internal/asset"
	"strings"
)

const finvizUrl = `https://finviz.com/quote.ashx?t=%s`

type finvizProvider struct {
}

func NewFinvizProvider() *finvizProvider {
	return &finvizProvider{}
}

func (p finvizProvider) Fetch(ticker string) (stock *asset.Stock, err error) {
	html, err := downloadHtml(fmt.Sprintf(finvizUrl, ticker))
	if err != nil {
		return nil, err
	}

	return parseHtml(html)
}

func downloadHtml(url string) (*string, error) {
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// run task list
	var res string
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) error {
			node, err := dom.GetDocument().Do(ctx)
			if err != nil {
				return err
			}
			res, err = dom.GetOuterHTML().WithNodeID(node.NodeID).Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, err
	}

	html := strings.TrimSpace(res)
	return &html, nil
	/*resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	html := string(content)
	return &html, nil*/
}

func parseHtml(html *string) (*asset.Stock, error) {
	//regexp.MustCompile()
	return nil, nil
}
