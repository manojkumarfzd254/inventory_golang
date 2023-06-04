package actions

import (
	"fmt"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v6"
	"github.com/gobuffalo/x/responder"

	"library/models"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Category)
// DB Table: Plural (categories)
// Resource: Plural (Categories)
// Path: Plural (/categories)
// View Template Folder: Plural (/templates/categories/)

// CategoriesResource is the resource for the Category model
type CategoriesResource struct {
	buffalo.Resource
}

// List gets all Categories. This function is mapped to the path
// GET /categories
func (v CategoriesResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	categories := &models.Categories{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Categories from the DB
	if err := q.All(categories); err != nil {
		return err
	}
	// jsonCategories := categories
	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		c.Set("pagination", q.Paginator)

		c.Set("categories", categories)
		return c.Render(http.StatusOK, r2.HTML("backend/categories/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		// for _, category := range jsonCategories {
		// 	// Access category properties and perform operations
		// 	fmt.Println(category.ID)
		// 	fmt.Println(category.category_name)
		// 	// ...
		// }
		return c.Render(200, r2.JSON(categories))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r2.XML(categories))
	}).Respond(c)
}

// Show gets the data for one Category. This function is mapped to
// the path GET /categories/{category_id}
func (v CategoriesResource) Show(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Category
	category := &models.Category{}

	// To find the Category the parameter category_id is used.
	if err := tx.Find(category, c.Param("category_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("category", category)

		return c.Render(http.StatusOK, r2.HTML("backend/categories/show.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r2.JSON(category))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r2.XML(category))
	}).Respond(c)
}

// New renders the form for creating a new Category.
// This function is mapped to the path GET /categories/new
func (v CategoriesResource) New(c buffalo.Context) error {
	c.Set("category", &models.Category{})

	return c.Render(http.StatusOK, r2.HTML("backend/categories/new.plush.html"))
}

// Create adds a Category to the DB. This function is mapped to the
// path POST /categories
func (v CategoriesResource) Create(c buffalo.Context) error {
	// Allocate an empty Category
	category := &models.Category{}

	// Bind category to the html form elements
	if err := c.Bind(category); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(category)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			c.Set("category", category)

			return c.Render(http.StatusUnprocessableEntity, r2.HTML("backend/categories/new.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r2.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r2.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "category.created.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/auth/categories/%v", category.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r2.JSON(category))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r2.XML(category))
	}).Respond(c)
}

// Edit renders a edit form for a Category. This function is
// mapped to the path GET /categories/{category_id}/edit
func (v CategoriesResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Category
	category := &models.Category{}

	if err := tx.Find(category, c.Param("category_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("category", category)
	return c.Render(http.StatusOK, r2.HTML("backend/categories/edit.plush.html"))
}

// Update changes a Category in the DB. This function is mapped to
// the path PUT /categories/{category_id}
func (v CategoriesResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Category
	category := &models.Category{}

	if err := tx.Find(category, c.Param("category_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Bind Category to the html form elements
	if err := c.Bind(category); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(category)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("category", category)

			return c.Render(http.StatusUnprocessableEntity, r2.HTML("backend/categories/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r2.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r2.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "category.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/auth/categories/%v", category.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r2.JSON(category))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r2.XML(category))
	}).Respond(c)
}

// Destroy deletes a Category from the DB. This function is mapped
// to the path DELETE /categories/{category_id}
func (v CategoriesResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Category
	category := &models.Category{}

	// To find the Category the parameter category_id is used.
	if err := tx.Find(category, c.Param("category_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := tx.Destroy(category); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "category.destroyed.success"))

		// Redirect to the index spage
		return c.Redirect(http.StatusSeeOther, "/auth/categories")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r2.JSON(category))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r2.XML(category))
	}).Respond(c)
}
