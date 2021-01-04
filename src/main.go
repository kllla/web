package main

import (
	"github.com/kllla/web/src/service/client"
	"github.com/kllla/web/src/service/invite"
	"github.com/kllla/web/src/service/post"
	"github.com/kllla/web/src/service/render"
	"github.com/kllla/web/src/service/render/ascii"
	"github.com/kllla/web/src/service/shorten"
	"log"
	"net/http"
	"strings"
)

var renderer = render.NewRender()

var clientSvc = client.NewService()
var shortenSvc = shorten.NewService()
var postSvc = post.NewService()
var inviteSvc = invite.NewService()

func main() {
	log.Printf("Listening on port %s", "8080")

	mux := http.NewServeMux()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/home", homeHandler)
	mux.HandleFunc("/post/", postsHandler)
	mux.HandleFunc("/shorten", shortenHandler)
	mux.HandleFunc("/invites", inviteHandler)
	mux.HandleFunc("/register/", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/logout", logoutHandler)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}

func inviteHandler(w http.ResponseWriter, r *http.Request) {
	isAuthForPosting := clientSvc.AuthenticationCheck(w, r)
	if isAuthForPosting {
		switch r.Method {
		case http.MethodGet:
			createdBy := clientSvc.GetSessionUsername(w, r)
			banner := ascii.RenderString(createdBy)
			invitesCreatedBy := inviteSvc.GetAllInvitesCreatedBy(createdBy)
			var intInvites []interface{}
			intInvites = append(intInvites, invitesCreatedBy)
			renderer.RenderPageWithOptions(w, render.Invites, &render.PageRenderOptions{
				Banner:        banner,
				AuthedNavBar:  isAuthForPosting,
				AuthedContent: true,
				Indata:        intInvites,
			})
			return
		case http.MethodPost:
			createdBy := clientSvc.GetSessionUsername(w, r)
			banner := ascii.RenderString(createdBy)
			inviteSvc.CreateInvite(w, r, createdBy)
			ints := inviteSvc.GetAllInvitesCreatedBy(createdBy)
			var intInvites []interface{}
			intInvites = append(intInvites, ints)
			renderer.RenderPageWithOptions(w, render.Invites, &render.PageRenderOptions{
				Banner:        banner,
				AuthedNavBar:  isAuthForPosting,
				AuthedContent: true,
				Indata:        intInvites,
			})
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	isAuthForPosting := clientSvc.AuthenticationCheck(w, r)
	if isAuthForPosting {
		createdBy := clientSvc.GetSessionUsername(w, r)
		banner := ascii.RenderString(createdBy)
		switch r.Method {
		case http.MethodGet:
			shortenedURLS := shortenSvc.GetAllShortenedURLsCreatedBy(createdBy)
			var intUrls []interface{}
			intUrls = append(intUrls, shortenedURLS)
			renderer.RenderPageWithOptions(w, render.ShortenedURLs, &render.PageRenderOptions{
				Banner:        banner,
				AuthedNavBar:  true,
				AuthedContent: true,
				Indata:        intUrls,
			})
			return
		case http.MethodPost:
			_, err := shortenSvc.CreatedShortenedURL(w, r, createdBy)
			if err != nil {
				renderer.RenderPageWithOptions(w, render.ShortenedURLs, &render.PageRenderOptions{
					Banner:        banner,
					AuthedNavBar:  true,
					AuthedContent: true,
					Indata:        nil,
				})
			}
			shortenedURLS := shortenSvc.GetAllShortenedURLsCreatedBy(createdBy)
			renderer.RenderPageWithOptions(w, render.ShortenedURLs, &render.PageRenderOptions{
				Banner:        banner,
				AuthedNavBar:  true,
				AuthedContent: true,
				Indata:        shorten.RenderWrap(shortenedURLS),
			})
			return
		}
	}
	return
}

func postsHandler(w http.ResponseWriter, r *http.Request) {
	isAuthForPosting := clientSvc.AuthenticationCheck(w, r)
	pathID := strings.TrimPrefix(r.URL.Path, "/post/")
	var postByID []*post.Post
	pathIDok := pathID != ""
	if pathIDok {
		postByID = postSvc.GetPostByID(pathID)
	}
	author := clientSvc.GetSessionUsername(w, r)
	banner := ascii.RenderString(author)
	switch r.Method {
	case http.MethodGet:
		if isAuthForPosting {
			if pathIDok {
				action := postSvc.GetActionFromFormData(w, r)
				if action == "delete" {
					postSvc.DeletePost(postByID[0].ID)
				} else {
					renderer.RenderPageWithOptions(w, render.Posts, &render.PageRenderOptions{
						Banner:        "",
						AuthedNavBar:  isAuthForPosting,
						AuthedContent: false,
						Indata:        post.RenderWrap(postByID),
					})
					return
				}
			}
			userPosts := postSvc.GetAllPostsForUsername(author)
			var intPosts []interface{}
			intPosts = append(intPosts, userPosts)
			renderer.RenderPageWithOptions(w, render.Posts, &render.PageRenderOptions{
				Banner:        banner,
				AuthedNavBar:  true,
				AuthedContent: true,
				Indata:        intPosts,
			})
			return
		}
		if pathIDok && (len(postByID) == 1 && postByID[0].Public) {
			renderer.RenderPageWithOptions(w, render.Posts, &render.PageRenderOptions{
				Banner:        "",
				AuthedNavBar:  isAuthForPosting,
				AuthedContent: false,
				Indata:        post.RenderWrap(postByID),
			})
			return
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	case http.MethodPost:
		if isAuthForPosting {
			formDataPost := postSvc.GetPostFromFormData(w, r, author)
			postSvc.CreatePost(formDataPost)
			userPosts := postSvc.GetAllPostsForUsername(author)
			renderer.RenderPageWithOptions(w, render.Posts, &render.PageRenderOptions{
				Banner:        banner,
				AuthedNavBar:  true,
				AuthedContent: true,
				Indata:        post.RenderWrap(userPosts),
			})
			return
		}
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	clientSvc.UnAuthentication(w, r)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		renderer.RenderPageWithOptions(w, render.Login, &render.PageRenderOptions{
			Banner:        "",
			AuthedNavBar:  false,
			AuthedContent: false,
			Indata:        nil,
		})
		return
	case http.MethodPost:
		if clientSvc.VerifyCredentialsAndAuthenticate(w, r) {
			http.Redirect(w, r, "/home", http.StatusTemporaryRedirect)
			return
		}
		renderer.RenderPageWithOptions(w, render.Login, &render.PageRenderOptions{
			Banner:        "",
			AuthedNavBar:  false,
			AuthedContent: false,
			Indata:        nil,
		})
		return
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if id := strings.TrimPrefix(r.URL.Path, "/"); id != "" {
		ok, shrtned := shortenSvc.GetUrlIfIDisValid(id)
		if ok {
			http.Redirect(w, r, shrtned.LongUrl, http.StatusTemporaryRedirect)
		}
	} else {
		isAuthed := clientSvc.AuthenticationCheck(w, r)
		renderer.RenderPageWithOptions(w, render.Index, &render.PageRenderOptions{
			Banner:        "",
			AuthedNavBar:  isAuthed,
			AuthedContent: false,
			Indata:        post.RenderWrap(postSvc.GetPublicPosts()),
		})
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	isAuthed := clientSvc.AuthenticationCheck(w, r)
	if isAuthed {
		username := clientSvc.GetSessionUsername(w, r)
		banner := ascii.RenderString(username)
		renderer.RenderPageWithOptions(w, render.UserHome, &render.PageRenderOptions{
			Banner:        banner,
			AuthedNavBar:  isAuthed,
			AuthedContent: false,
			Indata: append(invite.RenderWrap(inviteSvc.GetAllInvitesCreatedBy(username)),
				post.RenderWrap(postSvc.GetAllPostsForUsername(username))),
		})
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	return
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	inviteID := strings.TrimPrefix(r.URL.Path, "/register/")
	switch r.Method {
	case http.MethodGet:
		if inviteID == "" {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		} else
		if ok, _ := inviteSvc.GetInviteIfIDisValid(inviteID); ok {
			var inviteIDs []interface{}
			inviteIDs = append(inviteIDs, inviteID)
			renderer.RenderPageWithOptions(w, render.Register, &render.PageRenderOptions{
				Banner:        "",
				AuthedNavBar:  false,
				AuthedContent: false,
				Indata:        inviteIDs,
			})
			return
		}
	case http.MethodPost:
		if clientSvc.CreateCredentials(w, r) {
			r.ParseForm()
			inviteID := r.FormValue("invite")
			inviteSvc.DeleteInvite(inviteID)
			renderer.RenderPageWithOptions(w, render.Login, &render.PageRenderOptions{
				Banner:        "",
				AuthedNavBar:  false,
				AuthedContent: false,
				Indata:        nil,
			})
			return
		}
		return
	}
}
