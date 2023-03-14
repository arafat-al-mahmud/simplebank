package api

import (
	"log"

	db "github.com/arafat-al-mahmud/simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  *db.Store
	router *gin.Engine
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router

	return server
}

func (server *Server) Start(address string) error {
	err := server.router.Run(address)
	if err != nil {
		log.Fatal("Error running server ", err)
	}
	return err
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
