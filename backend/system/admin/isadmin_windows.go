//go:build windows

package isadmin

import (
	"golang.org/x/sys/windows"
)

// Check if the program has administrative privileges.
func Check() bool {
	var sid *windows.SID

	_ = windows.AllocateAndInitializeSid(
		&windows.SECURITY_NT_AUTHORITY,
		2,
		windows.SECURITY_BUILTIN_DOMAIN_RID,
		windows.DOMAIN_ALIAS_RID_ADMINS,
		0, 0, 0, 0, 0, 0,
		&sid)

	admin, _ := windows.Token(0).IsMember(sid)

	return admin
}
