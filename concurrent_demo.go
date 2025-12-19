package main

import (
	"fmt"
	"sync"
	"time"
	"goroutines-example/semaphore" // импорт пакета семафора
)

func main() {
	// Создаем счетный семафор с максимальным количеством разрешений 2
	// и таймаутом ожидания 5 секунд
	sem := semaphore.NewCountingSemaphore(2, 5*time.Second)

	fmt.Printf("Создан счетный семафор с максимальным количеством разрешений: %d\n", 2)
	fmt.Printf("Доступно разрешений: %d\n", sem.AvailablePermits())

	var wg sync.WaitGroup

	// Запускаем 5 горутин с небольшой задержкой между запусками
	// чтобы лучше видеть, как работает семафор
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			fmt.Printf("Горутина %d пытается захватить разрешение...\n", id)

			// Пытаемся захватить разрешение
			err := sem.Acquire()
			if err != nil {
				fmt.Printf("Горутина %d не смогла захватить разрешение: %v\n", id, err)
				return
			}

			fmt.Printf("Горутина %d захватила разрешение. Осталось доступных: %d\n", id, sem.AvailablePermits())

			// Имитируем работу с ресурсом
			time.Sleep(2 * time.Second)

			// Освобождаем разрешение
			sem.Release()
			fmt.Printf("Горутина %d освободила разрешение. Осталось доступных: %d\n", id, sem.AvailablePermits())
		}(i)
		
		// Небольшая задержка перед запуском следующей горутины
		time.Sleep(300 * time.Millisecond)
	}

	// Ждем завершения всех горутин
	wg.Wait()

	// После завершения всех горутин, ждем, пока освободятся все разрешения
	time.Sleep(3 * time.Second)
	fmt.Printf("\nПосле завершения всех горутин доступно разрешений: %d\n", sem.AvailablePermits())

	fmt.Println("\n--- Демонстрация метода TryAcquire ---")

	// Показываем использование TryAcquire
	for i := 0; i < 3; i++ {
		if sem.TryAcquire() {
			fmt.Printf("Успешно захвачено разрешение с помощью TryAcquire (попытка %d)\n", i+1)
		} else {
			fmt.Printf("Не удалось захватить разрешение с помощью TryAcquire (попытка %d) - нет доступных разрешений\n", i+1)
		}
	}

	// Проверяем, сколько разрешений мы можем освободить обратно
	availableAfterTry := sem.AvailablePermits()
	sem.ReleaseN(availableAfterTry)
	fmt.Printf("Все разрешения освобождены. Теперь доступно: %d\n", sem.AvailablePermits())

	fmt.Println("\n--- Демонстрация методов AcquireN и ReleaseN ---")

	// Демонстрация захвата нескольких разрешений
	fmt.Printf("Пытаемся захватить 2 разрешения сразу...\n")
	err := sem.AcquireN(2)
	if err != nil {
		fmt.Printf("Ошибка при захвате 2 разрешений: %v\n", err)
	} else {
		fmt.Printf("Успешно захвачено 2 разрешения. Осталось доступных: %d\n", sem.AvailablePermits())
		
		// Освобождаем захваченные разрешения
		sem.ReleaseN(2)
		fmt.Printf("Освобождены 2 разрешения. Теперь доступно: %d\n", sem.AvailablePermits())
	}

	fmt.Println("Демонстрация работы счетного семафора завершена.")
}