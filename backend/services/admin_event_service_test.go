package services

import (
	"errors"
	"testing"
	"time"

	"backend/domain"

	"gorm.io/gorm"
)

type fakeAdminEventRepository struct {
	event        *domain.Event
	createErr    error
	getErr       error
	saveErr      error
	createdEvent *domain.Event
	savedEvent   *domain.Event
}

func (repo *fakeAdminEventRepository) Create(event *domain.Event) error {
	repo.createdEvent = event
	return repo.createErr
}

func (repo *fakeAdminEventRepository) GetByIDForAdmin(id uint) (*domain.Event, error) {
	if repo.getErr != nil {
		return nil, repo.getErr
	}
	return repo.event, nil
}

func (repo *fakeAdminEventRepository) Save(event *domain.Event) error {
	repo.savedEvent = event
	return repo.saveErr
}

type fakeEventReportRepository struct {
	total        int64
	active       int64
	cancelled    int64
	err          error
	statusErr    error
	requestedIDs []uint
}

func (repo *fakeEventReportRepository) CountTicketsByEvent(eventID uint) (int64, error) {
	repo.requestedIDs = append(repo.requestedIDs, eventID)
	return repo.total, repo.err
}

func (repo *fakeEventReportRepository) CountTicketsByEventAndStatus(eventID uint, status string) (int64, error) {
	if repo.statusErr != nil {
		return 0, repo.statusErr
	}
	if status == domain.TicketStatusActive {
		return repo.active, nil
	}
	return repo.cancelled, nil
}

func newAdminEventTestService(repo *fakeAdminEventRepository, tickets *fakeEventReportRepository) *AdminEventService {
	service := NewAdminEventService(repo, tickets)
	service.now = func() time.Time {
		return time.Date(2026, 7, 12, 12, 0, 0, 0, time.UTC)
	}
	return service
}

func validCreateEventInput() CreateEventInput {
	return CreateEventInput{
		Title:       "Evento",
		Description: "Descripcion",
		Date:        time.Date(2026, 12, 10, 20, 0, 0, 0, time.UTC),
		Location:    "Cordoba",
		Capacity:    100,
		Price:       25000,
		ImageURL:    "https://example.com/event.jpg",
	}
}

func TestAdminCreateEventValid(t *testing.T) {
	repo := &fakeAdminEventRepository{}
	service := newAdminEventTestService(repo, &fakeEventReportRepository{})

	event, err := service.CreateEvent(validCreateEventInput())
	if err != nil {
		t.Fatalf("CreateEvent returned error: %v", err)
	}

	if event.Status != domain.EventStatusActive {
		t.Fatalf("expected active status, got %s", event.Status)
	}

	if event.AvailableCapacity != event.Capacity {
		t.Fatalf("expected available capacity %d, got %d", event.Capacity, event.AvailableCapacity)
	}
}

func TestAdminCreateEventInvalidInputs(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*CreateEventInput)
	}{
		{name: "empty title", mutate: func(input *CreateEventInput) { input.Title = " " }},
		{name: "empty location", mutate: func(input *CreateEventInput) { input.Location = " " }},
		{name: "past date", mutate: func(input *CreateEventInput) { input.Date = time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC) }},
		{name: "zero capacity", mutate: func(input *CreateEventInput) { input.Capacity = 0 }},
		{name: "negative capacity", mutate: func(input *CreateEventInput) { input.Capacity = -1 }},
		{name: "negative price", mutate: func(input *CreateEventInput) { input.Price = -1 }},
		{name: "invalid image url", mutate: func(input *CreateEventInput) { input.ImageURL = "ftp://example.com/image.jpg" }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := validCreateEventInput()
			tt.mutate(&input)
			service := newAdminEventTestService(&fakeAdminEventRepository{}, &fakeEventReportRepository{})

			_, err := service.CreateEvent(input)
			if !errors.Is(err, ErrInvalidEventData) {
				t.Fatalf("expected ErrInvalidEventData, got %v", err)
			}
		})
	}
}

func TestAdminCreateEventRepositoryError(t *testing.T) {
	repoErr := errors.New("repo error")
	service := newAdminEventTestService(&fakeAdminEventRepository{createErr: repoErr}, &fakeEventReportRepository{})

	_, err := service.CreateEvent(validCreateEventInput())
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}

func TestAdminUpdateEventValid(t *testing.T) {
	title := "Nuevo titulo"
	repo := &fakeAdminEventRepository{
		event: &domain.Event{ID: 1, Title: "Old", Date: validCreateEventInput().Date, Location: "Old", Capacity: 100, AvailableCapacity: 80, Status: domain.EventStatusActive},
	}
	service := newAdminEventTestService(repo, &fakeEventReportRepository{})

	event, err := service.UpdateEvent(UpdateEventInput{ID: 1, Title: &title})
	if err != nil {
		t.Fatalf("UpdateEvent returned error: %v", err)
	}

	if event.Title != title {
		t.Fatalf("expected updated title, got %s", event.Title)
	}
}

func TestAdminUpdateCapacityPreservesIssuedTickets(t *testing.T) {
	newCapacity := 150
	repo := &fakeAdminEventRepository{
		event: &domain.Event{ID: 1, Capacity: 100, AvailableCapacity: 80, Status: domain.EventStatusActive},
	}
	service := newAdminEventTestService(repo, &fakeEventReportRepository{})

	event, err := service.UpdateEvent(UpdateEventInput{ID: 1, Capacity: &newCapacity})
	if err != nil {
		t.Fatalf("UpdateEvent returned error: %v", err)
	}

	if event.AvailableCapacity != 130 {
		t.Fatalf("expected available capacity 130, got %d", event.AvailableCapacity)
	}
}

func TestAdminUpdateCapacityLowerThanIssuedReturnsConflict(t *testing.T) {
	newCapacity := 10
	repo := &fakeAdminEventRepository{
		event: &domain.Event{ID: 1, Capacity: 100, AvailableCapacity: 80, Status: domain.EventStatusActive},
	}
	service := newAdminEventTestService(repo, &fakeEventReportRepository{})

	_, err := service.UpdateEvent(UpdateEventInput{ID: 1, Capacity: &newCapacity})
	if !errors.Is(err, ErrEventCapacityConflict) {
		t.Fatalf("expected ErrEventCapacityConflict, got %v", err)
	}
}

func TestAdminUpdateMissingEventReturnsError(t *testing.T) {
	service := newAdminEventTestService(&fakeAdminEventRepository{getErr: gorm.ErrRecordNotFound}, &fakeEventReportRepository{})

	_, err := service.UpdateEvent(UpdateEventInput{ID: 1})
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestAdminUpdateCancelledEventReturnsConflict(t *testing.T) {
	repo := &fakeAdminEventRepository{
		event: &domain.Event{ID: 1, Status: domain.EventStatusCancelled},
	}
	service := newAdminEventTestService(repo, &fakeEventReportRepository{})

	_, err := service.UpdateEvent(UpdateEventInput{ID: 1})
	if !errors.Is(err, ErrEventCancelled) {
		t.Fatalf("expected ErrEventCancelled, got %v", err)
	}
}

func TestAdminCancelEventSuccess(t *testing.T) {
	repo := &fakeAdminEventRepository{
		event: &domain.Event{ID: 1, Status: domain.EventStatusActive, AvailableCapacity: 10},
	}
	service := newAdminEventTestService(repo, &fakeEventReportRepository{})

	event, err := service.CancelEvent(1)
	if err != nil {
		t.Fatalf("CancelEvent returned error: %v", err)
	}

	if event.Status != domain.EventStatusCancelled {
		t.Fatalf("expected cancelled status, got %s", event.Status)
	}

	if event.AvailableCapacity != 10 {
		t.Fatalf("expected available capacity unchanged, got %d", event.AvailableCapacity)
	}
}

func TestAdminCancelMissingEventReturnsError(t *testing.T) {
	service := newAdminEventTestService(&fakeAdminEventRepository{getErr: gorm.ErrRecordNotFound}, &fakeEventReportRepository{})

	_, err := service.CancelEvent(1)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestAdminCancelAlreadyCancelledReturnsConflict(t *testing.T) {
	repo := &fakeAdminEventRepository{event: &domain.Event{ID: 1, Status: domain.EventStatusCancelled}}
	service := newAdminEventTestService(repo, &fakeEventReportRepository{})

	_, err := service.CancelEvent(1)
	if !errors.Is(err, ErrEventCancelled) {
		t.Fatalf("expected ErrEventCancelled, got %v", err)
	}
}

func TestAdminEventReportSuccess(t *testing.T) {
	repo := &fakeAdminEventRepository{
		event: &domain.Event{ID: 1, Title: "Evento", Status: domain.EventStatusActive, Capacity: 100, AvailableCapacity: 80, Price: 25000},
	}
	service := newAdminEventTestService(repo, &fakeEventReportRepository{total: 25, active: 18, cancelled: 7})

	report, err := service.GetEventReport(1)
	if err != nil {
		t.Fatalf("GetEventReport returned error: %v", err)
	}

	if report.TicketsIssued != 25 || report.ActiveTickets != 18 || report.CancelledTickets != 7 {
		t.Fatalf("unexpected report counts: %+v", report)
	}

	if report.OccupancyPercentage != 18 {
		t.Fatalf("expected occupancy 18, got %f", report.OccupancyPercentage)
	}

	if report.EstimatedRevenue != 450000 {
		t.Fatalf("expected revenue 450000, got %f", report.EstimatedRevenue)
	}
}

func TestAdminEventReportWithZeroCapacityAvoidsDivisionByZero(t *testing.T) {
	repo := &fakeAdminEventRepository{
		event: &domain.Event{ID: 1, Capacity: 0, Price: 10},
	}
	service := newAdminEventTestService(repo, &fakeEventReportRepository{total: 1, active: 1})

	report, err := service.GetEventReport(1)
	if err != nil {
		t.Fatalf("GetEventReport returned error: %v", err)
	}

	if report.OccupancyPercentage != 0 {
		t.Fatalf("expected occupancy 0, got %f", report.OccupancyPercentage)
	}
}

func TestAdminEventReportRepositoryError(t *testing.T) {
	repoErr := errors.New("repo error")
	repo := &fakeAdminEventRepository{event: &domain.Event{ID: 1}}
	service := newAdminEventTestService(repo, &fakeEventReportRepository{err: repoErr})

	_, err := service.GetEventReport(1)
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
}
