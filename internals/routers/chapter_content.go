package router

type SaveContentPayload struct {
	ChapterID int    `validate:"required,gte=1"`
	Content   string `validate:"required"`
}

/* func (r Router) CreateContent(c echo.Context) error {
    validate := utils.NewValidator()
    payload := new(SaveContentPayload)
    idStr := c.Param("ChapterID")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        return r.badRequestError(err)
    }
    userId, err := r.JwtService.RetreiveUserIdFromContext(c)
    if err != nil {
        return r.forbiddenError(err)
    }
} */
