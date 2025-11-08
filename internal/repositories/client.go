package repositories

import (
	"errors"
	"fmt"

	"github.com/SaltaGet/NOA-GESTION-BACK/internal/models"
	"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *ClientRepository) ClientGetByID(id string) (*schemas.ClientResponse, error) {
	var client schemas.ClientResponse
	if err := r.DB.Where("id = ?", id).First(&client).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, schemas.ErrorResponse(404, "Client no encontrado", err)
		}
		return nil, schemas.ErrorResponse(500, "Error al buscar el cliente", err)
	}
	return &client, nil
}

func (r *ClientRepository) ClientGetByName(name string) (*[]schemas.ClientResponseDTO, error) {
	var client []schemas.ClientResponseDTO
	if err := r.DB.Where("last_name LIKE ? OR first_name LIKE ?", "%"+name+"%", "%"+name+"%").Find(&client).Error; err != nil {
    return nil, schemas.ErrorResponse(500, "Error al buscar el cliente", err)
	}
	return &client, nil
}

func (r *ClientRepository) ClientExist(email, dni, cuil string) (error) {
	var field string
	err := r.DB.Raw(`
		SELECT 
			CASE
				WHEN email = ? THEN 'email'
				WHEN dni = ? THEN 'dni'
				WHEN cuil = ? THEN 'cuil'
			END as field
		FROM clients
		WHERE email = ? OR dni = ? OR cuil = ?
		LIMIT 1
	`, email, dni, cuil, email, dni, cuil).Scan(&field).Error

	if err != nil {
		return schemas.ErrorResponse(500, "Error al corroborar el cliente", err)
	}

	if field == "" {
		return schemas.ErrorResponse(400, fmt.Sprintf("El campo %s ya existe, debe de ser único", field), err)
	}

	return nil
}

func (r *ClientRepository) ClientGetAll() (*[]schemas.ClientResponseDTO, error) {
	var clients []schemas.ClientResponseDTO
	if err := r.DB.Find(&clients).Error; err != nil {
		return nil, schemas.ErrorResponse(500, "Error al buscar los clientes", err)
	}
	return &clients, nil
}

func (r *ClientRepository) ClientCreate(client *schemas.ClientCreate) (string, error) {
	newClient := schemas.ClientResponseDTO{
		ID: uuid.NewString(),
		FirstName: client.FirstName,
		LastName:  client.LastName,
		Identifier: client.Identifier,
		Email:     client.Email,
	}
	if err := r.DB.Create(&newClient).Error; err != nil {
		return "", schemas.ErrorResponse(500, "Error al crear el cliente", err)
	}
	return newClient.ID, nil
}

func (r *ClientRepository) ClientUpdate(client *schemas.ClientUpdate) error {
	if err := r.DB.Where("id = ?", client.ID).Updates(&models.Client{
		FirstName: client.FirstName,
		LastName:  client.LastName,
		Email:     client.Email,
	}).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return schemas.ErrorResponse(404, "Cliente no encontrado", err)
		}	
		return schemas.ErrorResponse(500, "Error al obtener el cliente", err) 
	}
	
	return nil
}

func (r *ClientRepository) ClientDelete(id string) error {
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
