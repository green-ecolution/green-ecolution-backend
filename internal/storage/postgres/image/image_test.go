package image

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage"
	mapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/store"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/testutils"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

type RandomImage struct {
	ID        int32     `faker:"-"`
	CreatedAt time.Time `faker:"-"`
	UpdatedAt time.Time `faker:"-"`
	URL       string    `faker:"url"`
	Filename  *string   `faker:"word"`
	MimeType  *string   `faker:"oneof:image/png,image/jpeg"`
}

func TestMain(m *testing.M) {
	closeCon, _, err := testutils.SetupPostgresContainer()
	if err != nil {
		slog.Error("Error setting up postgres container", "error", err)
		os.Exit(1)
	}
	defer closeCon()

	os.Exit(m.Run())
}

func createStore(db *pgx.Conn) *store.Store {
	return store.NewStore(db)
}

func initMapper() ImageRepositoryMappers {
	return NewImageRepositoryMappers(&mapper.InternalImageRepoMapperImpl{})
}

func createImage(t *testing.T, str *store.Store) *entities.Image {
	var img entities.Image
	if err := faker.FakeData(&img); err != nil {
		t.Fatalf("error faking image data: %v", err)
	}
	mappers := initMapper()
	repo := NewImageRepository(str, mappers)

	got, err := repo.Create(context.Background(),
		WithURL(img.URL),
		WithFilename(img.Filename),
		WithMimeType(img.MimeType),
	)
	assert.NoError(t, err)

	assert.NotNil(t, got)
	assert.NotZero(t, got.ID)
	assert.Equal(t, img.URL, got.URL)
	assert.Equal(t, img.Filename, got.Filename)
	assert.Equal(t, img.MimeType, got.MimeType)
	assert.NotZero(t, got.CreatedAt)
	assert.NotZero(t, got.UpdatedAt)

	return got
}

func TestCreateImage(t *testing.T) {
	t.Parallel()
	t.Run("should create image", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			createImage(t, str)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Create(context.Background(),
				WithURL("https://example.com"),
			)
			assert.Error(t, err)
		})
	})
}

func TestGetAllImages(t *testing.T) {
	t.Parallel()
	t.Run("should get all images", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			createImage(t, str)
			createImage(t, str)
			createImage(t, str)

			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)

			assert.Len(t, got, 3)
		})
	})

	t.Run("should return empty list if no images found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			got, err := repo.GetAll(context.Background())
			assert.NoError(t, err)

			assert.Len(t, got, 0)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetAll(context.Background())
			assert.Error(t, err)
		})
	})
}

func TestGetImageByID(t *testing.T) {
	t.Parallel()
	t.Run("should get image by id", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			img := createImage(t, str)
			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			got, err := repo.GetByID(context.Background(), img.ID)
			assert.NoError(t, err)

			assert.NotNil(t, got)
			assertImage(t, img, got)
		})
	})

	t.Run("should return error if image not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			_, err := repo.GetByID(context.Background(), 999)
			assert.Error(t, err)
			assert.ErrorIs(t, err, storage.ErrImageNotFound)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func TestUpdateImage(t *testing.T) {
	t.Parallel()
	t.Run("should update image", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			prev := createImage(t, str)
			want := &entities.Image{
				ID:        prev.ID,
				URL:       "https://example.com",
				Filename:  utils.P("new-filename"),
				MimeType:  utils.P("image/jpeg"),
				CreatedAt: prev.CreatedAt,
			}

			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			got, err := repo.Update(context.Background(), prev.ID,
				WithURL("https://example.com"),
				WithFilename(utils.P("new-filename")),
				WithMimeType(utils.P("image/jpeg")),
			)
			assert.NoError(t, err)

			assert.NotNil(t, got)
			assert.NotEqual(t, prev, got)
			assertImage(t, want, got)
		})
	})

	t.Run("should only update filled image fields", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			prev := createImage(t, str)
			want := &entities.Image{
				ID:        prev.ID,
				URL:       prev.URL,
				Filename:  utils.P("new-filename"),
				MimeType:  prev.MimeType,
				CreatedAt: prev.CreatedAt,
			}

			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			got, err := repo.Update(context.Background(), prev.ID,
				WithFilename(utils.P("new-filename")),
			)
			assert.NoError(t, err)

			assert.NotNil(t, got)
			assert.NotEqual(t, prev, got)
			assert.Equal(t, want.Filename, got.Filename)
			assert.Equal(t, want.MimeType, got.MimeType)
			assert.Equal(t, prev.URL, got.URL)
		})
	})

	t.Run("should return error if image not found", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			img := createImage(t, str)
			img.ID = 999
			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			_, err := repo.Update(context.Background(), img.ID)
			assert.Error(t, err)
			assert.ErrorIs(t, err, storage.ErrImageNotFound)
		})
	})

	t.Run("should not update if all fields are nil", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			prev := createImage(t, str)
			want := &entities.Image{
				ID:        prev.ID,
				URL:       prev.URL,
				Filename:  prev.Filename,
				MimeType:  prev.MimeType,
				CreatedAt: prev.CreatedAt,
			}

			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			got, err := repo.Update(context.Background(), prev.ID)
			assert.NoError(t, err)
			assertImage(t, want, got)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewImageRepository(str, mappers)
      img := createImage(t, str)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			_, err = repo.Update(context.Background(), img.ID)
			assert.Error(t, err)
		})
	})
}

func TestDeleteImage(t *testing.T) {
	t.Parallel()
	t.Run("should delete image", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			img := createImage(t, str)

			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			err := repo.Delete(context.Background(), img.ID)
			assert.NoError(t, err)

			_, err = repo.GetByID(context.Background(), img.ID)
			assert.Error(t, err)
		})
	})

	t.Run("should return error if query fails", func(t *testing.T) {
		testutils.WithTx(t, func(db *pgx.Conn) {
			str := createStore(db)
			mappers := initMapper()
			repo := NewImageRepository(str, mappers)

			err := db.Close(context.Background())
			assert.NoError(t, err)

			err = repo.Delete(context.Background(), 1)
			assert.Error(t, err)
		})
	})
}

func assertImage(t *testing.T, want, got *entities.Image) {
	assert.Equal(t, want.ID, got.ID)
	assert.Equal(t, want.URL, got.URL)
	assert.Equal(t, want.Filename, got.Filename)
	assert.Equal(t, want.MimeType, got.MimeType)
	assert.Equal(t, want.CreatedAt, got.CreatedAt)
	assert.NotZero(t, got.UpdatedAt)
}
