package auth

func extractSession() string{

	return "loggedIn"
}

func GetSession() string{
	return extractSession()	
}
