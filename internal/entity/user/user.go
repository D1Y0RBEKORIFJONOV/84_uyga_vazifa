package userentity

type (
	User struct {
		FirstName  string `json:"first_name" bson:"first_name" redis:"first_name"`
		LastName   string `json:"last_name" bson:"last_name" redis:"last_name"`
		Email      string `json:"email" bson:"email" redis:"email"`
		Password   string `json:"password" bson:"password" redis:"password"`
		Role       string `json:"role" bson:"role" redis:"role"`
		UserID     string `json:"user_id" bson:"user_id" redis:"user_id"`
		SecretCode string `json:"secret_code" bson:"-" redis:"secret_code"`
	}
	CreateUser struct {
		FirstName string `json:"first_name" bson:"first_name"`
		LastName  string `json:"last_name" bson:"last_name"`
		Email     string `json:"email" bson:"email"`
		Password  string `json:"password" bson:"password"`
	}
	LoginRequest struct {
		Email    string `json:"email" bson:"email"`
		Password string `json:"password" bson:"password"`
	}
	VerifyRequest struct {
		Email  string `json:"email" bson:"email"`
		Secret string `json:"secret" bson:"secret"`
	}
	Status struct {
		UserID   string     `json:"user_id" bson:"user_id"`
		Messages []*Message `json:"messages" bson:"messages"`
	}
	Message struct {
		CreateAt string `json:"create_at" bson:"create_at"`
		Status   string `json:"status" bson:"status"`
	}
	Token struct {
		AccessToken  string `json:"access_token" bson:"access_token"`
		RefreshToken string `json:"refresh_token" bson:"refresh_token"`
	}
)
