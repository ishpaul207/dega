package space

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/factly/dega-server/service/core/model"
	"github.com/factly/dega-server/util"
	"github.com/factly/dega-server/util/timex"
	"github.com/factly/x/errorx"
	"github.com/factly/x/loggerx"
	"github.com/factly/x/middlewarex"
	"github.com/factly/x/renderx"
	"github.com/go-chi/chi"
	"github.com/spf13/viper"
)

func users(w http.ResponseWriter, r *http.Request) {
	uID, err := middlewarex.GetUser(r.Context())
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	applicationID, err := util.GetApplicationID(uint(uID), "dega")
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.Unauthorized()))
		return
	}

	sID := chi.URLParam(r, "space_id")
	spaceID, err := strconv.Atoi(sID)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InvalidID()))
		return
	}

	orgID, err := util.GetOrganisationIDfromSpaceID(uint(spaceID), uint(uID))
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	req, err := http.NewRequest("GET", viper.GetString("kavach_url")+fmt.Sprintf("/organisations/%d/applications/", orgID)+fmt.Sprintf("%d", applicationID)+"/spaces/"+fmt.Sprintf("%d", spaceID)+"/users", nil)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}
	req.Header.Set("X-User", fmt.Sprintf("%d", uID))
	req.Header.Set("Content-Type", "application/json")
	client := http.Client{Timeout: time.Second * time.Duration(timex.HTTP_TIMEOUT)}
	resp, err := client.Do(req)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	if resp.StatusCode != 200 {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.InternalServerError()))
		return
	}

	defer resp.Body.Close()

	var users []model.User

	err = json.NewDecoder(resp.Body).Decode(&users)
	if err != nil {
		loggerx.Error(err)
		errorx.Render(w, errorx.Parser(errorx.DecodeError()))
		return
	}

	renderx.JSON(w, http.StatusOK, users)
}
