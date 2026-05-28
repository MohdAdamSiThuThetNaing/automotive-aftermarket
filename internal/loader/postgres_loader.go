package loader

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/cache"
	"github.com/MohdAdamSiThuThetNaing/automotive-aftermarket/internal/models"
)

const (
	loadTranslationsQuery = `
SELECT
	entity_type,
	entity_id,
	locale,
	field_name,
	field_value,
	updated_at
FROM translation
WHERE entity_type = $1
AND entity_id = ANY($2)
AND locale = ANY($3)
`

	maxBatchSize = 1000
)

type PostgresLoader struct {
	db    *pgxpool.Pool
	cache *cache.TranslationCache
}

func NewPostgresLoader(
	db *pgxpool.Pool,
	cache *cache.TranslationCache,
) *PostgresLoader {

	return &PostgresLoader{
		db:    db,
		cache: cache,
	}
}

func (l *PostgresLoader) Load(
	ctx context.Context,
	entityType string,
	entityIDs []string,
	locales []string,
) (map[string][]models.Translation, error) {

	if len(entityIDs) == 0 {
		return map[string][]models.Translation{}, nil
	}

	start := time.Now()
	entityIDs = uniqueStrings(entityIDs)

	result :=
		make(
			map[string][]models.Translation,
			len(entityIDs),
		)

	missingEntityIDs :=
		l.collectCacheHits(
			entityType,
			entityIDs,
			result,
		)

	if len(missingEntityIDs) == 0 {

		log.Printf(
			"translation loader cache_hit_only entity_type=%s entities=%d duration=%s",
			entityType,
			len(entityIDs),
			time.Since(start),
		)

		return result, nil
	}

	batches :=
		chunkStrings(
			missingEntityIDs,
			maxBatchSize,
		)

	for _, batch := range batches {
		loadedTranslations, err :=
			l.loadBatch(
				ctx,
				entityType,
				batch,
				locales,
			)

		if err != nil {
			return nil, err
		}

		l.populateCache(
			loadedTranslations,
		)

		for key, translations :=
			range loadedTranslations {
			result[key] = translations
		}
	}

	log.Printf(
		"translation loader completed entity_type=%s requested=%d missing=%d duration=%s",
		entityType,
		len(entityIDs),
		len(missingEntityIDs),
		time.Since(start),
	)

	return result, nil
}

func (l *PostgresLoader) loadBatch(
	ctx context.Context,
	entityType string,
	entityIDs []string,
	locales []string,
) (map[string][]models.Translation, error) {

	rows, err := l.db.Query(
		ctx,
		loadTranslationsQuery,
		entityType,
		entityIDs,
		locales,
	)

	if err != nil {
		return nil, fmt.Errorf(
			"query translations: %w",
			err,
		)
	}

	defer rows.Close()

	result :=
		make(
			map[string][]models.Translation,
		)

	for rows.Next() {

		var translation models.Translation

		err := rows.Scan(
			&translation.EntityType,
			&translation.EntityID,
			&translation.Locale,
			&translation.FieldName,
			&translation.FieldValue,
			&translation.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf(
				"scan translation: %w",
				err,
			)
		}

		key := buildKey(
			translation.EntityType,
			translation.EntityID,
		)

		result[key] =
			append(
				result[key],
				translation,
			)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf(
			"iterate translations: %w",
			err,
		)
	}

	return result, nil
}

func (l *PostgresLoader) collectCacheHits(
	entityType string,
	entityIDs []string,
	result map[string][]models.Translation,
) []string {

	missing :=
		make(
			[]string,
			0,
			len(entityIDs),
		)

	for _, entityID := range entityIDs {

		key := buildKey(
			entityType,
			entityID,
		)

		translations, found :=
			l.cache.Get(key)

		if found {
			result[key] = translations
			continue
		}

		missing =
			append(
				missing,
				entityID,
			)
	}

	return missing
}

func (l *PostgresLoader) populateCache(
	translations map[string][]models.Translation,
) {

	for key, values :=
		range translations {
		l.cache.Set(
			key,
			values,
		)
	}
}

func (l *PostgresLoader) Invalidate(
	entityType string,
	entityID string,
) {

	key := buildKey(
		entityType,
		entityID,
	)

	l.cache.Invalidate(key)
}

func buildKey(
	entityType string,
	entityID string,
) string {

	return entityType + ":" + entityID
}

func chunkStrings(
	values []string,
	size int,
) [][]string {

	var chunks [][]string

	for size < len(values) {
		values, chunks =
			values[size:],
			append(
				chunks,
				values[0:size:size],
			)
	}

	chunks =
		append(
			chunks,
			values,
		)

	return chunks
}

func uniqueStrings(
	values []string,
) []string {

	seen := make(map[string]struct{})

	result :=
		make(
			[]string,
			0,
			len(values),
		)

	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result =
			append(
				result,
				value,
			)
	}

	return result
}
