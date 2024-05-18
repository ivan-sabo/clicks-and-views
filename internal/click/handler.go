package click

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// ClickDTO represents HTTP request/response model.
type ClickDTO struct {
	ID        uint   `json:"id,omitempty"`
	URL       string `json:"url" validate:"required,url"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// ToDomain maps DTO model into domain model.
func (c ClickDTO) ToDomain() Click {
	return Click{
		URL: c.URL,
	}
}

// ClickDTOCollection represents ClickDTO collection.
type ClickDTOCollection []ClickDTO

// NewClickDTOCollection maps domain models into DTO models.
func NewClickDTOCollection(clickCollection ClickCollection) ClickDTOCollection {
	clickDTOCollection := make(ClickDTOCollection, 0, len(clickCollection))

	for _, click := range clickCollection {
		clickDTOCollection = append(clickDTOCollection, NewClickDTO(click))
	}

	return clickDTOCollection
}

// FilterDTO represents HTTP request model.
type FilterDTO struct {
	URL    string    `query:"url"`
	Before time.Time `query:"before"`
	After  time.Time `query:"after"`
}

// ToDomain maps DTO model into domain model.
func (f *FilterDTO) ToDomain() Filter {
	return Filter{
		URL:    f.URL,
		Before: f.Before,
		After:  f.After,
	}
}

// Handler defines all API methods for Click.
type Handler struct {
	clickRepository Repository
}

// Create implements handler for Create Click HTTP request.
func (h *Handler) Create(c echo.Context) error {
	var clickDTO ClickDTO
	if err := c.Bind(&clickDTO); err != nil {
		return err
	}

	click, err := h.clickRepository.Create(c.Request().Context(), clickDTO.ToDomain())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, NewClickDTO(click))
}

// Filter implements handler for Filter Click HTTP request.
func (h *Handler) Filter(c echo.Context) error {
	var filterDTO FilterDTO
	if err := c.Bind(&filterDTO); err != nil {
		return err
	}

	clickCollection, err := h.clickRepository.Filter(c.Request().Context(), filterDTO.ToDomain())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, NewClickDTOCollection(clickCollection))
}

// NewHandler is a Handler constructor.
func NewHandler(clickRepository Repository) Handler {
	return Handler{
		clickRepository: clickRepository,
	}
}

// NewClickDTO is a ClickDTO constructor.
func NewClickDTO(c Click) ClickDTO {
	return ClickDTO{
		ID:        c.ID,
		URL:       c.URL,
		CreatedAt: c.CreatedAt.Format(time.DateTime),
	}
}
