#Foundation
## Components
Components are essential part of foundation which act as main functional part of the framework. Actually the components act as mini applications uses framework capabilities as its infrastructure to implement functionality of application.

The foundation could uses one or more components at same time depend on developer needs.

### Structure
Each component base on the application logic could consist several parts. 


   - [component.go](../components/dashboard/component.go) implements essential parts of component including views,paths,menus,permissions,installation/update/uninstall procedures and etc.
   - [view.go](../components/dashboard/component.go) Could implement views if required by the application logic. it will load html templates and render it by substitute variables as it is declared inside the view template the dispatch it to the user.
   - [models.go](../components/dashboard/component.go) Used to implement models.
   - [controller.go](../components/dashboard/component.go) Implements actions of component such as read from/write to database or compute things. 
   - assets Used to server static files and resources. assets should be introduced in component.go to be able to getting used.
   - views Used to serve views templates
   
File names are not mandatory but Foundation strongly suggest to use as mentioned for better reading.


#### Component.go
Each component should contain an struct that implements system.Component as below.
```go
type Component interface {
	Register()                         // Register the component
	Routers()                          // Initialize routers
	Menu()                             // Initialize the menus
	Install()                          // Install the component
	Uninstall()                        // Uninstall the component
	Update()                           // Update the component
	ComputeHash()                      // Compute hash of component for compare to updates
	ViewPath() string                  // return path of the view templates
	AssetsPath() string                // return path of the assets
	GetTemplates() *template.Template  // return templates
}
```

##### function Register()
```go
func (component *component) Register() {
	// put component as given name to list of registered components
	system.Components["dashboard"] = component
	
	// search for templates in views folder
	files, err := path.Dir(component.Views).Find("*.html")
	if err != nil {
		log.Critical(fmt.Errorf("unable to parse template html files %s", err.Error()))
	}

    // save list of templates somewhere to pass to GetTemplates()
	component.Templates, err = template.ParseFiles(files...)
	if err != nil {
		log.Critical(fmt.Errorf("unable to parse template html layouts. %s", err.Error()))
	}
}
```


##### function Routers()
Full version of [router doc](routers.md)
```go
func (component component) Routers() {
	//implements https://domain.xyz/hello
	system.Router.Match("hello","GET", func(req router.Request) {
		//call view of Hello from views
		View{}.Hello(req)
	})
	
	//Setup static route to component assets
	system.Router.Static(component.Assets,component.Assets,nil)

}

```


##### function Routers()
Full version of [menu doc](menu.md)
```go

func (component *component) Menu() {
	//Select Menu Based on template available positions
	adminlte.MainNav.Push(
		
		//Push root
		system.Menu{
			Name: "MainMenu", Title: "Home", Permission: "", URL: "dashboard", Icon: "fa-home", Class: "home",
			
			//Push sub menu
			SubMenu: []system.Menu{
				system.Menu{
					Name: "MainMenu", Title: "->1", Permission: "", URL: "dashboard", Icon: "fa-home", Class: "home",
				},
				system.Menu{
					Name: "MainMenu", Title: "->2", Permission: "", URL: "dashboard", Icon: "fa-home", Class: "home",
				},
				system.Menu{
					Name: "MainMenu", Title: "->3", Permission: "", URL: "dashboard", Icon: "fa-home", Class: "home",
				},
			}},
		system.Menu{
			Name: "MainMenu2", Title: "Home2", Permission: "", URL: "dashboard", Icon: "fa-home", Class: "home",
			SubMenu: []system.Menu{
				system.Menu{
					Name: "MainMenu", Title: "->1", Permission: "", URL: "dashboard", Icon: "fa-home", Class: "home",
				},
				system.Menu{
					Name: "MainMenu", Title: "->2", Permission: "", URL: "dashboard", Icon: "fa-home", Class: "home",
				},
				system.Menu{
					Name: "MainMenu", Title: "->3", Permission: "", URL: "dashboard", Icon: "fa-home", Class: "home",
				},
			}},
	)
}

```

##### function ViewPath() and AssetsPath()
```go
// return path of views, default: components/component/views
func (component component) ViewPath() string {
	return component.Views
}

// return path of assets, default: components/component/assets
func (component component) AssetsPath() string {
	return component.Assets
}
```

##### function GetTemplates()
```go
//return templates that is generated using template.ParseFiles. checkout Register()
func (component component) GetTemplates() *template.Template {
	return component.Templates
}
```