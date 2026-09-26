package repository

import (
	"dezzles-apps/build-engine/model"
	_ "embed"
	"fmt"
	"log"
)

//go:embed sql/update-portfolio.sql
var updateString string

//go:embed sql/insert-portfolio.sql
var insertString string

type PortfolioRepository struct {
	database *Database
}

func (r *PortfolioRepository) Initialise(database *Database) error {
	r.database = database
	return nil
}

func (r *PortfolioRepository) GetPortfolioId(projectKey string) (int, error) {
	rows, err := r.database.getDB().Query("SELECT portfolio_id from portfolio_entries where project_key = ?", projectKey)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var portfolioID int
	if rows.Next() {
		if err := rows.Scan(&portfolioID); err != nil {
			return 0, err
		}
		fmt.Println("Found portfolioID: ", portfolioID)
		return portfolioID, nil
	} else {
		return 0, nil
	}
}

func (r *PortfolioRepository) UpdatePortfolio(update model.PublishEvent) error {
	portfolioId, err := r.GetPortfolioId(update.Name)
	if err != nil {
		return err
	}
	if portfolioId == 0 {
		fmt.Println("Inserting new portfolio entry for projectKey: ", update.Name)
		r.database.getDB().Exec(
			insertString,
			update.Name,
			update.Config.Title,
			update.Config.Description,
			update.Version,
			update.Config.Url,
			"prod",
			true,
			update.Config.Category,
			update.Config.ShareRepository,
			update.RepositoryUrl,
		)
		portfolioId, err = r.GetPortfolioId(update.Name)
		if err != nil {
			return err
		}
	} else {
		fmt.Println("Updating portfolio entry for projectKey: ", update.Name)
		r.database.getDB().Exec(
			updateString,
			update.Config.Title,
			update.Config.Description,
			update.Version,
			update.Config.Url,
			"prod",
			true,
			update.Config.Category,
			update.Config.ShareRepository,
			update.RepositoryUrl,
			update.Name,
		)
		fmt.Println("Updated portfolio entry for projectKey: ", update.Name)
	}
	return r.UpdateTechnologies(portfolioId, update.Config.Technologies)
}

func (r *PortfolioRepository) UpdateTechnologies(portfolioId int, technologies []string) error {
	log.Printf("Removing technologies for %d: %s", portfolioId, technologies)
	_, err := r.database.getDB().Exec("DELETE FROM portfolio_entry_technologies WHERE portfolio_id = ?", portfolioId)
	if err != nil {
		return err
	}
	log.Printf("Removed technologies for %d: %s", portfolioId, technologies)

	if len(technologies) == 0 {
		return nil
	}

	rows, err := r.database.getDB().Query("SELECT technology_id, name FROM technologies")
	if err != nil {
		return err
	}
	defer rows.Close()

	existingTechs := make(map[string]int)
	for rows.Next() {
		var techID int
		var techName string
		if err := rows.Scan(&techID, &techName); err != nil {
			return err
		}
		existingTechs[techName] = techID
	}

	// Log warnings for missing technologies
	for _, tech := range technologies {
		if _, exists := existingTechs[tech]; !exists {
			fmt.Printf("WARNING: Technology '%s' does not exist in the technologies table\n", tech)
		}
	}

	// Build insert query for existing technologies only
	insertQuery := "INSERT INTO portfolio_entry_technologies (portfolio_id, technology_id) VALUES "
	insertArgs := make([]interface{}, 0)

	for i, tech := range technologies {
		if _, exists := existingTechs[tech]; !exists {
			continue
		}
		if i > 0 {
			insertQuery += ", "
		}
		insertQuery += "(?, ?)"
		insertArgs = append(insertArgs, portfolioId)
		insertArgs = append(insertArgs, existingTechs[tech])
	}
	// Execute insert for all technologies (non-existent ones will be silently skipped by the IN clause)
	_, err = r.database.getDB().Exec(insertQuery, insertArgs...)
	return err
}
