package entity

type UserInfoPost struct {
	UserId        int64 `db:"user_id"`
	ForwardPostId int64 `db:"forward_post_id"`
	// Description string    `db:"description"`
	// CreatedAt   time.Time `db:"created_at"`
	// UpdatedAt   time.Time `db:"updated_at"`
}
