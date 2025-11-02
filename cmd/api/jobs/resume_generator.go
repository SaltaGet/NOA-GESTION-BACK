package jobs

import (
	"github.com/DanielChachagua/GestionCar/cmd/api/logging"
	"github.com/DanielChachagua/GestionCar/pkg/database"
	"github.com/DanielChachagua/GestionCar/pkg/dependencies"
	"github.com/DanielChachagua/GestionCar/pkg/repositories"
	"github.com/DanielChachagua/GestionCar/pkg/utils"
)

func GenetateResume(deps *dependencies.Application) error {
	connections, err := deps.TenantController.TenantService.TenantGetConections()
	if err != nil {
		return err
	}
	
	first, end := utils.GetFirstAndLastDayTwoMonthsAgo(2)

	for _, connection := range *connections {
		db, err := database.GetTenantDB(connection)
		if err != nil {
			logging.ERROR("Error al obtener conexion a la db: %s", err.Error())
			return err
		}
		err = repositories.GenerateResume(db, first, end)
		if err != nil {
			logging.ERROR("Error al generar resumen: %s", err.Error())
			return err
		}
	}

	logging.INFO("✅ Resumen generado con exito")
	return nil
}