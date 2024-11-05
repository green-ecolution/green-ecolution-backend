package flowerbed

import (
	"context"
	"testing"

	"github.com/green-ecolution/green-ecolution-backend/internal/entities"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestFlowerbedRepository_Update(t *testing.T) {
	t.Run("should update flowerbed", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.Update(context.Background(), 1, WithDescription("updated"))

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "updated", got.Description)
	})

	t.Run("should return error when update flowerbed with non-existing id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.Update(context.Background(), 99, WithDescription("updated"))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error when update flowerbed with negative id", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.Update(context.Background(), -1, WithDescription("updated"))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.Update(ctx, 1, WithDescription("updated"))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should not update flowerbed when no changes", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)
		gotBefore, err := r.GetByID(context.Background(), 1)
		assert.NoError(t, err)

		// when
		got, err := r.Update(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, gotBefore, got)
	})
}

func TestFlowerbedRepository_UpdateWithImages(t *testing.T) {
	images := []*entities.Image{
		{
			ID:       1,
			URL:      "/test/url/to/image",
			Filename: utils.P("Screenshot"),
			MimeType: utils.P("png"),
		},
	}

	t.Run("should update flowerbed and update images", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.UpdateWithImages(
			context.Background(), 
			1, 
			WithDescription("Updated description"),
			WithImagesIDs([]int32{1}),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "Updated description", got.Description)

		linkedImages, err := suite.Store.GetAllImagesByFlowerbedID(context.Background(), got.ID)
		assert.NoError(t, err)
		for i, fb := range linkedImages {
			assert.Equal(t, images[i].ID, fb.ID)
			assert.Equal(t, images[i].Filename, fb.Filename)
			assert.Equal(t, images[i].MimeType, fb.MimeType)
			assert.Equal(t, images[i].URL, fb.Url)
		}
	})

	t.Run("should update flowerbed and link new images", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.UpdateWithImages(
			context.Background(), 
			2, 
			WithDescription("Updated description"),
			WithImagesIDs([]int32{1}),
		)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "Updated description", got.Description)

		linkedImages, err := suite.Store.GetAllImagesByFlowerbedID(context.Background(), got.ID)
		assert.NoError(t, err)
		for i, fb := range linkedImages {
			assert.Equal(t, images[i].ID, fb.ID)
			assert.Equal(t, images[i].Filename, fb.Filename)
			assert.Equal(t, images[i].MimeType, fb.MimeType)
			assert.Equal(t, images[i].URL, fb.Url)
		}
	})

	t.Run("should update flowerbed with no new images", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.UpdateWithImages(context.Background(), 2, WithDescription("Updated description"))

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, "Updated description", got.Description)
		assert.Empty(t, got.Images)
	})

	t.Run("should return error when updating non-existing flowerbed", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.UpdateWithImages(context.Background(), 99, WithDescription("Updated description"))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should only update other values when images parameter is nil", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)

		// when
		got, err := r.UpdateWithImages(context.Background(), 2, WithImages(nil))

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)

		linkedImages, err := suite.Store.GetAllImagesByFlowerbedID(context.Background(), got.ID)
		assert.NoError(t, err)
		assert.Empty(t, linkedImages)
	})

	t.Run("should return error if context is canceled", func(t *testing.T) {
		// given
		r := NewFlowerbedRepository(suite.Store, mappers)
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		// when
		got, err := r.UpdateWithImages(ctx, 1, WithDescription("Updated description"))

		// then
		assert.Error(t, err)
		assert.Nil(t, got)
	})

	t.Run("should not update flowerbed when no changes are made", func(t *testing.T) {
		// given
		suite.ResetDB(t)
		suite.InsertSeed(t, "internal/storage/postgres/seed/test/flowerbed")
		r := NewFlowerbedRepository(suite.Store, mappers)
		gotBefore, err := r.GetByID(context.Background(), 1)
		assert.NoError(t, err)

		// when
		got, err := r.UpdateWithImages(context.Background(), 1)

		// then
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, gotBefore, got)
	})
}
