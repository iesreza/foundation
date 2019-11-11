package system

var onMenuRenderCallbacks = []func(m *Menu){}

func OnMenuRender(f func(m *Menu))  {
	onMenuRenderCallbacks = append(onMenuRenderCallbacks,f)
}
