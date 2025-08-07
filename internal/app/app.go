package app

import (
	"dev_profiler/internal/controllers"
)

// Run starts the application
func Run() {
	mainController := controllers.NewMainController()
	mainController.Run()
}
