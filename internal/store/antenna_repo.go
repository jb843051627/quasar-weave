package store

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jb843051627/quasar-weave/internal/model"
)

func (s *Store) SaveAntenna(ctx context.Context, antenna model.Antenna) error {
	if err := validateEntity("antenna", antenna.ID); err != nil {
		return err
	}
	if !antenna.Status.Valid() {
		return fmt.Errorf("invalid antenna status %q", antenna.Status)
	}
	return s.Save(ctx, kindAntenna, antenna.ID, antenna)
}

func (s *Store) GetAntenna(ctx context.Context, id string) (model.Antenna, error) {
	return LoadJSON[model.Antenna](ctx, s, kindAntenna, id)
}

func (s *Store) ListAntennas(ctx context.Context) ([]model.Antenna, error) {
	items, err := listKind[model.Antenna](ctx, s, kindAntenna)
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	return items, err
}

func (s *Store) DeleteAntenna(ctx context.Context, id string) error {
	return s.Delete(ctx, kindAntenna, id)
}

func (s *Store) TouchAntenna(ctx context.Context, id string, status model.AntennaStatus, now time.Time) (model.Antenna, error) {
	antenna, err := s.GetAntenna(ctx, id)
	if err != nil {
		return antenna, err
	}
	antenna.Status = status
	antenna.LastHeartbeat = now
	antenna.UpdatedAt = now
	if err := s.SaveAntenna(ctx, antenna); err != nil {
		return antenna, err
	}
	return antenna, nil
}
