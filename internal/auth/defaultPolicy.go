package auth

type policySet struct {
	subject string
	object  string
	action  string
}

var DefaultPolicyList = []policySet{
	// User
	// api/me
	{
		subject: "user", object: "/api/me", action: "read",
	},
	{
		subject: "user", object: "/api/me", action: "update",
	},
	// Admin
	// api/me
	{
		subject: "admin", object: "/api/me", action: "read",
	},
	{
		subject: "admin", object: "/api/me", action: "create",
	},
	{
		subject: "admin", object: "/api/me", action: "update",
	},
	// api/users
	{
		subject: "admin", object: "/api/users", action: "create",
	},
	{
		subject: "admin", object: "/api/users", action: "read",
	},
	{
		subject: "admin", object: "/api/users", action: "update",
	},
	{
		subject: "admin", object: "/api/users", action: "delete",
	},
}
