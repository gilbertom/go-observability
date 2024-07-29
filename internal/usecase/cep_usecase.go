package usecase

import (
	"context"

	"github.com/gilbertom/go-temperatura-cep/internal/entity"
	"go.opentelemetry.io/otel/trace"
)

// CepUsecase represents a use case for working with CEP (Postal Code).
type CepUsecase struct {
    repo entity.CepRepository
}

// NewCepUsecase creates a new instance of CepUsecase.
func NewCepUsecase(repo entity.CepRepository) *CepUsecase {
    return &CepUsecase{repo: repo}
}

// GetLocalityByCep retrieves the locality information for a given CEP (Postal Code).
func (u *CepUsecase) GetLocalityByCep(ctx context.Context, cep string, tracer trace.Tracer ) (*entity.Cep, error) {
    return u.repo.GetLocalityByCep(ctx, cep, tracer)
}

// ValidateCep validates a CEP (Postal Code).
func (u *CepUsecase) ValidateCep(cep string) bool {
    if len(cep) != 8 {
        return false
    }
    return true
}