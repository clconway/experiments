package suggest

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/memcache"
)

type SuggestList []string

type Suggest struct {
	Keyword     string
	Suggestions SuggestList
}

func init() {
	http.HandleFunc("/suggest", suggest)
}

func suggest(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	keyword := r.FormValue("keyword")

	staticItem := &memcache.Item{
		Key:    "foo",
		Object: []string{"foo", "bar", "baz"},
	}
	if err := memcache.Gob.Set(c, staticItem); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	var suggestions SuggestList
	if _, err := memcache.Gob.Get(c, keyword, &suggestions); err == memcache.ErrCacheMiss {
		suggestions = []string{"cache miss", keyword}
	}

	if err := suggestTemplate.Execute(w, &Suggest{keyword, suggestions}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var suggestTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Suggestions</title>
  </head>
  <body>
    <b>{{.Keyword}}</b>
    <ul>
    {{range .Suggestions}}
      <li>{{.}}</li>
    {{end}}
    </ul>
    <form action="/suggest" method="post">
      <div><textarea name="keyword" rows="1" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`))
