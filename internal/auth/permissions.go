package auth

type Permission int32

const (
	PermRead  Permission = 1 << iota // 1
	PermWrite                        // 2
	PermAdmin                        // 4
)

const PermDefault = PermRead | PermWrite            // 3 — all new users
const PermAll     = PermRead | PermWrite | PermAdmin // 7 — admin

func HasPermission(userPerms, required Permission) bool {
	return userPerms&required == required
}
