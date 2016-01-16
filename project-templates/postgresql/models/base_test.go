package models

import (
	"fmt"
	"os"
	"$GO_BOOTSTRAP_REPO_NAME/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_PROJECT_NAME/libstring"
	"$GO_BOOTSTRAP_REPO_NAME/$GO_BOOTSTRAP_REPO_USER/$GO_BOOTSTRAP_PROJECT_NAME/libunix"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"testing"
)

func newEmailForTest() string {
	return fmt.Sprintf("user-%v@example.com", libstring.RandString(32))
}

func newDbForTest(t *testing.T) *sqlx.DB {
	var err error
	pguser, _, pghost, pgport, pgsslmode := os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGHOST"), os.Getenv("PGPORT"), os.Getenv("PGSSLMODE")
	if pguser == ""{
		pguser, err = libunix.CurrentUser()
		if err != nil {
			t.Fatalf("Getting current user should never fail. Error: %v", err)
		}
	}

	if pghost == "" {
		pghost = "localhost"
	}

	if pgport == "" {
		pgport = "5432"
	}

	if pgsslmode == "" {
		pgsslmode = "disable"
	}

	db, err := sqlx.Connect("postgres", fmt.Sprintf("postgres://%v@%v:%v/$GO_BOOTSTRAP_PROJECT_NAME-test?sslmode=%v", pguser, pghost, pgport, pgsslmode))
	if err != nil {
		t.Fatalf("Connecting to local postgres should never fail. Error: %v", err)
	}
	return db
}

func newBaseForTest(t *testing.T) *Base {
	base := &Base{}
	base.db = newDbForTest(t)

	return base
}

func TestNewTransactionIfNeeded(t *testing.T) {
	base := newBaseForTest(t)

	// New Transaction block
	tx, wrapInSingleTransaction, err := base.newTransactionIfNeeded(nil)
	if err != nil {
		t.Fatalf("Creating new transaction block should not fail. Error: %v", err)
	}
	if wrapInSingleTransaction != true {
		t.Fatalf("Creating new transaction block should set wrapInSingleTransaction == true.")
	}
	if tx == nil {
		t.Fatalf("Creating new transaction block should not fail. Error: %v", err)
	}

	// Existing Transaction block
	tx2, wrapInSingleTransaction, err := base.newTransactionIfNeeded(tx)
	if err != nil {
		t.Fatalf("Receiving existing transaction block should not fail. Error: %v", err)
	}
	if wrapInSingleTransaction != false {
		t.Fatalf("Receiving existing transaction block should set wrapInSingleTransaction == false.")
	}
	if tx2 == nil {
		t.Fatalf("Receiving existing transaction block should not fail. Error: %v", err)
	}
	if tx2 != tx {
		t.Fatalf("Receiving existing transaction block should not fail. Error: %v", err)
	}
}

func TestCreateDeleteGeneric(t *testing.T) {
	base := newBaseForTest(t)
	base.table = "users"

	// INSERT INTO users (name) VALUES (...)
	data := make(map[string]interface{})
	data["email"] = newEmailForTest()
	data["password"] = "abc123"

	result, err := base.InsertIntoTable(nil, data)
	if err != nil {
		t.Fatalf("Inserting new user should not fail. Error: %v", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Inserting new user should not fail. Error: %v", err)
	}

	// DELETE FROM users WHERE id=...
	where := fmt.Sprintf("id=%v", lastInsertedId)

	_, err = base.DeleteFromTable(nil, where)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}

func TestCreateDeleteById(t *testing.T) {
	base := newBaseForTest(t)
	base.table = "users"

	// INSERT INTO users (...) VALUES (...)
	data := make(map[string]interface{})
	data["email"] = newEmailForTest()
	data["password"] = "abc123"

	result, err := base.InsertIntoTable(nil, data)
	if err != nil {
		t.Fatalf("Inserting new user should not fail. Error: %v", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Inserting new user should not fail. Error: %v", err)
	}

	// DELETE FROM users WHERE id=...
	_, err = base.DeleteById(nil, lastInsertedId)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}

func TestCreateUpdateGenericDelete(t *testing.T) {
	base := newBaseForTest(t)
	base.table = "users"

	// INSERT INTO users (...) VALUES (...)
	data := make(map[string]interface{})
	data["email"] = newEmailForTest()
	data["password"] = "abc123"

	result, err := base.InsertIntoTable(nil, data)
	if err != nil {
		t.Fatalf("Inserting new user should not fail. Error: %v", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Inserting new user should not fail. Error: %v", err)
	}

	// UPDATE users SET email=$1 WHERE id=$2
	data["email"] = newEmailForTest()
	where := fmt.Sprintf("id=%v", lastInsertedId)

	_, err = base.UpdateFromTable(nil, data, where)
	if err != nil {
		t.Errorf("Updating existing user should not fail. Error: %v", err)
	}

	// DELETE FROM users WHERE id=...
	_, err = base.DeleteById(nil, lastInsertedId)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}

func TestCreateUpdateByIDDelete(t *testing.T) {
	base := newBaseForTest(t)
	base.table = "users"

	// INSERT INTO users (...) VALUES (...)
	data := make(map[string]interface{})
	data["email"] = newEmailForTest()
	data["password"] = "abc123"

	result, err := base.InsertIntoTable(nil, data)
	if err != nil {
		t.Fatalf("Inserting new user should not fail. Error: %v", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Inserting new user should not fail. Error: %v", err)
	}

	// UPDATE users SET name=$1 WHERE id=$2
	data["email"] = newEmailForTest()

	_, err = base.UpdateByID(nil, data, lastInsertedId)
	if err != nil {
		t.Errorf("Updating existing user should not fail. Error: %v", err)
	}

	// DELETE FROM users WHERE id=...
	_, err = base.DeleteById(nil, lastInsertedId)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}

func TestCreateUpdateByKeyValueStringDelete(t *testing.T) {
	base := newBaseForTest(t)
	base.table = "users"

	originalEmail := newEmailForTest()

	// INSERT INTO users (...) VALUES (...)
	data := make(map[string]interface{})
	data["email"] = newEmailForTest()
	data["password"] = originalEmail

	result, err := base.InsertIntoTable(nil, data)
	if err != nil {
		t.Fatalf("Inserting new user should not fail. Error: %v", err)
	}

	lastInsertedId, err := result.LastInsertId()
	if err != nil {
		t.Fatalf("Inserting new user should not fail. Error: %v", err)
	}

	// UPDATE users SET email=$1 WHERE id=$2
	data["email"] = newEmailForTest()

	_, err = base.UpdateByKeyValueString(nil, data, "email", originalEmail)
	if err != nil {
		t.Errorf("Updating existing user should not fail. Error: %v", err)
	}

	// DELETE FROM users WHERE id=...
	_, err = base.DeleteById(nil, lastInsertedId)
	if err != nil {
		t.Fatalf("Deleting user by id should not fail. Error: %v", err)
	}

}
