package core

// ResourceController is an interface for a controller
// that handles CRUD operations.
type ResourceController interface {
	// Index returns a list of resources.
	Index(ctx *CaesarCtx) error

	// Create creates a new resource.
	Create(ctx *CaesarCtx) error

	// Store creates a new resource.
	Store(ctx *CaesarCtx) error

	// Show returns a single resource.
	Show(ctx *CaesarCtx) error

	// Edit returns a form to edit a resource.
	Edit(ctx *CaesarCtx) error

	// Update updates a resource.
	Update(ctx *CaesarCtx) error

	// Delete deletes a resource.
	Delete(ctx *CaesarCtx) error
}

// BaseResourceController is a base controller for resources.
type BaseResourceController struct{}

func (c *BaseResourceController) Index(ctx *CaesarCtx) error {
	return nil
}

func (c *BaseResourceController) Create(ctx *CaesarCtx) error {
	return nil
}

func (c *BaseResourceController) Show(ctx *CaesarCtx) error {
	return nil
}

func (c *BaseResourceController) Store(ctx *CaesarCtx) error {
	return nil
}

func (c *BaseResourceController) Edit(ctx *CaesarCtx) error {
	return nil
}

func (c *BaseResourceController) Update(ctx *CaesarCtx) error {
	return nil
}

func (c *BaseResourceController) Delete(ctx *CaesarCtx) error {
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

	res.Routes[ResourceControllerMethodIndex] = r.Get(pattern+"/", controller.Index)
	res.Routes[ResourceControllerMethodCreate] = r.Get(pattern+"/create", controller.Create)
	res.Routes[ResourceControllerMethodShow] = r.Get(pattern+"/:id", controller.Show)
	res.Routes[ResourceControllerMethodEdit] = r.Get(pattern+"/:id/edit", controller.Edit)
	res.Routes[ResourceControllerMethodUpdate] = r.Put(pattern+"/:id", controller.Update)
	res.Routes[ResourceControllerMethodDelete] = r.Delete(pattern+"/:id", controller.Delete)

	return res
}

// ResourceControllerMethod is a type for the methods of a resource controller.
type ResourceControllerMethod string

const (
	ResourceControllerMethodIndex  ResourceControllerMethod = "Index"
	ResourceControllerMethodCreate ResourceControllerMethod = "Create"
	ResourceControllerMethodStore  ResourceControllerMethod = "Store"
	ResourceControllerMethodShow   ResourceControllerMethod = "Show"
	ResourceControllerMethodEdit   ResourceControllerMethod = "Edit"
	ResourceControllerMethodUpdate ResourceControllerMethod = "Update"
	ResourceControllerMethodDelete ResourceControllerMethod = "Delete"
)

// Exclude excludes methods from the resource.
func (res *Resource) Exclude(methods ...ResourceControllerMethod) {
	for _, method := range methods {
		delete(res.Routes, method)
	}
}
