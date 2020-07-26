package auth

// AuthenticateByAPIKey resolves API Key into a username.
func AuthenticateByAPIKey(apiKey string) (username string, found bool) {
	username, found = map[string]string{
		"user1key1": "user1",
		"user1key2": "user1",
		"user2key1": "user2",
		"user3key1": "user3",
		"user4key1": "unknown", // Not present in policy.
	}[apiKey]
	return
}
