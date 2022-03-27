package chaincode

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

//Badge 徽章
type Badge struct {
	ID          string `json:"ID"`
	PublisherID string `json:"PublisherID"`
	Description string `json:"Description"`
	Validity    int    `json:"Validity"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	badges := []Badge{
		{ID: "B032701", PublisherID: "T001", Description: "First Badge", Validity: 3},
		{ID: "B032702", PublisherID: "T002", Description: "Second Badge", Validity: 3},
		{ID: "B032703", PublisherID: "T003", Description: "Third Badge", Validity: 3},
		{ID: "B032704", PublisherID: "T001", Description: "Forth Badge", Validity: 3},
		{ID: "B032705", PublisherID: "T002", Description: "Fifth Badge", Validity: 3},
		{ID: "B032706", PublisherID: "T003", Description: "Sixth Badge", Validity: 3},
		{ID: "B032707", PublisherID: "T001", Description: "Seventh Badge", Validity: 3},
	}

	for _, badge := range badges {
		badgeJSON, err := json.Marshal(badge)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(badge.ID, badgeJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateBadge 创建徽章
func (s *SmartContract) CreateBadge(ctx contractapi.TransactionContextInterface, id string, pid string, description string, validity int) error {
	exists, err := s.BadgeExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	badge := Badge{
		ID:          id,
		PublisherID: pid,
		Description: description,
		Validity:    validity,
	}
	badgeJSON, err := json.Marshal(badge)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, badgeJSON)
}

//ReadBadge 读取徽章
func (s *SmartContract) ReadBadge(ctx contractapi.TransactionContextInterface, id string) (*Badge, error) {
	badgeJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if badgeJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var badge Badge
	err = json.Unmarshal(badgeJSON, &badge)
	if err != nil {
		return nil, err
	}

	return &badge, nil
}

//UpdateBadge 更新徽章
func (s *SmartContract) UpdateBadge(ctx contractapi.TransactionContextInterface, id string, pid string, description string, validity int) error {
	exists, err := s.BadgeExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the badge %s does not exist", id)
	}

	// overwriting original asset with new asset
	badge := Badge{
		ID:          id,
		PublisherID: pid,
		Description: description,
		Validity:    validity,
	}
	badgeJSON, err := json.Marshal(badge)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(id, badgeJSON)
}

// DeleteBadge 删除徽章
func (s *SmartContract) DeleteBadge(ctx contractapi.TransactionContextInterface, id string) error {
	exists, err := s.BadgeExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the badge %s does not exist", id)
	}

	return ctx.GetStub().DelState(id)
}

// BadgeExists 检测徽章是否存在
func (s *SmartContract) BadgeExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	badgeJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return badgeJSON != nil, nil
}

// GetAllBadge returns all assets found in world state
func (s *SmartContract) GetAllBadge(ctx contractapi.TransactionContextInterface) ([]*Badge, error) {

	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var badges []*Badge
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}

		var badge Badge
		err = json.Unmarshal(queryResponse.Value, &badge)
		if err != nil {
			return nil, err
		}
		badges = append(badges, &badge)
	}

	return badges, nil
}
