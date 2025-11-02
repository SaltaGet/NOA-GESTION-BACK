package repositories

import (
	"github.com/jinzhu/gorm"
	"github.com/Daniel-Vinassco/Gestion-Car/pkg/ports"
)

// transactionManager implementa ports.TransactionManager
type transactionManager struct {
	db *gorm.DB
}

// NewTransactionManager crea un nuevo gestor de transacciones.
func NewTransactionManager(db *gorm.DB) ports.TransactionManager {
	return &transactionManager{db: db}
}

// unitOfWork implementa ports.UnitOfWork
type unitOfWork struct {
	tx             *gorm.DB
	attendanceRepo ports.AttendanceRepository
	expenseRepo    ports.ExpenseRepository
}

// Execute envuelve la lógica de negocio en una transacción GORM.
func (tm *transactionManager) Execute(fn func(uow ports.UnitOfWork) error) error {
	return tm.db.Transaction(func(tx *gorm.DB) error {
		uow := &unitOfWork{tx: tx}
		return fn(uow)
	})
}

// --- Implementación de los métodos de UnitOfWork ---

// AttendanceRepository devuelve un repositorio de asistencia que usa la transacción.
func (uow *unitOfWork) AttendanceRepository() ports.AttendanceRepository {
	if uow.attendanceRepo == nil {
		// Crea una NUEVA instancia del repo, pasándole el MANEJADOR DE LA TRANSACCIÓN
		uow.attendanceRepo = NewAttendanceRepository(uow.tx)
	}
	return uow.attendanceRepo
}

// ExpenseRepository devuelve un repositorio de gastos que usa la transacción.
func (uow *unitOfWork) ExpenseRepository() ports.ExpenseRepository {
	if uow.expenseRepo == nil {
		uow.expenseRepo = NewExpenseRepository(uow.tx)
	}
	return uow.expenseRepo
}
