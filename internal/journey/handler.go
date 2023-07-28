package journey

import (
	"github.com/dportaluppi/journey-api/pkg/journey"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	getter  journey.Getter
	creator journey.Creator
	updater journey.Updater
	deleter journey.Deleter
}

func NewHandler(getter journey.Getter, creator journey.Creator, updater journey.Updater, deleter journey.Deleter) *Handler {
	return &Handler{getter: getter, creator: creator, updater: updater, deleter: deleter}
}

func (h *Handler) GetJourneys(c *gin.Context) {
	storefront := c.Query("storefront")
	audiencesQuery := c.QueryArray("audiences")
	channelsQuery := c.QueryArray("channels")
	date, _ := time.Parse(time.RFC3339, c.Query("date"))
	sort := c.Query("sort")

	var audiences []string
	for _, a := range audiencesQuery {
		audiences = append(audiences, strings.Split(a, ",")...)
	}

	var channels []string
	for _, ch := range channelsQuery {
		channels = append(channels, strings.Split(ch, ",")...)
	}

	journeys, err := h.getter.GetJourneysByCriteria(c.Request.Context(), storefront, audiences, channels, date, sort)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, journeys)
}

func (h *Handler) GetJourneyByID(c *gin.Context) {
	id := c.Param("id")

	j, err := h.getter.GetJourneyByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, j)
}

func (h *Handler) CreateJourney(c *gin.Context) {
	var j journey.Journey
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	journeyCreated, err := h.creator.CreateJourney(c.Request.Context(), &j)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, journeyCreated)
}

func (h *Handler) UpdateJourney(c *gin.Context) {
	id := c.Param("id")
	var j journey.Journey
	if err := c.ShouldBindJSON(&j); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	err := h.updater.UpdateJourney(c.Request.Context(), id, &j)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Journey updated successfully"})
}

func (h *Handler) DeleteJourney(c *gin.Context) {
	id := c.Param("id")

	err := h.deleter.DeleteJourney(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Journey deleted successfully"})
}
