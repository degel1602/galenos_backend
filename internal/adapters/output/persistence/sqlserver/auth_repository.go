package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) output.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Login(ctx context.Context, username, password string) (int, error) {
	var resultStr sql.NullString
	// Procedimiento almacenado usp_go_Login devuelve un string como salida
	// "SUCCESS;IdEmpleado;..." o "ERROR;Mensaje"
	_, err := r.db.ExecContext(ctx, "EXEC usp_go_Login @Usuario = @p1, @Password = @p2, @Resultado = @p3 OUTPUT",
		sql.Named("p1", username),
		sql.Named("p2", password),
		sql.Named("p3", sql.Out{Dest: &resultStr}),
	)
	if err != nil {
		return 0, fmt.Errorf("error ejecutando usp_go_Login: %w", err)
	}

	res := resultStr.String
	log.Printf("DB Login Response: %q", res)

	parts := strings.Split(res, ";")
	if len(parts) > 0 && parts[0] == "ERROR" {
		msg := "Credenciales inválidas"
		if len(parts) > 1 {
			msg = parts[1]
		}
		return 0, fmt.Errorf("%w: %s", domain.ErrInvalidCredentials, msg)
	}

	if len(parts) > 0 && parts[0] == "OK" {
		var idEmpleado int
		// Buscamos el IdEmpleado porque el SP no lo devuelve en el parámetro OUTPUT
		err = r.db.QueryRowContext(ctx, "SELECT IdEmpleado FROM dbo.Empleados WHERE Usuario = @p1", sql.Named("p1", username)).Scan(&idEmpleado)
		if err != nil {
			return 0, fmt.Errorf("error obteniendo IdEmpleado tras login exitoso: %w", err)
		}
		return idEmpleado, nil
	}

	return 0, fmt.Errorf("formato de respuesta de login inesperado: %s", res)
}

func (r *authRepository) GetMenus(ctx context.Context, idEmpleado int) ([]domain.Menu, error) {
	rows, err := r.db.QueryContext(ctx, "EXEC webMenuSeleccionarIdEmpleado @IdEmpleado = @p1", sql.Named("p1", idEmpleado))
	if err != nil {
		return nil, fmt.Errorf("error consultando menus: %w", err)
	}
	defer rows.Close()

	var menus []domain.Menu
	for rows.Next() {
		var m domain.Menu
		var texto, keyIconWeb, claveWeb sql.NullString
		if err := rows.Scan(&m.IdListGrupo, &texto, &keyIconWeb, &claveWeb, &m.Indice, &m.Estado, &m.NroSubMenu); err != nil {
			return nil, fmt.Errorf("error escaneando menu: %w", err)
		}
		m.Texto = texto.String
		m.KeyIconWeb = keyIconWeb.String
		m.ClaveWeb = claveWeb.String
		menus = append(menus, m)
	}
	return menus, nil
}

func (r *authRepository) GetMenuPermisos(ctx context.Context, idEmpleado int) ([]domain.MenuPermiso, error) {
	rows, err := r.db.QueryContext(ctx, "EXEC web_MenuPermisosIdempleado @IdUsuario = @p1", sql.Named("p1", idEmpleado))
	if err != nil {
		return nil, fmt.Errorf("error consultando menu permisos: %w", err)
	}
	defer rows.Close()

	var permisos []domain.MenuPermiso
	for rows.Next() {
		var p domain.MenuPermiso
		var opciones, texto, menu, keyIconWeb, claveWeb sql.NullString
		if err := rows.Scan(&opciones, &p.Indice, &texto, &menu, &p.IdListGrupo, &keyIconWeb, &p.Estado, &claveWeb, &p.Agregar, &p.Modificar, &p.Eliminar); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				continue
			}
		}
		p.Opciones = opciones.String
		p.Texto = texto.String
		p.Menu = menu.String
		p.KeyIconWeb = keyIconWeb.String
		p.ClaveWeb = claveWeb.String
		permisos = append(permisos, p)
	}
	return permisos, nil
}
