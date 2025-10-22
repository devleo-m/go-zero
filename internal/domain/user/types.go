package user

// Role representa o papel/função do usuário no sistema
type Role string

const (
	// RoleAdmin representa um administrador do sistema
	RoleAdmin Role = "admin"

	// RoleManager representa um gerente
	RoleManager Role = "manager"

	// RoleUser representa um usuário comum
	RoleUser Role = "user"

	// RoleGuest representa um usuário convidado
	RoleGuest Role = "guest"
)

// String retorna a representação string do role
func (r Role) String() string {
	return string(r)
}

// IsValid verifica se o role é válido
func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleManager, RoleUser, RoleGuest:
		return true
	default:
		return false
	}
}

// HasPermission verifica se o role tem uma permissão específica
func (r Role) HasPermission(permission string) bool {
	switch r {
	case RoleAdmin:
		// Admin tem todas as permissões
		return true
	case RoleManager:
		// Manager tem permissões limitadas
		managerPermissions := []string{
			"read_users",
			"create_users",
			"update_users",
			"read_reports",
			"manage_appointments",
		}
		for _, perm := range managerPermissions {
			if perm == permission {
				return true
			}
		}
		return false
	case RoleUser:
		// User tem permissões básicas
		userPermissions := []string{
			"read_own_profile",
			"update_own_profile",
			"create_appointments",
			"read_own_appointments",
		}
		for _, perm := range userPermissions {
			if perm == permission {
				return true
			}
		}
		return false
	case RoleGuest:
		// Guest tem permissões mínimas
		guestPermissions := []string{
			"read_public_info",
		}
		for _, perm := range guestPermissions {
			if perm == permission {
				return true
			}
		}
		return false
	default:
		return false
	}
}

// CanManage verifica se o role pode gerenciar outro role
func (r Role) CanManage(targetRole Role) bool {
	// Hierarquia de roles (maior para menor)
	hierarchy := map[Role]int{
		RoleAdmin:   4,
		RoleManager: 3,
		RoleUser:    2,
		RoleGuest:   1,
	}

	return hierarchy[r] > hierarchy[targetRole]
}

// Status representa o status do usuário
type Status string

const (
	// StatusActive representa um usuário ativo
	StatusActive Status = "active"

	// StatusInactive representa um usuário inativo
	StatusInactive Status = "inactive"

	// StatusPending representa um usuário pendente de ativação
	StatusPending Status = "pending"

	// StatusSuspended representa um usuário suspenso
	StatusSuspended Status = "suspended"
)

// String retorna a representação string do status
func (s Status) String() string {
	return string(s)
}

// IsValid verifica se o status é válido
func (s Status) IsValid() bool {
	switch s {
	case StatusActive, StatusInactive, StatusPending, StatusSuspended:
		return true
	default:
		return false
	}
}

// CanLogin verifica se o usuário pode fazer login com este status
func (s Status) CanLogin() bool {
	return s == StatusActive
}

// RequiresActivation verifica se o status requer ativação
func (s Status) RequiresActivation() bool {
	return s == StatusPending
}
