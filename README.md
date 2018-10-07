# My dev.to

> This is a learning project. In other words: I'm messing around, things get
> broken every ~~day~~ commit.

## The idea

> This is what I'm aiming to achieve. Maybe it isn't what is actually happening
> though. See for yourself I guess.

I'm trying to build a *reliable*, *performant* web application like [dev.to][]. The priority is
the user end-result.

### Performant

Server side rendering is good for SEO and first paint. It's pain once you're
browsing though. The answer is to combine both worlds:

> Server side rendering by simulating what a browser would do. And then, the
> client's browser does the job.

How it works? It's explained a bit below...

### Reliable

> Right now the tests are a mess. Please don't look at the code :laugh:

Testing. The user doesn't care if all the unit test pass, he just wants the
*result* to work.

This means that I test the *result*. Unit test are there to facilitate
*debugging*, not to *detect* bugs.

For example, to check whether the action `/api/posts/write` works (it
creates/updates a post in the database), we check that when we *then* call
`/api/posts/get?id=:id`, we do get the post.

This doesn't mean you can't check that `/api/posts/write` made the right SQL
request, whether it called `users.Current()` to check that the user was logged
in, etc... It just means that these are *extras* tests. Just to facilitate the
dev's job (which isn't the priority. The priority is the **user**).

## How it's organized

> Any battle is already won before it has even fought

Yep, I just put a life quote in my README. I know.

Anyway, I think it can be applied to programs as well. A well structured program
will perform much better that a clunky one, that's known. Therefore, this is what
I'm trying to figure out (and why the commit rate is so slow): the right
structure.

### Timeline

```
|<- request
|
|- get required components from tetsu.json
|
|-> send them off to NodeJS
|   use the /api to render stuff
|<- get the HTML back
|
|-> send the html to the client
|   First paint, DOMContentLoaded
|-> push to the client the components they are going to need
|   (based on the link there are on the page)

TODO: user interaction (shouldn't be to hard)
```

Here's the file structure I've decided to adopt for now.

#### `/`

This is where the `main.go` file lives. The rest is just meta data for the
project. This file starts *everything* it can start without knowing anything
about any connection. This globally means initiate every services (see below).

#### `/services`

Services are just external resources that the API uses.

Every service creates their own package (therefore, they're in their own folder).

A good rule of thumb to know whether a package should be a service is whether
it needs to be *initiated*. If it does, then it most likely needs to be a
service.

To be initialized, a service must implement the `init` method in the package
(Go will automatically call it *at most once*, when it is imported), and
*import it* in `main.go`, even if it doesn't use it straight away, like so:

```go
import _ "myservices" // just for side effects (init)
```

Services can be tiny, and they can do *very* different things. Have a look:

`/services/db`: manages the connection with database. Makes sure that there
is only one connection.

`/services/uli`: a sort of log wrapper. It's displays information about the
request.

#### `/api`

The API is the core of the application. It's a REST API, although this might
change if I wake up feeling like messing around with GraphQL, this'll change.

The API is split up into controllers, that are then split up into actions (the
names are stolen from CakePHP).

Therefore, this is reflected in the both the URL and the file structure.

Every controller is a sub-folder in `/api`, and every action is a handler in
this sub-folder.

Every controller must export a function (`Manage`), which will manage the
routing of URLs to actions. This function is usually defined in a file with the
same name as the controller.

```
/api
    /posts
        get.go
        list.go
        write.go
        posts.go // Manage is in this file
    /users
        auth.go
        current.go
        users.go
    ...
```

Controllers can *depend* on each other. For example, the controller `users` can
export a function `Current` that returns the current user's data, and other
controller can use it by doing `u := users.Current()`.

Those kind of function should be quite simple, and therefore defined in the
eponymous file (`posts.go` for the controller `posts`)

Various util functions can be stored in the `/api` package directly. This is
used for functions such as writing JSON from an object, internal error
messages, etc...

*Note: right now, those functions are in their own package, `resp` which is a
service. It's a bad idea. It's going to be moved*

### `/web`

> I haven't messed around with the front end yet. This is what I think I will
> do, and it should kind of work.

*Just some thoughts...*

This folder manages the views. Now, in order to be able to do both server side
rendering *and* client side rendering, I need one thing: a server side
identical to a client side.

NodeJS, here I come.

Views are actually *extremely* simple. They are just a list components. More on
that below.

The template they are plugged in is always the same for every dynamic page on
the website. It contains the HTML.

#### `/web/components`

Every component will live in its own `.js` files (yes, no preprocessor), under
the `components` folder. The filename should exactly match the name of the
component.

What I mean by concatenable components is that if you want to load 2
components, instead of making 2 requests, you add one after the other, and
serve it as a single file. No bundler. This means components should be in IIFE,
like so:

```js
;(function () {
    // your component code goes here.
})();
```

They have the role of fetching the data they need, and add it to the page.

Note that components may depend on other components, in which case they need to
be added *explicitely* to each view they occur in. See below.

#### `/web/<controller>.go`

Since views are very simple, they can all be grouped in one single file per
controller.

Each view is a function (an `http.Handler`) that is going to send the
components needed to

[dev.to]: https://dev.to
