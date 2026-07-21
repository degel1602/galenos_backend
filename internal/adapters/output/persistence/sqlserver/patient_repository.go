package sqlserver

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
	"github.com/galenos-pro/appointments-api/internal/ports/shared"
)

type patientRepository struct {
	db *sql.DB
}

// NewPatientRepository construye el adaptador que implementa el puerto de
// salida output.PatientRepository contra la tabla pacientes en SQL Server.
func NewPatientRepository(db *sql.DB) output.PatientRepository {
	return &patientRepository{db: db}
}

func (r *patientRepository) List(ctx context.Context, page shared.PageRequest) ([]domain.Patient, int, error) {
	const countQuery = `SELECT COUNT(1) FROM pacientes`

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("counting patients: %w", err)
	}

	const listQuery = `
		SELECT NroDocumento, ApellidoPaterno, ApellidoMaterno, PrimerNombre,
		       SegundoNombre, TercerNombre
		FROM pacientes 
		WHERE NroDocumento is not null 
		ORDER BY NroDocumento
		OFFSET @Offset ROWS FETCH NEXT @PageSize ROWS ONLY`

	rows, err := r.db.QueryContext(ctx, listQuery,
		sql.Named("Offset", page.Offset()),
		sql.Named("PageSize", page.PageSize),
	)
	if err != nil {
		return nil, 0, fmt.Errorf("querying patients: %w", err)
	}
	defer rows.Close()

	patients := make([]domain.Patient, 0, page.PageSize)
	for rows.Next() {
		var (
			patient         domain.Patient
			maternalSurname sql.NullString
			secondName      sql.NullString
			thirdName       sql.NullString
		)
		if err := rows.Scan(
			&patient.DocumentNumber,
			&patient.PaternalSurname,
			&maternalSurname,
			&patient.FirstName,
			&secondName,
			&thirdName,
		); err != nil {
			return nil, 0, fmt.Errorf("scanning patient row: %w", err)
		}
		patient.MaternalSurname = maternalSurname.String
		patient.SecondName = secondName.String
		patient.ThirdName = thirdName.String
		patients = append(patients, patient)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterating patient rows: %w", err)
	}

	return patients, total, nil
}

// GetByDocumentNumber invoca el procedimiento almacenado sp_listarPaciente,
// que recibe @NumDocumento varchar(18) y retorna ApellidoPaterno,
// ApellidoMaterno, PrimerNombre, SegundoNombre y TercerNombre del paciente
// cuyo número de documento coincide (el SP no repite el número de
// documento en el resultado, ya que es el propio parámetro de entrada).
// El nombre del procedimiento se pasa como texto de consulta y el
// parámetro va nombrado (@NumDocumento -> "NumDocumento"), lo que hace que
// el driver mssql lo despache como una llamada RPC en vez de concatenar
// SQL, evitando cualquier inyección.
func (r *patientRepository) GetByDocumentNumber(ctx context.Context, documentNumber string) (*domain.Patient, error) {
	const procedure = `sp_listarPaciente`

	rows, err := r.db.QueryContext(ctx, procedure, sql.Named("NumDocumento", documentNumber))
	if err != nil {
		return nil, fmt.Errorf("calling sp_listarPaciente: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("iterating sp_listarPaciente result: %w", err)
		}
		return nil, domain.ErrPatientNotFound
	}

	var (
		patient         domain.Patient
		maternalSurname sql.NullString
		secondName      sql.NullString
		thirdName       sql.NullString
	)
	if err := rows.Scan(
		&patient.PaternalSurname,
		&maternalSurname,
		&patient.FirstName,
		&secondName,
		&thirdName,
	); err != nil {
		return nil, fmt.Errorf("scanning sp_listarPaciente row: %w", err)
	}
	patient.DocumentNumber = documentNumber
	patient.MaternalSurname = maternalSurname.String
	patient.SecondName = secondName.String
	patient.ThirdName = thirdName.String

	return &patient, nil
}
