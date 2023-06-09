package api

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/betelgeusexru/golang-hotel-reservation/testutils"
	"github.com/betelgeusexru/golang-hotel-reservation/types"
)

func TestPostUser (t *testing.T) {
	tdb := testutils.SetupDatabase(t)
	defer tdb.Teardown(t)

	userHander := NewUserHandler(tdb.UserStore)

	app := testutils.SetupFiberApp()
	app.Post("/", userHander.HandlePostUser)

	params := types.CreateUserParams{
		Email: "some@foo.com",
		FirstName: "James",
		LastName: "Foo",
		Password: "12345678",
	}

	b, _ := json.Marshal(params)

	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	var user types.User
	json.NewDecoder(resp.Body).Decode(&user)

	if len(user.ID) == 0 {
		t.Errorf("expecting a user id to be set")
	}

	if len(user.EncryptedPassword) > 0 {
		t.Errorf("expecting the EncryptedPassword not to be included in the json")
	}

	if user.FirstName != params.FirstName {
		t.Errorf("expected firstname %s but got %s", params.FirstName, user.FirstName)
	}

	if user.LastName != params.LastName {
		t.Errorf("expected lastname %s but got %s", params.LastName, user.LastName)
	}

	if user.Email != params.Email {
		t.Errorf("expected email %s but got %s", params.Email, user.Email)
	}
}