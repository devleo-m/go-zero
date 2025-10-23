package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserPreferencesModel representa as preferências do usuário
// Separado para lazy loading e melhor organização
type UserPreferencesModel struct {
	// Campos base
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time      `gorm:"not null;index"`
	UpdatedAt time.Time      `gorm:"not null;index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Chave estrangeira
	UserID uuid.UUID `gorm:"type:uuid;not null;uniqueIndex;constraint:OnDelete:CASCADE"`

	// Preferências de localização
	Timezone *string `gorm:"type:varchar(50);default:'UTC'"`
	Language *string `gorm:"type:varchar(10);default:'pt-BR'"`
	Currency *string `gorm:"type:varchar(3);default:'BRL'"`

	// Preferências de notificação
	EmailNotifications     bool `gorm:"default:true;not null"`
	SMSNotifications       bool `gorm:"default:false;not null"`
	PushNotifications      bool `gorm:"default:true;not null"`
	MarketingEmails        bool `gorm:"default:false;not null"`
	SecurityAlerts         bool `gorm:"default:true;not null"`
	AppointmentReminders   bool `gorm:"default:true;not null"`
	NewsletterSubscription bool `gorm:"default:false;not null"`

	// Preferências de privacidade
	ProfileVisibility string `gorm:"type:varchar(20);default:'private';not null"` // public, private, friends
	ShowEmail         bool   `gorm:"default:false;not null"`
	ShowPhone         bool   `gorm:"default:false;not null"`
	ShowLastLogin     bool   `gorm:"default:false;not null"`

	// Preferências de interface
	Theme          string `gorm:"type:varchar(20);default:'light';not null"` // light, dark, auto
	DateFormat     string `gorm:"type:varchar(20);default:'DD/MM/YYYY';not null"`
	TimeFormat     string `gorm:"type:varchar(10);default:'24h';not null"` // 12h, 24h
	ItemsPerPage   int    `gorm:"default:20;not null"`
	AutoSaveDrafts bool   `gorm:"default:true;not null"`
	ShowTutorials  bool   `gorm:"default:true;not null"`

	// Preferências de acessibilidade
	HighContrast  bool `gorm:"default:false;not null"`
	LargeText     bool `gorm:"default:false;not null"`
	ScreenReader  bool `gorm:"default:false;not null"`
	KeyboardNav   bool `gorm:"default:false;not null"`
	ReducedMotion bool `gorm:"default:false;not null"`
}

// TableName define o nome da tabela
func (UserPreferencesModel) TableName() string {
	return "user_preferences"
}

// BeforeCreate hook executado antes de criar as preferências
func (p *UserPreferencesModel) BeforeCreate(tx *gorm.DB) error {
	// Validar UserID
	if p.UserID == uuid.Nil {
		return gorm.ErrInvalidData
	}

	// Gerar ID se não existir
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}

	// Definir timestamps
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}

	// Validar valores de enum
	if !isValidProfileVisibility(p.ProfileVisibility) {
		p.ProfileVisibility = "private"
	}
	if !isValidTheme(p.Theme) {
		p.Theme = "light"
	}
	if !isValidTimeFormat(p.TimeFormat) {
		p.TimeFormat = "24h"
	}

	return nil
}

// BeforeUpdate hook executado antes de atualizar as preferências
func (p *UserPreferencesModel) BeforeUpdate(tx *gorm.DB) error {
	// Atualizar timestamp
	p.UpdatedAt = time.Now()

	// Validar valores se foram alterados
	if tx.Statement.Changed("ProfileVisibility") && !isValidProfileVisibility(p.ProfileVisibility) {
		return gorm.ErrInvalidData
	}
	if tx.Statement.Changed("Theme") && !isValidTheme(p.Theme) {
		return gorm.ErrInvalidData
	}
	if tx.Statement.Changed("TimeFormat") && !isValidTimeFormat(p.TimeFormat) {
		return gorm.ErrInvalidData
	}

	return nil
}

// Validação de enums
func isValidProfileVisibility(visibility string) bool {
	valid := []string{"public", "private", "friends"}
	for _, v := range valid {
		if v == visibility {
			return true
		}
	}
	return false
}

func isValidTheme(theme string) bool {
	valid := []string{"light", "dark", "auto"}
	for _, t := range valid {
		if t == theme {
			return true
		}
	}
	return false
}

func isValidTimeFormat(format string) bool {
	valid := []string{"12h", "24h"}
	for _, f := range valid {
		if f == format {
			return true
		}
	}
	return false
}

// GetTimezone retorna o timezone ou UTC como padrão
func (p *UserPreferencesModel) GetTimezone() string {
	if p.Timezone != nil && *p.Timezone != "" {
		return *p.Timezone
	}
	return "UTC"
}

// GetLanguage retorna o idioma ou pt-BR como padrão
func (p *UserPreferencesModel) GetLanguage() string {
	if p.Language != nil && *p.Language != "" {
		return *p.Language
	}
	return "pt-BR"
}

// GetCurrency retorna a moeda ou BRL como padrão
func (p *UserPreferencesModel) GetCurrency() string {
	if p.Currency != nil && *p.Currency != "" {
		return *p.Currency
	}
	return "BRL"
}

// IsNotificationEnabled verifica se um tipo de notificação está habilitado
func (p *UserPreferencesModel) IsNotificationEnabled(notificationType string) bool {
	switch notificationType {
	case "email":
		return p.EmailNotifications
	case "sms":
		return p.SMSNotifications
	case "push":
		return p.PushNotifications
	case "marketing":
		return p.MarketingEmails
	case "security":
		return p.SecurityAlerts
	case "appointment":
		return p.AppointmentReminders
	case "newsletter":
		return p.NewsletterSubscription
	default:
		return false
	}
}

// SetNotification define o status de um tipo de notificação
func (p *UserPreferencesModel) SetNotification(notificationType string, enabled bool) {
	switch notificationType {
	case "email":
		p.EmailNotifications = enabled
	case "sms":
		p.SMSNotifications = enabled
	case "push":
		p.PushNotifications = enabled
	case "marketing":
		p.MarketingEmails = enabled
	case "security":
		p.SecurityAlerts = enabled
	case "appointment":
		p.AppointmentReminders = enabled
	case "newsletter":
		p.NewsletterSubscription = enabled
	}
}
