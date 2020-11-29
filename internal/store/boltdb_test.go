package store

import (
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/docker/distribution/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/biensupernice/krane/internal/constants"
)

const testBoltPath = "./krane.db"

func teardown() {
	os.Remove(testBoltPath)
}

func TestMain(m *testing.M) {
	CreateIfNotExist(testBoltPath)
	defer Instance().Shutdown()

	code := m.Run()

	teardown()
	os.Exit(code)
}

type Avenger struct {
	ID        uuid.UUID
	FirstName string
	LastName  string
	Age       int
	CreatedAt time.Time
}

func TestBoltGet(t *testing.T) {
	bkt := constants.DeploymentsCollectionName

	// Setup
	thor := &Avenger{
		ID:        uuid.Generate(),
		FirstName: "Thor",
		LastName:  "Odinson",
		Age:       43,
		CreatedAt: time.Now().UTC(),
	}

	bytes, err := json.Marshal(thor)
	assert.Nil(t, err)

	// Act
	err = Instance().Put(bkt, thor.ID.String(), bytes)
	assert.Nil(t, err)

	// Assert
	thorBytes, err := Instance().Get(bkt, thor.ID.String())
	assert.Nil(t, err)

	var hero Avenger
	err = json.Unmarshal(thorBytes, &hero)
	assert.Nil(t, err)

	assert.Equal(t, thor.ID, hero.ID)
	assert.Equal(t, thor.FirstName, hero.FirstName)
	assert.Equal(t, thor.LastName, hero.LastName)
	assert.Equal(t, thor.Age, hero.Age)
	assert.True(t, thor.CreatedAt.Equal(hero.CreatedAt))
}

func TestBoltPut(t *testing.T) {
	bkt := constants.DeploymentsCollectionName

	// Setup
	blackwidow := &Avenger{
		ID:        uuid.Generate(),
		FirstName: "Natasha",
		LastName:  "Romanova",
		Age:       23,
		CreatedAt: time.Now().UTC(),
	}

	bytes, err := json.Marshal(blackwidow)
	assert.Nil(t, err)

	// Act
	err = Instance().Put(bkt, blackwidow.ID.String(), bytes)
	assert.Nil(t, err)

	// Assert
	blackwidowBytes, err := Instance().Get(bkt, blackwidow.ID.String())

	var hero Avenger
	err = json.Unmarshal(blackwidowBytes, &hero)
	assert.Nil(t, err)

	assert.Equal(t, blackwidow.ID, hero.ID)
	assert.Equal(t, blackwidow.FirstName, hero.FirstName)
	assert.Equal(t, blackwidow.LastName, hero.LastName)
	assert.Equal(t, blackwidow.Age, hero.Age)
	assert.True(t, blackwidow.CreatedAt.Equal(hero.CreatedAt))
}

func TestBoltGetAll(t *testing.T) {
	bkt := constants.DeploymentsCollectionName

	// Setup
	avengers := make([]Avenger, 0)
	thor := &Avenger{
		ID:        uuid.Generate(),
		FirstName: "Thor",
		LastName:  "Odinson",
		Age:       43,
		CreatedAt: time.Now().UTC(),
	}

	avengers = append(avengers, *&Avenger{
		ID:        uuid.Generate(),
		FirstName: "Tony",
		LastName:  "Stark",
		Age:       30,
		CreatedAt: time.Now().UTC(),
	})

	avengers = append(avengers, *&Avenger{
		ID:        uuid.Generate(),
		FirstName: "Natasha",
		LastName:  "Romanova",
		Age:       23,
		CreatedAt: time.Now().UTC(),
	})

	// Act
	for _, avenger := range avengers {
		bytes, err := json.Marshal(thor)
		assert.Nil(t, err)

		err = Instance().Put(bkt, avenger.ID.String(), bytes)
		assert.Nil(t, err)
	}

	// Assert
	for _, avenger := range avengers {
		bytes, err := Instance().Get(bkt, avenger.ID.String())
		assert.Nil(t, err)

		var hero Avenger
		err = json.Unmarshal(bytes, &hero)
		assert.Nil(t, err)

		assert.NotNil(t, hero)
		assert.NotNil(t, hero.ID)
		assert.NotNil(t, hero.FirstName)
		assert.NotNil(t, hero.LastName)
		assert.NotNil(t, hero.CreatedAt)
	}
}
