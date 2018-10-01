CREATE TABLE IF NOT EXISTS users (
	id        SERIAL,
	token     VARCHAR(255),
	service   VARCHAR(1024) NOT NULL,
	email     VARCHAR(255) NOT NULL,
	username  VARCHAR(255),
	avatar    VARCHAR(255),
	name      VARCHAR(255),
	bio       VARCHAR(255),
	url       VARCHAR(255),
	updated   TIMESTAMPTZ DEFAULT  now(),
	PRIMARY KEY (id),
	UNIQUE (email, service)
);

CREATE TABLE IF NOT EXISTS posts (
	id       SERIAL,
	userid   INTEGER,
	title    VARCHAR(255),
	content  TEXT,
	written  TIMESTAMPTZ DEFAULT NOW(),
	updated  TIMESTAMPTZ DEFAULT NOW(),
	PRIMARY KEY (id),
	FOREIGN KEY (userid) REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS comments (
	id       SERIAL,
	userid   INTEGER,
	postid   INTEGER,
	content  TEXT,
	written  TIMESTAMPTZ DEFAULT NOW(),
	updated  TIMESTAMPTZ DEFAULT NOW(),
	PRIMARY KEY (id),
	FOREIGN KEY (postid) REFERENCES posts(id),
	FOREIGN KEY (userid) REFERENCES users(id)
);
