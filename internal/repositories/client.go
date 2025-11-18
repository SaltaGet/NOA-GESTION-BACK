package repositories

import (
	"errors"
	"fmt"
	"strings"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
)

func (r *ClientRepository) ClientGetByID(id int64) (*schemas.ClientResponse, error) {
	var client models.Client
	if err := r.DB.
	Preload("MemberCreate", func(db *gorm.DB) *gorm.DB { 
		return db.Select("id", "first_name", "last_name", "username") 
	}).
	Preload("Pay", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "income_sale_id", "client_id", "total", "method_pay").Where("method_pay = ?", "credit")
	}).
	Where("id = ?", id).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Client no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al buscar el cliente", err)
	}

	var clientResponse schemas.ClientResponse
	copier.Copy(&clientResponse, &client)


	return &clientResponse, nil
}

func (r *ClientRepository) ClientGetByFilter(search string) (*[]schemas.ClientResponseDTO, error) {
	var client []models.Client
	if err := r.DB.Limit(10).Where("last_name LIKE ? OR first_name LIKE ? OR identifier LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%").Find(&client).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error al buscar el cliente", err)
	}

	var clientResponse []schemas.ClientResponseDTO
	copier.Copy(&clientResponse, &client)

	return &clientResponse, nil
}

func (r *ClientRepository) ClientGetAll(limit, page int64, search *map[string]string) (*[]schemas.ClientResponseDTO, int64, error) {
	var clients []schemas.ClientResponseDTO

	query := r.DB.Model(&models.Client{})

	// Aplicar filtros dinámicos
	if search != nil {
		for key, value := range *search {
			if value == "" {
				continue
			}

			switch strings.ToLower(key) {
			case "identifier":
				query = query.Where("identifier LIKE ?", "%"+value+"%")
			case "first_name":
				query = query.Where("first_name LIKE ?", "%"+value+"%")
			case "last_name":
				query = query.Where("last_name LIKE ?", "%"+value+"%")
			case "email":
				query = query.Where("email LIKE ?", "%"+value+"%")
			}
		}
	}

	if limit > 0 {
		offset := (page - 1) * limit
		query = query.Limit(int(limit)).Offset(int(offset))
	}

	if err := query.Find(&clients).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "Error al buscar los clientes", err)
	}

	// Contar total de registros
	var total int64
	if err := r.DB.Model(&models.Client{}).Count(&total).Error; err != nil {
		return nil, 0, schemas.ErrorResponse(500, "Error al contar los clientes", err)
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
	}
	if err := r.DB.Create(&newClient).Error; err != nil {
		if schemas.IsDuplicateError(err) {
			msg := err.Error()
			switch {
			case strings.Contains(msg, "email"):
				return 0, schemas.ErrorResponse(409, fmt.Sprintf("Ya existe un cliente con el email %s", *client.Email), err)
			case strings.Contains(msg, "identifier"):
				return 0, schemas.ErrorResponse(409, fmt.Sprintf("Ya existe un cliente con el identificador %s", *client.Identifier), err)
			default:
				return 0, schemas.ErrorResponse(409, "El cliente ya existe", err)
			}
		}
		return 0, schemas.ErrorResponse(500, "Error al crear el cliente", err)
	}
	return newClient.ID, nil
}

func (r *ClientRepository) ClientUpdate(client *schemas.ClientUpdate) error {
	if err := r.DB.Where("id = ?", client.ID).Updates(&models.Client{
		FirstName: client.FirstName,
		LastName:  client.LastName,
		CompanyName: client.CompanyName,
		Identifier: client.Identifier,
		Email:     client.Email,
		Phone:     client.Phone,
		Address:   client.Address,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Cliente no encontrado", err)
		}
		if schemas.IsDuplicateError(err) {
			msg := err.Error()
			switch {
			case strings.Contains(msg, "email"):
				return schemas.ErrorResponse(409, fmt.Sprintf("Ya existe un cliente con el email %s", *client.Email), err)
			case strings.Contains(msg, "identifier"):
				return schemas.ErrorResponse(409, fmt.Sprintf("Ya existe un cliente con el identificador %s", *client.Identifier), err)
			default:
				return schemas.ErrorResponse(409, "El cliente ya existe", err)
			}
		}
		return schemas.ErrorResponse(500, "Error al obtener el cliente", err)
	}

	return nil
}

func (r *ClientRepository) ClientDelete(id int64) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("id = ?", id).Delete(&models.Client{}).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return schemas.ErrorResponse(404, "Cliente no encontrado", err)
			}
			return schemas.ErrorResponse(500, "Error al eliminar el cliente", err)
		}

		return nil
	})
}
