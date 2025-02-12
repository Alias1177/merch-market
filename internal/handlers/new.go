package handlers

type Handler struct {
	userU UserUsecase
}

func New(userU UserUsecase) *Handler {
	return &Handler{
		userU: userU,
	}
}
