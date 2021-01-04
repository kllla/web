package render

import (
	"fmt"
	"github.com/kllla/web/src/service/client/credentials"
	"github.com/kllla/web/src/service/invite"
	"github.com/kllla/web/src/service/post"
	"github.com/kllla/web/src/service/shorten"
	"html/template"
)

const (
	constBanner = `
██   ██ ██      ██         ██       █████  
██  ██  ██      ██         ██      ██   ██ 
█████   ██      ██         ██      ███████ 
██  ██  ██      ██         ██      ██   ██ 
██   ██ ███████ ███████ ██ ███████ ██   ██ 
`
	baseTemplate = `<!DOCTYPE html>
<html lang="en">
<constBanner>
    <link rel="stylesheet" type="text/css" href="https://storage.googleapis.com/infra-person.appspot.com/kllpw.css">
    <script src="https://storage.googleapis.com/infra-person.appspot.com/kllpw.js"></script>
</constBanner>
<head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
</head>
    <div class="fixed-header">
		<div class="container">
			<div class="ascii-art"><pre>{{.Banner}}</pre></div>
		</div>
		<div class="container">
            {{.NavBar}}
        </div>
    </div>
<div class="body-container">
{{ .Body }}
</div>
</body>
</html>`
)

var Pages map[string]*pageImpl

func MakePages() map[string]Page {
	Pages := make(map[string]Page)
	Pages[Index.TemplateTitle] = Index
	Pages[Login.TemplateTitle] = Login
	Pages[Register.TemplateTitle] = Register
	Pages[UserHome.TemplateTitle] = UserHome
	Pages[Posts.TemplateTitle] = Posts
	Pages[ShortenedURLs.TemplateTitle] = ShortenedURLs
	Pages[Invites.TemplateTitle] = Invites
	return Pages
}

type navEntry struct {
	path string
	text string
}

var defaultNavBar = makeDefaultNavBar()
var authedNavBar = makeAuthedNavBar()

func makeDefaultNavBar() []*navEntry {
	navbr := []*navEntry{
		{path: "/", text: "/"},
		{path: "/login", text: "/login"},
	}
	return navbr
}

func makeAuthedNavBar() []*navEntry {
	navbr := []*navEntry{
		{path: "/", text: "/"},
		{path: "/home", text: "/home"},
		{path: "/post", text: "/post"},
		{path: "/shorten", text: "/shorten"},
		{path: "/invites", text: "/invites"},
		{path: "/logout", text: "/logout"},
	}
	return navbr
}

type Page interface {
	GetTemplateTitle() string
	RenderPageDataWithOptions(options *PageRenderOptions) map[string]interface{}
	ToTemplate() *template.Template
}

func (p *pageImpl) ToTemplate() *template.Template {
	t := template.New(p.TemplateTitle) // Create a template.
	t, err := t.Parse(p.Template)      // Parse template file.
	if err != nil {
		fmt.Print("render err", err)
	}
	return t
}

func (p *pageImpl) GetTemplateTitle() string {
	return p.TemplateTitle
}

func (p *pageImpl) RenderPageDataWithOptions(options *PageRenderOptions) map[string]interface{} {
	p.BuildBanner(options.Banner)
	p.BuildBody(options.AuthedContent, options.Indata)
	p.BuildNavBar(options.AuthedNavBar)
	return map[string]interface{}{
		"Title":  p.Title,
		"Banner": p.Banner,
		"NavBar": p.NavBar,
		"Body":   p.Body,
	}
}
func (p *pageImpl) BuildBanner(indata string) {
	if indata == "" {
		p.Banner = constBanner
	} else {
		p.Banner = indata
	}
}

func (p *pageImpl) BuildBody(authed bool, indata ...interface{}) bool {
	body := p.ContentToBody(authed, indata)
	p.Body = template.HTML(body)
	return authed
}

func (p *pageImpl) BuildNavBar(authed bool) {
	str := "<div class=\"nav-container\">" +
		"<nav>"
	if authed {
		for _, navEntry := range authedNavBar {
			str += fmt.Sprintf("\n<a href=\"%s\">%s</a>", navEntry.path, navEntry.text)
		}
	} else {
		for _, navEntry := range defaultNavBar {
			str += fmt.Sprintf("\n<a href=\"%s\">%s</a>", navEntry.path, navEntry.text)
		}
	}
	str += "\n</nav></div>"
	p.NavBar = template.HTML(str)
}

type pageImpl struct {
	TemplateTitle string
	Template      string

	Title     string
	Banner    string
	NavBarRaw []*navEntry
	NavBar    template.HTML

	BodyRaw       string
	Body          template.HTML
	ContentToBody func(authed bool, inData ...interface{}) string
}

type NavBar map[string]string

var Index = &pageImpl{
	Title:         "kllla",
	NavBarRaw:     makeDefaultNavBar(),
	Banner:        constBanner,
	TemplateTitle: "index",
	Template:      baseTemplate,
	ContentToBody: func(authed bool, inData ...interface{}) string {

		str := processIndata(authed, true, "", inData)
		return str
	},
}

var Login = &pageImpl{
	Title:         "kllla",
	Banner:        constBanner,
	TemplateTitle: "login",
	NavBarRaw:     makeDefaultNavBar(),
	Template:      baseTemplate,
	ContentToBody: func(authed bool, inData ...interface{}) string {
		return "<div class=\"login-container\">" +
			"<form action=\"/login\" method=\"POST\">" +
			"<h1>login:</h1>" +
			"<input type=\"text\" placeholder=\"username\" id=\"username\" name=\"username\"></input>" +
			"<input type=\"password\" placeholder=\"password\" id=\"password\" name=\"password\"></input>" +
			"<input type=\"submit\" value=\"login\"></input>" +
			"</form>" +
			"</div>"
	},
}
var Invites = &pageImpl{
	Title:         "invites",
	Template:      baseTemplate,
	NavBarRaw:     authedNavBar,
	Banner:        constBanner,
	TemplateTitle: "invites",
	ContentToBody: func(authed bool, inData ...interface{}) string {
		formStr := ""
		if authed {
			formStr += "<div class=\"post-container\">" +
				"<div class=\"post-container-form\">" +
				"<h1>generate invite</h1>" +
				"<form action=\"/invites\" method=\"post\">" +
				"<input type=\"submit\" name=\"submit\" value=\"generate\"/>" +
				"</form>" +
				"</div>" +
				"</div>"
		}
		body := processIndata(authed, false, "", inData)
		return formStr + body
	},
}

var ShortenedURLs = &pageImpl{
	Title:         "shortened urls",
	Template:      baseTemplate,
	NavBarRaw:     makeAuthedNavBar(),
	Banner:        constBanner,
	TemplateTitle: "shorten",
	ContentToBody: func(authed bool, inData ...interface{}) string {
		body := processIndata(authed, false, "", inData)
		formStr := ""
		if authed {
			formStr += "<div class=\"container\">" +
				"<div class=\"post-container-form\">" +
				"<h1>shorten url</h1>" +
				"<form action=\"/shorten\" method=\"post\">" +
				"<input type=\"text\" name=\"longURL\" placeholder=\"long url\"/>" +
				"<input type=\"submit\" name=\"submit\" value=\"shorten\"/>" +
				"</form>" +
				"</div>" +
				"</div>"
		}
		return formStr + body
	},
}
var Posts = &pageImpl{
	Title:         "posts",
	Template:      baseTemplate,
	NavBarRaw:     makeAuthedNavBar(),
	Banner:        constBanner,
	TemplateTitle: "posts",
	ContentToBody: func(authed bool, inData ...interface{}) string {
		htmlStr := ""
		short := false
		for _, data := range inData {
			switch data.(type) {
			case []interface{}:
				ds := data.([]interface{})
				for _, nestedData := range ds {
					htmlStr = processIndata(authed, short, htmlStr, nestedData)
				}
			case []*post.Post:
				posts := data.([]*post.Post)
				htmlStr += printPostsPost(posts, authed, short)
			}
		}
		body := htmlStr
		formStr := ""
		if authed {
			formStr += "<div class=\"container\">" +
				"<div class=\"post-container-form\">" +
				"<h1>new post:</h1>" +
				"<form action=\"/post/\" method=\"post\">" +
				"<input type=\"text\" name=\"title\" placeholder=\"title\"/>" +
				"<textarea name=\"content\" placeholder=\"content\"></textarea>" +
				"<label for=\"public\">public:</label>" +
				"<input type=\"checkbox\" id=\"public\" name=\"public\" value=\"true\"/>" +
				"<input type=\"submit\" name=\"submit\" value=\"post	\"/>" +
				"</form>" +
				"</div>" +
				"</div>" +
				"<div class=\"post-title\"><h1>your posts:</h1></div>"
		}
		formStr += body
		return formStr
	},
}

func processIndata(authed bool, short bool, htmlStr string, inData ...interface{}) string {
	for _, data := range inData {
		switch data.(type) {
		case []interface{}:
			ds := data.([]interface{})
			for _, nestedData := range ds {
				htmlStr = processIndata(authed, short, htmlStr, nestedData)
			}
		case string:
			posts := data.(string)
			htmlStr += posts
		case []*post.Post:
			posts := data.([]*post.Post)
			htmlStr += printPosts(posts, authed, short)
		case *credentials.Credentials:
			creds := data.(*credentials.Credentials)
			htmlStr += fmt.Sprintf("<div class=\"post-container\">"+
				"</h1> %s </h1>"+
				"</div>", creds.Username)
		case []*shorten.ShortenedURL:
			urls := data.([]*shorten.ShortenedURL)
			htmlStr += printURLS(urls, authed, short)
		case []*invite.Invite:
			invites := data.([]*invite.Invite)
			htmlStr += printInvites(invites, authed, short)
		}
	}
	return htmlStr
}

func printInvites(invites []*invite.Invite, authed bool, short bool) string {

	htmlStr := ""
	for pos, invite := range invites {
		htmlStr += fmt.Sprintf("<div class=\"post-container\">"+
			"<h1>invite %d</h1>"+
			"<a href=\"https://www.kll.la/register/%s\">https://www.kll.la/register/%s</a>"+
			"<div class=\"post-container-author\">%s</div> "+
			"<div class=\"post-container-date\">%s</div>"+
			"</div>", pos+1, invite.InviteID, invite.InviteID, invite.CreatedBy, invite.ExpiryTime.Format("2006/01/02"))
	}
	return htmlStr

}

func printURLS(urls []*shorten.ShortenedURL, authed bool, short bool) string {
	str := ""
	for _, url := range urls {
		str += fmt.Sprintf("<div class=\"post-container\">"+
			"<p><a href=\"%s\">%s</p></a>"+
			"<p><a href=\"https://%s\">%s</p></a>"+
			"<div class=\"post-container-author\">%s</div> "+
			"<div class=\"post-container-date\">%s</div>"+
			"</div>", url.LongUrl, url.LongUrl, url.ShortenedURL, url.ShortenedURL, url.CreatedBy, url.ExpiryTime.Format("2006/01/02"))
	}
	return str
}

func printPosts(posts []*post.Post, authed bool, short bool) string {
	str := ""
	for _, post := range posts {
		authedControls := ""
		if authed {
			authedControls = fmt.Sprintf("<div class=\"editactions\">"+
				"<form action=\"/post/%s\" method=\"get\">"+
				"<input type=\"submit\" value=\"delete\"></input>"+
				"<input type=\"text\" hidden name=\"action\" value=\"delete\"></input>"+
				"</form>"+
				"</div>", post.ID)
		}
		if short {
			str += fmt.Sprintf("<div class=\"post-short-container\">"+
				"<a href=\"/post/%s\">"+
				"<h1>%s</h1></a>"+
				"<div class=\"post-container-date\">%s</div>"+
				"<div class=\"post-container-author\">%s</div>"+
				authedControls+
				"</div>", post.ID, post.Title, post.Author, post.Date.Format("2006/01/02"))
		} else {
			str += fmt.Sprintf("<div class=\"post-container\">"+
				authedControls+
				"<a href=\"%s\">"+
				"<h1>%s</h1></a>"+
				"<p>%s</p>"+
				"<div class=\"post-container-author\">%s</div> "+
				"<div class=\"post-container-date\">%s</div>"+
				"</div>", post.ID, post.Title, post.Content, post.Author, post.Date.Format("2006/01/02"))
		}
	}
	return str
}

func printPostsPost(posts []*post.Post, authed bool, short bool) string {
	str := ""
	for _, post := range posts {
		authedControls := ""
		if authed {
			authedControls = fmt.Sprintf("<div class=\"editactions\">"+
				"<form action=\"/%s\" method=\"get\">"+
				"<input type=\"submit\" value=\"delete\"></input>"+
				"<input type=\"text\" hidden name=\"action\" value=\"delete\"></input>"+
				"</form>"+
				"</div>", post.ID)
		}
		if short {
			str += fmt.Sprintf("<div class=\"post-short-container\">"+
				"<a href=\"/%s\">"+
				"<h1>%s</h1></a>"+
				"<div class=\"post-container-date\">%s</div>"+
				"<div class=\"post-container-author\">%s</div>"+
				authedControls+
				"</div>", post.ID, post.Title, post.Author, post.Date.Format("2006/01/02"))
		} else {
			str += fmt.Sprintf("<div class=\"post-container\">"+
				authedControls+
				"<a href=\"%s\">"+
				"<h1>%s</h1></a>"+
				"<p>%s</p>"+
				"<div class=\"post-container-author\">%s</div> "+
				"<div class=\"post-container-date\">%s</div>"+
				"</div>", post.ID, post.Title, post.Content, post.Author, post.Date.Format("2006/01/02"))
		}
	}
	return str
}

var Register = &pageImpl{
	Title:         "kllla",
	Banner:        constBanner,
	NavBarRaw:     makeDefaultNavBar(),
	TemplateTitle: "register",
	Template:      baseTemplate,
	ContentToBody: func(authed bool, inData ...interface{}) string {
		value := processIndata(authed, false, "", inData)
		return "<div class=\"login-container\">" +
			"<form action=\"/register/\" method=\"POST\">" +
			"<h1>register:</h1>" +
			"<input type=\"text\" placeholder=\"username\" id=\"username\" name=\"username\"></input>" +
			"<input type=\"password\" placeholder=\"password\" id=\"password\" name=\"password\"></input>" +
			"<input type=\"text\" hidden placeholder=\"\" id=\"invite\" name=\"invite\" value=\"" + value + "\"></input>" +
			"<input type=\"submit\" value=\"register\"></input>" +
			"</form>" +
			"</div>"
	},
}

var UserHome = &pageImpl{
	Title:         "kllla",
	Banner:        constBanner,
	TemplateTitle: "userhome",
	NavBarRaw:     makeAuthedNavBar(),
	Template:      baseTemplate,
	ContentToBody: func(authed bool, inData ...interface{}) string {
		str := "<div class=\"post-container\"><div class=\"post-title\"><h1>your activity:</h1></div></div>"
		str += processIndata(authed, true, "", inData)
		return str
	},
}
