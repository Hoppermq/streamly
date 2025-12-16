package organization

type Handler struct {
}

type HandlerOption func(h *Handler)

func (h *Handler) CreateOrganization() (int, error) {
	return 0, nil
}
