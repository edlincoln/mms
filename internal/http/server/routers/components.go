package routers

import "github.com/go-chi/chi"

type Router interface {
	Route(r chi.Router)
	GetPath() string
}

type RouterComponent struct {
	MmsPairController   Router
	ExtractorController Router
}

func NewRouterComponent() RouterComponent {
	return RouterComponent{
		MmsPairController:   GetInstance().MmsPairController(),
		ExtractorController: GetInstance().ExtractorController(),
	}
}
