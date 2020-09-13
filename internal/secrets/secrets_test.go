package secrets

import (
	"os"
	"testing"

	"github.com/biensupernice/krane/internal/store"
)

const boltpath = "./krane.db"

func teardown() { os.Remove(boltpath) }

func TestMain(m *testing.M) {
	store.New((boltpath))
	defer store.Instance().Shutdown()

	code := m.Run()

	teardown()
	os.Exit(code)
}

func TestAddNewSecret(t *testing.T) {

}

func TestAliasCreation(t *testing.T) {
}
