package dto_dappradar

import (
	"encoding/json"
	"review-service/pkg/db"
	"review-service/pkg/utils"
)

type EndpointDappRepo struct {
	EndpointDappList []EndpointDapp
}

type EndpointDapp struct {
	Endpoint       string
	BlockchainName string

	DetailDapp *DetailDapp
}

func (endpointDapp *EndpointDapp) InsertDB() error {
	query :=
		`
		INSERT INTO dApp 
			(dAppId, dAppSrc, dAppCode, dAppName, category, 
				subCategory, dAppLogo, description, socials, chains,
					sourceUrl, sourceName, volume24h, user24h, balance, createdDate, updatedDate)
		VALUES
			($1, $2, $3,
				$4, $5, $6,
				$7, $8, $9,
				$10, $11, $12,
				$13, $14, $15, $16, $17);
		`

	var socialsJSONB any //default nil
	var err error
	if endpointDapp.DetailDapp.Social != nil {
		socialsJSONB, err = json.Marshal(endpointDapp.DetailDapp.Social)
		if err != nil {
			return err
		}
	}

	endpointDapp.DetailDapp.CreatedDate = utils.Timestamp()
	endpointDapp.DetailDapp.UpdatedDate = utils.Timestamp()

	subCategories := ``
	for index, category := range endpointDapp.DetailDapp.SubCategories {
		subCategories += category
		if index == len(endpointDapp.DetailDapp.SubCategories)-1 {
			continue
		}
		subCategories += `,`
	}

	chainsJSONB, err := json.Marshal(endpointDapp.DetailDapp.Chains)
	if err != nil {
		return err
	}

	_, err = db.PSQL.Exec(query, endpointDapp.DetailDapp.DAppId, endpointDapp.DetailDapp.DAppSrc,
		endpointDapp.DetailDapp.DAppCode, endpointDapp.DetailDapp.DAppName, endpointDapp.DetailDapp.Category,
		subCategories, endpointDapp.DetailDapp.Image,
		endpointDapp.DetailDapp.Description, socialsJSONB, chainsJSONB,
		endpointDapp.DetailDapp.SourceUrl, endpointDapp.DetailDapp.SourceName, endpointDapp.DetailDapp.Volume24h,
		endpointDapp.DetailDapp.User24h, endpointDapp.DetailDapp.Balance, endpointDapp.DetailDapp.CreatedDate,
		endpointDapp.DetailDapp.UpdatedDate)
	if err != nil {
		return err
	}
	return nil
}
