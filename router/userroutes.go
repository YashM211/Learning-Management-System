package router

import (
	"lms/private"
	"lms/util"
)

// SetupUserRoutes func sets up all the user routes
func SetupUserRoutes() {
	USER.Post("/signup", CreateUser)              // Sign Up a user
	USER.Post("/signin", LoginUser)               // Sign In a user
	USER.Get("/logout", Logout)                   // Logs out a user
	USER.Get("/get-access-token", GetAccessToken) // returns a new access_token
	// USER.Get("/googlelogin", GoogleLogin)
	// USER.Get("/googlecallback", GoogleCallback)
	// USER.Get("/getbookstock", GetBookStock)
	USER.Get("/getusername", GetUsername)
	// USER.Post("/createcatagorycookie", CreateCatagoryCookie)
	// USER.Get("/getfilteredbooks", Getfilteredbooks)

	// privUser handles all the private user routes that requires authentication
	privUser := USER.Group("/private")
	privUser.Use(util.SecureAuth()) // middleware to secure all routes for this group
	// privUser.Get("/user", GetUserData)
	//privUser.Post("/logout", private.LogOut)
	privUser.Post("/createcourse", private.Createcourse)
	privUser.Get("/availablecourse", private.AvailableCourses)
	// privUser.Post("/deleteentry", private.DeleteEntry)
	// privUser.Post("/addtocart", private.AddtoCart)
	privUser.Post("/deletecourse", private.DeleteCourse)
	// privUser.Post("/placecartorders", private.PlaceCartOrders)
	privUser.Get("/getyourcourse", private.GetYourCourses)
	// privUser.Get("/getpurchasedata", private.GetPurchaseData)
	// privUser.Post("/createbookcookie", private.CreateBookCookie)
	// privUser.Get("/showbook", private.ShowBook)
	// privUser.Post("/trackpage", private.TrackPackage)

}
