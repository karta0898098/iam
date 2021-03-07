package domain

type SignupOption struct {
	Nickname  string // Nickname user nickname
	FirstName string // FirstName user first name
	LastName  string // LastName user last name
	Email     string // Email user email address
	Avatar    string // Avatar means user profile picture URL
}

type UpdateProfileOption struct {
	Nickname  *string // Nickname user nickname
	FirstName *string // FirstName user first name
	LastName  *string // LastName user last name
	Email     *string // Email user email address
	Avatar    *string // Avatar means user profile picture URL
}
