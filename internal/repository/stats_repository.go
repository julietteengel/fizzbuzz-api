package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/julietteengel/fizzbuzz-api/internal/config"
	"github.com/julietteengel/fizzbuzz-api/internal/model"
)

type IStatsRepository interface {
	RecordRequest(ctx context.Context, request model.FizzBuzzRequest) error
	GetMostFrequent(ctx context.Context) (*model.StatsResponse, error)
}

type statsRepository struct {
	db        *gorm.DB
	memStats  map[string]*model.StatsEntry // ⚠️ PROBLÈME: Map grandit indéfiniment = fuite mémoire
	memMutex  sync.RWMutex
	useMemory bool

	// AMÉLIORATION: Ajouter ces champs pour éviter les fuites mémoire
	// maxEntries int                    // Limite max d'entrées (ex: 10000)
	// cleanupTicker *time.Ticker        // Nettoyage périodique des anciennes entrées
	// entryTTL time.Duration            // TTL pour expirer les entrées (ex: 24h)
}

func NewStatsRepository(database *gorm.DB, cfg *config.Config) IStatsRepository {
	useMemory := cfg.Database.StatsStorage == "memory"
	return &statsRepository{
		db:        database,
		memStats:  make(map[string]*model.StatsEntry),
		useMemory: useMemory,

		// AMÉLIORATION: Initialiser la protection contre les fuites mémoire
		// maxEntries:    10000,                    // Limite à 10k entrées
		// entryTTL:      24 * time.Hour,           // Expirer après 24h
		// cleanupTicker: time.NewTicker(1 * time.Hour), // Cleanup toutes les heures
	}

	// AMÉLIORATION: Démarrer le nettoyage périodique
	// if useMemory {
	//     go repo.startPeriodicCleanup()
	// }
}

func (r *statsRepository) RecordRequest(ctx context.Context, request model.FizzBuzzRequest) error {
	if r.useMemory {
		return r.recordInMemory(request)
	}
	return r.recordInDatabase(ctx, request)
}

func (r *statsRepository) GetMostFrequent(ctx context.Context) (*model.StatsResponse, error) {
	if r.useMemory {
		return r.getMostFrequentFromMemory()
	}
	return r.getMostFrequentFromDatabase(ctx)
}

func (r *statsRepository) recordInMemory(request model.FizzBuzzRequest) error {
	r.memMutex.Lock()         //Exclusif, bloque TOUT (lecteurs + écrivains): L'enregistrement des stats bloque temporairement les lectures
	defer r.memMutex.Unlock() // S'exécute automatiquement à la fin, même si une erreur survient, unlock() sera appelé

	// AMÉLIORATION: Vérifier la limite avant d'ajouter une nouvelle entrée
	// if len(r.memStats) >= r.maxEntries {
	//     r.evictOldestEntry() // Supprimer la plus ancienne entrée
	// }

	key := r.generateKey(request)
	if entry, exists := r.memStats[key]; exists {
		entry.HitCount++
		entry.UpdatedAt = time.Now()
	} else {
		// PROBLÈME: Nouvelle entrée sans limite = fuite mémoire potentielle
		r.memStats[key] = &model.StatsEntry{
			Int1:      request.Int1,
			Int2:      request.Int2,
			Limit:     request.Limit,
			Str1:      request.Str1,
			Str2:      request.Str2,
			HitCount:  1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	return nil
}

// 💡 AMÉLIORATION: Fonctions pour éviter les fuites mémoire
// func (r *statsRepository) cleanupExpiredEntries() {
//     r.memMutex.Lock()
//     defer r.memMutex.Unlock()
//
//     now := time.Now()
//     for key, entry := range r.memStats {
//         if now.Sub(entry.UpdatedAt) > r.entryTTL {
//             delete(r.memStats, key) // Supprimer les entrées expirées
//         }
//     }
// }

func (r *statsRepository) getMostFrequentFromMemory() (*model.StatsResponse, error) {
	r.memMutex.RLock() //Partagé entre lecteurs, mais bloqué par écrivains, plusieurs utilisateurs peuvent consulter /stats en même temps
	defer r.memMutex.RUnlock()

	var mostFrequent *model.StatsEntry
	for _, entry := range r.memStats {
		if mostFrequent == nil || entry.HitCount > mostFrequent.HitCount {
			mostFrequent = entry
		}
	}

	if mostFrequent == nil {
		return nil, nil
	}

	return &model.StatsResponse{
		Request: model.FizzBuzzRequest{
			Int1:  mostFrequent.Int1,
			Int2:  mostFrequent.Int2,
			Limit: mostFrequent.Limit,
			Str1:  mostFrequent.Str1,
			Str2:  mostFrequent.Str2,
		},
		HitCount: mostFrequent.HitCount,
	}, nil
}

func (r *statsRepository) recordInDatabase(ctx context.Context, request model.FizzBuzzRequest) error {
	// PB sans transaction: si 2 requêtes simultanées avec les mêmes paramètres int1=3, int2=5, limit=15, str1="fizz", str2="buzz" :
	// Problème : Les deux threads lisent la même ancienne valeur avant que l'autre ait fini sa mise à jour.
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// DÉBUT DE TRANSACTION - PostgreSQL pose un LOCK, donc requete B doit attendre la fin de A donc pas de pb de concurrence
		entry := model.StatsEntry{
			Int1:  request.Int1,
			Int2:  request.Int2,
			Limit: request.Limit,
			Str1:  request.Str1,
			Str2:  request.Str2,
		}

		result := tx.Where(&entry).First(&entry)
		if result.Error == gorm.ErrRecordNotFound {
			entry.HitCount = 1
			return tx.Create(&entry).Error
		}
		if result.Error != nil {
			return result.Error
		}

		entry.HitCount++
		return tx.Save(&entry).Error
	})
}

func (r *statsRepository) getMostFrequentFromDatabase(ctx context.Context) (*model.StatsResponse, error) {
	var entry model.StatsEntry
	result := r.db.WithContext(ctx).Order("hit_count DESC").First(&entry)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}

	return &model.StatsResponse{
		Request: model.FizzBuzzRequest{
			Int1:  entry.Int1,
			Int2:  entry.Int2,
			Limit: entry.Limit,
			Str1:  entry.Str1,
			Str2:  entry.Str2,
		},
		HitCount: entry.HitCount,
	}, nil
}

func (r *statsRepository) generateKey(request model.FizzBuzzRequest) string {
	return fmt.Sprintf("%d_%d_%d_%s_%s", request.Int1, request.Int2, request.Limit, request.Str1, request.Str2)
}
