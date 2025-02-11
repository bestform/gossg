package page

import (
	"os"
	"strings"

	"github.com/bestform/gossg/config"
	"github.com/bestform/gossg/markdown"
)

type RenderedPage struct {
	Title    string
	Filename string
	Content  []byte
	SubPages []RenderedPage
}

func CreatePage(pageConf config.Page, headerPath, footerPath string) RenderedPage {
	header, err := os.ReadFile(headerPath)
	if err != nil {
		panic(err)
	}
	footer, err := os.ReadFile(footerPath)
	if err != nil {
		panic(err)
	}

	switch pageConf.ContentType {
	case config.CONTENT_TYPE_HTML:
		return createHtmlPage(pageConf, header, footer)
	case config.CONTENT_TYPE_BLOG_INDEX:
		indexPage := createBlogIndexPage(pageConf, header, footer)
		var blogPages []RenderedPage
		for _, blogPageConf := range pageConf.Pages {
			blogPage := CreatePage(blogPageConf, headerPath, footerPath)
			blogPages = append(blogPages, blogPage)
		}
		indexPage.SubPages = blogPages
		return indexPage
	case config.CONTENT_TYPE_MARKDOWN:
		return createMarkdownPage(pageConf, header, footer)
	default:
		panic("Unknown content type or not implemented yet")
	}
}

func createMarkdownPage(pageConf config.Page, header, footer []byte) RenderedPage {
	mdContent, err := os.ReadFile(pageConf.ContentPath)
	if err != nil {
		panic(err)
	}
	lexer := markdown.NewLexer(string(mdContent))
	tokens := lexer.Tokenize()
	parser := markdown.NewParser(tokens)
	nodes := parser.Parse()

	renderedMarkdown := nodes.Render()
	renderedMarkdown = "<h1>" + pageConf.Title + "</h1>" + renderedMarkdown
	markdownTemplate, err := os.ReadFile("templates/markdown_template.html")
	if err != nil {
		panic(err)
	}
	renderedMarkdownInTemplate := strings.ReplaceAll(string(markdownTemplate), "%%CONTENT%%", renderedMarkdown)
	var page []byte
	page = append(page, header...)
	page = append(page, []byte(renderedMarkdownInTemplate)...)
	page = append(page, footer...)

	return RenderedPage{Title: pageConf.Title, Filename: pageConf.TargetFilename, Content: page}
}

func createBlogIndexPage(pageConf config.Page, header, footer []byte) RenderedPage {
	content, err := os.ReadFile("templates/blog_index.html")
	if err != nil {
		panic(err)
	}
	var page []byte
	page = append(page, header...)
	page = append(page, content...)
	page = append(page, footer...)

	return RenderedPage{Title: pageConf.Title, Filename: pageConf.TargetFilename, Content: page}
}

func createHtmlPage(pageConf config.Page, header, footer []byte) RenderedPage {
	content, err := os.ReadFile(pageConf.ContentPath)
	if err != nil {
		panic(err)
	}
	var page []byte
	page = append(page, header...)
	page = append(page, content...)
	page = append(page, footer...)

	return RenderedPage{Title: pageConf.Title, Filename: pageConf.TargetFilename, Content: page}
}
