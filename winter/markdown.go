package winter

// example for https://blog.kowalczyk.info/article/cxn3/advanced-markdown-processing-in-go.html

import (
	"fmt"
	"io"

	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

// renderImage overrides the standard HTML renderer to make the image clickable for a gallery view.
func renderImage(w io.Writer, img *ast.Image, entering bool) error {
	if entering {
		if _, err := io.WriteString(
			w,
			fmt.Sprintf(`
				<label class="gallery-item">
			    <input type="checkbox" />
				  <img alt="%s" src="%s" title="%s" />
				`,
				img.Children[0].AsLeaf().Literal,
				img.Destination,
				img.Title,
			),
		); err != nil {
			return err
		}
	} else {
		if _, err := io.WriteString(w, `</label>`); err != nil {
			return err
		}
	}
	return nil
}

func markdownRenderHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	if img, ok := node.(*ast.Image); ok {
		if err := renderImage(w, img, entering); err != nil {
			panic(err)
		}
		// Alt text is a "child" of ast.Image,
		// but we handle it inside the tag in renderImage.
		return ast.SkipChildren, true
	}
	return ast.GoToNext, false
}

func newCustomizedRender() *html.Renderer {
	opts := html.RendererOptions{
		Flags:          html.FlagsNone,
		RenderNodeHook: markdownRenderHook,
	}
	return html.NewRenderer(opts)
}
