package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"go-rest-api/config"
	"go-rest-api/models"
)

func CreateExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var exp models.Expense
	err := json.NewDecoder(r.Body).Decode(&exp)
	if err != nil {
		http.Error(w, "Yanlış JSON formatı", http.StatusBadRequest)
		return
	}

	if exp.Amount <= 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(exp.Category) == "" {
		http.Error(w, "Category is required", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(exp.SpentOn) == "" {
		http.Error(w, "SpentOn date is required", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO expenses (amount, category, note, spent_on)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err = config.DB.QueryRow(query, exp.Amount, exp.Category, exp.Note, exp.SpentOn).Scan(&exp.ID, &exp.CreatedAt)
	if err != nil {
		http.Error(w, "Məlumat bazaya yazılarkən xəta baş verdi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(exp)
}

func GetExpenses(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `SELECT id, amount, category, note, spent_on, created_at FROM expenses ORDER BY created_at DESC`
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Xərcləri oxuyarkən xəta baş verdi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	expenses := []models.Expense{}
	for rows.Next() {
		var exp models.Expense
		var note sql.NullString

		err := rows.Scan(&exp.ID, &exp.Amount, &exp.Category, &note, &exp.SpentOn, &exp.CreatedAt)
		if err != nil {
			http.Error(w, "Məlumat oxunarkən xəta baş verdi", http.StatusInternalServerError)
			return
		}
		if note.Valid {
			exp.Note = note.String
		}

		expenses = append(expenses, exp)
	}

	json.NewEncoder(w).Encode(expenses)
}

func GetExpenseByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Yanlış ID formatı", http.StatusBadRequest)
		return
	}

	query := `SELECT id, amount, category, note, spent_on, created_at FROM expenses WHERE id = $1`
	var exp models.Expense
	var note sql.NullString

	err = config.DB.QueryRow(query, id).Scan(&exp.ID, &exp.Amount, &exp.Category, &note, &exp.SpentOn, &exp.CreatedAt)
	if err == sql.ErrNoRows {
		http.Error(w, "Xərc tapılmadı", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Baza xətası", http.StatusInternalServerError)
		return
	}

	if note.Valid {
		exp.Note = note.String
	}

	json.NewEncoder(w).Encode(exp)
}

func UpdateExpense(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Yanlış ID formatı", http.StatusBadRequest)
		return
	}

	var updateData struct {
		Amount   *float64 `json:"amount"`
		Category *string  `json:"category"`
		Note     *string  `json:"note"`
	}

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Yanlış JSON formatı", http.StatusBadRequest)
		return
	}

	if updateData.Amount != nil && *updateData.Amount <= 0 {
		http.Error(w, "Amount must be greater than 0", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE expenses 
		SET amount = COALESCE($1, amount),
		    category = COALESCE($2, category),
		    note = COALESCE($3, note)
		WHERE id = $4
		RETURNING id, amount, category, note, spent_on, created_at`

	var exp models.Expense
	var note sql.NullString

	err = config.DB.QueryRow(query, updateData.Amount, updateData.Category, updateData.Note, id).
		Scan(&exp.ID, &exp.Amount, &exp.Category, &note, &exp.SpentOn, &exp.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Yenilənəcək xərc tapılmadı", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Yenilənmə zamanı xəta baş verdi", http.StatusInternalServerError)
		return
	}

	if note.Valid {
		exp.Note = note.String
	}

	json.NewEncoder(w).Encode(exp)
}

func DeleteExpense(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Yanlış ID formatı", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM expenses WHERE id = $1`
	res, err := config.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Silinmə zamanı xəta baş verdi", http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "Silinəcək xərc tapılmadı", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetSummaryByCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := `SELECT category, SUM(amount) FROM expenses GROUP BY category`
	rows, err := config.DB.Query(query)
	if err != nil {
		http.Error(w, "Hesablama xətası baş verdi", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	type Summary struct {
		Category    string  `json:"category"`
		TotalAmount float64 `json:"total_amount"`
	}

	summaries := []Summary{}
	for rows.Next() {
		var s Summary
		if err := rows.Scan(&s.Category, &s.TotalAmount); err != nil {
			http.Error(w, "Məlumat oxunarkən xəta baş verdi", http.StatusInternalServerError)
			return
		}
		summaries = append(summaries, s)
	}

	json.NewEncoder(w).Encode(summaries)
}