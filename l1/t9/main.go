package main

import "fmt"

type User interface {
	GetUsername() string
	HasPermission(string) bool
	GetRole() string
}

type BasicUser struct {
	username    string
	role        string
	permissions map[string]struct{}
}

func NewBasicUser(username string) *BasicUser {
	return &BasicUser{
		username: username,
		role:     "BasicUser",
		permissions: map[string]struct{}{
			"read": {},
		},
	}
}

func (b *BasicUser) GetUsername() string {
	return b.username
}

func (b *BasicUser) HasPermission(permission string) bool {
	_, ok := b.permissions[permission]
	return ok
}

func (b *BasicUser) GetRole() string {
	return b.role
}

type Moderator struct {
	BasicUser
}

func NewModerator(username string) *Moderator {
	m := &Moderator{
		BasicUser: *NewBasicUser(username),
	}
	m.permissions["edit"] = struct{}{}
	m.permissions["ban_user"] = struct{}{}
	return m
}

func (m *Moderator) GetUsername() string {
	return m.username
}

func (m *Moderator) HasPermission(permission string) bool {
	_, ok := m.permissions[permission]
	return ok
}

func (m *Moderator) GetRole() string {
	return m.role
}

type Admin struct {
	Moderator
}

func NewAdmin(username string) *Admin {
	a := &Admin{
		Moderator: *NewModerator(username),
	}
	a.permissions["delete"] = struct{}{}
	a.permissions["manage_roles"] = struct{}{}
	return a
}

func (a *Admin) GetUsername() string {
	return a.username
}

func (a *Admin) HasPermission(permission string) bool {
	_, ok := a.permissions[permission]
	return ok
}

func (a *Admin) GetRole() string {
	return a.role
}

func main() {
	users := []User{
		NewBasicUser("Ivan1"),
		NewModerator("Ivan2"),
		NewAdmin("Ivan3"),
	}

	for _, u := range users {
		fmt.Printf("User: %v, Role: %v\n", u.GetUsername(), u.GetRole())
		fmt.Printf("permissions: \n")
		fmt.Printf("\tread: %v\n", u.HasPermission("read"))
		fmt.Printf("\tban_user: %v\n", u.HasPermission("ban_user"))
		fmt.Printf("\tedit: %v\n", u.HasPermission("edit"))
		fmt.Printf("\tdelete: %v\n", u.HasPermission("delete"))
		fmt.Printf("\tmanage_roles: %v\n\n", u.HasPermission("manage_roles"))
	}
}
