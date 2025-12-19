package main

import (
	"fmt"
	"time"
	"goroutines-example/semaphore" // импорт пакета семафора
)

func main() {
	// Создаем счетный семафор с максимальным количеством разрешений 3
	// и таймаутом ожидания 5 секунд
	sem := semaphore.NewCountingSemaphore(3, 5*time.Second)

	fmt.Printf("Создан счетный семафор с максимальным количеством разрешений: %d\n", 3)
	fmt.Printf("Доступно разрешений: %d\n", sem.AvailablePermits())

	// Демонстрируем захват одного разрешения
	fmt.Println("\n--- Демонстрация захвата одного разрешения ---")
	
	// Захватываем 2 разрешения
	for i := 0; i < 2; i++ {
		err := sem.Acquire()
		if err != nil {
			fmt.Printf("Ошибка при захвате разрешения %d: %v\n", i+1, err)
		} else {
			fmt.Printf("Успешно захвачено разрешение %d. Осталось доступных: %d\n", i+1, sem.AvailablePermits())
		}
	}

	fmt.Printf("После захвата 2 разрешений доступно: %d\n", sem.AvailablePermits())

	// Пытаемся захватить разрешение без блокировки
	fmt.Println("\n--- Демонстрация метода TryAcquire ---")
	if sem.TryAcquire() {
		fmt.Printf("Успешно захвачено разрешение с помощью TryAcquire. Осталось доступных: %d\n", sem.AvailablePermits())
	} else {
		fmt.Println("Не удалось захватить разрешение с помощью TryAcquire - нет доступных разрешений")
	}

	if sem.TryAcquire() {
		fmt.Printf("Успешно захвачено разрешение с помощью TryAcquire. Осталось доступных: %d\n", sem.AvailablePermits())
	} else {
		fmt.Println("Не удалось захватить разрешение с помощью TryAcquire - нет доступных разрешений")
	}

	fmt.Printf("После попыток захвата доступно: %d\n", sem.AvailablePermits())

	// Освобождаем 2 разрешения
	fmt.Println("\n--- Демонстрация освобождения разрешений ---")
	for i := 0; i < 2; i++ {
		err := sem.Release()
		if err != nil {
			fmt.Printf("Ошибка при освобождении разрешения %d: %v\n", i+1, err)
		} else {
			fmt.Printf("Успешно освобождено разрешение %d. Осталось доступных: %d\n", i+1, sem.AvailablePermits())
		}
	}

	// Демонстрация захвата и освобождения нескольких разрешений
	fmt.Println("\n--- Демонстрация методов AcquireN и ReleaseN ---")
	
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

	// Демонстрируем ситуацию, когда пытаемся захватить больше разрешений, чем доступно
	fmt.Printf("\nПытаемся захватить 5 разрешений при максимальном количестве 3...\n")
	err = sem.AcquireN(5)
	if err != nil {
		fmt.Printf("Как и ожидалось, ошибка при захвате 5 разрешений: %v\n", err)
	} else {
		fmt.Printf("Успешно захвачено 5 разрешений. Осталось доступных: %d\n", sem.AvailablePermits())
		sem.ReleaseN(5)
	}

	// Проверим захват всех разрешений и их освобождение
	fmt.Println("\n--- Проверка захвата всех разрешений и освобождения ---")
	fmt.Printf("Доступно разрешений до захвата: %d\n", sem.AvailablePermits())
	
	err = sem.AcquireN(3) // Захватываем все 3 разрешения
	if err != nil {
		fmt.Printf("Ошибка при захвате 3 разрешений: %v\n", err)
	} else {
		fmt.Printf("Успешно захвачено 3 разрешения. Осталось доступных: %d\n", sem.AvailablePermits())
		
		// Проверим, что TryAcquire не сработает
		if sem.TryAcquire() {
			fmt.Println("НЕПРАВИЛЬНО: TryAcquire сработал, хотя не должно быть доступных разрешений")
		} else {
			fmt.Println("Правильно: TryAcquire не сработал, так как нет доступных разрешений")
		}
		
		// Освобождаем все разрешения
		sem.ReleaseN(3)
		fmt.Printf("Освобождены 3 разрешения. Теперь доступно: %d\n", sem.AvailablePermits())
	}

	fmt.Println("\nДемонстрация работы счетного семафора завершена.")
	fmt.Printf("Финальное количество доступных разрешений: %d\n", sem.AvailablePermits())
}