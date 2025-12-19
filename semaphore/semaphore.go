package semaphore

import (
	"fmt"
	"sync"
	"time"
)

// CountingSemaphore — структура счетного семафора
// В отличие от двоичного семафора, счетный может иметь значение больше 1,
// что позволяет контролировать доступ к нескольким одинаковым ресурсам
type CountingSemaphore struct {
	// Канал для хранения состояния семафора
	sem chan struct{}
	// Максимальное количество разрешений
	maxPermits int
	// Текущее количество доступных разрешений
	currentPermits int
	// Защита переменной currentPermits при многопоточном доступе
	mutex sync.RWMutex
	// Время ожидания основных операций с семафором, чтобы не 
	// блокировать операции с ним навечно
	timeout time.Duration
}

// Acquire — метод захвата одного разрешения у семафора
// Уменьшает счетчик доступных разрешений на 1
func (cs *CountingSemaphore) Acquire() error {
	select {
	case _ = <-cs.sem:
		cs.mutex.Lock()
		defer cs.mutex.Unlock()
		cs.currentPermits--
		return nil
	case <-time.After(cs.timeout):
		return fmt.Errorf("Не удалось захватить разрешение у семафора")
	}
}

// TryAcquire — метод попытки захвата разрешения без блокировки
// Возвращает true, если удалось захватить разрешение, иначе false
func (cs *CountingSemaphore) TryAcquire() bool {
	select {
	case _ = <-cs.sem:
		cs.mutex.Lock()
		defer cs.mutex.Unlock()
		cs.currentPermits--
		return true
	default:
		return false
	}
}

// Release — метод освобождения одного разрешения у семафора
// Увеличивает счетчик доступных разрешений на 1
func (cs *CountingSemaphore) Release() error {
	select {
	case cs.sem <- struct{}{}:
		cs.mutex.Lock()
		defer cs.mutex.Unlock()
		cs.currentPermits++
		return nil
	case <-time.After(cs.timeout):
		return fmt.Errorf("Не удалось освободить разрешение у семафора")
	}
}

// AvailablePermits — метод получения количества доступных разрешений
func (cs *CountingSemaphore) AvailablePermits() int {
	cs.mutex.RLock()
	defer cs.mutex.RUnlock()
	return cs.currentPermits
}

// AcquireN — метод захвата N разрешений у семафора
// Важно: для корректной работы с несколькими разрешениями используйте
// эту функцию вместо вызова Acquire несколько раз
func (cs *CountingSemaphore) AcquireN(n int) error {
	if n > cs.maxPermits {
		return fmt.Errorf("запрошено больше разрешений (%d), чем максимально доступно (%d)", n, cs.maxPermits)
	}

	// Проверяем, достаточно ли доступных разрешений
	if cs.AvailablePermits() < n {
		return fmt.Errorf("недостаточно разрешений: доступно %d, требуется %d", cs.AvailablePermits(), n)
	}

	for i := 0; i < n; i++ {
		err := cs.Acquire()
		if err != nil {
			// Если не удалось получить все разрешения, возвращаем уже захваченные
			cs.ReleaseN(i)
			return err
		}
	}
	return nil
}

// ReleaseN — метод освобождения N разрешений у семафора
func (cs *CountingSemaphore) ReleaseN(n int) error {
	cs.mutex.RLock()
	availableToRelease := cs.maxPermits - cs.currentPermits
	cs.mutex.RUnlock()

	if n > availableToRelease {
		return fmt.Errorf("попытка освободить больше разрешений (%d), чем захвачено (%d)", n, availableToRelease)
	}

	for i := 0; i < n; i++ {
		err := cs.Release()
		if err != nil {
			return err
		}
	}
	return nil
}

// NewCountingSemaphore — функция создания счетного семафора
// initialPermits — начальное количество разрешений (должно быть <= maxPermits)
func NewCountingSemaphore(maxPermits int, timeout time.Duration) *CountingSemaphore {
	sem := make(chan struct{}, maxPermits)
	
	// Заполняем канал начальными разрешениями
	for i := 0; i < maxPermits; i++ {
		sem <- struct{}{}
	}

	return &CountingSemaphore{
		sem:            sem,
		maxPermits:     maxPermits,
		currentPermits: maxPermits,
		timeout:        timeout,
	}
}