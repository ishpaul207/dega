package rating

import (
	"encoding/json"
	"net/http"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/factcheck/model"
	"github.com/factly/dega-server/util/render"
)

// create - Create rating
// @Summary Create rating
// @Description Create rating
// @Tags Rating
// @ID add-rating
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param Rating body rating true "Rating Object"
// @Success 201 {object} model.Rating
// @Router /factcheck/ratings [post]
func create(w http.ResponseWriter, r *http.Request) {

	rating := &rating{}

	json.NewDecoder(r.Body).Decode(&rating)

	result := &model.Rating{
		Name:         rating.Name,
		Slug:         rating.Slug,
		Description:  rating.Description,
		MediumID:     rating.MediumID,
		SpaceID:      rating.SpaceID,
		NumericValue: rating.NumericValue,
	}

	err := config.DB.Model(&model.Rating{}).Create(&result).Error

	if err != nil {
		return
	}

	config.DB.Model(&model.Rating{}).Preload("Medium").First(&result)

	render.JSON(w, http.StatusCreated, result)
}
