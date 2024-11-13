-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens(
	token,
	created_at,
	updated_at,
	user_id,
	expires_at,
	revoked_at
	)VALUES(
	$1,
	NOW(),
	NOW(),
	$2,
	$3,
	$4
	)RETURNING *;

-- name: GetRefreshToken :one
Select * FROM refresh_tokens WHERE token= $1 ;

-- name: GetUserFromRefreshToken :one
SELECT * From refresh_tokens
INNER JOIN users ON refresh_tokens.user_id = users.id
WHERE users.id = refresh_tokens.user_id ;


-- name: RevokeRefreshToken :one
UPDATE refresh_tokens
SET revoked_at = NOW() , updated_at = NOW()
WHERE token = $1 RETURNING * ;


-- name: DeleteRefreshToken :exec
SELECT * FROM refresh_tokens WHERE token = $1;
