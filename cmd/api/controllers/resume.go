package controllers

// import (
// 	"github.com/gofiber/fiber/v2"
// )

// func (r *ResumeController) ExpenseResumeCreate(c *fiber.Ctx) error {
// 	logging.INFO("Crear Resumen de gastos")
// 	resume := &models.ResumeExpenseCreate{}
// 	if err := c.BodyParser(resume); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Bad Request",
// 		})
// 	}
// 	if err := resume.Validate(); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: err.Error(),
// 		})
// 	}

// 	id, err := r.ResumeExpenseService.ResumeExpenseCreate(resume)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Resumen creado exitosamente")
// 	return c.Status(fiber.StatusOK).JSON(models.Response{
// 		Status:  true,
// 		Body:    id,
// 		Message: "Resumen creado exitosamente",
// 	})
// }

// func (r *ResumeController) ExpenseResumeGetByDateBetween(c *fiber.Ctx) error {
// 	logging.INFO("Obtener resumen de gastos")
// 	fromDate := c.Query("fromDate")
// 	toDate := c.Query("toDate")
// 	if fromDate == "" || toDate == "" {
// 		logging.ERROR("Error: fromDate y toDate son requiridos")
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "fromDate and toDate are required",
// 		})
// 	}

// 	resumes, err := r.ResumeExpenseService.ResumeExpenseGetByDateBetween(fromDate, toDate)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Resumen obtenido exitosamente")
// 	return c.Status(fiber.StatusOK).JSON(models.Response{
// 		Status:  true,
// 			Body:    resumes,
// 			Message: "Resumen obtenido exitosamente",
// 	})
// }

// func (r *ResumeController) ExpenseResumeGetByID(c *fiber.Ctx) error {
// 	logging.INFO("Obtener resumen de gastos")
// 	id := c.Params("id")
// 	if id == "" {
// 		logging.ERROR("Error: id is required")
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "id is required",
// 		})
// 	}

// 	resume, err := r.ResumeExpenseService.ResumeExpenseGetByID(id)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Resumen obtenido exitosamente")
// 	return c.Status(fiber.StatusOK).JSON(models.Response{
// 		Status:  true,
// 			Body:    resume,
// 			Message: "Resumen obtenido exitosamente",
// 	})
// }

// func (r *ResumeController) ExpenseResumeUpdate(c *fiber.Ctx) error {
// 	logging.INFO("Actualizar resumen de gastos")
// 	resume := &models.ResumeExpenseUpdate{}
// 	if err := c.BodyParser(resume); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Bad Request",
// 		})
// 	}
// 	if err := resume.Validate(); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: err.Error(),
// 		})
// 	}

// 	err := r.ResumeExpenseService.ResumeExpenseUpdate(resume)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Resumen actualizado exitosamente")
// 	return c.Status(fiber.StatusOK).JSON(models.Response{
// 		Status:  true,
// 			Body:    nil,
// 			Message: "Resumen actualizado exitosamente",
// 	})
// }


// // INCOME

// func (r *ResumeController) IncomeResumeCreate(c *fiber.Ctx) error {
// 	logging.INFO("Crear resumen de ingresos")
// 	resume := &models.ResumeIncomeCreate{}
// 	if err := c.BodyParser(resume); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Bad Request",
// 		})
// 	}
// 	if err := resume.Validate(); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: err.Error(),
// 		})
// 	}

// 	id, err := r.ResumeIncomeService.ResumeIncomeCreate(resume)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Resumen creado exitosamente")
// 	return c.Status(fiber.StatusOK).JSON(models.Response{
// 		Status:  true,
// 		Body:    id,
// 		Message: "Resumen creado exitosamente",
// 	})
// }

// func (r *ResumeController) IncomeResumeGetByDateBetween(c *fiber.Ctx) error {
// 	logging.INFO("Obtener resumen de ingresos")
// 	fromDate := c.Query("from_date")
// 	toDate := c.Query("to_date")
// 	if fromDate == "" || toDate == "" {
// 		logging.ERROR("Error: from_date y to_date son requiridos")
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "from_date and to_date are required",
// 		})
// 	}

// 	resume, err := r.ResumeIncomeService.ResumeIncomeGetByDateBetween(fromDate, toDate)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Resumen obtenido exitosamente")
// 	return c.Status(fiber.StatusOK).JSON(models.Response{
// 		Status:  true,
// 			Body:    resume,
// 			Message: "Resumen obtenido exitosamente",
// 	})
// }

// func (r *ResumeController) IncomeResumeGetByID(c *fiber.Ctx) error {
// 	logging.INFO("Obtener resumen de ingresos")
// 	id := c.Params("id")
// 	if id == "" {
// 		logging.ERROR("Error: id is required")
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "id is required",
// 		})
// 	}

// 	resume, err := r.ResumeIncomeService.ResumeIncomeGetByID(id)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Resumen obtenido exitosamente")
// 	return c.Status(fiber.StatusOK).JSON(models.Response{
// 		Status:  true,
// 			Body:    resume,
// 			Message: "Resumen obtenido exitosamente",
// 	})
// }

// func (r *ResumeController) IncomeResumeUpdate(c *fiber.Ctx) error {
// 	logging.INFO("Actualizar resumen de ingresos")
// 	resume := &models.ResumeIncomeUpdate{}
// 	if err := c.BodyParser(resume); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Bad Request",
// 		})
// 	}
// 	if err := resume.Validate(); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: err.Error(),
// 		})
// 	}

// 	err := r.ResumeIncomeService.ResumeIncomeUpdate(resume)
// 	if err != nil {
// 		if errResp, ok := err.(*models.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(models.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Resumen actualizado exitosamente")
// 	return c.Status(fiber.StatusOK).JSON(models.Response{
// 		Status:  true,
// 			Body:    nil,
// 			Message: "Resumen actualizado exitosamente",
// 	})
// }
