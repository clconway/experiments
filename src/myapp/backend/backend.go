package backend

import (
	"html/template"
	"net/http"

	"appengine"
	"appengine/memcache"
	"myapp/suggest"
)

func init() {
	http.HandleFunc("/populate", populate)
}

func populate(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	keyword := r.FormValue("key")
	suggest := &Suggest{keyword, []string{keyword, keyword, keyword, keyword}}

	//Put the requested item into memcache
	item := &memcache.Item{
		Key:    suggest.Keyword,
		Object: suggest.Suggestions,
	}
	if err := memcache.Gob.Set(c, item); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	if err := populateTemplate.Execute(w, suggest); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var populateTemplate = template.Must(template.New("add_item").Parse(`
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
  </body>
</html>
`))
