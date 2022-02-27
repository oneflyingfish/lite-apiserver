package ServerHandlers

import "net/http"

func Hello(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello to see you, LiteKube is here!"))
}
