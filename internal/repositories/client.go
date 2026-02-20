package repositories

import (
	"fmt"
	"math"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/database"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (r *ClientRepository) ClientGetByID(id int64) (*schemas.ClientResponse, error) {
	var client models.Client
	if err := r.DB.
		Preload("MemberCreate", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "first_name", "last_name", "username").Unscoped()
		}).
		Preload("Pay", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "income_sale_id", "client_id", "total", "method_pay", "created_at").Where("method_pay = ?", "credit")
		}).
		Where("id = ?", id).First(&client).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Cliente", schemas.Read)
	}

	var clientResponse schemas.ClientResponse
	copier.Copy(&clientResponse, &client)

	return &clientResponse, nil
}

func (r *ClientRepository) ClientGetByFilter(search string) (*[]schemas.ClientResponseDTO, error) {
	var client []models.Client
	if err := r.DB.Limit(10).Where("last_name LIKE ? OR first_name LIKE ? OR identifier LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").Find(&client).Error; err != nil {
		return nil, schemas.HandlerErrorGorm(err, "Cliente", schemas.Read)
	}

	var clientResponse []schemas.ClientResponseDTO
	copier.Copy(&clientResponse, &client)

	return &clientResponse, nil
}

func (r *ClientRepository) ClientGetAll(limit, page int64, search *map[string]string, filterDrbt bool) (*[]schemas.ClientResponseDTO, int64, error) {
	var clients []schemas.ClientResponseDTO

	// Base query con join opcional si hay que calcular deuda
	query := r.DB.Table("clients c")

	if filterDrbt {
		// query = query.
		// 	Select(`
		// 		c.id, c.first_name, c.last_name, c.company_name, c.identifier, c.responsability_front_iva,
		// 		c.email, c.phone, c.address, c.member_create_id, c.created_at, c.updated_at,
		// 		COALESCE(SUM(CASE WHEN p.method_pay = 'credit' THEN p.total ELSE 0 END), 0) AS debt
		// 	`).
		// 	Joins("LEFT JOIN pay_incomes p ON p.client_id = c.id").
		// 	Group("c.id").
		// 	Having("debt > 0")
		// Definimos la fórmula de la deuda para reutilizarla
    debtFormula := "COALESCE(SUM(CASE WHEN p.method_pay = 'credit' THEN p.total ELSE 0 END), 0)"
    
    query = query.
      Select(fmt.Sprintf(`
        c.id, c.first_name, c.last_name, c.company_name, c.identifier, c.responsability_front_iva,
        c.email, c.phone, c.address, c.member_create_id, c.created_at, c.updated_at,
        %s AS debt
      `, debtFormula)).
      Joins("LEFT JOIN pay_incomes p ON p.client_id = c.id").
      Group("c.id").
      // REGLA POSTGRES: Usar la función agregada, no el alias
      Having(fmt.Sprintf("%s > 0", debtFormula))
	} else {
		query = query.Select("c.*")
	}

	// Aplicar filtros dinámicos
	if search != nil {
		for key, value := range *search {
			if value == "" {
				continue
			}

			switch strings.ToLower(key) {
			case "identifier":
				query = query.Where("c.identifier LIKE ?", "%"+value+"%")
			case "first_name":
				query = query.Where("c.first_name LIKE ?", "%"+value+"%")
			case "last_name":
				query = query.Where("c.last_name LIKE ?", "%"+value+"%")
			case "email":
				query = query.Where("c.email LIKE ?", "%"+value+"%")
			}
		}
	}

	// Paginación
	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(int(limit)).Offset(int(offset))
	}

	// Ejecutar consulta
	if err := query.Scan(&clients).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Cliente", schemas.Read)
	}

	// Contar total
	var total int64
	countQuery := r.DB.Model(&models.Client{})

	// Si solo queremos los deudores, aplicar el mismo HAVING al count
	if filterDrbt {
		countQuery = countQuery.
			Select("clients.id").
			Joins("LEFT JOIN pay_incomes p ON p.client_id = clients.id").
			Group("clients.id").
			Having("SUM(CASE WHEN p.method_pay = 'credit' THEN p.total ELSE 0 END) > 0")
	}

	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, schemas.HandlerErrorGorm(err, "Cliente", schemas.Read)
	}

	return &clients, total, nil
}

func (r *ClientRepository) ClientCreate(memberID int64, client *schemas.ClientCreate) (int64, error) {
	newClient := models.Client{
		FirstName:      client.FirstName,
		LastName:       client.LastName,
		CompanyName:    client.CompanyName,
		Identifier:     client.Identifier,
		Email:          client.Email,
		Phone:          client.Phone,
		Address:        client.Address,
		MemberCreateID: memberID,
		ResponsabilityFrontIVA: client.ResponsabilityFrontIVA,
	}
	if err := r.DB.Create(&newClient).Error; err != nil {
		return 0, schemas.HandlerErrorGorm(err, "Cliente", schemas.Create)
	}

	go database.SaveAuditAsync(r.DB, models.AuditLog{
		MemberID: memberID,
		Method:   "create",
		Path:     "client",
	}, nil, newClient)

	return newClient.ID, nil
}

func (r *ClientRepository) ClientUpdate(memberID int64, client *schemas.ClientUpdate) error {
	// obtener estado anterior
	var oldClient models.Client
	if err := r.DB.First(&oldClient, client.ID).Error; err != nil {
		return schemas.HandlerErrorGorm(err, "Cliente", schemas.Read)
	}

	saveClient := oldClient

	oldClient.FirstName = client.FirstName
	oldClient.LastName = client.LastName
	oldClient.CompanyName = client.CompanyName
	oldClient.Identifier = client.Identifier
	oldClient.Email = client.Email
	oldClient.Phone = client.Phone
	oldClient.Address = client.Address
	oldClient.ResponsabilityFrontIVA = client.ResponsabilityFrontIVA

	// ejecutar actualización
	res := r.DB.Model(&models.Client{}).Where("id = ?", client.ID).Save(&oldClient)

	if res.Error != nil {
		return schemas.HandlerErrorGorm(res.Error, "Cliente", schemas.Update)
	}

	if res.RowsAffected == 0 {
		return schemas.ErrorResponse(404, "Cliente no encontrado", fmt.Errorf("cliente no encontrado"))
	}

	go database.SaveAuditAsync(r.DB, models.AuditLog{
		MemberID: memberID,
		Method:   "update",
		Path:     "client",
	}, saveClient, oldClient)

	return nil
}

func (r *ClientRepository) ClientDelete(memberID, id int64) error {
	var client models.Client
	if err := r.DB.First(&client, id).Error; err != nil {
		return schemas.HandlerErrorGorm(err, "Cliente", schemas.Read)
	}

	if err := r.DB.Delete(&client).Error; err != nil {
		return schemas.HandlerErrorGorm(err, "Cliente", schemas.Delete)
	}

	go database.SaveAuditAsync(r.DB, models.AuditLog{
		MemberID: memberID,
		Method:   "delete",
		Path:     "client",
	}, client, nil)

	return nil
}

// func (r *ClientRepository) ClientUpdateCredit(pointSaleID int64, clientUpdateCredit *schemas.ClientUpdateCredit) error {
// 	err := r.DB.Transaction(func(tx *gorm.DB) error {
// 		var register models.CashRegister
// 		if err := tx.
// 			Select("id").
// 			Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
// 			Order("hour_open DESC").
// 			First(&register).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return schemas.ErrorResponse(400, "No hay caja abierta para este punto de venta", err)
// 			}
// 			return schemas.ErrorResponse(500, "Error al obtener la apertura de caja", err)
// 		}

// 		var client models.Client
// 		if err := tx.Select("id").First(&client, clientUpdateCredit.ID).Error; err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return schemas.ErrorResponse(404, "Cliente no encontrado", err)
// 			}
// 			return schemas.ErrorResponse(500, "Error al obtener el cliente", err)
// 		}

// 		total := 0.0
// 		for _, p := range clientUpdateCredit.PayCredit {
// 			var payCredit models.PayIncome
// 			if err := tx.Where("client_id = ?", client.ID).First(&payCredit, p.CreditID).Error; err != nil {
// 				if errors.Is(err, gorm.ErrRecordNotFound) {
// 					return schemas.ErrorResponse(404, "Credito no encontrado", err)
// 				}
// 				return schemas.ErrorResponse(500, "Error al obtener el credito", err)
// 			}

// 			payCredit.MethodPay = p.MethodPay
// 			payCredit.CashRegisterID = &register.ID

// 			if err := tx.Save(&payCredit).Error; err != nil {
// 				return schemas.ErrorResponse(500, "Error al actualizar el credito", err)
// 			}

// 			total += payCredit.Total
// 		}

// 		if math.Abs(total-clientUpdateCredit.Total) > 1 {
// 			message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total del cliente (%.2f) no  puede ser mayor que 1", total, clientUpdateCredit.Total)
// 			return schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
// 		}

// 		return nil
// 	})

// 	return err
// }

func (r *ClientRepository) ClientUpdateCredit(memberID, pointSaleID int64, clientUpdateCredit *schemas.ClientUpdateCredit) error {
	var oldCredits []models.PayIncome
	var newCredits []models.PayIncome
	err := r.DB.Transaction(func(tx *gorm.DB) error {

		// Obtener caja abierta
		var register models.CashRegister
		if err := tx.
			Select("id").
			Where("is_close = ? AND point_sale_id = ?", false, pointSaleID).
			Order("hour_open DESC").
			First(&register).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Caja", schemas.Read)
		}

		// Obtener cliente
		var client models.Client
		if err := tx.Select("id").First(&client, clientUpdateCredit.ID).Error; err != nil {
			return schemas.HandlerErrorGorm(err, "Cliente", schemas.Read)
		}

		total := 0.0

		for _, p := range clientUpdateCredit.PayCredit {

			// Obtener crédito
			var payCredit models.PayIncome
			if err := tx.Where("client_id = ?", client.ID).First(&payCredit, p.CreditID).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Credito", schemas.Read)
			}

			// Copia para auditoría (estado ANTERIOR)
			oldCredits = append(oldCredits, payCredit)

			// Aplicar actualizaciones
			payCredit.MethodPay = p.MethodPay
			payCredit.CashRegisterID = &register.ID

			if err := tx.Save(&payCredit).Error; err != nil {
				return schemas.HandlerErrorGorm(err, "Credito", schemas.Update)
			}

			// Copia para auditoría (estado NUEVO)
			newCredits = append(newCredits, payCredit)

			total += payCredit.Total
		}

		// Validación negocio
		if math.Abs(total-clientUpdateCredit.Total) > 1 {
			message := fmt.Sprintf("la diferencia entre la suma de pagos (%.2f) y el total del cliente (%.2f) no puede ser mayor que 1", total, clientUpdateCredit.Total)
			return schemas.ErrorResponse(400, message, fmt.Errorf("%s", message))
		}

		return nil
	})

	if err == nil {
		go database.SaveAuditAsync(r.DB, models.AuditLog{
			MemberID: memberID,
			Method:   "update",
			Path:     "client",
		}, oldCredits, newCredits)
	}

	return err
}
