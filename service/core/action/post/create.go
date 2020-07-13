package post

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service/core/action/author"
	"github.com/factly/dega-server/service/core/model"
	factcheckModel "github.com/factly/dega-server/service/factcheck/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/util/slug"
	"github.com/factly/x/errorx"
	"github.com/factly/x/renderx"
	"github.com/factly/x/validationx"
)

// create - Create post
// @Summary Create post
// @Description Create post
// @Tags Post
// @ID add-post
// @Consume json
// @Produce json
// @Param X-User header string true "User ID"
// @Param X-Space header string true "Space ID"
// @Param Post body post true "Post Object"
// @Success 201 {object} postData
// @Router /core/posts [post]
func create(w http.ResponseWriter, r *http.Request) {

	sID, err := util.GetSpace(r.Context())
	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	post := post{}
	result := &postData{}
	result.Categories = make([]model.Category, 0)
	result.Tags = make([]model.Tag, 0)
	result.Authors = make([]model.Author, 0)
	result.Claims = make([]factcheckModel.Claim, 0)

	err = json.NewDecoder(r.Body).Decode(&post)

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	validationError := validationx.Check(post)

	if validationError != nil {
		errorx.Render(w, validationError)
		return
	}

	post.SpaceID = uint(sID)

	var postSlug string
	if post.Slug != "" && slug.Check(post.Slug) {
		postSlug = post.Slug
	} else {
		postSlug = slug.Make(post.Title)
	}

	result.Post = model.Post{
		Title:            post.Title,
		Slug:             slug.Approve(postSlug, sID, config.DB.NewScope(&model.Post{}).TableName()),
		Status:           post.Status,
		Subtitle:         post.Subtitle,
		Excerpt:          post.Excerpt,
		Description:      post.Description,
		IsFeatured:       post.IsFeatured,
		IsHighlighted:    post.IsHighlighted,
		IsSticky:         post.IsSticky,
		FeaturedMediumID: post.FeaturedMediumID,
		FormatID:         post.FormatID,
		PublishedDate:    post.PublishedDate,
		SpaceID:          post.SpaceID,
	}

	// check categories, tags & medium belong to same space or not
	err = post.CheckSpace(config.DB)
	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	err = config.DB.Model(&model.Post{}).Create(&result.Post).Error

	if err != nil {
		errorx.Render(w, errorx.Parser(errorx.DBError()))
		return
	}

	config.DB.Model(&model.Post{}).Preload("Medium").Preload("Format").First(&result.Post)

	if result.Format.Slug == "factcheck" {
		// create post claim
		for _, id := range post.ClaimIDs {
			postClaim := &factcheckModel.PostClaim{}
			postClaim.ClaimID = uint(id)
			postClaim.PostID = result.ID

			err = config.DB.Model(&factcheckModel.PostClaim{}).Create(&postClaim).Error
			if err != nil {
				errorx.Render(w, errorx.Parser(errorx.DBError()))
				return
			}
		}

		// fetch all post claims
		postClaims := []factcheckModel.PostClaim{}
		config.DB.Model(&factcheckModel.PostClaim{}).Where(&factcheckModel.PostClaim{
			PostID: result.ID,
		}).Preload("Claim").Preload("Claim.Claimant").Preload("Claim.Claimant.Medium").Preload("Claim.Rating").Preload("Claim.Rating.Medium").Find(&postClaims)

		// appending all post claims
		for _, postClaim := range postClaims {
			result.Claims = append(result.Claims, postClaim.Claim)
		}
	}

	// create post category & fetch categories
	for _, id := range post.CategoryIDs {
		postCategory := &model.PostCategory{}
		postCategory.CategoryID = uint(id)
		postCategory.PostID = result.ID

		err = config.DB.Model(&model.PostCategory{}).Create(&postCategory).Error

		if err != nil {
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	// fetch all post categories
	postCategories := []model.PostCategory{}
	config.DB.Model(&model.PostCategory{}).Where(&model.PostCategory{
		PostID: result.ID,
	}).Preload("Category").Preload("Category.Medium").Find(&postCategories)

	// appending post categories to result
	for _, postCategory := range postCategories {
		result.Categories = append(result.Categories, postCategory.Category)
	}

	// create post tag
	for _, id := range post.TagIDs {
		postTag := &model.PostTag{}
		postTag.TagID = uint(id)
		postTag.PostID = result.ID

		err = config.DB.Model(&model.PostTag{}).Create(&postTag).Error

		if err != nil {
			errorx.Render(w, errorx.Parser(errorx.DBError()))
			return
		}
	}

	// fetch all post tags
	postTags := []model.PostTag{}
	config.DB.Model(&model.PostTag{}).Where(&model.PostTag{
		PostID: result.ID,
	}).Preload("Tag").Find(&postTags)

	// appending previous post tags to result
	for _, postTag := range postTags {
		result.Tags = append(result.Tags, postTag.Tag)
	}

	// Adding author
	authors, err := author.All(r.Context())
	for _, id := range post.AuthorIDs {
		aID := fmt.Sprint(id)
		if authors[aID].Email != "" {
			author := model.PostAuthor{
				AuthorID: id,
				PostID:   result.Post.ID,
			}
			err := config.DB.Model(&model.PostAuthor{}).Create(&author).Error
			if err == nil {
				result.Authors = append(result.Authors, authors[aID])
			}
		}
	}

	renderx.JSON(w, http.StatusCreated, result)
}
