package handlers

import (
	"net/http"
)

func (apiCfg *ApiConfig) GetCreatorsAvailable(w http.ResponseWriter, r *http.Request) {
	//get creators from DB
	dbCreatorNamesList, err := apiCfg.DB.GetCreators(r.Context())
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Couldn't get creators names", err)
		return
	}

	//put only the actual names in a list
	creatorNames := []string{}
	for _, creator := range dbCreatorNamesList {
		if creator.Valid {
			creatorNames = append(creatorNames, creator.String)
		}
	}

	RespondWithJSON(w, http.StatusAccepted, creatorNames)
}