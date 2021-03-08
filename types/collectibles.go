package types

type (
	Collection struct {
		Id           string `json:"id"`
		Name         string `json:"name"`
		ImageUrl     string `json:"image_url"`
		Description  string `json:"description"`
		ExternalLink string `json:"external_link"`
		Total        int    `json:"total"`
		Address      string `json:"address"`
		Coin         uint   `json:"coin"`
		Type         string `json:"-"`
	}

	CollectionPage []Collection

	CollectibleTransferFee struct {
		Asset  string `json:"asset"`
		Amount string `json:"amount"`
	}

	Collectible struct {
		ID              string                  `json:"id"`
		CollectionID    string                  `json:"collection_id"`
		TokenID         string                  `json:"token_id"`
		ContractAddress string                  `json:"contract_address"`
		Category        string                  `json:"category"`
		ImageUrl        string                  `json:"image_url"`
		ExternalLink    string                  `json:"external_link"`
		ProviderLink    string                  `json:"provider_link"`
		Type            string                  `json:"type"`
		Description     string                  `json:"description"`
		Coin            uint                    `json:"coin"`
		Name            string                  `json:"name"`
		Version         string                  `json:"nft_version"`
		TransferFee     *CollectibleTransferFee `json:"transfer_fee,omitempty"`
	}

	CollectiblePage []Collectible
)
