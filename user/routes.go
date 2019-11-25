package user

import "github.com/zhughes3/website/router"

var Routes = []router.Route{
	router.Route{
		"UsersIndex",
		"GET",
		"/users",
		IndexHandler,
		false,
	},
	router.Route{
		"UsersShow",
		"GET",
		"/users/{userId}",
		ShowHandler,
		true,
	},
	router.Route{
		"UsersCreate",
		"POST",
		"/users",
		CreateHandler,
		false,
	},
	router.Route{
		"UsersLogin",
		"POST",
		"/users/login",
		LoginHandler,
		false,
	},
	router.Route{
		"UsersDelete",
		"DELETE",
		"/users/{userId}",
		DeleteHandler,
		true,
	},
	router.Route{
		"UsersUpdate",
		"PUT",
		"/users/{userId}",
		UpdateHandler,
		true,
	},
}
