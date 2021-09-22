package render

import (
	"github.com/russross/blackfriday"
	"regexp"
	"strings"
	"wiblog/pkg/conf"
	"wiblog/pkg/model"
	"wiblog/tools"
)

// blackfriday 配置
const (
	commonHtmlFlags = 0 |
		blackfriday.HTML_TOC |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES |
		blackfriday.HTML_NOFOLLOW_LINKS

	commonExtensions = 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS
)

var (
	// 渲染markdown操作和截取摘要操作
	regIdentifier = regexp.MustCompile(conf.Conf.WiBlogApp.General.Identifier)
	// header
	regHeader = regexp.MustCompile("</nav></div>")
)

// RenderPage 渲染markdown
func RenderPage(md []byte) []byte {
	renderer := blackfriday.HtmlRenderer(commonHtmlFlags, "", "")
	return blackfriday.Markdown(md, renderer, commonExtensions)
}

// GenerateExcerptMarkdown 生成预览和描述
func GenerateExcerptMarkdown(article *model.Article) {
	blogapp := conf.Conf.WiBlogApp

	if strings.HasPrefix(article.Content, blogapp.General.DescPrefix) {
		index := strings.Index(article.Content, "\r\n")
		prefix := article.Content[len(blogapp.General.DescPrefix):index]

		article.Desc = tools.IgnoreHtmlTag(prefix)
		article.Content = article.Content[index:]
	}

	// 查找目录
	content := RenderPage([]byte(article.Content))
	index := regHeader.FindIndex(content)
	if index != nil {
		article.Header = string(content[0:index[1]])
		article.Content = string(content[index[1]:])
	} else {
		article.Content = string(content)
	}

	// excerpt
	index = regIdentifier.FindStringIndex(article.Content)
	if index != nil {
		article.Excerpt = tools.IgnoreHtmlTag(article.Content[:index[0]])
		return
	}
	uc := []rune(article.Content)
	length := blogapp.General.Length
	if len(uc) < length {
		length = len(uc)
	}
	article.Excerpt = tools.IgnoreHtmlTag(string(uc[0:length]))
}
