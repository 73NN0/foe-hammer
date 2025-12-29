package ports

import "github.com/73NN0/foe-hammer/domain"

// ToolChecker vérifie que les outils externes sont disponibles
type ToolChecker interface {
	// Check retourne nil si l'outil est disponible, une erreur sinon
	Check(tool string) error
	// Suggest retourne comment installer l'outil manquant
	Suggest(tool string, host domain.Host) string
}
