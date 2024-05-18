package view

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// ViewDTO represents HTTP request/response model.
type ViewDTO struct {
	ID        uint   `json:"id,omitempty"`
	URL       string `json:"url" validate:"required,url"`
	CreatedAt string `json:"createdAt,omitempty"`
}

// ToDomain maps DTO model into domain model.
func (c ViewDTO) ToDomain() View {
	return View{
		URL: c.URL,
	}
}

// NewViewDTO is a ClickDTO constructor.
func NewViewDTO(c View) ViewDTO {
	return ViewDTO{
		ID:        c.ID,
		URL:       c.URL,
		CreatedAt: c.CreatedAt.Format(time.DateTime),
	}
}

// ViewDTOCollection represents ViewDTO collection.
type ViewDTOCollection []ViewDTO

// NewViewDTOCollection maps domain models into DTO models.
func NewViewDTOCollection(viewCollection ViewCollection) ViewDTOCollection {
	viewDTOCollection := make(ViewDTOCollection, 0, len(viewCollection))

	for _, view := range viewCollection {
		viewDTOCollection = append(viewDTOCollection, NewViewDTO(view))
	}

	return viewDTOCollection
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

// Handler defines all API methods for View.
type Handler struct {
	viewRepository Repository
}

// Create implements handler for Create View HTTP request.
func (h *Handler) Create(c echo.Context) error {
	var viewDTO ViewDTO
	if err := c.Bind(&viewDTO); err != nil {
		return err
	}

	click, err := h.viewRepository.Create(c.Request().Context(), viewDTO.ToDomain())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, NewViewDTO(click))
}

// Filter implements handler for Filter View HTTP request.
func (h *Handler) Filter(c echo.Context) error {
	var filterDTO FilterDTO
	if err := c.Bind(&filterDTO); err != nil {
		return err
	}

	viewCollection, err := h.viewRepository.Filter(c.Request().Context(), filterDTO.ToDomain())
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, NewViewDTOCollection(viewCollection))
}

// NewHandler is a Handler constructor.
func NewHandler(viewRepository Repository) Handler {
	return Handler{
		viewRepository: viewRepository,
	}
}
