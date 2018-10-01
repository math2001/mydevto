WITH test1 AS (
	-- create a fake user
	INSERT INTO users (username, name, bio, email, service)
	VALUES (
		'test1',
		'Test NO1',
		'I''m just a test',
		'test1@tests.com',
		'fake'
	) RETURNING (id)
	), test2 AS (
	-- create an other fake user
	INSERT INTO users (username, name, bio, email, service)
	VALUES (
		'test2',
		'Test NO2',
		'I''m just an other test',
		'test2@tests.com',
		'fake'
	) RETURNING (id)
	), post1 AS (
	-- create a post
	INSERT INTO posts (userid, title, content)
	SELECT test1.id,
	'The first posts is out!',
	'Indeed, on the Monday 01 October 2018, the first post was inserted!'
	FROM test1
	RETURNING (id)
	), post2 AS (
	-- create an other post
	INSERT INTO posts (userid, title, content)
	SELECT
	test1.id,
	'An other one',
	'Yep, that''s right. An other post. Crazy, right?'
	FROM test1
	RETURNING (id)
	), post3 AS (
	-- create an other post
	INSERT INTO posts (userid, title, content)
	SELECT
		test1.id,
		'It''s me again!',
		'My last posts though...'
	FROM test1
	RETURNING(ID)
)

-- create a comment on the first post by test1
INSERT INTO comments (userid, postid, content) SELECT
	test1.id,
	post1.id,
	'Yay!'
FROM test1, post1
