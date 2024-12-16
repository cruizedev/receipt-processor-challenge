package process

import (
	"errors"
	"log"
	"math"
	"strings"
	"sync"
	"unicode"

	"github.com/cruizedev/receipt-processor-challenge/internal/model"
	"github.com/google/uuid"
)

type Processor struct {
	receipts map[string]int
	lock     sync.RWMutex
}

var (
	ErrItemNotFound        = errors.New("item with the given id does not exist")
	ErrItemStillInProgress = errors.New("item with the given id is still processing")
)

func New() *Processor {
	return &Processor{
		receipts: make(map[string]int),
		lock:     sync.RWMutex{},
	}
}

func (p *Processor) New(r model.Receipt) string {
	p.lock.Lock()
	defer p.lock.Unlock()

	id := uuid.New().String()
	p.receipts[id] = -1

	go p.process(r, id)

	return id
}

func (p *Processor) process(r model.Receipt, id string) {
	score := 0

	for _, c := range r.Retailer {
		if unicode.IsLetter(c) {
			score += 1
		}
	}
	log.Printf("retailer name score %d", score)

	if math.Round(r.Total) == r.Total {
		score += 50
	}
	log.Printf("total is integer score %d", score)

	if math.Round(r.Total*4) == r.Total {
		score += 25
	}
	log.Printf("total is dividable by 0.25 score %d", score)

	score += (len(r.Items) / 2) * 5
	log.Printf("items pairs score %d", score)

	if r.PurchaseTime.Day()%2 == 1 {
		score += 6
	}
	log.Printf("day is odd score %d", score)

	if r.PurchaseTime.Hour() > 14 && r.PurchaseTime.Hour() < 16 {
		score += 10
	}
	log.Printf("hour is in 14 - 16 score %d", score)

	for _, item := range r.Items {
		if len(strings.TrimSpace(item.ShortDescription))%3 == 0 {
			score += int(math.Ceil(item.Price * 0.2))
		}
		log.Printf("item description score %d", score)
	}

	p.lock.Lock()
	defer p.lock.Unlock()

	p.receipts[id] = score
}

func (p *Processor) Get(id string) (int, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	score, ok := p.receipts[id]
	if !ok {
		return 0, ErrItemNotFound
	}

	if score == -1 {
		return 0, ErrItemStillInProgress
	}

	return score, nil
}
