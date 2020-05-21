package category

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/model"
	"github.com/factly/dega-server/util/render"
	"github.com/go-chi/chi"
)

// update - Update category by id
// @Summary Update a category by id
// @Description Update category by ID
// @Tags Category
// @ID update-category-by-id
// @Produce json
// @Consume json
// @Param X-User header string true "User ID"
// @Param category_id path string true "Category ID"
// @Param Category body category false "Category"
// @Success 200 {object} model.Category
// @Router /core/categories/{category_id} [put]
func update(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "category_id")
	id, err := strconv.Atoi(categoryID)

	if err != nil {
		return
	}

	category := &category{}
	json.NewDecoder(r.Body).Decode(&category)

	result := &model.Category{}
	result.ID = uint(id)

	config.DB.Model(&result).Updates(model.Category{
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
		MediumID:    category.MediumID,
	}).Preload("Medium").First(&result)

	render.JSON(w, http.StatusOK, result)
}
