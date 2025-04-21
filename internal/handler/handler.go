package handler

type Handler struct {
	userService      UserService
	pvzService       PvzService
	receptionService ReceptionService
}

func New(userService UserService, pvzService PvzService, receptionService ReceptionService) *Handler {
	return &Handler{
		userService:      userService,
		pvzService:       pvzService,
		receptionService: receptionService,
	}
}
