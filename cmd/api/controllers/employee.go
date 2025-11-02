package controllers

// import (
// 	"github.com/SaltaGet/NOA-GESTION-BACK/cmd/api/logging"
// 	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
// 	"github.com/gofiber/fiber/v2"
// )

// // GetEmployeeByID godoc
// //	@Summary		Get Employee By ID
// //	@Description	Get Employee By ID
// //	@Tags			Employee
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			id	path		string									true	"ID of Employee"
// //	@Success		200	{object}	schemas.Response{body=schemas.Employee}	"Employee obtained successfully"
// //	@Failure		400	{object}	schemas.Response							"Bad Request"
// //	@Failure		401	{object}	schemas.Response							"Auth is required"
// //	@Failure		403	{object}	schemas.Response							"Not Authorized"
// //	@Failure		404	{object}	schemas.Response							"Employee not found"
// //	@Failure		500	{object}	schemas.Response
// //	@Router			/employee/{id} [get]
// func (e *EmployeeController) GetEmployeeByID(c *fiber.Ctx) error {
// 	logging.INFO("Obtener un empleado por ID")
// 	id := c.Params("id")
// 	if id == "" {
// 		logging.ERROR("Error: ID is required")
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "ID is required",
// 		})
// 	}

// 	employee, err := e.EmployeeService.EmployeeGetByID(id)
// 	if err != nil {
// 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Empleado obtenido con éxito")
// 	return c.Status(200).JSON(schemas.Response{
// 		Status:  true,
// 		Body:    employee,
// 		Message: "Empleado obtenido con éxito",
// 	})
// }

// // GetAllEmployees godoc
// //	@Summary		Get all employees
// //	@Description	Fetches all employees from the specified tenant.
// //	@Tags			Employee
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Success		200	{object}	schemas.Response{body=[]schemas.Employee}	"List of employees"
// //	@Failure		400	{object}	schemas.Response							"Bad request"
// //	@Failure		401	{object}	schemas.Response							"Auth is required"
// //	@Failure		403	{object}	schemas.Response							"Not Authorized"
// //	@Failure		500	{object}	schemas.Response							"Internal server error"
// //	@Router			/employee/get_all [get]
// func (e *EmployeeController) GetAllEmployees(c *fiber.Ctx) error {
// 	logging.INFO("Obtener todos los empleados")
// 	employees, err := e.EmployeeService.EmployeeGetAll()
// 	if err != nil {
// 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Empleados obtenidos con éxito")
// 	return c.Status(200).JSON(schemas.Response{
// 		Status:  true,
// 		Body:    employees,
// 		Message: "Empleados obtenidos con éxito",
// 	})
// }

// // GetEmployeeByName godoc
// //	@Summary		Get Employee By Name
// //	@Description	Fetches employees from either laundry or workshop based on the provided name and workplace.
// //	@Tags			Employee
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			name	query		string									true	"Name of the Employee"
// //	@Success		200		{object}	schemas.Response{body=[]schemas.Employee}	"List of laundry employees"
// //	@Failure		400		{object}	schemas.Response							"Bad request"
// //	@Failure		401		{object}	schemas.Response							"Auth is required"
// //	@Failure		403		{object}	schemas.Response							"Not Authorized"
// //	@Failure		500		{object}	schemas.Response							"Internal server error"
// //	@Router			/employee/get_by_name [get]
// func (e *EmployeeController) GetEmployeeByName(c *fiber.Ctx) error {
// 	logging.INFO("Obtener un empleado por nombre")
// 	name := c.Query("name")
// 	if name == "" || len(name) < 3 {
// 		logging.ERROR("Error: El valor no debe de ser vacio o menor a 3 caracteres")
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "El valor no debe de ser vacio o menor a 3 caracteres",
// 		})
// 	}

// 	employees, err := e.EmployeeService.EmployeeGetByName(name)
// 	if err != nil {
// 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Empleados obtenidos con éxito")
// 	return c.Status(200).JSON(schemas.Response{
// 		Status:  true,
// 		Body:    employees,
// 		Message: "Empleados obtenidos con éxito",
// 	})
// }

// // CreateEmployee godoc
// //	@Summary		Create Employee
// //	@Description	Creates an employee for either laundry or workshop based on the provided information.
// //	@Tags			Employee
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			employeeCreate	body		schemas.EmployeeCreate			true	"Employee information"
// //	@Success		200				{object}	schemas.Response{body=string}	"Employee created"
// //	@Failure		400				{object}	schemas.Response					"Bad request"
// //	@Failure		401				{object}	schemas.Response					"Auth is required"
// //	@Failure		403				{object}	schemas.Response					"Not Authorized"
// //	@Failure		422				{object}	schemas.Response					"Model Invalid"
// //	@Failure		500				{object}	schemas.Response					"Internal server error"
// //	@Router			/employee/create [post]
// func (e *EmployeeController) CreateEmployee(c *fiber.Ctx) error {
// 	logging.INFO("Crear un empleado")
// 	var employeeCreate schemas.EmployeeCreate
// 	if err := c.BodyParser(&employeeCreate); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Invalid request" + err.Error(),
// 		})
// 	}
// 	if err := employeeCreate.Validate(); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: err.Error(),
// 		})
// 	}

// 	id, err := e.EmployeeService.EmployeeCreate(&employeeCreate)
// 	if err != nil {
// 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Empleado creado con éxito")
// 	return c.Status(200).JSON(schemas.Response{
// 		Status:  true,
// 		Body:    id,
// 		Message: "Empleado creado con éxito",
// 	})
// }

// // UpdateEmployee godoc
// //	@Summary		Update Employee
// //	@Description	Updates the details of an employee based on the provided data.
// //	@Tags			Employee
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			employeeUpdate	body		schemas.EmployeeUpdate	true	"Employee data to update"
// //	@Success		200				{object}	schemas.Response			"Empleado editado con éxito"
// //	@Failure		400				{object}	schemas.Response			"Invalid request or Workplace is required"
// //	@Failure		401				{object}	schemas.Response			"Auth is required"
// //	@Failure		403				{object}	schemas.Response			"Not Authorized"
// //	@Failure		404				{object}	schemas.Response			"Not Found"
// //	@Failure		422				{object}	schemas.Response			"Model Invalid"
// //	@Failure		500				{object}	schemas.Response			"Error interno"
// //	@Router			/employee/update [put]
// func (e *EmployeeController) UpdateEmployee(c *fiber.Ctx) error {
// 	logging.INFO("Actualizar un empleado")
// 	var employeeUpdate schemas.EmployeeUpdate
// 	if err := c.BodyParser(&employeeUpdate); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Invalid request" + err.Error(),
// 		})
// 	}
// 	if err := employeeUpdate.Validate(); err != nil {
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: err.Error(),
// 		})
// 	}

// 	err := e.EmployeeService.EmployeeUpdate(&employeeUpdate)
// 	if err != nil {
// 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Empleado editado con éxito")
// 	return c.Status(200).JSON(schemas.Response{
// 		Status:  true,
// 		Body:    nil,
// 		Message: "Empleado editado con éxito",
// 	})
// }

// // DeleteEmployee godoc
// //	@Summary		Delete Employee
// //	@Description	Removes an employee from the database based on the provided ID and tenant context.
// //	@Tags			Employee
// //	@Accept			json
// //	@Produce		json
// //	@Security		BearerAuth
// //	@Param			id	path		string			true	"ID of the employee"
// //	@Success		200	{object}	schemas.Response	"Empleado eliminado con éxito"
// //	@Failure		400	{object}	schemas.Response	"Bad Request"
// //	@Failure		401	{object}	schemas.Response	"Auth is required"
// //	@Failure		403	{object}	schemas.Response	"Not Authorized"
// //	@Failure		404	{object}	schemas.Response	"Not Found"
// //	@Failure		500	{object}	schemas.Response	"Error interno"
// //	@Router			/employee/delete/{id} [delete]
// func (e *EmployeeController) DeleteEmployee(c *fiber.Ctx) error {
// 	logging.INFO("Eliminar un empleado")
// 	id := c.Params("id")
// 	if id == "" {
// 		logging.ERROR("Error: ID is required")
// 		return c.Status(fiber.StatusBadRequest).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "ID is required",
// 		})
// 	}

// 	err := e.EmployeeService.EmployeeDelete(id)
// 	if err != nil {
// 		if errResp, ok := err.(*schemas.ErrorStruc); ok {
// 			logging.ERROR("Error: %s", errResp.Err.Error())
// 			return c.Status(errResp.StatusCode).JSON(schemas.Response{
// 				Status:  false,
// 				Body:    nil,
// 				Message: errResp.Message,
// 			})
// 		}
// 		logging.ERROR("Error: %s", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(schemas.Response{
// 			Status:  false,
// 			Body:    nil,
// 			Message: "Error interno",
// 		})
// 	}

// 	logging.INFO("Empleado eliminado con éxito")
// 	return c.Status(200).JSON(schemas.Response{
// 		Status:  true,
// 		Body:    nil,
// 		Message: "Empleado eliminado con éxito",
// 	})
// }

