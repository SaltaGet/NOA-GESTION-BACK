package services

import (
"github.com/SaltaGet/NOA-GESTION-BACK/internal/schemas"
)

func (c *ClientService) ClientGetAll(limit, page int64, search *map[string]string) (*[]schemas.ClientResponseDTO, int64, error) {
	return c.ClientRepository.ClientGetAll(limit, page, search)
}

func (c *ClientService) ClientGetByID(id int64) (*schemas.ClientResponse, error) {
	return c.ClientRepository.ClientGetByID(id)
}

func (c *ClientService) ClientGetByFilter(search string) (*[]schemas.ClientResponseDTO, error) {
	client, err := c.ClientRepository.ClientGetByFilter(search)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *ClientService) ClientCreate(memberID int64,clientCreate *schemas.ClientCreate) (int64, error) {
	client, err := c.ClientRepository.ClientCreate(memberID, clientCreate)
	if err != nil {
		return 0, err
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

func (c *ClientService) ClientDelete(id int64) (error) {
	err := c.ClientRepository.ClientDelete(id)
	if err != nil {
		return err
	}

	return nil
}