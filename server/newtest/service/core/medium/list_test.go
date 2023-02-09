package medium

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/factly/dega-server/config"
	"github.com/factly/dega-server/service"
	"github.com/factly/dega-server/service/core/model"
	"github.com/gavv/httpexpect/v2"
	"gopkg.in/h2non/gock.v1"
)

func TestMediumList(t *testing.T) {
	defer gock.DisableNetworking()
	testServer := httptest.NewServer(service.RegisterRoutes())
	gock.New(testServer.URL).EnableNetworking().Persist()
	defer gock.DisableNetworking()
	defer testServer.Close()
	//delete all entries from the db and insert some data
	config.DB.Exec("DELETE FROM media")
	config.DB.Exec("DELETE FROM space_permissions")

	e := httpexpect.New(t, testServer.URL)

	t.Run("get empty list of media", func(t *testing.T) {
		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 0})
	})

	t.Run("get non-empty list of media", func(t *testing.T) {
		var insertData = &model.Medium{
			Name:        "Create Medium Test 1",
			Slug:        "create-medium-test-1",
			Description: TestDescription,
			Type:        TestType,
			Title:       TestTitle,
			Caption:     TestCaption,
			AltText:     TestAltText,
			FileSize:    TestFileSize,
			URL:         TestUrl,
			Dimensions:  TestDimensions,
			MetaFields:  TestMetaFields,
			SpaceID:     TestSpaceID,
		}
		config.DB.Create(insertData)
		insertData = &model.Medium{
			Name:        "Create Medium Test 2",
			Slug:        "create-medium-test-2",
			Description: TestDescription,
			Type:        TestType,
			Title:       TestTitle,
			Caption:     TestCaption,
			AltText:     TestAltText,
			FileSize:    TestFileSize,
			URL:         TestUrl,
			Dimensions:  TestDimensions,
			MetaFields:  TestMetaFields,
			SpaceID:     TestSpaceID,
		}
		config.DB.Create(insertData)
		insertSpacePermission := &model.SpacePermission{
			SpaceID:   1,
			FactCheck: true,
			Media:     10,
			Posts:     10,
			Podcast:   true,
			Episodes:  10,
			Videos:    10,
		}
		config.DB.Create(insertSpacePermission)

		e.GET(basePath).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 2})

	})
	// should run test after running 'get non-empty list of media'
	t.Run("get media with pagination", func(t *testing.T) {

		res := e.GET(basePath).
			WithQueryObject(map[string]string{
				"limit": "1",
				"page":  "2",
				"sort":  "asc",
			}).
			WithHeaders(headers).
			Expect().
			Status(http.StatusOK).
			JSON().
			Object().
			ContainsMap(map[string]interface{}{"total": 2}).
			Value("nodes").
			Array().
			Element(0).
			Object()

		Data["name"] = "Create Medium Test 2"
		Data["slug"] = "create-medium-test-2"

		res.ContainsMap(Data)

	})
}
