package layer

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// BlockedUser is Layer's representation for a blocked user
type BlockedUser struct {
	UserID string `json:"user_id"`
}

// Block reresents a block operation on Layer
type Block struct {
	Operation string `json:"operation"`
	Property  string `json:"property"`
	Value     string `json:"value"`
}

// AddUserToBlockList adds one member to a user's block list
func (l *Layer) AddUserToBlockList(userID, blocked string) (ok bool, err error) {
	b := BlockedUser{UserID: blocked}
	body, err := json.Marshal(&b)
	if err != nil {
		return false, err
	}

	p := Parameters{Path: fmt.Sprintf("users/%s/blocks", userID), Body: body}
	resp, err := l.request("POST", &p)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return false, fmt.Errorf("Responded with Error Code %d", resp.StatusCode)
	}
	return true, nil
}

// GetUserBlockList Returns an array of all blocked users for the specified
func (l *Layer) GetUserBlockList(userID string) ([]BlockedUser, error) {
	p := Parameters{Path: fmt.Sprintf("users/%s/blocks", userID)}
	resp, err := l.request("GET", &p)
	if err != nil {
		return []BlockedUser{}, err
	}
	defer resp.Body.Close()

	b := []BlockedUser{}
	json.NewDecoder(resp.Body).Decode(&b)
	return b, nil
}

// UnblockUser Removes a blocked user from the Block List of the specified user
func (l *Layer) UnblockUser(userID, blockID string) (ok bool, err error) {
	p := Parameters{Path: fmt.Sprintf("users/%s/blocks/%s", userID, blockID)}
	resp, err := l.request("DELETE", &p)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 204 {
		return false, fmt.Errorf("Responded with Error Code %d", resp.StatusCode)
	}
	return true, nil
}

// BulkModifyBlockList supports bulk operations on a user's blocklist
func (l *Layer) BulkModifyBlockList(userID string, blockIDs, unBlockIDs []string) (ok bool, err error) {

	blocks := make(chan []Block)
	go buildBulkBlockOperation("add", blockIDs, blocks)
	go buildBulkBlockOperation("remove", unBlockIDs, blocks)

	b1, b2 := <-blocks, <-blocks

	b1 = append(b1, b2...)

	body, err := json.Marshal(&b1)
	if err != nil {
		return false, err
	}

	p := Parameters{Path: fmt.Sprintf("users/%s", userID), Body: body}
	resp, err := l.request("PATCH", &p)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return false, fmt.Errorf("Responded with Error Code %d", resp.StatusCode)
	}
	return true, nil
}

func buildBulkBlockOperation(op string, ids []string, ch chan<- []Block) {
	resp := []Block{}
	for _, v := range ids {
		resp = append(resp, Block{Operation: op, Property: "blocks", Value: v})
	}

	ch <- resp
}
