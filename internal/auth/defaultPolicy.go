package auth

type policySet struct {
	Subject string
	Object  string
	Action  string
}

type Policy struct {
	PType   string
	Subject string
	Object  string
	Action  string
}

type GroupingPolicy struct {
	PType string
	User  string
	Role  string
}

var DefaultPolicyList = []policySet{
	// User
	// api/me
	{
		Subject: "user", Object: "/api/me", Action: "read",
	},
	{
		Subject: "user", Object: "/api/me", Action: "update",
	},
	// Admin
	// api/me
	{
		Subject: "admin", Object: "/api/me", Action: "read",
	},
	{
		Subject: "admin", Object: "/api/me", Action: "create",
	},
	{
		Subject: "admin", Object: "/api/me", Action: "update",
	},
	// api/users
	{
		Subject: "admin", Object: "/api/users", Action: "create",
	},
	{
		Subject: "admin", Object: "/api/users", Action: "read",
	},
	{
		Subject: "admin", Object: "/api/users", Action: "update",
	},
	{
		Subject: "admin", Object: "/api/users", Action: "delete",
	},
}
