package auth

type policySet struct {
	subject string
	object  string
	action  string
}

var DefaultPolicyList = []policySet{
	// User
	{
		subject: "user", object: "/api/me", action: "read",
	},
	// Admin
	{
		subject: "admin", object: "/api/me", action: "read",
	},
	{
		subject: "admin", object: "/api/me", action: "create",
	},
	{
		subject: "admin", object: "/api/me", action: "update",
	},
}
