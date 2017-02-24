// +build appengine

package employees

import "net/http"

func init() {
	store := AppEngineStore{}
	handler := GetHandler(&store)
	http.Handle("/", handler)
}
