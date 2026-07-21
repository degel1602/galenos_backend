package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type appointmentRepository struct {
	db *sql.DB
}

// NewAppointmentRepository construye el adaptador que implementa el puerto
// de salida output.AppointmentRepository contra SQL Server.
func NewAppointmentRepository(db *sql.DB) output.AppointmentRepository {
	return &appointmentRepository{db: db}
}

// Create persiste la cita dentro de una transacción: primero verifica
// disponibilidad del médico con un bloqueo de lectura (UPDLOCK, HOLDLOCK)
// para evitar condiciones de carrera con inserciones concurrentes sobre el
// mismo médico/horario, y luego inserta. Todos los parámetros se envían
// como parámetros nombrados (@Nombre), nunca concatenados, para prevenir
// inyección SQL.
func (r *appointmentRepository) Create(ctx context.Context, appointment *domain.Appointment) error {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("beginning transaction: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // no-op si ya hubo Commit

	const conflictQuery = `
		SELECT COUNT(1)
		FROM dbo.appointments WITH (UPDLOCK, HOLDLOCK)
		WHERE doctor_id = @DoctorID
		  AND status <> 'CANCELLED'
		  AND starts_at < @EndsAt
		  AND ends_at > @StartsAt`

	var conflicts int
	err = tx.QueryRowContext(ctx, conflictQuery,
		sql.Named("DoctorID", string(appointment.DoctorID)),
		sql.Named("StartsAt", appointment.Slot.StartsAt),
		sql.Named("EndsAt", appointment.Slot.EndsAt),
	).Scan(&conflicts)
	if err != nil {
		return fmt.Errorf("checking schedule conflicts: %w", err)
	}
	if conflicts > 0 {
		return domain.ErrDoctorNotAvailable
	}

	const insertQuery = `
		INSERT INTO dbo.appointments
			(id, patient_id, doctor_id, starts_at, ends_at, status, reason, created_at, updated_at)
		VALUES
			(@ID, @PatientID, @DoctorID, @StartsAt, @EndsAt, @Status, @Reason, @CreatedAt, @UpdatedAt)`

	_, err = tx.ExecContext(ctx, insertQuery,
		sql.Named("ID", string(appointment.ID)),
		sql.Named("PatientID", string(appointment.PatientID)),
		sql.Named("DoctorID", string(appointment.DoctorID)),
		sql.Named("StartsAt", appointment.Slot.StartsAt),
		sql.Named("EndsAt", appointment.Slot.EndsAt),
		sql.Named("Status", string(appointment.Status)),
		sql.Named("Reason", appointment.Reason),
		sql.Named("CreatedAt", appointment.CreatedAt),
		sql.Named("UpdatedAt", appointment.UpdatedAt),
	)
	if err != nil {
		return fmt.Errorf("inserting appointment: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("committing transaction: %w", err)
	}

	return nil
}

// GetByID recupera una cita por id usando una consulta parametrizada.
func (r *appointmentRepository) GetByID(ctx context.Context, id domain.AppointmentID) (*domain.Appointment, error) {
	const query = `
		SELECT id, patient_id, doctor_id, starts_at, ends_at, status, reason, created_at, updated_at
		FROM dbo.appointments
		WHERE id = @ID`

	row := r.db.QueryRowContext(ctx, query, sql.Named("ID", string(id)))

	var (
		appointment domain.Appointment
		rawID       string
		patientID   string
		doctorID    string
		status      string
	)

	err := row.Scan(
		&rawID,
		&patientID,
		&doctorID,
		&appointment.Slot.StartsAt,
		&appointment.Slot.EndsAt,
		&status,
		&appointment.Reason,
		&appointment.CreatedAt,
		&appointment.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrAppointmentNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scanning appointment row: %w", err)
	}

	appointment.ID = domain.AppointmentID(rawID)
	appointment.PatientID = domain.PatientID(patientID)
	appointment.DoctorID = domain.DoctorID(doctorID)
	appointment.Status = domain.AppointmentStatus(status)

	return &appointment, nil
}
