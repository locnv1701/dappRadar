package dao

import (
	"review-service/pkg/db"
)

type ChainList struct {
	List []Chain
}

type Chain struct {
	ChainId         string `json:"chainid"`
	Chainname       string `json:"chainname"`
	ExplorerWebsite string `json:"explorerWebsite"`
	Path            string `json:"path"`
	Image           string `json:"image"`
}

func (chainList *ChainList) GetChainList() error {
	query := "select distinct on (chainId, chainname) chainId, chainName from chain_list"

	rows, err := db.PSQL.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		chain := Chain{}
		err := rows.Scan(&chain.ChainId, &chain.Chainname)
		if err != nil {
			return err
		}
		chainList.List = append(chainList.List, chain)
	}

	return nil

}
