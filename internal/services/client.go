package services

import (
"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (c *ClientService) ClientGetAll() (*[]schemas.Client, error) {
	clients, err := c.ClientRepository.ClientGetAll()
	if err != nil {
		return nil, err
	}

	return clients, nil
}

func (c *ClientService) ClientGetByID(id string) (*schemas.Client, error) {
	client, err := c.ClientRepository.ClientGetByID(id)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ClientService) ClientGetByName(name string) (*[]schemas.Client, error) {
	client, err := c.ClientRepository.ClientGetByName(name)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ClientService) ClientCreate(clientCreate *schemas.ClientCreate) (string, error) {
	client, err := c.ClientRepository.ClientCreate(clientCreate)
	if err != nil {
		return "", err
	}

	return client, nil
}

func (c *ClientService) ClientUpdate(clientUpdate *schemas.ClientUpdate) (error) {
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