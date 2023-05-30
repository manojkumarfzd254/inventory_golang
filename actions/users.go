package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/x/responder"
	"github.com/pkg/errors"

	"library/models"
)

func UserCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(http.StatusOK, r2.HTML("backend/users/create.plush.html"))
	}

	c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", "Welcome to library!")

	return c.Redirect(http.StatusFound, "/")
}

func UserEdit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Book
	user := &models.User{}

	if err := tx.Find(user, c.Param("ID")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("user", user)
	c.Set("checkID", c.Param("ID"))
	return c.Render(http.StatusOK, r2.HTML("backend/users/edit.plush.html"))
}

func UserShow(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Book
	user := &models.User{}

	if err := tx.Find(user, c.Param("ID")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("user", user)
	return c.Render(http.StatusOK, r2.HTML("backend/users/show.plush.html"))
}

func UserUpdate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Book
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	if err := tx.Find(user, c.Param("ID")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind User to the html form elements
	verrs, err := user.Update(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("user", user)

			return c.Render(http.StatusUnprocessableEntity, r2.HTML("backend/users/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r2.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r2.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", "User successfully Updated.")

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/auth/users/")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(user))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(user))
	}).Respond(c)
}

func UserList(c buffalo.Context) error {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	Users := []models.User{}
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Users from the DB
	if err := q.All(&Users); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		c.Set("pagination", q.Paginator)

		c.Set("users", Users)
		return c.Render(http.StatusOK, r2.HTML("backend/users/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(Users))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(Users))
	}).Respond(c)
	// return c.Render(http.StatusOK, r.HTML("users/new.plush.html"))
}

func UserSave(c buffalo.Context) error {
	user := &models.User{}
	if err := c.Bind(user); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := user.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)
			fmt.Println(verrs)
			// Render again the new.html template that the user can
			// correct the input.
			c.Set("user", user)

			return c.Render(http.StatusUnprocessableEntity, r2.HTML("backend/users/create.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r2.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r2.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "User successfully added."))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/auth/users/")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(user))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(user))
	}).Respond(c)

}
func UserDelete(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Book
	user := &models.User{}

	// To find the Book the parameter book_id is used.
	if err := tx.Find(user, c.Param("ID")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(user); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "User successfully deleted."))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/auth/users/")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(user))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(user))
	}).Respond(c)
}

///////////////////////////////////////////////////////////////////////////////////////////

// UsersNew renders the users form
func UsersNew(c buffalo.Context) error {
	u := models.User{}
	c.Set("user", u)
	return c.Render(http.StatusOK, r.HTML("users/new.plush.html"))
}

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(http.StatusOK, r.HTML("users/new.plush.html"))
	}

	c.Session().Set("current_user_id", u.ID)
	c.Flash().Add("success", "Welcome to library!")

	return c.Redirect(http.StatusFound, "/")
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid != nil {
			u := &models.User{}
			tx := c.Value("tx").(*pop.Connection)
			err := tx.Find(u, uid)
			if err != nil {
				c.Logger().Warnf("user attempted to access with current_user_id '%v' that is not found: %v", uid, err)

				c.Session().Delete("current_user_id")
				c.Session().Set("redirectURL", c.Request().URL.String())
				c.Flash().Add("danger", "You must be authorized with a correct user to see that page")
				return c.Redirect(http.StatusFound, "/auth/new")
			}
			c.Set("current_user", u)
		}
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(http.StatusFound, "/auth/new")
		}
		return next(c)
	}
}
