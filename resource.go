package core

// ResourceController is an interface for a controller
// that handles CRUD operations.
type ResourceController interface {
	// Index returns a list of resources.
	Index(ctx *Context) error

	// Create creates a new resource.
	Create(ctx *Context) error

	// Store creates a new resource.
	Store(ctx *Context) error

	// Show returns a single resource.
	Show(ctx *Context) error

	// Edit returns a form to edit a resource.
	Edit(ctx *Context) error

	// Update updates a resource.
	Update(ctx *Context) error

	// Delete deletes a resource.
	Delete(ctx *Context) error
}

// BaseResourceController is a base controller for resources.
type BaseResourceController struct{}

func (c *BaseResourceController) Index(ctx *Context) error {
	return nil
}

func (c *BaseResourceController) Create(ctx *Context) error {
	return nil
}

func (c *BaseResourceController) Show(ctx *Context) error {
	return nil
}

func (c *BaseResourceController) Store(ctx *Context) error {
	return nil
}

func (c *BaseResourceController) Edit(ctx *Context) error {
	return nil
}

func (c *BaseResourceController) Update(ctx *Context) error {
	return nil
}

func (c *BaseResourceController) Delete(ctx *Context) error {
	return nil
}

// Resource is the struct that holds the routes for a resource controller.
type Resource struct {
	Routes map[ResourceControllerMethod]*Route
}

// Resource adds a set of routes to the router for a resource controller.
func (r *Router) Resource(pattern string, controller ResourceController) *Resource {
	res := &Resource{
		Routes: make(map[ResourceControllerMethod]*Route),
	}

	res.Routes[MethodIndex] = r.Get(pattern+"/", controller.Index)
	res.Routes[MethodCreate] = r.Get(pattern+"/create", controller.Create)
	res.Routes[MethodShow] = r.Get(pattern+"/:id", controller.Show)
	res.Routes[MethodEdit] = r.Get(pattern+"/:id/edit", controller.Edit)
	res.Routes[MethodUpdate] = r.Put(pattern+"/:id", controller.Update)
	res.Routes[MethodDelete] = r.Delete(pattern+"/:id", controller.Delete)

	return res
}

// ResourceControllerMethod is a type for the methods of a resource controller.
type ResourceControllerMethod string

const (
	MethodIndex  ResourceControllerMethod = "Index"
	MethodCreate ResourceControllerMethod = "Create"
	MethodStore  ResourceControllerMethod = "Store"
	MethodShow   ResourceControllerMethod = "Show"
	MethodEdit   ResourceControllerMethod = "Edit"
	MethodUpdate ResourceControllerMethod = "Update"
	MethodDelete ResourceControllerMethod = "Delete"
)

// Exclude excludes methods from the resource.
func (res *Resource) Exclude(methods ...ResourceControllerMethod) *Resource {
	for _, method := range methods {
		delete(res.Routes, method)
	}
	return res
}

// Use adds middleware to the resource.
func (res *Resource) Use(middleware ...Handler) *Resource {
	for _, route := range res.Routes {
		for _, m := range middleware {
			route.Use(m)
		}
	}
	return res
}
