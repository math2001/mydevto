package posts

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/math2001/mydevto/api"
	"github.com/math2001/mydevto/services/db"
	"github.com/math2001/mydevto/services/uli"
	"github.com/math2001/sibu"
)

func list(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	q := r.URL.Query()
	var (
		limit  = -1
		userid = -1
		err    error
	)
	if lim := q.Get("limit"); lim != "" {
		limit, err = strconv.Atoi(lim)
		if err != nil {
			uli.Printf(ctx, "Invalid limit @ postsGet: %s", err)
			api.Error(w, r, http.StatusBadRequest,
				"Invalid `limit`. Should be a number")
			return
		}
	}
	if ui := q.Get("userid"); ui != "" {
		userid, err = strconv.Atoi(ui)
		if err != nil {
			uli.Printf(ctx, "Invalid userid @ postsGet: %s", err)
			api.Error(w, r, http.StatusBadRequest, "Invalid `userid`. Should be a number")
			return
		}
	}
	b := sibu.Sibu{}
	b.Add(`SELECT p.id, p.title, p.content, p.written, p.updated, u.name, u.username,
	u.avatar, u.bio, u.url, u.email, u.location FROM posts p JOIN users u
	ON p.userid=u.id`)
	where := sibu.OpClause{}
	if userid != -1 {
		where.Add("AND", "p.userid={{ p }}", userid)
	}
	if !where.Empty() {
		b.AddClause("WHERE", where)
	}
	if limit != -1 {
		b.Add("LIMIT {{ p }}", limit)
	}
	sql, args, err := b.Query()
	if err != nil {
		uli.Printf(ctx, "Errored building sql request @ postIndex: %s", err)
		api.InternalError(w, r)
		return
	}
	rows, err := db.QueryContext(ctx, sql, args...)
	if err != nil {
		fmt.Println("query good", err)
		uli.Printf(ctx, "Errored querying @ postIndex: %s", err)
		api.InternalError(w, r)
		return
	}
	var posts []db.Post
	for rows.Next() {
		p := db.Post{}
		u := &p.User
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Written, &p.Updated, &u.Name,
			&u.Username, &u.Avatar, &u.Bio, &u.URL, &u.Email, &u.Location)
		if err != nil {
			uli.Printf(ctx, "Errored scanning rows @ postIndex: %s", err)
			api.InternalError(w, r)
			return
		}
		posts = append(posts, p)
	}
	if err := rows.Err(); err != nil {
		uli.Printf(ctx, "Errored during iteration @ postIndex: %s", err)
		api.InternalError(w, r)
		return
	}
	api.Encode(w, r, posts, http.StatusOK)
}
