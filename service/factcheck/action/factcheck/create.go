package factcheck

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/factcheck/model"
	"github.com/factly/dega-server/util/render"
	"github.com/factly/dega-server/validation"
	"github.com/go-chi/chi"
	"github.com/go-playground/validator/v10"
)

// create - Create factcheck
// @Summary Create factcheck
// @Description Create factcheck
// @Tags Factcheck
// @ID add-factcheck
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param space_id path string true "Space ID"
// @Param Factcheck body factcheck true "Factcheck Object"
// @Success 201 {object} factcheckData
// @Failure 400 {array} string
// @Router /{space_id}/factcheck/factchecks [post]
func create(w http.ResponseWriter, r *http.Request) {

	spaceID := chi.URLParam(r, "space_id")
	sid, err := strconv.Atoi(spaceID)

	factcheck := factcheck{}
	result := &factcheckData{}

	json.NewDecoder(r.Body).Decode(&factcheck)

	validate := validator.New()

	err = validate.Struct(factcheck)

	if err != nil {
		msg := err.Error()
		validation.ValidErrors(w, r, msg)
		return
	}

	factcheck.SpaceID = uint(sid)

	result.Factcheck = model.Factcheck{
		Title:            factcheck.Title,
		Slug:             factcheck.Slug,
		Status:           factcheck.Status,
		Subtitle:         factcheck.Subtitle,
		Excerpt:          factcheck.Excerpt,
		Updates:          factcheck.Updates,
		Description:      factcheck.Description,
		IsFeatured:       factcheck.IsFeatured,
		IsHighlighted:    factcheck.IsHighlighted,
		IsSticky:         factcheck.IsSticky,
		FeaturedMediumID: factcheck.FeaturedMediumID,
		PublishedDate:    factcheck.PublishedDate,
		SpaceID:          factcheck.SpaceID,
	}

	// check claims, categories, tags & medium belong to same space or not
	err = factcheck.CheckSpace(config.DB)
	if err != nil {
		validation.Error(w, r, err.Error())
		return
	}

	err = config.DB.Model(&model.Factcheck{}).Create(&result.Factcheck).Error

	if err != nil {
		return
	}

	config.DB.Model(&model.Factcheck{}).Preload("Medium").Find(&result.Factcheck)

	// create factcheck category & fetch categories
	for _, id := range factcheck.CategoryIDS {
		factcheckCategory := &model.FactcheckCategory{}

		factcheckCategory.CategoryID = uint(id)
		factcheckCategory.FactcheckID = result.ID

		err = config.DB.Model(&model.FactcheckCategory{}).Create(&factcheckCategory).Error

		if err != nil {
			return
		}
		config.DB.Model(&model.FactcheckCategory{}).Preload("Category").Preload("Category.Medium").First(&factcheckCategory)
		result.Categories = append(result.Categories, factcheckCategory.Category)
	}
	// create factcheck tag & fetch tags
	for _, id := range factcheck.TagIDS {
		factcheckTag := &model.FactcheckTag{}

		factcheckTag.TagID = uint(id)
		factcheckTag.FactcheckID = result.ID

		err = config.DB.Model(&model.FactcheckTag{}).Create(&factcheckTag).Error

		if err != nil {
			return
		}
		config.DB.Model(&model.FactcheckTag{}).Preload("Tag").First(&factcheckTag)
		result.Tags = append(result.Tags, factcheckTag.Tag)
	}

	// create factcheck claim & fetch claims
	for _, id := range factcheck.ClaimIDS {
		factcheckClaim := &model.FactcheckClaim{}

		factcheckClaim.ClaimID = uint(id)
		factcheckClaim.FactcheckID = result.ID

		err = config.DB.Model(&model.FactcheckClaim{}).Create(&factcheckClaim).Error

		if err != nil {
			return
		}

		config.DB.Model(&model.FactcheckClaim{}).Preload("Claim").Preload("Claim.Claimant").Preload("Claim.Claimant.Medium").Preload("Claim.Rating").Preload("Claim.Rating.Medium").First(&factcheckClaim)
		result.Claims = append(result.Claims, factcheckClaim.Claim)
	}

	render.JSON(w, http.StatusCreated, result)
}
