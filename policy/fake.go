package policy

type Policy int

const (
	Default Policy = iota
	Smooth
	Spiky
)

func FindByUsername(username string) (policy Policy, found bool) {
	policy, found = map[string]Policy{
		"user1": Spiky,
		"user2": Spiky,
		"user3": Smooth,
	}[username]
	return
}
