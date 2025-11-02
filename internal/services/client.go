package services

import (
	"github.com/DanielChachagua/GestionCar/pkg/models"
)

func (c *ClientService) ClientGetAll() (*[]models.Client, error) {
	clients, err := c.ClientRepository.ClientGetAll()
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (c *ClientService) ClientGetByID(id string) (*models.Client, error) {
	client, err := c.ClientRepository.ClientGetByID(id)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ClientService) ClientGetByName(name string) (*[]models.Client, error) {
	client, err := c.ClientRepository.ClientGetByName(name)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ClientService) ClientCreate(clientCreate *models.ClientCreate) (string, error) {
	client, err := c.ClientRepository.ClientCreate(clientCreate)
	if err != nil {
		return "", err
	}

	return client, nil
}

func (c *ClientService) ClientUpdate(clientUpdate *models.ClientUpdate) (error) {
	err := c.ClientRepository.ClientUpdate(clientUpdate)
	if err != nil {
		return err
	}

	return nil
}

func (c *ClientService) ClientDelete(id string) (error) {
	err := c.ClientRepository.ClientDelete(id)
	if err != nil {
		return err
	}

	return nil
}