package main

import (
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/molejs/mole-server/mole"
	"os"
	"time"
)

func getenv(k, defaultVal string) string {
	v := os.Getenv(k)
	if v == "" {
		v = defaultVal
	}

	return v
}

func main() {
	gin.SetMode(gin.ReleaseMode)

	var (
		r            = gin.Default()
		addr         = getenv("MOLE_ADDR", ":8080")
		cert         = os.Getenv("MOLE_CERT")
		key          = os.Getenv("MOLE_KEY")
		mongoAddr    = getenv("MOLE_MONGO_ADDR", "127.0.0.1:27017")
		dbName       = getenv("MOLE_DB_NAME", "mole")
		dbMiddleware = mole.DatabaseMiddleware(mongoAddr, dbName)
	)

	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "PUT, POST",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

	r.POST("/logs", dbMiddleware, mole.ReportHandler)
	r.GET("/logs", dbMiddleware, mole.RetrieveHandler)
	r.GET("/log/:id", dbMiddleware, mole.SingleLogHandler)

	if cert != "" && key != "" {
		r.RunTLS(addr, cert, key)
	} else {
		r.Run(addr)
	}
}
