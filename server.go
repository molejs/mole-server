package main

import (
	"github.com/gin-gonic/gin"
	"os"
"./mole"
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

	r.POST("/logs", dbMiddleware, mole.ReportHandler)
	r.GET("/logs", dbMiddleware, mole.RetrieveHandler)

	if cert != "" && key != "" {
		r.RunTLS(addr, cert, key)
	} else {
		r.Run(addr)
	}
}
