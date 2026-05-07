package main

import (
	"database/sql"
	"math/rand"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	randSource = rand.NewSource(time.Now().UnixNano())
	randRange  = rand.New(randSource)
)

// openTestDB использует отдельную SQLite БД на каждый тест.
func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	dsn := filepath.Join(t.TempDir(), "test.db")
	db, err := sql.Open("sqlite", dsn)
	require.NoError(t, err)
	require.NoError(t, Migrate(db))
	t.Cleanup(func() { _ = db.Close() })
	return db
}

func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

func TestAddGetDelete(t *testing.T) {
	db := openTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	var err error
	parcel.Number, err = store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	stored, err := store.Get(parcel.Number)

	require.NoError(t, err)
	require.Equal(t, parcel, stored)

	err = store.Delete(parcel.Number)

	stored, err = store.Get(parcel.Number)
	require.Equal(t, sql.ErrNoRows, err)
}

func TestSetAddress(t *testing.T) {
	db := openTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	var err error
	parcel.Number, err = store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	newAddress := "new test address"
	err = store.SetAddress(parcel.Number, newAddress)

	require.NoError(t, err)

	stored, err := store.Get(parcel.Number)

	require.NoError(t, err)
	require.Equal(t, newAddress, stored.Address)
}

func TestSetStatus(t *testing.T) {
	db := openTestDB(t)
	store := NewParcelStore(db)
	parcel := getTestParcel()

	var err error
	parcel.Number, err = store.Add(parcel)

	require.NoError(t, err)
	require.NotEmpty(t, parcel.Number)

	err = store.SetStatus(parcel.Number, ParcelStatusSent)

	require.NoError(t, err)

	stored, err := store.Get(parcel.Number)

	require.NoError(t, err)
	require.Equal(t, ParcelStatusSent, stored.Status)
}

func TestGetByClient(t *testing.T) {
	db := openTestDB(t)
	store := NewParcelStore(db)

	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])

		require.NoError(t, err)
		require.NotEmpty(t, id)

		parcels[i].Number = id
		parcelMap[id] = parcels[i]
	}

	storedParcels, err := store.GetByClient(client)

	require.NoError(t, err)
	require.Len(t, storedParcels, len(parcels))

	for _, parcel := range storedParcels {
		expectedParcel, ok := parcelMap[parcel.Number]

		require.True(t, ok)
		require.Equal(t, expectedParcel, parcel)
	}
}
