package users

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gbrlsnchs/jwt"
	"github.com/gorilla/mux"
	"github.com/math2001/mydevto/api"
	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/services/uli"
	"github.com/mitchellh/mapstructure"
)

var jwtsigner jwt.Signer

// Manage is delegated the charges of mapping routes to functions by the main
// package
func Manage(r *mux.Router) {
	jwtsecret := os.Getenv("JWTSECRET")
	if jwtsecret == "" {
		log.Fatal("$JWTSECRET must be set")
	}
	jwtsigner = jwt.NewHS256(jwtsecret)

	r.Handle("/", api.ListRoutes(r)).Methods("GET")
	// in there documentation, github say they send a POST request, but they
	// actually send a GET... :(
	r.HandleFunc("/auth", auth).Methods("GET")
	r.HandleFunc("/current", current).Methods("GET")
}

// Current returns the current user's information from the sessions. It returns
// nil if he isn't connected
func Current(ctx context.Context, r *http.Request) *db.User {
	// session, err := sess.Store().Get(r, api.SessionAuth)
	uli.Printf(ctx, "getting user information from cookie...")
	cookie, err := r.Cookie(api.JWT)
	if err == http.ErrNoCookie {
		uli.Printf(ctx, "cookie not found.")
		return nil
	} else if err != nil {
		uli.Printf(ctx, "could not get cookie for unexpected reason: %s", err)
		return nil
	}
	uli.Printf(ctx, "parsing information from %q", cookie.Value)
	payload, sig, err := jwt.Parse(cookie.Value)
	if err != nil {
		uli.Printf(ctx, "could not parse payload from jwt: %s", err)
		return nil
	}
	if err = jwtsigner.Verify(payload, sig); err != nil {
		uli.Security(ctx) // someone's probably messing aroud with the JWT
		uli.Printf(ctx, "could not verify payload signature from jwt: %s", err)
		return nil
	}
	var jot = db.JWTToken{
		User: &db.User{},
		JWT:  &jwt.JWT{},
	}
	if err = jwt.Unmarshal(payload, &jot); err != nil {
		uli.Printf(ctx, "could not unmarshal payload from jwt: %s", err)
		return nil
	}
	u := &db.User{}
	mapstructure.Decode(jot.User, u)
	return u
}
