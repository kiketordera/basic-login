package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/timshannon/bolthold"
)

// Start the web server for a Server object
func (db *DB) Start(wg *sync.WaitGroup) {
	// When the process is end, warning the group that is end
	defer wg.Done()

	// Declare a web server in the port declared
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: db.createRouter(),
	}
	// Initialice a web server already declare, in background due to GO instruction
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Graceful process shutdown from interrupts
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	// We wait here and do nothing until the server is told to shut down
	<-quit

	// Here we shut down the server SRV and give 5 seconds for it to shut down or it would be killed
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
}

// This method renders the HTML / JSON template and export the HTML / JSON to the browser
func render(c *gin.Context, data gin.H, templateName string) {
	switch c.Request.Header.Get("Accept") {
	case "application/json":
		// Respond with JSON
		c.JSON(http.StatusOK, data["payload"])
	default:
		// Respond with HTML
		c.HTML(http.StatusOK, templateName, data)
	}
}

// CreateRouter creates our router with the URLs
func (db *DB) createRouter() *gin.Engine {
	router := gin.Default()
	// Route for the static content (images, SVG...)
	router.Static("/ui", BasePath+"/ui")
	// Path to the HTML templates
	router.LoadHTMLGlob(BasePath + "/ui/html/*/*.html")
	// Redirects when there is wrong route
	router.NoRoute(redirect)

	// Checks user cookie for each request
	router.Use(checkToken(), gzip.Gzip(gzip.DefaultCompression))

	home := router.Group("/")
	{
		home.GET("/", ensureNotLoggedIn(), showLogin)
		home.POST("/", ensureNotLoggedIn(), db.login)
		home.GET("/users", ensureLoggedIn(), db.showUsers)
		home.GET("/log-out", ensureLoggedIn(), db.logout)
	}
	return router
}

// InitDatabase returns the DataBase of the project or creates a new DataBase if does not exists
func InitDatabase() DB {
	// We open or create the DataBase and the Directory
	db, err := bolthold.Open(BasePath+"/data/"+projectName+".db", 0666, nil)
	if err != nil {
		if os.Getenv("GOPATH") == "" {
			fmt.Println("Try writing the GOPATH with: ")
			fmt.Println("export GOPATH=$HOME/go")
		}
		os.MkdirAll(BasePath+"/data/", os.ModePerm)
		db, err = bolthold.Open(BasePath+"/data/"+projectName+".db", 0666, nil)
		if err != nil {
			fmt.Print("We try to make the directory and did not work, here is the error: ")
			fmt.Print(err)
		}
	}
	DataBase := DB{
		DataBase: db,
	}
	return DataBase
}
