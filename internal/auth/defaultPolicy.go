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
	// api/user
	{
		subject: "admin", object: "/api/user", action: "create",
	},
	{
		subject: "admin", object: "/api/user", action: "read",
	},
	{
		subject: "admin", object: "/api/user", action: "update",
	},
	{
		subject: "admin", object: "/api/user", action: "delete",
	},
}
