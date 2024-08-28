package image

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image/mapper"
	imgMapper "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/image/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/test"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

var (
	querier sqlc.Querier
)

func TestMain(m *testing.M) {
	rootDir := utils.RootDir()
	seedPath := fmt.Sprintf("%s/internal/storage/postgres/test/seed/image", rootDir)
	close, db, err := test.SetupPostgresContainer(seedPath)
	if err != nil {
		slog.Error("Error setting up postgres container", "error", err)
		panic(err)
	}
	defer close()

	querier = sqlc.New(db)

	os.Exit(m.Run())
}

func TestImageRepository(t *testing.T) {
	t.Run("GetAll", func(t *testing.T) {
		// given
		repo := NewImageRepository(querier, NewImageRepositoryMappers(&imgMapper.InternalImageRepoMapperImpl{}))
		expected := []*entities.Image{
			{
				ID:       1,
				URL:      "https://avatars.githubusercontent.com/u/165842746?s=96&v=4",
				Filename: nil,
				MimeType: nil,
			},
			{
				ID:       2,
				URL:      "https://app.dev.green-ecolution.de/api/v1/images/avatar.png",
				Filename: utils.P("avatar.png"),
				MimeType: utils.P("image/png"),
			},
		}

		// when
		actual, err := repo.GetAll(context.Background())

		// then
		assert.NoError(t, err)
		assert.Len(t, actual, len(expected))

		for i, actual := range actual {
			assert.NotNil(t, actual.CreatedAt)
			assert.NotNil(t, actual.UpdatedAt)

			assert.EqualValues(t, expected[i].ID, actual.ID)
			assert.EqualValues(t, expected[i].URL, actual.URL)
			assert.EqualValues(t, expected[i].Filename, actual.Filename)
			assert.EqualValues(t, expected[i].MimeType, actual.MimeType)
		}
	})

	t.Run("GetByID", func(t *testing.T) {
		// given
		repo := NewImageRepository(querier, NewImageRepositoryMappers(&imgMapper.InternalImageRepoMapperImpl{}))
		expectedID1 := entities.Image{
			ID:       1,
			URL:      "https://avatars.githubusercontent.com/u/165842746?s=96&v=4",
			Filename: nil,
			MimeType: nil,
		}
		expectedID2 := entities.Image{
			ID:       2,
			URL:      "https://app.dev.green-ecolution.de/api/v1/images/avatar.png",
			Filename: utils.P("avatar.png"),
			MimeType: utils.P("image/png"),
		}

		// when
		actualID1, err1 := repo.GetByID(context.Background(), 1)
		actualID2, err2 := repo.GetByID(context.Background(), 2)

		// then
		assert.NoError(t, err1)
		assert.NotNil(t, actualID1)
		assert.NotNil(t, actualID1.CreatedAt)
		assert.NotNil(t, actualID1.UpdatedAt)
		assert.EqualValues(t, expectedID1.ID, actualID1.ID)
		assert.EqualValues(t, expectedID1.URL, actualID1.URL)
		assert.EqualValues(t, expectedID1.Filename, actualID1.Filename)
		assert.EqualValues(t, expectedID1.MimeType, actualID1.MimeType)

		assert.NoError(t, err2)
		assert.NotNil(t, actualID2)
		assert.NotNil(t, actualID2.CreatedAt)
		assert.NotNil(t, actualID2.UpdatedAt)
		assert.EqualValues(t, expectedID2.ID, actualID2.ID)
		assert.EqualValues(t, expectedID2.URL, actualID2.URL)
		assert.EqualValues(t, expectedID2.Filename, actualID2.Filename)
		assert.EqualValues(t, expectedID2.MimeType, actualID2.MimeType)
	})

	t.Run("Create", func(t *testing.T) {
		// given
		repo := NewImageRepository(querier, NewImageRepositoryMappers(&imgMapper.InternalImageRepoMapperImpl{}))
		newImage := &entities.Image{
			URL:      "http://example.com/image.jpg",
			Filename: utils.P("image.jpg"),
			MimeType: utils.P("image/jpeg"),
		}

		// when
		actual, err := repo.Create(context.Background(), newImage)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.NotNil(t, actual.CreatedAt)
		assert.NotNil(t, actual.UpdatedAt)
		assert.EqualValues(t, actual.ID, 3)
		assert.EqualValues(t, newImage.URL, actual.URL)
		assert.EqualValues(t, newImage.Filename, actual.Filename)
		assert.EqualValues(t, newImage.MimeType, actual.MimeType)
	})

	t.Run("create with null values for filename and mimetype", func(t *testing.T) {
		// given
		repo := NewImageRepository(querier, NewImageRepositoryMappers(&imgMapper.InternalImageRepoMapperImpl{}))
		newImage := &entities.Image{
			URL: "http://example.com/image.jpg",
		}

		// when
		actual, err := repo.Create(context.Background(), newImage)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.NotNil(t, actual.CreatedAt)
		assert.NotNil(t, actual.UpdatedAt)
		assert.EqualValues(t, actual.ID, 4)
		assert.EqualValues(t, newImage.URL, actual.URL)
		assert.EqualValues(t, newImage.Filename, actual.Filename)
		assert.EqualValues(t, newImage.MimeType, actual.MimeType)
	})

	t.Run("Update", func(t *testing.T) {
		// given
		repo := NewImageRepository(querier, NewImageRepositoryMappers(&imgMapper.InternalImageRepoMapperImpl{}))
		image := &entities.Image{
			ID:       1,
			URL:      "http://example.com/image.jpg",
			Filename: utils.P("image.jpg"),
			MimeType: utils.P("image/jpeg"),
		}

		// when
		actual, err := repo.Update(context.Background(), image)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.NotNil(t, actual.CreatedAt)
		assert.NotNil(t, actual.UpdatedAt)
		assert.EqualValues(t, actual.ID, 1)
		assert.EqualValues(t, actual.URL, image.URL)
		assert.EqualValues(t, actual.Filename, image.Filename)
		assert.EqualValues(t, actual.MimeType, image.MimeType)
	})

	t.Run("Delete", func(t *testing.T) {
		// given
		repo := NewImageRepository(querier, NewImageRepositoryMappers(&imgMapper.InternalImageRepoMapperImpl{}))

		// when
		err := repo.Delete(context.Background(), 1)

		// then
		assert.NoError(t, err)

		// check if image was deleted
		_, err = repo.GetByID(context.Background(), 1)
		assert.Error(t, err)
	})
}

func TestImageRepositoryErrors(t *testing.T) {
	t.Run("call get by id with not existing id should return error", func(t *testing.T) {
		// given
		repo := NewImageRepository(querier, NewImageRepositoryMappers(&imgMapper.InternalImageRepoMapperImpl{}))

		// when
		_, err := repo.GetByID(context.Background(), 999)

		// then
		assert.Error(t, err)
	})

	t.Run("call update with not existing id should return error", func(t *testing.T) {
		// given
		repo := NewImageRepository(querier, NewImageRepositoryMappers(&imgMapper.InternalImageRepoMapperImpl{}))
		image := &entities.Image{
			ID:       999,
			URL:      "http://example.com/image.jpg",
			Filename: utils.P("image.jpg"),
			MimeType: utils.P("image/jpeg"),
		}

		// when
		_, err := repo.Update(context.Background(), image)

		// then
		assert.Error(t, err)
	})

	t.Run("call delete with not existing id should not return error", func(t *testing.T) {
		// given
		repo := NewImageRepository(querier, NewImageRepositoryMappers(&imgMapper.InternalImageRepoMapperImpl{}))

		// when
		err := repo.Delete(context.Background(), 999)

		// then
		assert.NoError(t, err)
	})
}

func TestNewImageRepositoryMappers(t *testing.T) {
	type args struct {
		iMapper mapper.InternalImageRepoMapper
	}
	tests := []struct {
		name string
		args args
		want ImageRepositoryMappers
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewImageRepositoryMappers(tt.args.iMapper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewImageRepositoryMappers() = %v, want %v", got, tt.want)
			}
		})
	}
}
