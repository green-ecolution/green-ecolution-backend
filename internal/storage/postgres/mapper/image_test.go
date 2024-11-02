package mapper_test

import (
	"testing"
	"time"

	sqlc "github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/_sqlc"
	"github.com/green-ecolution/green-ecolution-backend/internal/storage/postgres/mapper/generated"
	"github.com/green-ecolution/green-ecolution-backend/internal/utils"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestImageMapper_FromSql(t *testing.T) {
	imageMapper := &generated.InternalImageRepoMapperImpl{}

	t.Run("should convert from sql to entity", func(t *testing.T) {
		// given
		src := allTestImages[0]

		// when
		got := imageMapper.FromSql(src)

		// then
		assert.NotNil(t, got)
		assert.Equal(t, src.ID, got.ID)
		assert.Equal(t, src.CreatedAt.Time, got.CreatedAt)
		assert.Equal(t, src.UpdatedAt.Time, got.UpdatedAt)
		assert.Equal(t, src.Url, got.URL)
		assert.Equal(t, src.Filename, got.Filename)
		assert.Equal(t, src.MimeType, got.MimeType)
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src *sqlc.Image = nil

		// when
		got := imageMapper.FromSql(src)

		// then
		assert.Nil(t, got)
	})
}

func TestImageMapper_FromSqlList(t *testing.T) {
	imageMapper := &generated.InternalImageRepoMapperImpl{}

	t.Run("should convert from sql slice to entity slice", func(t *testing.T) {
		// given
		src := allTestImages

		// when
		got := imageMapper.FromSqlList(src)

		// then
		assert.NotNil(t, got)
		assert.Len(t, got, 2)

		for i, src := range src {
			assert.NotNil(t, got)
			assert.Equal(t, src.ID, got[i].ID)
			assert.Equal(t, src.CreatedAt.Time, got[i].CreatedAt)
			assert.Equal(t, src.UpdatedAt.Time, got[i].UpdatedAt)
			assert.Equal(t, src.Url, got[i].URL)
			assert.Equal(t, src.Filename, got[i].Filename)
			assert.Equal(t, src.MimeType, got[i].MimeType)
		}
	})

	t.Run("should return nil for nil input", func(t *testing.T) {
		// given
		var src []*sqlc.Image = nil

		// when
		got := imageMapper.FromSqlList(src)

		// then
		assert.Nil(t, got)
	})
}

var allTestImages = []*sqlc.Image{
	{
		ID:             1,
		CreatedAt:      pgtype.Timestamp{Time: time.Now()},
		UpdatedAt:      pgtype.Timestamp{Time: time.Now()},
		Url:           	"/test/path/to/image",
		Filename:    	utils.StringPointer("Screenshot"),
		MimeType: 		utils.StringPointer("pdf"),
	},
	{
		ID:             2,
		CreatedAt:      pgtype.Timestamp{Time: time.Now()},
		UpdatedAt:      pgtype.Timestamp{Time: time.Now()},
		Url:           	"/test/path/to/image",
		Filename:    	utils.StringPointer("Screenshot 02"),
		MimeType: 		utils.StringPointer("png"),
	},
}