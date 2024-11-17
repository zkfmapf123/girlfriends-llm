package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zkfmapf123/go-llm/chains"
)

type GenerateVacationIdeaRequest struct {
	FavoriteSeason string   `json:"favorite_season"`
	Hobbies        []string `json:"hobbies"`
	Budget         int      `json:"budget"`
}

type GenerateVacationIdeaResponse struct {
	Id        uuid.UUID `json:"id"`
	Completed bool      `json:"completed"`
}

type GetVacationIdeaResponse struct {
	Id        uuid.UUID `json:"id"`
	Completed bool      `json:"completed"`
	Idea      string    `json:"idea"`
}

func generateVacation(r GenerateVacationIdeaRequest) GenerateVacationIdeaResponse {
	// First, generate a new UUID for the idea
	id := uuid.New()

	// Then invoke the GeneateVacationIdeaChange method of the chains package
	// passing through all of the parameters from the user
	go chains.GeneateVacationIdeaChange(id, r.Budget, r.FavoriteSeason, r.Hobbies)
	return GenerateVacationIdeaResponse{Id: id, Completed: false}
}

func getVacation(id uuid.UUID) (GetVacationIdeaResponse, error) {
	// Search the chains database for the ID requested by the user
	v, err := chains.GetVacationFromDb(id)

	// If the ID didn't exist, handle the error
	if err != nil {
		return GetVacationIdeaResponse{}, err
	}

	// Otherwise, return the vacation idea to the caller
	return GetVacationIdeaResponse{Id: v.Id, Completed: v.Completed, Idea: v.Idea}, nil
}

func GetVacationRouter(router *gin.Engine) *gin.Engine {

	// Add a new router group to the gin router
	registrationRoutes := router.Group("/vacation")

	// Handle the POST to /create
	registrationRoutes.POST("/create", func(c *gin.Context) {
		var req GenerateVacationIdeaRequest
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Bad Request",
			})
		} else {
			c.JSON(http.StatusOK, generateVacation(req))
		}
	})

	// Handle the GET to /:id
	registrationRoutes.GET("/:id", func(c *gin.Context) {
		id, err := uuid.Parse(c.Param("id"))

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Bad Request",
			})
		} else {
			resp, err := getVacation(id)
			if err != nil {
				c.JSON(http.StatusNotFound, gin.H{
					"message": "Id Not Found",
				})
			} else {
				c.JSON(http.StatusOK, resp)
			}
		}
	})

	// Return the updated router
	return router
}
