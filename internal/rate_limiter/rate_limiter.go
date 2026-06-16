package rate_limiter

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type RateLimiter struct {
	mu      sync.Mutex
	request map[int]RequestUserInfo
}

type RequestUserInfo struct {
	Counter     int
	RequestedAt time.Time
}

func New() *RateLimiter {
	return &RateLimiter{
		request: make(map[int]RequestUserInfo),
	}
}

func (r *RateLimiter) Allow(ctx context.Context, id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	info, ok := r.request[id]
	if !ok {
		r.request[id] = RequestUserInfo{
			Counter:     1,
			RequestedAt: time.Now(),
		}
		return true
	}
	if info.Counter >= 5 {
		return false
	}
	info.Counter++
	r.request[id] = info
	return true
}

func (r *RateLimiter) WorkerClear(ctx context.Context) {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ctx.Done():
			fmt.Println("ctx is done")
			return
		case <-ticker.C:
			r.mu.Lock()
			for id, val := range r.request {
				if time.Since(val.RequestedAt) > time.Minute {
					delete(r.request, id)
				}
			}
			r.mu.Unlock()
		}
	}
}

/*
func main(){
	counter := 0
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done();
			counter++
		}()
	}
	wg.Wait()
	fmt.Println(counter) // ?
}
*/
/*
func handle(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id != "" {
		n, err := strconv.Atoi(id)
		if err == nil {
			u, err := db.Find(n)
			if err == nil {
				if u.Active {
					json.NewEncoder(w).Encode(u)
				} else {
					w.WriteHeader(403)
				}
			} else {
				w.WriteHeader(500)
			}
		} else {
			w.WriteHeader(400)
		}
	} else {
		w.WriteHeader(400)
	}
}

func handler (w http.HandlerFunc, r *http.Request) {
	n, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = service(n)
	if err != nil{
		if errors.Is(err, errs.ErrFromValidateID){
			http.Error(w, "error from validate", http.StatusBadRequest)
			return
		}
		fmt.Errorf(w, err.Error(), http.StatusInternalServerError)
	}
}

func service(id int) error{
	if id < 1{
		return errs.ErrFromValidateID
	}
	u, err := repo(id)
	if err != nil{
		return fmt.Errorf("error from repo")
	}

	if u.IsActive {

	} else{
		re
	}
}

func repo(id int) (models.User, error){
	const query = `SELECT name, isActive FROM users WHERE id = $1`

	var u models.User

	err := db.QueryRow(query, id).Scan(&u.Name, &u.IsActive)
	if err != nil{
		return models.User{}, fmt.Errorf("not found %w", err)
	}

	return u, nil
}
*/
/*
func main() {
	// ParseRange("1-5") -> 1, 5, nil
}

func ParseRange(s string) (start, end int, err error) {
	parts := strings.Split(s, "-")
	start, _ = strconv.Atoi(parts[0])
	end, _ = strconv.Atoi(parts[1])
	return start, end, nil
}
*/
